package permissions

import (
	"github.com/gin-gonic/gin"
)

// PermissionManager mendefinisikan kontrak untuk sistem permission
type PermissionManager interface {
	// Role Management
	CreateRole(name, guardName string) (interface{}, error)
	FindRole(id uint) (interface{}, error)
	FindRoleByName(name string) (interface{}, error)
	DeleteRole(id uint) error
	GetAllRoles() ([]interface{}, error)

	// Permission Management
	CreatePermission(name, guardName string) (interface{}, error)
	FindPermission(id uint) (interface{}, error)
	FindPermissionByName(name string) (interface{}, error)
	DeletePermission(id uint) error
	GetAllPermissions() ([]interface{}, error)

	// Assignment
	AssignRole(userID uint, roleName string) error
	AssignRoleToModel(model interface{}, role interface{}) error
	AssignPermissionToRole(role interface{}, permission interface{}) error
	AssignDirectPermissionToModel(model interface{}, permission interface{}) error
	SyncRoles(model interface{}, roles []interface{}) error
	SyncPermissions(model interface{}, permissions []interface{}) error

	// Revocation
	RevokeRole(userID uint, roleName string) error
	RevokePermissionFromRole(role interface{}, permission interface{}) error
	RevokePermissionFromModel(model interface{}, permission interface{}) error
	RemoveAllRolesFromModel(model interface{}) error
	RemoveAllPermissionsFromModel(model interface{}) error

	// Checking
	HasRole(model interface{}, role interface{}) (bool, error)
	HasAnyRole(model interface{}, roles []interface{}) (bool, error)
	HasAllRoles(model interface{}, roles []interface{}) (bool, error)
	HasPermission(model interface{}, permission interface{}) (bool, error)
	HasAnyPermission(model interface{}, permissions []interface{}) (bool, error)
	HasAllPermissions(model interface{}, permissions []interface{}) (bool, error)

	// Middleware
	RoleMiddleware(roles ...string) gin.HandlerFunc
	PermissionMiddleware(permissions ...string) gin.HandlerFunc
	AnyPermissionMiddleware(permissions ...string) gin.HandlerFunc

	// Utility
	GetModelRoles(model interface{}) ([]interface{}, error)
	GetModelPermissions(model interface{}) ([]interface{}, error)
	GetRolePermissions(role interface{}) ([]interface{}, error)
}
