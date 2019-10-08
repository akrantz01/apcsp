package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func init() {
	// Set default log config
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(false)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp:       false,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})
	logrus.WithField("app", "initialization").Trace("Configured base logger settings")

	// Setup configuration file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	logrus.WithField("app", "initialization").Trace("Set configuration file path and name")

	// Initialize configuration default
	viper.SetDefault("http.host", "127.0.0.1")
	viper.SetDefault("http.port", 8080)
	viper.SetDefault("http.domain", "http://127.0.0.1:8080")
	viper.SetDefault("http.reset_files", false)
	viper.SetDefault("email.host", "127.0.0.1")
	viper.SetDefault("email.port", 25)
	viper.SetDefault("email.ssl", false)
	viper.SetDefault("email.username", "postmaster")
	viper.SetDefault("email.password", "postmaster")
	viper.SetDefault("logging.format", "text")
	viper.SetDefault("logging.report_caller", false)
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.username", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.database", "postgres")
	viper.SetDefault("database.ssl", "disable")
	viper.SetDefault("database.reset", false)
	logrus.WithField("app", "initialization").Trace("Set defaults for configuration keys")

	// Allow loading config from environment variables
	viper.SetEnvPrefix("chat")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Parse configuration file
	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
		case *os.PathError:
			logrus.WithField("app", "initialization").Info("Configuration file not found, continuing with defaults and environment variables")
		default:
			logrus.WithField("app", "initialization").WithError(err).Fatal("Failed to parse configuration file")
		}
	}

	// Validate ssl mode
	if mode := viper.GetString("database.ssl"); mode != "disable" && mode != "allow" && mode != "prefer" && mode != "require" && mode != "verify-ca" && mode != "verify-full" {
		logrus.WithFields(logrus.Fields{"app": "initialization", "key": "database.ssl", "value": viper.GetString("database.ssl"), "options": []string{"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}}).Fatal("Invalid value for database ssl mode")
	}
	logrus.WithField("app", "initialization").Trace("Validated database ssl connection mode")

	// Delete all uploaded files
	if viper.GetBool("http.reset_files") {
		if err := os.RemoveAll("./uploaded"); err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{"app": "initialization", "key": "http.reset_files"}).Fatal("Failed to reset uploaded files")
		}
		logrus.WithField("app", "initialization").Info("Successfully deleted all uploaded files")
	}
	logrus.WithField("app", "initialization").Trace("Optionally deleted all uploaded files")

	// Create directory if not exist
	if _, err := os.Stat("./uploaded"); os.IsNotExist(err) {
		if err := os.Mkdir("./uploaded", os.ModePerm); err != nil {
			logrus.WithField("app", "initialization").WithError(err).Fatal("Failed to create uploads directory")
		}
		logrus.WithField("app", "initialization").Debug("Created uploads directory")
	} else if err != nil {
		logrus.WithField("app", "initialization").WithError(err).Fatal("Failed to stat uploads directory")
	}
	logrus.WithField("app", "initialization").Trace("Created uploads directory if it did not exist")

	// Validate and set log format
	if format := viper.GetString("logging.format"); format != "text" && format != "json" {
		logrus.WithFields(logrus.Fields{"app": "initialization", "key": "logging.format", "value": format, "options": []string{"text", "json"}}).Fatal("Invalid output format specified")
	} else if format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.WithField("app", "initialization").Debug("Set output format to JSON")
	} else if format == "text" {
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableTimestamp:       false,
			FullTimestamp:          true,
			DisableLevelTruncation: true,
		})
		logrus.WithField("app", "initialization").Debug("Set output format to text")
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
	logrus.WithField("app", "initialization").Tracef("Set minimum log level to %s", viper.GetString("logging.level"))

	// Set reporting caller
	if viper.GetBool("logging.report_caller") {
		logrus.SetReportCaller(viper.GetBool("logging.report_caller"))
		logrus.WithField("app", "initialization").Trace("Enabled caller reporting")
	}
}
