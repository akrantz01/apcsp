package main

import (
	"bytes"
	"context"
	"github.com/akrantz01/apcsp/api/authentication"
	"github.com/akrantz01/apcsp/api/chats"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/files"
	"github.com/akrantz01/apcsp/api/messages"
	"github.com/akrantz01/apcsp/api/users"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/akrantz01/apcsp/api/websockets"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var shutdown = make(chan os.Signal, 1)
var mail = make(chan *gomail.Message)

func main() {
	logger := logrus.WithField("app", "main")

	// Initialize file embedding
	box := packr.New("static", "./static")

	// Connect to the database
	db := database.SetupDatabase()

	// Create websocket hub
	hub := websockets.NewHub()
	logger.Trace("Created websocket hub for connection management")

	// Setup routes
	router := mux.NewRouter()
	logger.Trace("Created HTTP router")

	// Enable authentication middleware
	router.Use(authMiddleware(db))
	logger.Trace("Initialized API authentication middleware")

	// API sub-router
	api := router.PathPrefix("/api").Subrouter()
	logger.Trace("Initialized API subrouter")

	// Authentication routes
	api.HandleFunc("/auth/login", authentication.Login(db))
	api.HandleFunc("/auth/logout", authentication.Logout(db))
	api.HandleFunc("/auth/forgot-password", authentication.ForgotPassword(db, mail, box))
	api.HandleFunc("/auth/reset-password", authentication.ResetPassword(db, mail, box))
	api.HandleFunc("/auth/verify-email", authentication.VerifyEmail(db))
	logger.Trace("Add authentication routes")

	// User routes
	api.HandleFunc("/users", users.AllUsers(db, mail, box))
	api.HandleFunc("/users/{user}", users.SpecificUser(db))
	logger.Trace("Add user management routes")

	// Chat routes
	api.HandleFunc("/chats", chats.AllChats(db))
	api.HandleFunc("/chats/{chat}", chats.SpecificChat(db))
	logger.Trace("Add chat management routes")

	// Messages routes
	api.HandleFunc("/chats/{chat}/messages", messages.AllMessages(hub, db))
	api.HandleFunc("/chats/{chat}/messages/{message}", messages.SpecificMessage(db))
	logger.Trace("Add chat message management routes")

	// Files routes
	api.HandleFunc("/files/{file}", files.Files(hub, db))
	logger.Trace("Add file management routes")

	// Websocket routes
	api.HandleFunc("/ws", websockets.Websockets(hub, db))
	logger.Trace("Add websocket routes")

	// Add static HTML routes
	router.HandleFunc("/reset-password", func(w http.ResponseWriter, r *http.Request) {
		// Get file from box
		resetPassword, err := box.Find("reset-password.html")
		if err != nil {
			logrus.WithField("app", "static-files").WithError(err).Error("Failed to load file from box")
			util.Responses.Error(w, http.StatusInternalServerError, "failed to load file from box")
			return
		}
		buffer := bytes.NewBuffer(resetPassword)

		// Write headers
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.FormatInt(int64(len(resetPassword)), 10))
		w.WriteHeader(http.StatusOK)

		// Copy to client
		if _, err := io.Copy(w, buffer); err != nil {
			logrus.WithField("app", "static-files").WithError(err).Error("Failed to copy file data to client")
			util.Responses.Error(w, http.StatusInternalServerError, "failed to copy to client")
			return
		}
	})

	// Register router with http and enable cors
	http.Handle("/", loggingHandler{handler: cors.AllowAll().Handler(router)})
	logger.Trace("Register router with http handler")

	// Wait for OS shutdown signals
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	logger.Trace("Begin listening for SIGINT and SIGTERM")

	// Setup HTTP server
	server := &http.Server{
		Addr:         viper.GetString("http.host") + ":" + viper.GetString("http.port"),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      nil,
	}

	// Start mail daemon
	go func() {
		logger := logrus.WithField("app", "mailer")

		// Initialize mail server connector
		d := gomail.NewPlainDialer(viper.GetString("email.host"), viper.GetInt("email.port"), viper.GetString("email.username"), viper.GetString("email.password"))
		d.SSL = viper.GetBool("email.ssl")
		logger.Trace("Configured dialer")

		// Connect and create writer
		logger.WithFields(logrus.Fields{"host": viper.GetString("email.host"), "port": viper.GetInt("email.port"), "ssl": viper.GetBool("email.ssl")}).Info("Connecting to email server...")
		var s gomail.SendCloser
		var err error
		if s, err = d.Dial(); err != nil {
			logrus.WithError(err).Fatal("Failed to connect to email server")
		}
		logger.Info("Successfully connected to email server")

		open := true

		logger.Trace("Started mail handling loop")
		for {
			select {
			// Wait for mail
			case m, ok := <-mail:
				if !ok {
					logger.Trace("Invalid message received")
					return
				}
				logger.Trace("Got new message to send")

				// Connect if not open
				if !open {
					logger.Trace("Sender not open, opening...")
					if s, err = d.Dial(); err != nil {
						logrus.WithError(err).Fatal("Failed to re-connect to email server")
					}
					open = true
					logger.Trace("Successfully re-connected to mail server")
				}
				logger.Trace("Ensured sender connection was open")

				// Send the message
				logger.Trace("Sending message...")
				if err := gomail.Send(s, m); err != nil {
					logger.WithError(err).Error("Failed to send email")
				}
				logger.Trace("Message sent successfully")

			// Close connection after 4 seconds (fixes error with Amazon SES)
			case <-time.After(4 * time.Second):
				if open {
					logger.Trace("Connection open for more than 4 seconds, closing")
					if err := s.Close(); err != nil {
						logger.WithError(err).Fatal("Failed to close mail server connection")
					}
					open = false
					logger.Trace("Successfully closed sender")
				}
			}
		}
	}()

	// Start websocket server
	go hub.Run()
	logger.Trace("Started websocket server in separate goroutine")

	// Start http server
	go func() {
		logrus.WithFields(logrus.Fields{"app": "http-server", "host": viper.GetString("http.host"), "port": viper.GetInt("http.port")}).Info("Starting API listener...")
		if err := server.ListenAndServe(); err != nil {
			if strings.Contains(err.Error(), "Server closed") {
				return
			}
			logrus.WithError(err).WithField("app", "http-server").Fatal("Failure while server was running")
		}
	}()
	logger.Trace("Start server in separate goroutine")

	// Wait for server shutdown
	<-shutdown
	logger.Trace("Shutdown signal received. Shutdown sequence started")

	// Create shutdown context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	logger.Trace("Create shutdown context, forcing after 5 seconds")

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		logrus.WithError(err).WithField("app", "http-server").Fatal("Failed to shutdown server")
	}
	logrus.WithField("app", "http-server").Info("Gracefully shutdown API listener")
}
