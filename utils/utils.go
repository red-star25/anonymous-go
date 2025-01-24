package utils

import (
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

const (
	JWT_SECRET = "JWT_SECRET"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	valid := true
	msg := ""
	if err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword)); err != nil {
		msg = "Login or Password is incorrect"
		valid = false
	}
	return valid, msg
}

func GenerateToken(userName string) (token, refreshToken string, err error) {
	claims := jwt.MapClaims{
		"username": userName,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"username": userName,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv(JWT_SECRET)))
	if err != nil {
		log.Fatal(err)
		return
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv(JWT_SECRET)))
	if err != nil {
		log.Fatal(err)
		return
	}

	return token, refreshToken, nil
}

func ValidateTranslator(err error) string {
	if err != nil {
		validationErr := err.(validator.ValidationErrors)
		for _, e := range validationErr {
			if e.Tag() == "min" {
				return e.Field() + " must be at least " + e.Param() + " characters long"
			}
			if e.Tag() == "required" {
				return e.Field() + " is required"
			}
		}
	}
	return err.Error()
}
