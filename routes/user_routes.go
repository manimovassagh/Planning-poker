// routes/user_routes.go

package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/manimovassagh/Planning-poker/controllers"
)

func UserRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/register", controllers.RegisterUser) // From register_controller.go
	api.Post("/login", controllers.LoginUser)       // From login_controller.go
}
