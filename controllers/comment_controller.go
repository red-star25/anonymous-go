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

func AddComment(c *gin.Context) {
	var body models.Comments
	if err := c.ShouldBindJSON(&body); err != nil {
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
	body.ID = primitive.NewObjectID()
	body.Created_At = time.Now()

	filter := bson.M{"_id": postID}
	update := bson.M{"$push": bson.M{
		"comments": body,
	}}
	_, err = database.PostCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	// TODO: add comment to the current users collection

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment added successfully",
		"comment": body,
	})

}

func GetComment(c *gin.Context) {

}

func DeleteComment(c *gin.Context) {

}
