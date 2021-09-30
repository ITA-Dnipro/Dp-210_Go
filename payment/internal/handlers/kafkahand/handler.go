package kafkahand

import (
	"context"
	"encoding/json"
	"errors"
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
		Dr:         dataRepository,
		Logger:     log,
		statClient: sc,
	}
}

type Handler struct {
	Dr         DataRepository
	Logger     *zap.Logger
	statClient st.StatClient
}

func (h *Handler) InsertToDb(bill *entity.Bill) error {
	return h.Dr.InsertDataToDb(bill)
}

func (h *Handler) SendMonthlyReport(ctx context.Context) ([]byte, error) {
	docPat, err := h.Dr.GetAllDataFromDb()
	if err != nil {
		return nil, err
	}
	if len(docPat.Patients) == 0 || len(docPat.Doctors) == 0 {
		return nil, errors.New("lists from db empty")
	}

	doctorsBytes, _ := json.Marshal(docPat.Doctors)
	_, _ = h.statClient.DocStat(ctx, &st.DocRequest{DocsBytesArr: doctorsBytes})

	report, _ := json.Marshal(docPat.Patients)

	return report, nil
}
