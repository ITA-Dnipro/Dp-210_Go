package appointmen

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	"github.com/google/uuid"
)

// AppointmentsRepository represent doctor repository.
type AppointmentsRepository interface {
	GetWithFilter(ctx context.Context, filter entity.AppointmentFilter) ([]entity.Appointment, error)
	Create(ctx context.Context, a *entity.Appointment) error
	Delete(ctx context.Context, id string) error
}

// DoctorsClient represent doctor grpc client.
type DoctorsClient interface {
	GetByID(ctx context.Context, id string) (entity.Doctor, error)
}

//Producer represent Kafka producer.
type Producer interface {
	SendAppointment(a *entity.Appointment) error
	SendBill(a *entity.Appointment) error
	SendNotification(a *entity.Appointment) error
}

// NewUsecases create appoinments usecases.
func NewUsecases(ar AppointmentsRepository, dr DoctorsClient, producer Producer) *Usecases {
	return &Usecases{
		ar:       ar,
		dr:       dr,
		producer: producer,
	}
}

// Usecases represent a appointment usecases.
type Usecases struct {
	ar       AppointmentsRepository
	dr       DoctorsClient
	producer Producer
}

func (uc *Usecases) CreateRequest(ctx context.Context, a *entity.Appointment) error {
	// if a.From.Before(time.Now().UTC()) {
	// 	return fmt.Errorf("can't create appointment in past %s", a.From)
	// }
	// d, err := uc.dr.GetByID(ctx, a.DoctorID)
	// if err != nil {
	// 	return fmt.Errorf("can't find a doctor with %v id", a.DoctorID)
	// }
	a.To = a.From.Add(time.Minute * 30)
	// if a.To.After(d.EndAt) || a.From.Before(d.StartAt) {
	// 	return fmt.Errorf("doctor doesn't work %v - %v", a.From, a.To)
	// }
	a.ID = uuid.New().String()
	return uc.producer.SendAppointment(a)
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
	if err := uc.producer.SendNotification(a); err != nil {
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
	if err := uc.producer.SendBill(a); err != nil {
		return fmt.Errorf("create bill event")
	}
	return nil
}

func (uc *Usecases) GetWithFilter(ctx context.Context, f entity.AppointmentFilter) ([]entity.Appointment, error) {
	return uc.ar.GetWithFilter(ctx, f)
}
