// routes/ws_routes.go

package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/manimovassagh/Planning-poker/controllers"
)

func WebSocketRoutes(app *fiber.App) {
	app.Get("/ws", websocket.New(controllers.WebSocketHandler))
}