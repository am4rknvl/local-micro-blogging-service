package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=amar password=mystrongpass dbname=microblog port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	DB = db

	// Enable UUID extension for PostgreSQL
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	// Auto-migrate the database schema
	// Note: Models will be migrated in main.go to avoid import cycles

	fmt.Println("Connected to microblog!")
}
