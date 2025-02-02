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
	var userID models.UserID

	if err := c.BindJSON(&userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
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

	// check if post exists
	var post models.Post
	filter := bson.M{"_id": postID}
	err = database.PostCollection().FindOne(ctx, filter).Decode(&post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Post not found",
		})
		return
	}

	filter = bson.M{
		"_id": postID,
	}
	update := bson.M{
		"$push": bson.M{
			"likes": userID,
		},
	}

	var likeResult bson.M
	err = database.PostCollection().FindOne(ctx, bson.M{"_id": postID, "likes": userID}).Decode(&likeResult)
	if err == nil {
		update = bson.M{
			"$pull": bson.M{
				"likes": userID,
			},
		}

		_, err = database.PostCollection().UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to remove like",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Like removed successfully",
		})
		return
	} else {
		result, err := database.PostCollection().UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to update like",
			})
			return
		}

		if result.MatchedCount == 0 {
			filter = bson.M{"_id": postID}
			update = bson.M{
				"$push": bson.M{
					"likes": userID,
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

}
