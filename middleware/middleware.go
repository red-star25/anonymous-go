package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/jwtauth/v5"
	"github.com/red-star25/anonymous-go/utils"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

func Protected() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}
		token, err := utils.ValidateToken(ClientToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("token", token)

		c.Next()
	}
}
