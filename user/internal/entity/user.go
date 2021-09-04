package entity

import "github.com/ITA-Dnipro/Dp-210_Go/user/internal/role"

// User struct (Model)
type User struct {
	ID             string    `json:"id" readonly:"true"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PermissionRole role.Role `json:"roles"`
	PasswordHash   []byte    `json:"-"`
}
