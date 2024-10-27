// main.go

package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/manimovassagh/Planning-poker/database"
	"github.com/manimovassagh/Planning-poker/routes"
)

func main() {
	// Initialize the database connection
	database.InitializeDatabase()

	// Set up Fiber app
	app := fiber.New()

	// Register routes
	routes.UserRoutes(app)
	routes.ProfileRoutes(app) // Add profile routes

	// Start the server
	if err := app.Listen(":3000"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
