// routes/profile_routes.go

package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/manimovassagh/Planning-poker/controllers"
	"github.com/manimovassagh/Planning-poker/middleware"
)

func ProfileRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/user/me", middleware.JWTProtected, controllers.GetUserProfile)
}