package main

import (
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"your_project_path/models" // Replace with the actual path to your models
)

func main() {
	// Connect to SQLite
	db, err := gorm.Open(sqlite.Open("planning_poker_dev.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Run migrations
	err = db.AutoMigrate(&models.User{}, &models.Session{}, &models.Task{}, &models.Vote{}, &models.SessionParticipant{})
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Database migration completed successfully with SQLite.")
}