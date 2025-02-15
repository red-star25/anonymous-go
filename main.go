package main

import (
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/red-star25/anonymous-go/config"
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

	database.SetupDatabase()

	redisConfig := config.NewRedisConfig()
	redisConfig.StartRedis()

	store, err := redis.NewStore(10, "tcp", os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"), []byte(os.Getenv("REDIS_SECRET")))
	if err != nil {
		log.Fatal(err)
	}

	r.Use(sessions.Sessions("redis-session", store))
	routes.SetupRoutes(r)

	log.Fatal(r.Run("0.0.0.0:" + os.Getenv("PORT")))
}
