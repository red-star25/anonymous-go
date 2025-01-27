package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/red-star25/anonymous-go/controllers"
	"github.com/red-star25/anonymous-go/middleware"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)

	r.POST("/createPost", middleware.Protected(), controllers.CreatePost)
}
