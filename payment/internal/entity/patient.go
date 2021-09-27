package entity

import "github.com/google/uuid"

type Patient struct {
	PatientId    uuid.UUID
	PatientTotal int64
}
