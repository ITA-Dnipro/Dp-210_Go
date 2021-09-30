package entity

import "github.com/google/uuid"

// User struct (Model)
type User struct {
	ID             uuid.UUID `json:"id" readonly:"true"`
	Name           string    `json:"name,omitempty" validate:"omitempty"`
	Email          string    `json:"email" validate:"required,email"`
	PermissionRole string    `json:"role" validate:"required"`
}
