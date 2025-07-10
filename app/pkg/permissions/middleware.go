package permissions

import (
	"response-std/app/models/entities"
	"response-std/app/pkg/response"

	"github.com/gin-gonic/gin"
)

// Middleware untuk permission check
func (s *Spatie) Middleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			response.Unauthorized(c, "Unauthenticated", nil, "[Permission Middleware]")
			c.Abort()
			return
		}

		u, ok := user.(entities.User)
		if !ok {
			response.Unauthorized(c, "Invalid user data", nil, "[Permission Middleware]")
			c.Abort()
			return
		}

		hasPermission, err := s.CheckPermission(u.ID, permission)
		if err != nil || !hasPermission {
			response.Forbidden(c, "Permission denied", err, "[Permission Middleware]")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RoleMiddleware memeriksa apakah user memiliki salah satu role yang dibutuhkan
func (s *Spatie) RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			response.Unauthorized(c, "Unauthenticated", nil, "[Role Middleware]")
			c.Abort()
			return
		}

		u, ok := user.(entities.User)
		if !ok {
			response.Unauthorized(c, "Invalid user data", nil, "[Role Middleware]")
			c.Abort()
			return
		}

		userRoles, err := s.GetUserRoles(u.ID)
		if err != nil {
			response.Error(c, 500, "Failed to check roles", err, "[Role Middleware]")
			c.Abort()
			return
		}

		hasRole := false
		for _, userRole := range userRoles {
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
			response.Forbidden(c, "Role access denied", nil, "[Role Middleware]")
			c.Abort()
			return
		}

		c.Next()
	}
}

// AnyPermissionMiddleware memeriksa apakah user memiliki salah satu permission yang dibutuhkan
func (s *Spatie) AnyPermissionMiddleware(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			response.Unauthorized(c, "Unauthenticated", nil, "[AnyPermission Middleware]")
			c.Abort()
			return
		}

		u, ok := user.(entities.User)
		if !ok {
			response.Unauthorized(c, "Invalid user data", nil, "[AnyPermission Middleware]")
			c.Abort()
			return
		}

		for _, perm := range permissions {
			hasPermission, err := s.CheckPermission(u.ID, perm)
			if err == nil && hasPermission {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Permission denied", nil, "[AnyPermission Middleware]")
		c.Abort()
	}
}

// AllPermissionsMiddleware memeriksa apakah user memiliki semua permission yang dibutuhkan
func (s *Spatie) AllPermissionsMiddleware(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			response.Unauthorized(c, "Unauthenticated", nil, "[AllPermissions Middleware]")
			c.Abort()
			return
		}

		u, ok := user.(entities.User)
		if !ok {
			response.Unauthorized(c, "Invalid user data", nil, "[AllPermissions Middleware]")
			c.Abort()
			return
		}

		for _, perm := range permissions {
			hasPermission, err := s.CheckPermission(u.ID, perm)
			if err != nil || !hasPermission {
				mssg := "Permission denied: " + perm
				response.Forbidden(c, mssg, err, "[AllPermissions Middleware]")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
