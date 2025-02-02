package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/jwtauth/v5"
	"github.com/red-star25/anonymous-go/config"
	"github.com/red-star25/anonymous-go/utils"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

func Protected() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		tokenInterface := session.Get("token")

		if tokenInterface == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: no token found in session"})
			c.Abort()
			return
		}

		tokenStr, ok := tokenInterface.(string)
		if !ok || tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: token format invalid"})
			c.Abort()
			return
		}

		token, err := tokenAuth.Decode(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: token invalid"})
			c.Abort()
			return
		}

		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: token invalid"})
			c.Abort()
			return
		}

		userIDInterface, ok := token.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user_id claim missing"})
			c.Abort()
			return
		}

		userID, ok := userIDInterface.(string)
		if !ok || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid user_id claim"})
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		redisToken, err := config.RedisClient.Get(ctx, userID).Result()
		if redisToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: token not found in redis"})
			c.Abort()
			return
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired, please log in again"})
			c.Abort()
			return
		}

		if redisToken != tokenStr {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: token mismatch"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()

	}
}
