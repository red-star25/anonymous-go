package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/red-star25/anonymous-go/controllers"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/signup", controllers.SignUp())
	app.Post("/login", controllers.Login())
}
