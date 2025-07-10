package responses

import (
	"response-std/app/models/entities"
	"time"
)

type UserResponse struct {
	ID              uint             `json:"id,omitempty"`
	Name            string           `json:"name"`
	Email           string           `json:"email"`
	EmailVerifiedAt *time.Time       `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time        `json:"created_at,omitempty"`
	UpdatedAt       time.Time        `json:"updated_at,omitempty"`
	Roles           []entities.Roles `json:"roles,omitempty"`
}

func UserToResponse(u *entities.User) UserResponse {
	return UserResponse{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		EmailVerifiedAt: u.EmailVerifiedAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
		Roles:           u.Roles,
	}
}
