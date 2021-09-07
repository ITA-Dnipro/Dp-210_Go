package entity

import "time"

// Doctor struct (Model) and request
type Doctor struct {
	ID         string    `json:"id"         validate:"required"`
	FirstName  string    `json:"name"       validate:"required"`
	LastName   string    `json:"last_name"  validate:"omitempty"`
	Speciality string    `json:"speciality" validate:"required"`
	StartAt    time.Time `json:"start_at"   validate:"required"`
	EndAt      time.Time `json:"end_at"     validate:"required"`
}
