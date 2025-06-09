// models/model_has_roles.go
package models

type ModelHasRoles struct {
	RoleID    uint   `gorm:"primaryKey"`
	ModelType string `gorm:"size:255;primaryKey"`
	ModelID   uint   `gorm:"primaryKey"` // Ini akan menjadi foreign key ke user
	Roles      Roles
}
