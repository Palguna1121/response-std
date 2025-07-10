package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"response-std/app/models/entities"
	"response-std/app/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ---------------------------
// AUTH MIDDLEWARE (Token Verification)
// ---------------------------
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Check if Authorization header exists and has Bearer token
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(c, "Authorization header required", nil, "[Auth Middleware]")
			c.Abort()
			return
		}

		// Extract token from header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse token format: ID|plain_token
		parts := strings.SplitN(tokenString, "|", 2)
		if len(parts) != 2 {
			response.Unauthorized(c, "Invalid token format", nil, "[Auth Middleware]")
			c.Abort()
			return
		}

		tokenID, plainToken := parts[0], parts[1]

		// Convert token ID to integer
		id, err := strconv.Atoi(tokenID)
		if err != nil {
			response.Unauthorized(c, "Invalid token ID", err, "[Auth Middleware]")
			c.Abort()
			return
		}

		// Hash the plain token to compare with stored hash
		hashedToken := sha256.Sum256([]byte(plainToken))
		hashedTokenHex := hex.EncodeToString(hashedToken[:])

		// Find token in database
		var token entities.PersonalAccessTokens
		err = db.Where("id = ? AND token = ?", id, hashedTokenHex).First(&token).Error
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token", err, "[Auth Middleware]")
			c.Abort()
			return
		}

		// Check if token is expired
		if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now()) {
			// Delete expired token
			db.Delete(&token)
			response.Unauthorized(c, "Token has expired", nil, "[Auth Middleware]")
			c.Abort()
			return
		}

		// Get user associated with the token
		var user entities.User
		err = db.Preload("Roles.Permissions").Preload("Permissions").
			Where("id = ?", token.TokenableID).First(&user).Error
		if err != nil {
			response.Unauthorized(c, "User not found", err, "[Auth Middleware]")
			c.Abort()
			return
		}

		// Store user in context for use in handlers
		c.Set("user", user)
		c.Set("token", token)

		c.Next()
	}
}

// ---------------------------
// ROLE MIDDLEWARE (Role-based Access Control)
// ---------------------------
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by AuthMiddleware)
		userInterface, exists := c.Get("user")
		if !exists {
			response.Unauthorized(c, "User not authenticated", nil, "[Role Middleware]")
			c.Abort()
			return
		}

		user, ok := userInterface.(entities.User)
		if !ok {
			response.Unauthorized(c, "Invalid user data", nil, "[Role Middleware]")
			c.Abort()
			return
		}

		// Check if user has the required role
		hasRole := false
		for _, role := range user.Roles {
			if role.Name == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			response.Forbidden(c, fmt.Sprintf("Access denied. Required role: %s", requiredRole), nil, "[Role Middleware]")
			c.Abort()
			return
		}

		c.Next()
	}
}

// ---------------------------
// PERMISSION MIDDLEWARE (Permission-based Access Control)
// ---------------------------
func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context
		userInterface, exists := c.Get("user")
		if !exists {
			response.Unauthorized(c, "User not authenticated", nil, "[Permission Middleware]")
			c.Abort()
			return
		}

		user, ok := userInterface.(entities.User)
		if !ok {
			response.Unauthorized(c, "Invalid user data", nil, "[Permission Middleware]")
			c.Abort()
			return
		}

		// Check direct permissions
		hasPermission := false
		for _, permission := range user.Permissions {
			if permission.Name == requiredPermission {
				hasPermission = true
				break
			}
		}

		// If not found in direct permissions, check role permissions
		if !hasPermission {
			for _, role := range user.Roles {
				for _, permission := range role.Permissions {
					if permission.Name == requiredPermission {
						hasPermission = true
						break
					}
				}
				if hasPermission {
					break
				}
			}
		}

		if !hasPermission {
			response.Forbidden(c, fmt.Sprintf("Access denied. Required permission: %s", requiredPermission), nil, "[Permission Middleware]")
			c.Abort()
			return
		}

		c.Next()
	}
}
