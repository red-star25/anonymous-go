package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/red-star25/anonymous-go/database"
	"github.com/red-star25/anonymous-go/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()
	routes.SetupRoutes(app)

	database.SetupDatabase()

	log.Fatal(app.Listen(":3000"))
}
