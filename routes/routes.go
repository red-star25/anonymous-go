package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/red-star25/anonymous-go/controllers"
	"github.com/red-star25/anonymous-go/middleware"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)

	r.POST("/post", middleware.Protected(), controllers.CreatePost)
	r.GET("/post", middleware.Protected(), controllers.GetPosts)
	r.GET("/post/:id", middleware.Protected(), controllers.GetPost)
	r.PUT("post/:id", middleware.Protected(), controllers.UpdatePost)
}
