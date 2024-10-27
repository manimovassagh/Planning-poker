// routes/session_routes.go

package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/manimovassagh/Planning-poker/controllers"
	"github.com/manimovassagh/Planning-poker/middleware"
)

func SessionRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/sessions", middleware.JWTProtected, controllers.CreateSession) // Protected session creation
}