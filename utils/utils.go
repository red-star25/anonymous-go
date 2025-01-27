package utils

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt/v5"
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

	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(givenPassword)); err != nil {
		msg = "Login or Password is incorrect"
		valid = false
	}

	return valid, msg
}

func GenerateToken(userName string, userID string) (token, refreshToken string, err error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"expires_at": time.Now().Add(time.Hour * 24),
	})

	refreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"expires_at": time.Now().Add(time.Hour * 24 * 7),
	})

	token, err = claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Fatal("could not generate token")
	}

	refreshToken, err = refreshClaims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Fatal("could not generate refresh token")
	}

	return token, refreshToken, nil
}

func ValidateTranslator(err error) string {
	if err != nil {
		validationErr := err.(validator.ValidationErrors)
		for _, e := range validationErr {
			if e.Tag() == "min" {
				return strings.ToLower(e.Field()) + " must be at least " + e.Param() + " characters long"
			}
			if e.Tag() == "required" {
				return strings.ToLower(e.Field()) + " is required"
			}
		}
	}
	return err.Error()
}
func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(JWT_SECRET)), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Get UserID from JWT Claims

// func ParseJWT(tokenStr string) (map[string]interface{}, error) {
// 	secret := os.Getenv("JWT_SECRET")
// 	if secret == "" {
// 		return nil, errors.New("JWT_SECRET is not set in the environment")
// 	}

// 	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(secret), nil
// 	})

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse token: %w", err)
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		return claims, nil
// 	}

// 	return nil, errors.New("invalid token")
// }

// func GetUserIDFromToken(token string) (*string, error) {
// 	claims, err := ParseJWT(token)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid token")
// 	}
// 	userID, ok := claims["user_id"].(string)
// 	if !ok {
// 		return nil, fmt.Errorf("user_id not found in token claims")
// 	}
// 	return &userID, nil
// }
