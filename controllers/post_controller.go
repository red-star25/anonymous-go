package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/red-star25/anonymous-go/database"
	"github.com/red-star25/anonymous-go/models"
	"github.com/red-star25/anonymous-go/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreatePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var userPost models.Post
	if err := c.BindJSON(&userPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON",
		})
		return
	}
	validate := validator.New()
	if validationErr := validate.Struct(userPost); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.ValidateTranslator(validationErr),
		})
		return
	}

	userPost.Created_At = time.Now()
	userPost.Updated_At = time.Now()
	userPost.Likes = make([]models.Likes, 0)
	userPost.ID = primitive.NewObjectID()
	userPost.Comments = make([]models.Comments, 0)

	_, insertError := database.PostCollection().InsertOne(ctx, userPost)
	if insertError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error creating post",
		})
		return
	}

	id, err := primitive.ObjectIDFromHex(userPost.User_ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error creating post",
		})
		return
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{primitive.E{Key: "$push", Value: bson.D{primitive.E{Key: "user_posts", Value: userPost}}}}

	_, updateError := database.UserCollection().UpdateOne(ctx, filter, update)
	if updateError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error creating post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post created successfully",
		"post":    userPost,
	})
}
