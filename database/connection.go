// database/connection.go

package database

import (
	"log"

	"github.com/manimovassagh/Planning-poker/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitializeDatabase initializes the database and assigns it to the global DB variable.
func InitializeDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("planning_poker_dev.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Automatically migrate database tables
	err = DB.AutoMigrate(&models.User{}, &models.Session{}, &models.Task{}, &models.Vote{}, &models.SessionParticipant{})
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Database connection and migration successful.")
}