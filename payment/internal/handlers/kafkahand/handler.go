package kafkahand

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/entity"
	st "github.com/ITA-Dnipro/Dp-210_Go/payment/internal/proto/statistics"

	"go.uber.org/zap"
)

type DataRepository interface {
	InsertDataToDb(*entity.Bill) error
	GetAllDataFromDb() (*entity.DocPat, error)
}

func NewHandler(dataRepository DataRepository, log *zap.Logger, sc st.StatClient) *Handler {
	return &Handler{
		dr:         dataRepository,
		Logger:     log,
		statClient: sc,
	}
}

type Handler struct {
	dr         DataRepository
	Logger     *zap.Logger
	statClient st.StatClient
}

func (h *Handler) InsertToDb(bill *entity.Bill) error {
	return h.dr.InsertDataToDb(bill)
}

func (h *Handler) SendMonthlyReport() ([]byte, error) {
	docPat, err := h.dr.GetAllDataFromDb()
	if err != nil {
		return nil, err
	}
	if len(docPat.Patients) == 0 || len(docPat.Doctors) == 0 {
		return nil, nil
	}

	doctorsBytes, err := json.Marshal(docPat.Doctors)
	if err != nil {
		return nil, fmt.Errorf("doctors marshal error > %s", err)
	}
	_, _ = h.statClient.DocStat(context.Background(), &st.DocRequest{DocsBytesArr: doctorsBytes})

	report, err := json.Marshal(docPat.Patients)
	if err != nil {
		return nil, fmt.Errorf("patients marshal error > %s", err)
	}

	return report, nil
}
