package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
	"log"
)

// Connect to the database and create the schema
func SetupDatabase() *gorm.DB {
	log.Println("Connecting to database...")
	// Connect to database
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", viper.GetString("database.host"), viper.GetString("database.port"), viper.GetString("database.username"), viper.GetString("database.password"), viper.GetString("database.database"), viper.GetString("database.ssl")))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")

	log.Println("Building database schema...")
	// Create schema if not exist
	for _, model := range []interface{}{} {
		if viper.GetBool("database.reset") {
			db.DropTableIfExists(model)
			db.CreateTable(model)
		} else if !db.HasTable(model) {
			db.CreateTable(model)
		}
	}
	log.Println("Built database schema")

	return db
}
