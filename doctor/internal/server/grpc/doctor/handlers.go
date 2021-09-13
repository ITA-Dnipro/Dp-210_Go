package doctor

import (
	"context"

	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/entity"
	ds "github.com/ITA-Dnipro/Dp-210_Go/doctor/proto/doctors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DoctorsUsecases interface {
	GetByID(ctx context.Context, id string) (entity.Doctor, error)
}

type Handlers struct {
	usecases DoctorsUsecases
	logger   *zap.Logger
}

func NewHandlers(uc DoctorsUsecases, log *zap.Logger) *Handlers {
	return &Handlers{usecases: uc, logger: log}
}

func (h *Handlers) GetByID(ctx context.Context, req *ds.GetByIDReq) (*ds.GetByIDRes, error) {
	doctor, err := h.usecases.GetByID(ctx, req.GetDoctorID())
	if err != nil {
		return nil, err
	}
	doctorGRPC := &ds.Doctor{
		DoctorID:   doctor.ID,
		FirstName:  doctor.FirstName,
		LastName:   doctor.LastName,
		Speciality: doctor.Speciality,
		StartAt:    timestamppb.New(doctor.StartAt),
		EndAt:      timestamppb.New(doctor.EndAt),
	}

	return &ds.GetByIDRes{Doctor: doctorGRPC}, nil
}
