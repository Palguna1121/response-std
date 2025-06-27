package permissions

import (
	"response-std/core/models/entities"
	"response-std/core/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Spatie struct {
	db *gorm.DB
}

func NewSpatie(db *gorm.DB) *Spatie {
	return &Spatie{db: db}
}

// AssignRole assigns a role to a user
func (s *Spatie) AssignRole(userID uint, roleName string) error {
	var role entities.Roles
	if err := s.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}

	return s.db.Model(&entities.User{ID: userID}).Association("Roles").Append(&role)
}

// RevokeRole revokes a role from a user
func (s *Spatie) RevokeRole(userID uint, roleName string) error {
	var role entities.Roles
	if err := s.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}

	return s.db.Model(&entities.User{ID: userID}).Association("Roles").Delete(&role)
}

// AssignPermission assigns a permission to a role
func (s *Spatie) AssignPermission(roleName string, permissionName string) error {
	var role entities.Roles
	if err := s.db.Where("name = ?", roleName).Preload("Permissions").First(&role).Error; err != nil {
		return err
	}

	var permission entities.Permission
	if err := s.db.Where("name = ?", permissionName).First(&permission).Error; err != nil {
		return err
	}

	return s.db.Model(&role).Association("Permissions").Append(&permission)
}

// CheckPermission checks if user has a permission
func (s *Spatie) CheckPermission(userID uint, permissionName string) (bool, error) {
	var user entities.User
	if err := s.db.Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
		return false, err
	}

	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			if permission.Name == permissionName {
				return true, nil
			}
		}
	}

	return false, nil
}

// CreateRole creates a new role
func (s *Spatie) CreateRole(name string, guardName string) (*entities.Roles, error) {
	role := entities.Roles{
		Name:      name,
		GuardName: guardName,
	}

	if err := s.db.Create(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// CreatePermission creates a new permission
func (s *Spatie) CreatePermission(name string, guardName string) (*entities.Permission, error) {
	permission := entities.Permission{
		Name:      name,
		GuardName: guardName,
	}

	if err := s.db.Create(&permission).Error; err != nil {
		return nil, err
	}

	return &permission, nil
}

// GetUserRoles returns all roles for a user
func (s *Spatie) GetUserRoles(userID uint) ([]entities.Roles, error) {
	var user entities.User
	if err := s.db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}

	return user.Roles, nil
}

// GetRolePermissions returns all permissions for a role
func (s *Spatie) GetRolePermissions(roleName string) ([]entities.Permission, error) {
	var role entities.Roles
	if err := s.db.Preload("Permissions").Where("name = ?", roleName).First(&role).Error; err != nil {
		return nil, err
	}

	return role.Permissions, nil
}

// Middleware untuk permission check
func (s *Spatie) Middleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			response.Unauthorized(c, "Unauthenticated", nil)
			c.Abort()
			return
		}

		u, ok := user.(entities.User)
		if !ok {
			response.Unauthorized(c, "Invalid user data", nil)
			c.Abort()
			return
		}

		hasPermission, err := s.CheckPermission(u.ID, permission)
		if err != nil || !hasPermission {
			response.Forbidden(c, "Permission denied", err, "[Spatie Middleware]")
			c.Abort()
			return
		}

		c.Next()
	}
}
