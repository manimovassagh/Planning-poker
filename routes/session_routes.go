// routes/session_routes.go

package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/manimovassagh/Planning-poker/controllers"
	"github.com/manimovassagh/Planning-poker/middleware"
)

func SessionRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/sessions", middleware.JWTProtected, controllers.CreateSession)           // Protected session creation
	api.Post("/sessions/:id/join", middleware.JWTProtected, controllers.JoinSession)     // Protected join session
	api.Get("/sessions/:id/participants", middleware.JWTProtected, controllers.GetSessionParticipants) // Protected view participants
	api.Post("/sessions/:id/tasks", middleware.JWTProtected, controllers.CreateTask)     // Protected task creation
	api.Get("/sessions/:id/tasks", middleware.JWTProtected, controllers.GetTasks)        // Protected view tasks
	api.Post("/sessions/:session_id/tasks/:task_id/vote", middleware.JWTProtected, controllers.SubmitVote) // Protected submit vote
	api.Get("/sessions/:session_id/tasks/:task_id/votes", middleware.JWTProtected, controllers.GetVotes)   // Protected view votes
	api.Post("/sessions/:session_id/tasks/:task_id/reveal", middleware.JWTProtected, controllers.RevealVotes) // Protected reveal votes
	api.Post("/sessions/:id/close", middleware.JWTProtected, controllers.CloseSession)   // Protected close session
	api.Get("/sessions/:id/summary", middleware.JWTProtected, controllers.GetSessionSummary) // Protected view session summary
}