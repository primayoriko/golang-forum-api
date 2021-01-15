package models

import (
	"errors"
	"time"

	validator "github.com/asaskevich/govalidator"
	"gitlab.com/hydra/forum-api/api/utils"
)

// User is model for users table in the db
type User struct {
	ID        uint32    `gorm:"serial;" json:"id"`
	Username  string    `gorm:"size:255;not null;unique" json:"username"`
	Email     string    `gorm:"size:255;not null;unique" json:"email"`
	Password  string    `gorm:"size:255;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// RegistrationRequest is struct that represent field that could be included on register request
type RegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserUpdateRequest is struct that represent field that could be included on update request
type UserUpdateRequest struct {
	ID       uint32 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse is struct that represent field that could be shown on response
type UserResponse struct {
	ID        uint32    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"registered_at"`
}

// InjectToModel used to inject request into corresponding model
func (rr *RegistrationRequest) InjectToModel(target *User) error {
	if !utils.IsNonEmpty(rr.Username, rr.Email, rr.Password) {
		return errors.New("username/email/password can't be left empty/blank")
	}

	if !validator.IsEmail(rr.Email) || !validator.IsPrintableASCII(rr.Username) ||
		!validator.IsPrintableASCII(rr.Password) {
		return errors.New("username/email/password have bad format or character")
	}

	target.Email = rr.Email
	target.Username = rr.Username
	target.Password = rr.Password

	return nil
}

// InjectToModel used to inject request into corresponding model
func (ur *UserUpdateRequest) InjectToModel(target *User) error {
	if ur.ID == 0 {
		return errors.New("id cannot be left blank/empty")
	}

	if !validator.IsPrintableASCII(ur.Password) || !validator.IsEmail(ur.Email) {
		return errors.New("email/password have bad format or character")
	}

	if ur.Email != "" {
		target.Email = ur.Email
	}

	if ur.Password != "" {
		target.Password = ur.Password
	}

	target.ID = ur.ID

	return nil
}

// InjectToResponse used to inject model into corresponding response structure
func (u *User) InjectToResponse(target *UserResponse) error {
	target.ID = u.ID
	target.Username = u.Username
	target.Email = u.Email
	target.CreatedAt = u.CreatedAt
	return nil
}

// InsertFromModel used to insert response struct data form It's corresponding model
func (ur *UserResponse) InsertFromModel(user User) error {
	ur.ID = user.ID
	ur.Username = user.Username
	ur.Email = user.Email
	ur.CreatedAt = user.CreatedAt
	return nil
}

// type UsersResponse struct {
// 	Data []User `json:"users" example:"[{\"id\":100, \"name\":\"Mikun\"}]" description:"User identity"`
// }
