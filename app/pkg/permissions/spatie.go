package permissions

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"response-std/app/models/entities"

	"gorm.io/gorm"
)

type Spatie struct {
	db          *gorm.DB
	repo        *Repository
	cache       map[string]interface{}
	cacheMutex  sync.RWMutex
	cacheExpiry time.Duration
}

func NewSpatie(db *gorm.DB) *Spatie {
	return &Spatie{
		db:          db,
		repo:        NewRepository(db),
		cache:       make(map[string]interface{}),
		cacheExpiry: 5 * time.Minute,
	}
}

// ==================== Role Management ====================

func (s *Spatie) CreateRole(name, guardName string) (*entities.Roles, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("role name cannot be empty")
	}

	role := &entities.Roles{
		Name:      name,
		GuardName: guardName,
	}

	if err := s.repo.CreateRole(context.Background(), role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *Spatie) FindRole(id uint) (*entities.Roles, error) {
	return s.repo.FindRoleByID(context.Background(), id)
}

func (s *Spatie) FindRoleByName(name string) (*entities.Roles, error) {
	return s.repo.FindRoleByName(context.Background(), name)
}

func (s *Spatie) DeleteRole(id uint) error {
	return s.repo.DeleteRole(context.Background(), id)
}

func (s *Spatie) GetAllRoles() ([]entities.Roles, error) {
	var roles []entities.Roles
	if err := s.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// ==================== Permission Management ====================

func (s *Spatie) CreatePermission(name, guardName string) (*entities.Permission, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("permission name cannot be empty")
	}

	perm := &entities.Permission{
		Name:      name,
		GuardName: guardName,
	}

	if err := s.repo.CreatePermission(context.Background(), perm); err != nil {
		return nil, err
	}

	return perm, nil
}

func (s *Spatie) FindPermission(id uint) (*entities.Permission, error) {
	return s.repo.FindPermissionByID(context.Background(), id)
}

func (s *Spatie) FindPermissionByName(name string) (*entities.Permission, error) {
	return s.repo.FindPermissionByName(context.Background(), name)
}

func (s *Spatie) DeletePermission(id uint) error {
	return s.repo.DeletePermission(context.Background(), id)
}

func (s *Spatie) GetAllPermissions() ([]entities.Permission, error) {
	var perms []entities.Permission
	if err := s.db.Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

func (s *Spatie) CheckPermission(userID uint, permissionName string) (bool, error) {
	if userID == 0 || permissionName == "" {
		return false, errors.New("invalid user ID or permission name")
	}

	hasPermission, err := s.repo.CheckUserPermission(context.Background(), userID, permissionName)
	if err != nil {
		return false, fmt.Errorf("error checking permission: %w", err)
	}

	return hasPermission, nil
}

// ==================== Assignment ====================

func (s *Spatie) AssignRole(userID uint, roleName string) error {
	role, err := s.FindRoleByName(roleName)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	return s.repo.AssignRoleToUser(context.Background(), userID, role.ID)
}

func (s *Spatie) AssignPermissionToRole(roleID, permissionID uint) error {
	return s.repo.AssignPermissionToRole(context.Background(), roleID, permissionID)
}

func (s *Spatie) AssignDirectPermissionToUser(userID, permissionID uint) error {
	return s.repo.AssignDirectPermissionToUser(context.Background(), userID, permissionID)
}

// ==================== Checking ====================

func (s *Spatie) HasRole(userID uint, roleName string) (bool, error) {
	roles, err := s.GetUserRoles(userID)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if role.Name == roleName {
			return true, nil
		}
	}

	return false, nil
}

func (s *Spatie) HasAnyRole(userID uint, roleNames []string) (bool, error) {
	roles, err := s.GetUserRoles(userID)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		for _, requiredRole := range roleNames {
			if role.Name == requiredRole {
				return true, nil
			}
		}
	}

	return false, nil
}

func (s *Spatie) HasPermission(userID uint, permissionName string) (bool, error) {
	return s.repo.CheckUserPermission(context.Background(), userID, permissionName)
}

// ==================== Utility ====================

func (s *Spatie) GetUserRoles(userID uint) ([]entities.Roles, error) {
	var user entities.User
	if err := s.db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user.Roles, nil
}

func (s *Spatie) GetRolePermissions(roleID uint) ([]entities.Permission, error) {
	var role entities.Roles
	if err := s.db.Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, err
	}
	return role.Permissions, nil
}
