package stathand

import (
	"encoding/json"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/entity"
	"go.uber.org/zap"
	"sort"
)

func NewHandler(log *zap.Logger) *Handler {
	return &Handler{
		Logger: log,
	}
}

type Handler struct {
	Logger *zap.Logger
}

func (h *Handler) DoctorsUnmarshal(bytes []byte) (doctors []entity.Doctor, err error) {
	err = json.Unmarshal(bytes, &doctors)
	if err != nil {
		return
	}
	return
}

func (h *Handler) GetBest(doctors []entity.Doctor) (doctorsOut []entity.Doctor) {
	sort.Slice(doctors, func(i, j int) bool {
		return doctors[i].DoctorTotal > doctors[j].DoctorTotal
	})

	comparisonPoint := doctors[0].DoctorTotal
	for _, doctor := range doctors {
		if doctor.DoctorTotal == comparisonPoint {
			doctorsOut = append(doctorsOut, doctor)
		}
	}
	return
}
