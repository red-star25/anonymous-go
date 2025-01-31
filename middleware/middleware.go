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
		ClientToken, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}

		if ClientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token not found"})
			c.Abort()
			return
		}
		isValid, err := utils.ValidateToken(ClientToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		if isValid {
			c.SetCookie("token", ClientToken, 3600*24, "", "", true, true)
			c.Next()
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

	}
}
