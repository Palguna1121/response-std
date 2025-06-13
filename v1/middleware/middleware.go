package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"response-std/core/models"
	"response-std/core/response"

	"github.com/gin-gonic/gin"
)

// ---------------------------
// OWNER MIDDLEWARE (Resource Ownership Check)
// ---------------------------
func OwnerMiddleware(resourceIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context
		userInterface, exists := c.Get("user")
		if !exists {
			response.Unauthorized(c, "User not authenticated")
			c.Abort()
			return
		}

		user, ok := userInterface.(models.User)
		if !ok {
			response.Unauthorized(c, "Invalid user data")
			c.Abort()
			return
		}

		// Get resource ID from URL parameter
		resourceID := c.Param(resourceIDParam)
		userIDStr := fmt.Sprintf("%d", user.ID)

		// Check if user owns the resource (basic implementation)
		if resourceID != userIDStr {
			// Check if user has admin role (admins can access any resource)
			isAdmin := false
			for _, role := range user.Roles {
				if role.Name == "admin" || role.Name == "super_admin" {
					isAdmin = true
					break
				}
			}

			if !isAdmin {
				response.Forbidden(c, "Access denied. You can only access your own resources")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// ---------------------------
// VALIDATION MIDDLEWARE (Request Validation)
// ---------------------------
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add request size limit (10MB)
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)

		// Validate Content-Type for POST/PUT requests
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") &&
				!strings.Contains(contentType, "multipart/form-data") {
				response.BadRequest(c, "Content-Type must be application/json or multipart/form-data")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// ---------------------------
// ADMIN MIDDLEWARE (Super Admin Check)
// ---------------------------
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			response.Unauthorized(c, "User not authenticated")
			c.Abort()
			return
		}

		user, ok := userInterface.(models.User)
		if !ok {
			response.Unauthorized(c, "Invalid user data")
			c.Abort()
			return
		}

		// Check if user has admin or super_admin role
		isAdmin := false
		for _, role := range user.Roles {
			if role.Name == "admin" {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			response.Forbidden(c, "Access denied. Admin privileges required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// ---------------------------
// CUSTOM MIDDLEWARE HELPERS
// ---------------------------

// Multiple roles middleware
func MultipleRolesMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			response.Unauthorized(c, "User not authenticated")
			c.Abort()
			return
		}

		user, ok := userInterface.(models.User)
		if !ok {
			response.Unauthorized(c, "Invalid user data")
			c.Abort()
			return
		}

		hasRole := false
		for _, userRole := range user.Roles {
			for _, requiredRole := range roles {
				if userRole.Name == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			response.Forbidden(c, "Access denied. Insufficient privileges")
			c.Abort()
			return
		}

		c.Next()
	}
}

// IP Whitelist middleware
func IPWhitelistMiddleware(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		allowed := false
		for _, ip := range allowedIPs {
			if ip == clientIP {
				allowed = true
				break
			}
		}

		if !allowed {
			response.Forbidden(c, "IP address not whitelisted")
			c.Abort()
			return
		}

		c.Next()
	}
}
