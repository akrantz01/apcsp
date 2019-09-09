package main

import (
	"context"
	"github.com/akrantz01/apcsp/api/authentication"
	"github.com/akrantz01/apcsp/api/chats"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/files"
	"github.com/akrantz01/apcsp/api/messages"
	"github.com/akrantz01/apcsp/api/users"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var shutdown = make(chan os.Signal, 1)

func main() {
	logger := logrus.WithField("app", "main")

	// Connect to the database
	db := database.SetupDatabase()

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
	logger.Trace("Add authentication routes")

	// User routes
	api.HandleFunc("/users", users.AllUsers(db))
	api.HandleFunc("/users/{user}", users.SpecificUser(db))
	logger.Trace("Add user management routes")

	// Chat routes
	api.HandleFunc("/chats", chats.AllChats(db))
	api.HandleFunc("/chats/{chat}", chats.SpecificChat(db))
	logger.Trace("Add chat management routes")

	// Messages routes
	api.HandleFunc("/chats/{chat}/messages", messages.AllMessages(db))
	api.HandleFunc("/chats/{chat}/messages/{message}", messages.SpecificMessage(db))
	logger.Trace("Add chat message management routes")

	// Files routes
	api.HandleFunc("/files/{file}", files.Files(db))
	logger.Trace("Add file management routes")

	// Register router with http and enable cors
	http.Handle("/", handlers.LoggingHandler(os.Stdout, cors.AllowAll().Handler(router)))
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
