package entities

import (
	"time"
)

type PersonalAccessTokens struct {
	ID            uint       `gorm:"primaryKey" json:"id,omitempty"`
	TokenableID   uint       `gorm:"index" json:"-"`
	TokenableType string     `gorm:"size:255" json:"-"`
	Name          string     `gorm:"size:255" json:"name,omitempty"`
	Token         string     `gorm:"size:64;unique" json:"-"`
	Abilities     *string    `gorm:"type:text" json:"abilities,omitempty"`
	LastUsedAt    *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at,omitempty"`
}
