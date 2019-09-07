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
	// Connect to the database
	db := database.SetupDatabase()

	// Setup routes
	router := mux.NewRouter()

	// Enable authentication middleware
	router.Use(authMiddleware(db))

	// API sub-router
	api := router.PathPrefix("/api").Subrouter()

	// Authentication routes
	api.HandleFunc("/auth/login", authentication.Login(db)).Methods("POST")
	api.HandleFunc("/auth/logout", authentication.Logout(db)).Methods("GET")

	// User routes
	api.HandleFunc("/users", users.AllUsers(db)).Methods("GET", "POST")
	api.HandleFunc("/users/{user}", users.SpecificUser(db)).Methods("GET", "PUT", "DELETE")

	// Chat routes
	api.HandleFunc("/chats", chats.AllChats(db)).Methods("GET", "POST")
	api.HandleFunc("/chats/{chat}", chats.SpecificChat(db)).Methods("GET", "DELETE")

	// Messages routes
	api.HandleFunc("/chats/{chat}/messages", messages.AllMessages(db)).Methods("GET", "POST")
	api.HandleFunc("/chats/{chat}/messages/{message}", messages.SpecificMessage(db)).Methods("GET", "PUT", "DELETE")

	// Files routes
	api.HandleFunc("/files/{file}", files.Files(db)).Methods("GET", "POST")

	// Register router with http and enable cors
	http.Handle("/", handlers.LoggingHandler(os.Stdout, cors.AllowAll().Handler(router)))

	// Wait for OS shutdown signals
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

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

	// Wait for server shutdown
	<-shutdown

	// Create shutdown context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		logrus.WithError(err).WithField("app", "http-server").Fatal("Failed to shutdown server")
	}
	logrus.WithField("app", "http-server").Info("Gracefully shutdown API listener")
}
