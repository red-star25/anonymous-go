package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/red-star25/anonymous-go/database"
	"github.com/red-star25/anonymous-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddComment(c *gin.Context) {
	var body models.Comment
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if body.Comment_Body == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Comment text is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	postID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	var post models.Post
	filter := bson.M{"_id": postID}
	err = database.PostCollection().FindOne(ctx, filter).Decode(&post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Post not found",
		})
		return
	}

	body.ID = primitive.NewObjectID()
	body.Created_At = time.Now()

	filter = bson.M{"_id": postID}
	update := bson.M{"$push": bson.M{
		"comments": body,
	}}
	_, err = database.PostCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment added successfully",
		"comment": body,
	})

}

func GetComments(c *gin.Context) {
	id := c.Param("id")
	postHexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid post id",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	var comments []models.Comment
	filter := bson.M{"_id": postHexID}
	projection := bson.M{"comments": 1}
	err = database.PostCollection().FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&comments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})

}

func DeleteComment(c *gin.Context) {
	id := c.Query("id")
	commentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	postid := c.Query("postid")
	postID, err := primitive.ObjectIDFromHex(postid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	filter := bson.M{"_id": postID}
	update := bson.M{
		"$pull": bson.M{
			"comments": bson.M{"_id": commentID},
		},
	}
	res, err := database.PostCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	if res.ModifiedCount == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No comment deleted",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Comment deleted successfully",
		})
	}

}
