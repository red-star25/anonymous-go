package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/red-star25/anonymous-go/controllers"
	"github.com/red-star25/anonymous-go/middleware"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/hello", controllers.Hello)
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)

	p := r.Group("/post")
	{
		p.POST("/", middleware.Protected(), controllers.CreatePost)
		p.GET("/", middleware.Protected(), controllers.GetPosts)
		p.GET("/:id", middleware.Protected(), controllers.GetPost)
		p.PUT("/:id", middleware.Protected(), controllers.UpdatePost)
		p.DELETE("/:id", middleware.Protected(), controllers.DeletePost)
	}

	c := r.Group("/comment")
	{
		c.POST("/:id", middleware.Protected(), controllers.AddComment)
		c.GET("/:id", middleware.Protected(), controllers.GetComments)
		c.DELETE("/", middleware.Protected(), controllers.DeleteComment)
	}

	r.POST("/like/:id", middleware.Protected(), controllers.Like)

	r.POST("/logout", middleware.Protected(), controllers.Logout)

}
