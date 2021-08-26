package appointmen

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/kafka"
	"github.com/google/uuid"
)

// DoctorsRepository represent doctor repository.
type AppointmentsRepository interface {
	GetBeforeTime(ctx context.Context, t time.Time) ([]entity.Appointment, error)
	GetByPatientID(ctx context.Context, id string) ([]entity.Appointment, error)
	GetByDoctorID(ctx context.Context, id string) ([]entity.Appointment, error)
	GetByUserID(ctx context.Context, id string) ([]entity.Appointment, error)
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

//Events represent Kafka events.
type Events interface {
	Emit(topic string, payload interface{}) error
}

// NewUsecases create appoinments usecases.
func NewUsecases(ar AppointmentsRepository, dr DoctorsRepository, pr PatientsRepository, events Events) *Usecases {
	return &Usecases{
		ar:     ar,
		dr:     dr,
		pr:     pr,
		events: events,
	}
}

// Usecases represent a appointment usecases.
type Usecases struct {
	ar     AppointmentsRepository
	dr     DoctorsRepository
	pr     PatientsRepository
	events Events
}

func (uc *Usecases) CreateRequest(ctx context.Context, a *entity.Appointment) error {
	d, err := uc.dr.GetByID(ctx, a.DoctorID)
	if err != nil {
		return fmt.Errorf("can't find a doctor with %v id", a.DoctorID)
	}
	a.To = a.From.Add(time.Minute * 30)
	// if a.To.After(d.EndAt) || a.From.Before(d.StartAt) {
	// 	return fmt.Errorf("doctor doesn't work %v - %v", a.From, a.To)
	// }
	_ = d
	a.ID = uuid.New().String()
	return uc.events.Emit(kafka.AppoinmentTopic, a)
}
func (uc *Usecases) CreateFromEvent(payload []byte) error {
	var a entity.Appointment
	if err := json.Unmarshal(payload, &a); err != nil {
		return fmt.Errorf("marshaling appointment:%w", err)
	}
	return uc.Create(context.Background(), &a)
}

// Create Add new appointment.
func (uc *Usecases) Create(ctx context.Context, a *entity.Appointment) error {

	if err := uc.ar.Create(ctx, a); err != nil {
		return fmt.Errorf("creating appointment:%w", err)
	}
	if err := uc.events.Emit(kafka.MailTopic, "sent mail to user email"); err != nil {
		return fmt.Errorf("create mail event")
	}
	return nil
}

// Delete deletes a appointment from storage.
func (uc *Usecases) Delete(ctx context.Context, id string) error {
	return uc.ar.Delete(ctx, id)
}

// Delete deletes a appointment from storage.
func (uc *Usecases) DeleteWithBilling(ctx context.Context, a *entity.Appointment) error {
	if err := uc.ar.Delete(ctx, a.ID); err != nil {
		return fmt.Errorf("deleting appointment:%w", err)
	}
	if err := uc.events.Emit(kafka.BillTopic, "bill to user id"); err != nil {
		return fmt.Errorf("create bill event")
	}
	return nil
}
func (u *Usecases) DeleteOld(ctx context.Context, t time.Time) error {
	appoinments, err := u.ar.GetBeforeTime(ctx, t)
	if err != nil {
		return fmt.Errorf("geting appointment before time:%w", err)
	}
	for _, a := range appoinments {
		if err := u.DeleteWithBilling(ctx, &a); err != nil {
			return fmt.Errorf("delet appoinments:%w", err)
		}
	}
	return nil
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
func (uc *Usecases) GetByUser(ctx context.Context, id string) ([]entity.Appointment, error) {
	return uc.ar.GetByUserID(ctx, id)
}

// GetAll get all appointments.
func (uc *Usecases) GetAll(ctx context.Context) (res []entity.Appointment, err error) {
	return uc.ar.GetAll(ctx)
}
