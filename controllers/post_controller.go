package controllers

import (
	"context"
	"log"
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
	userPost.Likes = make([]models.Like, 0)
	userPost.ID = primitive.NewObjectID()
	userPost.Comments = make([]models.Comment, 0)

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

func GetPosts(c *gin.Context) {
	var posts []models.Post

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	cursor, err := database.PostCollection().Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching posts",
		})
		return
	}

	for cursor.Next(ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
	}

	if !(len(posts) > 0) {
		c.JSON(http.StatusOK, gin.H{
			"message": "No posts found",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"posts": posts,
		})
	}
}

func GetPost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	postID := c.Param("id")
	id, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid post ID",
		})
		return
	}

	var post models.Post
	err = database.PostCollection().FindOne(ctx, bson.M{"_id": id}).Decode(&post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Post not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func UpdatePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	var post models.Post

	if err := c.Bind(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON",
		})
	}

	postID := c.Param("id")
	id, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid post ID",
		})
		return
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"title":      post.Title,
		"body":       post.Body,
		"updated_at": time.Now(),
	}}

	_, updateError := database.PostCollection().UpdateOne(ctx, filter, update)
	if updateError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error updating post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
	})
}

func DeletePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	id := c.Param("id")
	postID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid id",
		})
		return
	}

	_, err = database.PostCollection().DeleteOne(ctx, bson.M{"_id": postID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Post id not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post deleted successfully",
	})

}
