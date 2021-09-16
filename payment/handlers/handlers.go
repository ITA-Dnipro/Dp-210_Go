package handlers

import (
	"encoding/json"
	"errors"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/entity"

	"go.uber.org/zap"
)

type DataRepository interface {
	InsertDataToDb(*entity.Bill) error
	GetAllPatients() ([]entity.Patient, error)
}

func NewHandler(dataRepository DataRepository, log *zap.Logger) *Handler {
	return &Handler{
		dr:     dataRepository,
		Logger: log,
	}
}

type Handler struct {
	dr     DataRepository
	Logger *zap.Logger
}

func (h *Handler) InsertToDb(bill *entity.Bill) error {
	return h.dr.InsertDataToDb(bill)
}

func (h *Handler) SendMonthlyReportToTopic() (report []byte, err error) {
	patientsArr, err := h.dr.GetAllPatients()
	if err != nil {
		return nil, err
	}
	if len(patientsArr) == 0 {
		return nil, nil
	}
	bytesArr, err := json.Marshal(patientsArr)
	if err != nil {
		return nil, errors.New("marshal error")
	}
	return bytesArr, nil
}
