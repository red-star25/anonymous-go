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
)

func Like(c *gin.Context) {
	var body models.Like
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	id := c.Query("id") // Get post id from query parameter
	postID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// check if the user has already liked the post
	filter := bson.M{
		"_id":           postID,
		"likes.user_id": body.User_ID,
	}
	// update the like if the user has already liked the posts
	update := bson.M{
		"$set": bson.M{
			"likes.$.is_liked": body.Is_Liked,
		},
	}
	result, err := database.PostCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update like",
		})
		return
	}

	// if the user has not liked the post, add the like
	if result.MatchedCount == 0 {
		filter = bson.M{"_id": postID}
		update = bson.M{
			"$push": bson.M{
				"likes": body,
			},
		}
		_, err = database.PostCollection().UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to add like",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Like added successfully",
	})

}
