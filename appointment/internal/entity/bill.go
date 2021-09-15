package entity

import (
	"github.com/google/uuid"
)

type Bill struct {
	DoctorID  uuid.UUID `json:"doctor_id"`
	PatientID uuid.UUID `json:"patient_id"`
	Price     int64     `json:"price"`
}
