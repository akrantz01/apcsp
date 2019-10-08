package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"reflect"
)

var logger = logrus.WithField("app", "database")

// Connect to the database and create the schema
func SetupDatabase() *gorm.DB {
	logger.WithFields(logrus.Fields{"host": viper.GetString("database.host"), "port": viper.GetInt("database.port"), "ssl": viper.GetString("database.ssl")}).Info("Connecting to database...")
	// Connect to database
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", viper.GetString("database.host"), viper.GetString("database.port"), viper.GetString("database.username"), viper.GetString("database.password"), viper.GetString("database.database"), viper.GetString("database.ssl")))
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	logger.Info("Successfully connected to database")

	logger.Info("Building database schema if nonexistent...")
	// Create schema if not exist
	for _, model := range []interface{}{&User{}, &Token{}, &Chat{}, &Message{}, &File{}} {
		if viper.GetBool("database.reset") {
			db.DropTableIfExists(model)
			db.CreateTable(model)
			logger.WithField("model", reflect.TypeOf(model).Name()).Trace("Dropped model if it existed and re-created in database")
		} else if !db.HasTable(model) {
			db.CreateTable(model)
			logger.WithField("model", reflect.TypeOf(model).Name()).Trace("Created model in database")
		}
	}
	logger.Info("Successfully built database schema if nonexistent")

	// Enable struct preloading (for relationships)
	db.Set("gorm:auto_preload", true)
	logger.Trace("Enable automatically preloading table relationships")

	return db
}
