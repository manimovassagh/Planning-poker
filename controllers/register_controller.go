// controllers/register_controller.go

package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/manimovassagh/Planning-poker/database"
	"github.com/manimovassagh/Planning-poker/models"
	"golang.org/x/crypto/bcrypt"
)

// Struct to capture registration input fields
type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

// RegisterUser handles user registration and saves a hashed password
func RegisterUser(c *fiber.Ctx) error {
	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Log the received password for debugging
	log.Printf("Received Password: %s", input.Password)

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Password encryption failed"})
	}

	// Create a new user model and set the hashed password
	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Nickname:     input.Nickname,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Log the hashed password for verification
	log.Printf("Hashed Password (to store): %s", user.PasswordHash)

	// Save the user to the database
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "User registered successfully", "user": user})
}
