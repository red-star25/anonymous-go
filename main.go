package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/red-star25/anonymous-go/database"
	"github.com/red-star25/anonymous-go/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.New()
	r.Use(gin.Logger())

	routes.SetupRoutes(r)

	database.SetupDatabase()

	log.Fatal(r.Run(":3000"))
}
