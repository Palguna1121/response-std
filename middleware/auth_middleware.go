package middleware

import (
	"response-std/config"
	"response-std/helper"
	"response-std/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			helper.Unauthorized(c, "Authorization header is missing or invalid")
			c.Abort()
			return
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")
		parts := strings.Split(rawToken, "|")
		if len(parts) != 2 {
			helper.Unauthorized(c, "Invalid token format")
			c.Abort()
			return
		}

		userID := parts[0]
		jwtToken := parts[1]

		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			return config.ENV.JWT_SECRET, nil // Menggunakan secret dari config
		})

		if err != nil || !token.Valid {
			helper.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		var user models.User
		db := config.DB
		if err := db.First(&user, userID).Error; err != nil {
			helper.Unauthorized(c, "User not found")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
