package entity

import "github.com/ITA-Dnipro/Dp-210_Go/user/internal/role"

// User struct (Model)
type User struct {
	ID             string    `json:"id" readonly:"true"`
	Name           string    `json:"name,omitempty" validate:"omitempty"`
	Email          string    `json:"email" validate:"required,email"`
	PermissionRole role.Role `json:"role" validate:"required"`
	PasswordHash   []byte    `json:"-"`
}

// UserList struct
type UserList struct {
	Users  []User `json:"data"`
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}
