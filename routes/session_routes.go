// routes/session_routes.go

package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/manimovassagh/Planning-poker/controllers"
)

func SessionRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/sessions", controllers.CreateSession)
}
