package patientdata

import (
	"context"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
)

type PtDataRepository interface {
	CreateRecord(ctx context.Context, pc *entity.PatientData) error
}

type Usecases struct {
	pdRepo PtDataRepository
}

func NewUsecases(pdRepo PtDataRepository) *Usecases {
	return &Usecases{
		pdRepo: pdRepo,
	}
}

func (ucs *Usecases) CreateRecord(ctx context.Context, pd *entity.PatientData) error {
	return ucs.pdRepo.CreateRecord(ctx, pd)
}
