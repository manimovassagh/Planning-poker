// controllers/profile_controller.go

package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/manimovassagh/Planning-poker/database"
	"github.com/manimovassagh/Planning-poker/models"
)

// GetUserProfile retrieves the authenticated user's profile information
func GetUserProfile(c *fiber.Ctx) error {
	// Get user_id from context set by the JWT middleware
	userID := c.Locals("user_id").(uint)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Return user profile details (excluding sensitive info like PasswordHash)
	return c.JSON(fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"nickname": user.Nickname,
	})
}
