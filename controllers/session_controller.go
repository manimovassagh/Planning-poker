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

// JoinSession allows an authenticated user to join an existing session
func JoinSession(c *fiber.Ctx) error {
	// Get session ID from the route parameter
	sessionID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid session ID"})
	}

	// Check if session exists
	var session models.Session
	if err := database.DB.First(&session, sessionID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Session not found"})
	}

	// Get user ID from context (set by JWT middleware)
	userID := c.Locals("user_id").(uint)

	// Check if the user is already a participant in the session
	var participant models.SessionParticipant
	if err := database.DB.Where("session_id = ? AND user_id = ?", sessionID, userID).First(&participant).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User already joined this session"})
	}

	// Add the user as a participant in the session
	newParticipant := models.SessionParticipant{
		SessionID: uint(sessionID),
		UserID:    userID,
		Role:      "participant",
		JoinedAt:  time.Now(),
	}

	if err := database.DB.Create(&newParticipant).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not join session"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User successfully joined the session",
		"participant": newParticipant,
	})
}

// GetSessionParticipants retrieves all participants of a given session
func GetSessionParticipants(c *fiber.Ctx) error {
	// Get session ID from the route parameter
	sessionID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid session ID"})
	}

	// Verify that the session exists
	var session models.Session
	if err := database.DB.First(&session, sessionID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Session not found"})
	}

	// Get user ID from the JWT token (set in middleware)
	userID := c.Locals("user_id").(uint)

	// Verify that the user is a participant in the session
	var participant models.SessionParticipant
	if err := database.DB.Where("session_id = ? AND user_id = ?", sessionID, userID).First(&participant).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User is not a participant in this session"})
	}

	// Retrieve all participants in the session
	var participants []models.SessionParticipant
	if err := database.DB.Where("session_id = ?", sessionID).Find(&participants).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve participants"})
	}

	return c.JSON(fiber.Map{"participants": participants})
}


// Struct to capture task creation input
type CreateTaskInput struct {
	TaskName        string `json:"task_name"`
	TaskDescription string `json:"task_description"`
}

// CreateTask allows the session admin to add a new task to the session
func CreateTask(c *fiber.Ctx) error {
	// Get session ID from the route parameter
	sessionID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid session ID"})
	}

	// Verify that the session exists
	var session models.Session
	if err := database.DB.First(&session, sessionID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Session not found"})
	}

	// Get user ID from JWT token (set by middleware)
	userID := c.Locals("user_id").(uint)

	// Ensure the user is the admin of the session
	if session.AdminID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only the admin can create tasks"})
	}

	// Parse request body
	var input CreateTaskInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Create the new task associated with this session
	task := models.Task{
		SessionID:       uint(sessionID),
		TaskName:        input.TaskName,
		TaskDescription: input.TaskDescription,
		CreatedAt:       time.Now(),
		Status:          "pending",
	}

	// Save the task to the database
	if err := database.DB.Create(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create task"})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}


// GetTasks retrieves all tasks for a specific session
func GetTasks(c *fiber.Ctx) error {
	// Get session ID from the route parameter
	sessionID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid session ID"})
	}

	// Verify that the session exists
	var session models.Session
	if err := database.DB.First(&session, sessionID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Session not found"})
	}

	// Get user ID from JWT token (set by middleware)
	userID := c.Locals("user_id").(uint)

	// Check if the user is a participant in the session
	var participant models.SessionParticipant
	if err := database.DB.Where("session_id = ? AND user_id = ?", sessionID, userID).First(&participant).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User is not a participant in this session"})
	}

	// Retrieve all tasks associated with this session
	var tasks []models.Task
	if err := database.DB.Where("session_id = ?", sessionID).Find(&tasks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve tasks"})
	}

	return c.JSON(fiber.Map{"tasks": tasks})
}