package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
)

func init() {
	// Set default log config
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp:       false,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})

	// Setup configuration file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	// Initialize configuration default
	viper.SetDefault("http.host", "127.0.0.1")
	viper.SetDefault("http.port", 8080)
	viper.SetDefault("http.domain", "http://127.0.0.1:8080")
	viper.SetDefault("http.reset_files", false)
	viper.SetDefault("logging.format", "text")
	viper.SetDefault("logging.level", "info")
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
			logrus.WithField("app", "initialization").WithError(err).Fatal("Configuration file not found")
		default:
			logrus.WithField("app", "initialization").WithError(err).Fatal("Failed to parse configuration file")
		}
	}

	// Validate ssl mode
	if mode := viper.GetString("database.ssl"); mode != "disable" && mode != "allow" && mode != "prefer" && mode != "require" && mode != "verify-ca" && mode != "verify-full" {
		logrus.WithFields(logrus.Fields{"app": "initialization", "key": "database.ssl", "value": viper.GetString("database.ssl"), "options": []string{"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}}).Fatal("Invalid value for database ssl mode")
	}

	// Delete all uploaded files
	if viper.GetBool("http.reset_files") {
		if err := os.RemoveAll("./uploaded"); err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{"app": "initialization", "key": "http.reset_files"}).Fatal("Failed to reset uploaded files")
		}
	}

	// Create directory if not exist
	if _, err := os.Stat("./uploaded"); os.IsNotExist(err) {
		if err := os.Mkdir("./uploaded", os.ModePerm); err != nil {
			log.Printf("Failed to create uploads directory: %v", err)
			logrus.WithField("app", "initialization").WithError(err).Fatal("Failed to create uploads directory")
		}
	} else if err != nil {
		logrus.WithField("app", "initialization").WithError(err).Fatal("Failed to stat uploads directory")
	}

	// Validate and set log format
	if format := viper.GetString("logging.format"); format != "text" && format != "json" {
		logrus.WithFields(logrus.Fields{"app": "initialization", "key": "logging.format", "value": format, "options": []string{"text", "json"}}).Fatal("Invalid output format specified")
	} else if format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else if format == "text" {
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableTimestamp:       false,
			FullTimestamp:          true,
			DisableLevelTruncation: true,
		})
	}

	// Validate and set log level
	switch viper.GetString("logging.level") {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.WithFields(logrus.Fields{"app": "initialization", "key": "logging.level", "value": viper.GetString("logging.level"), "options": []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}}).Fatal("Invalid level for minimum log level")
	}
}
