// controllers/session_controller.go

package controllers

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/manimovassagh/Planning-poker/database"
	"github.com/manimovassagh/Planning-poker/models"
)

func CreateSession(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint) // Assuming JWT middleware sets this

	var session models.Session
	if err := c.BodyParser(&session); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	session.CreatedBy = userId
	session.AdminID = userId
	session.CreatedAt = time.Now()
	session.IsActive = true

	if err := database.DB.Create(&session).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create session"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Session created successfully", "session": session})
}
