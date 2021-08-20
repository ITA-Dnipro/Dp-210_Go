package entity

import "github.com/ITA-Dnipro/Dp-210_Go/internal/role"

// User struct (Model)
type User struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PermissionRole role.Role `json:"roles"`
	PasswordHash   []byte    `json:"-"`
}

// NewUser represent new user in request.
type NewUser struct {
	ID              string `json:"id" validate:"omitempty"`
	Name            string `json:"name,omitempty" validate:"omitempty"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"omitempty,eqfield=Password"`
}

type PasswordsRequest struct {
	UserID          string `json:"id"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type UserNewPassword struct {
	UserID          string `json:"id"`
	OldPassword     string `json:"password_old"`
	Password        string `json:"password_new"`
	PasswordConfirm string `json:"password_new_confirm"`
}
