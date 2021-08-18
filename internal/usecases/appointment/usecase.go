package appointmen

import (
	"context"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
	"github.com/google/uuid"
)

// DoctorsRepository represent doctor repository.
type AppointmentsRepository interface {
	GetByPatientID(ctx context.Context, id string) ([]entity.Appointment, error)
	GetByDoctorID(ctx context.Context, id string) ([]entity.Appointment, error)
	Create(ctx context.Context, a *entity.Appointment) error
	GetAll(ctx context.Context) ([]entity.Appointment, error)
	Delete(ctx context.Context, id string) error
}

// DoctorsRepository represent doctor repository.
type DoctorsRepository interface {
	GetByID(ctx context.Context, id string) (entity.Doctor, error)
}

// PatientsRepository represent parient repository.
type PatientsRepository interface {
	GetByID(ctx context.Context, id string) (entity.Patient, error)
}

// NewUsecases create appoinments usecases.
func NewUsecases(ar AppointmentsRepository, dr DoctorsRepository, pr PatientsRepository) *Usecases {
	return &Usecases{
		ar: ar,
		dr: dr,
		pr: pr,
	}
}

// Usecases represent a appointment usecases.
type Usecases struct {
	ar AppointmentsRepository
	dr DoctorsRepository
	pr PatientsRepository
}

// Create Add new appointment.
func (uc *Usecases) Create(ctx context.Context, a *entity.Appointment) error {
	a.ID = uuid.New().String()
	a.To = a.From.Add(time.Minute * 30)
	return uc.ar.Create(ctx, a)
}

// Delete deletes a appointment from storage.
func (uc *Usecases) Delete(ctx context.Context, id string) error {
	return uc.ar.Delete(ctx, id)
}

// GetByDoctorID get appointmens by doctor id.
func (uc *Usecases) GetByDoctorID(ctx context.Context, id string) ([]entity.Appointment, error) {
	return uc.ar.GetByDoctorID(ctx, id)
}

// GetByPatientID get appointmens by patient id.
func (uc *Usecases) GetByPatientID(ctx context.Context, id string) ([]entity.Appointment, error) {
	return uc.ar.GetByDoctorID(ctx, id)
}

// GetByUser get appointmens by user.
func (uc *Usecases) GetByUser(ctx context.Context, userID string, userRole role.Role) ([]entity.Appointment, error) {
	if userRole == role.Doctor {
		return uc.ar.GetByDoctorID(ctx, userID)
	}
	if userRole == role.Patient {
		return uc.ar.GetByPatientID(ctx, userID)
	}
	return []entity.Appointment{}, nil
}

// GetAll get all appointments.
func (uc *Usecases) GetAll(ctx context.Context) (res []entity.Appointment, err error) {
	return uc.ar.GetAll(ctx)
}
