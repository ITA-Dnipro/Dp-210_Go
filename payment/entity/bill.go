package entity

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func NewBill(bill []byte) (*Bill, error) {
	var b Bill
	if err := json.Unmarshal(bill, &b); err != nil {
		return nil, fmt.Errorf("unmarshal error > %s", err)
	}
	return &b, nil
}

type Bill struct {
	DoctorID  uuid.UUID `json:"doctor_id"`
	PatientID uuid.UUID `json:"patient_id"`
	Price     int64     `json:"price"`
}
