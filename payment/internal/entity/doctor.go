package entity

import "github.com/google/uuid"

type Doctor struct {
	DoctorId    uuid.UUID `json:"doctor_id"`
	DoctorTotal int64     `json:"doctor_total"`
}
