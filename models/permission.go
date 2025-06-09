// models/permission.go
package models

import "time"

type Permission struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;uniqueIndex:permission_name_guard_name"`
	GuardName string `gorm:"size:255;uniqueIndex:permission_name_guard_name"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Roles     []Roles `gorm:"many2many:role_has_permissions;foreignKey:ID;joinForeignKey:permission_id;joinReferences:role_id"`
}
