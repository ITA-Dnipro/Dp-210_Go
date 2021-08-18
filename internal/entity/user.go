package entity

import "github.com/ITA-Dnipro/Dp-210_Go/internal/role"

// User struct (Model)
type User struct {
	ID             string    `json:"id" readonly:"true"`
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
