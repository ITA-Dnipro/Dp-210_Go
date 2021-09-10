package patientdata

import (
	"context"
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/user"
	"golang.org/x/crypto/bcrypt"
)

type PatientRepository interface {
	CreateRecord(ctx context.Context, pc *entity.Patient) error
	GetByEmail(ctx context.Context, email string) *entity.Patient
}

type Usecases struct {
	pdRepo   PatientRepository
	userRepo user.UsersRepository
}

func NewUsecases(patientRepo PatientRepository, ur user.UsersRepository) *Usecases {
	return &Usecases{
		pdRepo:   patientRepo,
		userRepo: ur,
	}
}

func (ucs *Usecases) CreateRecord(ctx context.Context, patient *entity.Patient) error {
	return ucs.pdRepo.CreateRecord(ctx, patient)
}

func (ucs *Usecases) GetByEmail(ctx context.Context, email string) *entity.Patient {
	return ucs.pdRepo.GetByEmail(ctx, email)
}

func (ucs *Usecases) CreateUserFromPatient(ctx context.Context, nu entity.NewUser, patient *entity.Patient) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generate password hash:%w", err)
	}

	u := entity.User{
		ID:             patient.Id.String(),
		Name:           patient.Name,
		Email:          patient.Email,
		PermissionRole: role.Patient,
		PasswordHash:   hash,
	}
	return ucs.userRepo.Create(ctx, u)
}
