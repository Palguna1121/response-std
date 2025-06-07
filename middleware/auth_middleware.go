package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"response-std/helper"
	"response-std/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(tokenHeader, "Bearer ") {
			helper.Unauthorized(c, "Token tidak valid")
			c.Abort()
			return
		}

		parts := strings.SplitN(strings.TrimPrefix(tokenHeader, "Bearer "), "|", 2)
		if len(parts) != 2 {
			helper.Unauthorized(c, "Token tidak lengkap")
			c.Abort()
			return
		}

		id := parts[0]
		plain := parts[1]
		hashed := sha256.Sum256([]byte(plain))
		hashedHex := hex.EncodeToString(hashed[:])

		var token models.PersonalAccessToken
		if err := db.Where("id = ? AND token = ?", id, hashedHex).First(&token).Error; err != nil {
			helper.Unauthorized(c, "Token tidak dikenali")
			c.Abort()
			return
		}

		if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now()) {
			helper.Unauthorized(c, "Token expired")
			c.Abort()
			return
		}

		var user models.User
		if err := db.Preload("Roles.Permissions").Preload("Permissions").First(&user, token.TokenableID).Error; err != nil {
			helper.InternalServerError(c, "User tidak ditemukan")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func RoleMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			helper.Unauthorized(c, "Unauthenticated")
			c.Abort()
			return
		}

		u := user.(models.User)
		if !hasRole(u.Roles, roles) {
			helper.Forbidden(c, "Forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}

func hasRole(userRoles []models.Role, roles []string) bool {
	roleMap := make(map[string]bool)
	for _, r := range userRoles {
		roleMap[r.Name] = true
	}

	for _, r := range roles {
		if !roleMap[r] {
			return false
		}
	}

	return true
}
