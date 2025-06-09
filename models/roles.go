// models/role.go
package models

import "time"

type Roles struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:255;uniqueIndex:role_name_guard_name"`
	GuardName   string `gorm:"size:255;uniqueIndex:role_name_guard_name"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Permissions []Permission `gorm:"many2many:role_has_permissions;foreignKey:ID;joinForeignKey:role_id;joinReferences:permission_id"`
	Users       []User       `gorm:"many2many:model_has_roles;foreignKey:ID;joinForeignKey:role_id;joinReferences:model_id"`
}
