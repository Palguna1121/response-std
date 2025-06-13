// models/user.go
package models

import (
	"time"

	"response-std/core/helper"

	"gorm.io/gorm"
)

type User struct {
	ID                   uint   `gorm:"primaryKey"`
	Name                 string `gorm:"size:50"`
	Email                string `gorm:"size:50;unique"`
	EmailVerifiedAt      *time.Time
	Password             string
	RememberToken        *string `gorm:"size:100"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt         `gorm:"index"`
	Roles                []Roles                `gorm:"many2many:model_has_roles;foreignKey:ID;joinForeignKey:model_id;joinReferences:role_id"`
	Permissions          []Permission           `gorm:"many2many:model_has_permissions;foreignKey:ID;joinForeignKey:model_id;joinReferences:permission_id"`
	PersonalAccessTokens []PersonalAccessTokens `gorm:"foreignKey:TokenableID;constraint:OnDelete:CASCADE"`
}

func (u *User) CheckPassword(pw string) bool {
	if u.Password == "" {
		return false
	}
	return helper.CheckPasswordHash(pw, u.Password)
}
