package controllers

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/red-star25/anonymous-go/database"
	"github.com/red-star25/anonymous-go/models"
	"github.com/red-star25/anonymous-go/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SignUp() fiber.Handler {
	return func(c fiber.Ctx) error {

		// First Parse the JSON request body into the User model
		var user models.User
		if err := c.Bind().JSON(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid JSON",
			})
		}

		// Check validation errors from the User model
		if validationErr := validator.New().Struct(user); validationErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": utils.ValidateTranslator(validationErr),
			})
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Check if the user already exists
		count, err := database.UserCollection().CountDocuments(ctx, bson.M{"username": user.UserName})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error checking if user exists",
			})
		}

		if count > 0 {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "User already exists",
			})
		}

		// Hash the password and assign it back to the user
		password := utils.HashPassword(*user.Password)
		user.Password = &password

		// fill other fields
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()

		// Generate token
		token, refreshToken, _ := utils.GenerateToken(user.UserID)
		user.Token = &token
		user.RefreshToken = &refreshToken
		user.UserPosts = make([]models.Post, 0)

		// Insert the user into the database
		_, insertErr := database.UserCollection().InsertOne(ctx, user)
		if insertErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"meesage": "Error inserting user",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User created successfully",
			"user":    user,
		})
	}
}

func Login() fiber.Handler {
	return func(c fiber.Ctx) error {
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		var user models.User
		if err := c.Bind().JSON(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid JSON",
			})
		}

		if validationErr := validator.New().Struct(user); validationErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": utils.ValidateTranslator(validationErr),
			})
		}

		var dbUser models.User
		err := database.UserCollection().FindOne(ctx, bson.M{"username": user.UserName}).Decode(&dbUser)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Login or Password is incorrect",
			})
		}

		valid, msg := utils.VerifyPassword(*dbUser.Password, *user.Password)
		if !valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": msg,
			})
		}

		token, refreshToken, _ := utils.GenerateToken(dbUser.UserID)
		dbUser.Token = &token
		dbUser.RefreshToken = &refreshToken

		_, updateErr := database.UserCollection().UpdateOne(ctx, bson.M{"username": user.UserName}, bson.M{"$set": bson.M{"token": token, "refresh_token": refreshToken}})
		if updateErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error updating user",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Login successful",
			"user":    dbUser,
		})

	}
}
