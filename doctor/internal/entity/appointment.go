package entity

import (
	"time"
)

type Appointment struct {
	ID        string    `json:"id" readonly:"true"`
	DoctorID  string    `json:"doctor_id"           validate:"required"`
	PatientID string    `json:"patient_id"          validate:"required"`
	Reason    string    `json:"reason"              validate:"omitempty"`
	From      time.Time `json:"from"                validate:"required"`
	To        time.Time `json:"to" readonly:"true"`
}

type AppointmentFilter struct {
	DoctorID  *string
	PatientID *string
	From      *time.Time
	To        *time.Time
}
