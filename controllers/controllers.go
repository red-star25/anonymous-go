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
				"message": validationErr.Error(),
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

		return nil
	}
}
