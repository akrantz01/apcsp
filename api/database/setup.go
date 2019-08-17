package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

const (
	host     = "172.17.0.2"
	port     = "5432"
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
	sslmode  = "disable"
)

// Connect to the database and create the schema
func SetupDatabase() *gorm.DB {
	log.Println("Connecting to database...")
	// Connect to database
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")

	log.Println("Building database schema...")
	// Create schema if not exist
	for _, model := range []interface{}{} {
		if !db.HasTable(model) {
			db.CreateTable(model)
		}
	}
	log.Println("Built database schema")

	return db
}
