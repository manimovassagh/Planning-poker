// controllers/session_controller.go

package controllers

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/manimovassagh/Planning-poker/database"
	"github.com/manimovassagh/Planning-poker/models"
)

// Struct to capture session creation input
type CreateSessionInput struct {
	SessionName string `json:"session_name"`
}

// CreateSession handles the creation of a new planning poker session
func CreateSession(c *fiber.Ctx) error {
	// Get user_id from JWT middleware
	userID := c.Locals("user_id").(uint)

	// Parse request body
	var input CreateSessionInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Create a new session with the authenticated user as the admin
	session := models.Session{
		SessionName: input.SessionName,
		CreatedBy:   userID,
		AdminID:     userID,
		CreatedAt:   time.Now(),
		IsActive:    true,
	}

	// Save the session to the database
	if err := database.DB.Create(&session).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create session"})
	}

	// Return the created session details
	return c.Status(fiber.StatusCreated).JSON(session)
}