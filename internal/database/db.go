package db

import (
	"fmt"
	"log"

	"github.com/am4rknvl/local-micro-blogging-service.git/internal/models"
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

	// Auto-migrate the database schema
	DB.AutoMigrate(&models.Post{})

	fmt.Println("Connected to microblog!")
}
