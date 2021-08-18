package entity

import (
	"time"
)

type Appointment struct {
	ID        string    `json:"id" readonly:"true"`
	DoctorID  string    `json:"doctor_id"`
	PatientID string    `json:"patient_id"`
	Reason    string    `json:"reason"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to" readonly:"true"`
}
