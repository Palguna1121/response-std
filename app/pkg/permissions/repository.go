package permissions

import (
	"context"
	"response-std/app/models/entities"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Role Repository Methods
func (r *Repository) FindRoleByID(ctx context.Context, id uint) (*entities.Roles, error) {
	var role entities.Roles
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repository) FindRoleByName(ctx context.Context, name string) (*entities.Roles, error) {
	var role entities.Roles
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repository) CreateRole(ctx context.Context, role *entities.Roles) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *Repository) DeleteRole(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Roles{}, id).Error
}

// Permission Repository Methods
func (r *Repository) FindPermissionByID(ctx context.Context, id uint) (*entities.Permission, error) {
	var perm entities.Permission
	if err := r.db.WithContext(ctx).First(&perm, id).Error; err != nil {
		return nil, err
	}
	return &perm, nil
}

func (r *Repository) FindPermissionByName(ctx context.Context, name string) (*entities.Permission, error) {
	var perm entities.Permission
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&perm).Error; err != nil {
		return nil, err
	}
	return &perm, nil
}

func (r *Repository) CreatePermission(ctx context.Context, perm *entities.Permission) error {
	return r.db.WithContext(ctx).Create(perm).Error
}

func (r *Repository) DeletePermission(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Permission{}, id).Error
}

// Model-Role-Permission Relationship Methods
func (r *Repository) AssignRoleToUser(ctx context.Context, userID uint, roleID uint) error {
	return r.db.WithContext(ctx).Table("model_has_roles").Create(&entities.ModelHasRoles{
		ModelID:   userID,
		ModelType: "App\\Models\\User",
		RoleID:    roleID,
	}).Error
}

func (r *Repository) AssignDirectPermissionToUser(ctx context.Context, userID uint, permissionID uint) error {
	return r.db.WithContext(ctx).Table("model_has_permissions").Create(&entities.ModelHasPermissions{
		ModelID:      userID,
		ModelType:    "App\\Models\\User",
		PermissionID: permissionID,
	}).Error
}

func (r *Repository) AssignPermissionToRole(ctx context.Context, roleID uint, permissionID uint) error {
	role := entities.Roles{ID: roleID}
	perm := entities.Permission{ID: permissionID}
	return r.db.WithContext(ctx).Model(&role).Association("Permissions").Append(&perm)
}

func (r *Repository) CheckUserPermission(ctx context.Context, userID uint, permission string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.ModelHasPermissions{}).
		Joins("JOIN permissions ON permissions.id = model_has_permissions.permission_id").
		Where("model_has_permissions.model_id = ? AND model_has_permissions.model_type = ? AND permissions.name = ?",
			userID, "User", permission).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	// Check permissions via roles
	err = r.db.WithContext(ctx).
		Model(&entities.ModelHasRoles{}).
		Joins("JOIN roles ON roles.id = model_has_roles.role_id").
		Joins("JOIN role_has_permissions ON role_has_permissions.role_id = roles.id").
		Joins("JOIN permissions ON permissions.id = role_has_permissions.permission_id").
		Where("model_has_roles.model_id = ? AND model_has_roles.model_type = ? AND permissions.name = ?",
			userID, "User", permission).
		Count(&count).Error

	return count > 0, err
}
