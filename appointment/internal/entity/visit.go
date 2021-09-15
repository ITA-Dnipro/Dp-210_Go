package entity

import (
	"github.com/google/uuid"
)

type Visit struct {
	AppointmentID uuid.UUID `json:"appointment_id"`
	DoctorID      uuid.UUID `json:"doctor_id"`
	PatientID     uuid.UUID `json:"patient_id"`
	Price         int64     `json:"price"`
	Result        string    `json:"result" validate:"required"`
}
