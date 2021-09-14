package entity

import (
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/role"
	"github.com/google/uuid"
)

// User struct (Model)
type User struct {
	ID             uuid.UUID `json:"id" readonly:"true"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PermissionRole role.Role `json:"roles"`
	PasswordHash   []byte    `json:"-"`
}
