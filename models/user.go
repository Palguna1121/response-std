package models

import (
	"time"
	// "gorm.io/gorm"
)

type User struct {
	ID              uint   `gorm:"primaryKey"`
	Name            string `gorm:"size:50;not null"`
	Email           string `gorm:"size:50;unique;not null"`
	RoleID          uint   `gorm:"not null"`
	EmailVerifiedAt *time.Time
	Password        string `gorm:"not null"`
	RememberToken   string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	// DeletedAt       gorm.DeletedAt `gorm:"index"`
	Role          Role `gorm:"foreignKey:RoleID;references:ID"`
}
