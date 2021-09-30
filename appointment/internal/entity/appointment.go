package entity

import (
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	ID        uuid.UUID `json:"id" readonly:"true"`
	DoctorID  uuid.UUID `json:"doctor_id"           validate:"required"`
	PatientID uuid.UUID `json:"patient_id"          validate:"required"`
	Reason    string    `json:"reason"              validate:"omitempty"`
	From      time.Time `json:"from"                validate:"required"`
	To        time.Time `json:"to" readonly:"true"`
}

type AppointmentsParam struct {
	Limits int       `json:"limit"`
	Cursor string    `json:"cursor"`
	From   time.Time `json:"-"`
	To     time.Time `json:"-"`
}
