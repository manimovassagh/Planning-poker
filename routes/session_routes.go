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
}