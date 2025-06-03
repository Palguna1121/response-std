package models

import (
	"time"
)

type Role struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:50;not null;unique"`
	Description string    `gorm:"size:255"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Users       []User    `gorm:"foreignKey:RoleID"`
}
