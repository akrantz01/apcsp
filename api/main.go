package main

import (
	"context"
	"github.com/akrantz01/apcsp/api/authentication"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/users"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"log"
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

	// API sub-router
	api := router.PathPrefix("/api").Subrouter()

	// Authentication routes
	api.HandleFunc("/auth/login", authentication.Login(db)).Methods("POST")
	api.HandleFunc("/auth/logout", authentication.Logout(db)).Methods("GET")

	// User routes
	api.HandleFunc("/users", users.AllUsers(db)).Methods("POST")
	api.HandleFunc("/users/{user}", users.SpecificUser(db)).Methods("GET", "PUT", "DELETE")

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
		log.Printf("API listening on %s:%s...", viper.GetString("http.host"), viper.GetString("http.port"))
		if err := server.ListenAndServe(); err != nil {
			if strings.Contains(err.Error(), "Server closed") {
				return
			}
			log.Fatalf("failure while running server: %s", err)
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
		log.Fatalf("Failed to shutdown server: %v", err)
	}
}
