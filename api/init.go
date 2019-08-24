package main

import (
	"github.com/spf13/viper"
	"log"
)

func init() {
	// Setup configuration file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	// Initialize configuration default
	viper.SetDefault("http.host", "127.0.0.1")
	viper.SetDefault("http.port", 8080)
	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.username", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.database", "postgres")
	viper.SetDefault("database.ssl", "disable")
	viper.SetDefault("database.reset", false)

	// Parse configuration file
	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Fatalf("Configuration file not found")
		default:
			log.Fatalf("Failed to parse configuration file: %s", err)
		}
	}

	// Validate ssl mode
	if mode := viper.GetString("database.ssl"); mode != "disable" && mode != "allow" && mode != "prefer" && mode != "require" && mode != "verify-ca" && mode != "verify-full" {
		log.Fatal("invalid value for ssl, must be one of: disable, allow, prefer, require, verify-ca, verify-full")
	}
}
