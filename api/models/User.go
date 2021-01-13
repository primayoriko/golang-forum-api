package models

import (
	"time"
)

// User is model for users table in the db
type User struct {
	ID        uint32    `gorm:"serial;" json:"id"`
	Username  string    `gorm:"size:255;not null;unique" json:"username"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// UserResponse is struct that represent field that could be shown on response
type UserResponse struct {
	ID        uint32    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"registered_at"`
}

// UserUpdateRequest is struct that represent field that could be included on request
type UserUpdateRequest struct {
	ID       uint32 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// InjectToResponse used to inject model into corresponding response structure
func (u *User) InjectToResponse(target *UserResponse) error {
	target.ID = u.ID
	target.Username = u.Username
	target.Email = u.Email
	target.CreatedAt = u.CreatedAt
	return nil
}
