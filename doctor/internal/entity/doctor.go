package entity

import (
	"time"

	"github.com/google/uuid"
)

// Doctor struct (Model) and request
type Doctor struct {
	ID         uuid.UUID `json:"id" validate:"required"`
	FirstName  string    `json:"name" validate:"required"`
	LastName   string    `json:"last_name" validate:"omitempty"`
	Speciality string    `json:"speciality" validate:"required"`
	StartAt    time.Time `json:"start_at"`
	EndAt      time.Time `json:"end_at"`
}
