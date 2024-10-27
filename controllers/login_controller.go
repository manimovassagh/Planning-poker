// controllers/login_controller.go

package controllers

import (
	"log"
	"net/http"
	"time"
	"strings"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"github.com/gofiber/fiber/v2"
	"github.com/manimovassagh/Planning-poker/database"
	"github.com/manimovassagh/Planning-poker/models"
)

var jwtSecret = []byte("your_secret_key") // Replace with a secure key in production

// LoginUser handles user login and generates a JWT token on successful authentication
func LoginUser(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Fetch user record based on the provided username
	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Trim whitespace from the password (for debugging only)
	inputPassword := strings.TrimSpace(input.Password)
	storedHash := strings.TrimSpace(user.PasswordHash)

	// Log trimmed password and stored hash for debugging
	log.Printf("Trimmed Stored Hash: '%s'", storedHash)
	log.Printf("Trimmed Input Password: '%s'", inputPassword)

	// Perform password comparison with bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(inputPassword)); err != nil {
		log.Printf("Password comparison failed: %v", err)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	// Return the token on successful login
	return c.JSON(fiber.Map{"token": tokenString})
}