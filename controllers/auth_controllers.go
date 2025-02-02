package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/red-star25/anonymous-go/config"
	"github.com/red-star25/anonymous-go/database"
	"github.com/red-star25/anonymous-go/models"
	"github.com/red-star25/anonymous-go/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SignUp(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON",
		})
		return
	}

	validate := validator.New()
	if validationErr := validate.Struct(user); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.ValidateTranslator(validationErr),
		})
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	count, err := database.UserCollection().CountDocuments(ctx, bson.M{"username": user.Username})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error checking if user exists",
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User already exists",
		})
		return
	}

	password := utils.HashPassword(*user.Password)
	user.Password = &password

	user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_ID = user.ID.Hex()

	token, refreshToken, _ := utils.GenerateToken(*user.Username, user.User_ID)
	user.Token = &token
	user.Refresh_Token = &refreshToken
	user.User_Posts = make([]models.Post, 0)

	_, insertErr := database.UserCollection().InsertOne(ctx, user)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error creating user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

func Login(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON",
		})
		return
	}

	validate := validator.New()
	if validationErr := validate.Struct(user); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.ValidateTranslator(validationErr),
		})
		return
	}

	var dbUser models.User
	err := database.UserCollection().FindOne(ctx, bson.M{"username": user.Username}).Decode(&dbUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User does not exist",
		})
		return
	}

	valid, msg := utils.VerifyPassword(*dbUser.Password, *user.Password)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": msg,
		})
		return
	}

	token, refreshToken, _ := utils.GenerateToken(*dbUser.Username, dbUser.User_ID)
	dbUser.Token = &token
	dbUser.Refresh_Token = &refreshToken

	_, updateErr := database.UserCollection().UpdateOne(ctx, bson.M{"username": user.Username}, bson.M{"$set": bson.M{"token": token, "refresh_token": refreshToken}})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error updating user",
		})
		return
	}

	session := sessions.Default(c)
	session.Set("token", token)
	session.Save()

	config.SetRedisToken(dbUser.User_ID, token)

	// Respond with the created user
	c.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
		"user":    dbUser,
		"token":   token,
	})
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
