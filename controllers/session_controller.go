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
		"message":     "User successfully joined the session",
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

// Struct to capture vote submission input
type SubmitVoteInput struct {
	VoteValue int `json:"vote_value"`
}

// SubmitVote allows a participant to submit or update their vote for a task
func SubmitVote(c *fiber.Ctx) error {
	// Get session and task IDs from the route parameters
	sessionID, err := c.ParamsInt("session_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid session ID"})
	}
	taskID, err := c.ParamsInt("task_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	// Verify that the session exists
	var session models.Session
	if err := database.DB.First(&session, sessionID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Session not found"})
	}

	// Verify that the task exists and belongs to the session
	var task models.Task
	if err := database.DB.Where("id = ? AND session_id = ?", taskID, sessionID).First(&task).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found in this session"})
	}

	// Get user ID from JWT token (set by middleware)
	userID := c.Locals("user_id").(uint)

	// Check if the user is a participant in the session
	var participant models.SessionParticipant
	if err := database.DB.Where("session_id = ? AND user_id = ?", sessionID, userID).First(&participant).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User is not a participant in this session"})
	}

	// Parse the vote input
	var input SubmitVoteInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Check if the user has already voted on this task
	var vote models.Vote
	if err := database.DB.Where("task_id = ? AND user_id = ?", taskID, userID).First(&vote).Error; err == nil {
		// If a vote exists, update it
		vote.VoteValue = input.VoteValue
		vote.CreatedAt = time.Now()
		if err := database.DB.Save(&vote).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update vote"})
		}
		BroadcastUpdate("New vote submitted for task")
		return c.JSON(fiber.Map{"message": "Vote updated successfully", "vote": vote})
	}

	// If no previous vote exists, create a new vote
	newVote := models.Vote{
		TaskID:    uint(taskID),
		UserID:    userID,
		VoteValue: input.VoteValue,
		CreatedAt: time.Now(),
	}
	if err := database.DB.Create(&newVote).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not submit vote"})
	}

	return c.JSON(fiber.Map{"message": "Vote submitted successfully", "vote": newVote})
}

// GetVotes retrieves all votes for a specific task in a session
func GetVotes(c *fiber.Ctx) error {
	// Get session and task IDs from the route parameters
	sessionID, err := c.ParamsInt("session_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid session ID"})
	}
	taskID, err := c.ParamsInt("task_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	// Verify that the session exists
	var session models.Session
	if err := database.DB.First(&session, sessionID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Session not found"})
	}

	// Verify that the task exists and belongs to the session
	var task models.Task
	if err := database.DB.Where("id = ? AND session_id = ?", taskID, sessionID).First(&task).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found in this session"})
	}

	// Get user ID from JWT token (set by middleware)
	userID := c.Locals("user_id").(uint)

	// Ensure the requesting user is the session admin
	if session.AdminID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only the admin can view votes"})
	}

	// Retrieve all votes for the task
	var votes []models.Vote
	if err := database.DB.Where("task_id = ?", taskID).Find(&votes).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve votes"})
	}

	return c.JSON(fiber.Map{"votes": votes})
}

// RevealVotes allows the admin to reveal votes for a specific task
func RevealVotes(c *fiber.Ctx) error {
	// Get session and task IDs from the route parameters
	sessionID, err := c.ParamsInt("session_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid session ID"})
	}
	taskID, err := c.ParamsInt("task_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	// Verify that the session exists
	var session models.Session
	if err := database.DB.First(&session, sessionID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Session not found"})
	}

	// Verify that the task exists and belongs to the session
	var task models.Task
	if err := database.DB.Where("id = ? AND session_id = ?", taskID, sessionID).First(&task).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found in this session"})
	}

	// Get user ID from JWT token (set by middleware)
	userID := c.Locals("user_id").(uint)

	// Ensure the requesting user is the session admin
	if session.AdminID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only the admin can reveal votes"})
	}

	// Update task status to "revealed"
	task.Status = "revealed"
	// After setting task status to "revealed" in RevealVotes
	BroadcastUpdate("Task votes have been revealed")
	if err := database.DB.Save(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not reveal votes"})
	}

	// Retrieve all votes for the task
	var votes []models.Vote
	if err := database.DB.Where("task_id = ?", taskID).Find(&votes).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve votes"})
	}

	return c.JSON(fiber.Map{
		"message":     "Votes revealed successfully",
		"task_status": task.Status,
		"votes":       votes,
	})
}

// CloseSession allows the session admin to close the session
func CloseSession(c *fiber.Ctx) error {
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

	// Ensure the requesting user is the session admin
	if session.AdminID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only the admin can close this session"})
	}

	// Update session's IsActive status to false
	session.IsActive = false
	if err := database.DB.Save(&session).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not close session"})
	}

	return c.JSON(fiber.Map{
		"message":        "Session closed successfully",
		"session_status": "closed",
	})
}

// GetSessionSummary provides a summary of all tasks and votes for a closed session
func GetSessionSummary(c *fiber.Ctx) error {
	// Get session ID from the route parameter
	sessionID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid session ID"})
	}

	// Verify that the session exists and is closed
	var session models.Session
	if err := database.DB.First(&session, sessionID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Session not found"})
	}
	if session.IsActive {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Session is still active. Summary is available only for closed sessions"})
	}

	// Get user ID from JWT token (set by middleware)
	userID := c.Locals("user_id").(uint)

	// Ensure the requesting user is the session admin
	if session.AdminID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only the admin can view the session summary"})
	}

	// Retrieve all tasks for the session
	var tasks []models.Task
	if err := database.DB.Where("session_id = ?", sessionID).Find(&tasks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve tasks"})
	}

	// Retrieve votes for each task and build the summary
	summary := make([]fiber.Map, len(tasks))
	for i, task := range tasks {
		var votes []models.Vote
		if err := database.DB.Where("task_id = ?", task.ID).Find(&votes).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve votes"})
		}
		summary[i] = fiber.Map{
			"task_id":          task.ID,
			"task_name":        task.TaskName,
			"task_description": task.TaskDescription,
			"task_status":      task.Status,
			"votes":            votes,
		}
	}

	return c.JSON(fiber.Map{
		"session_name": session.SessionName,
		"summary":      summary,
	})
}
