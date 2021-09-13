package appointmen

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/customerrors"
	"github.com/google/uuid"
)

// AppointmentsRepository represent doctor repository.
type AppointmentsRepository interface {
	GetByPatientID(ctx context.Context, id uuid.UUID, al *entity.AppointmentList) error
	GetByDoctorID(ctx context.Context, id uuid.UUID, al *entity.AppointmentList) error
	GetByID(ctx context.Context, id uuid.UUID) (entity.Appointment, error)
	GetAll(ctx context.Context, al *entity.AppointmentList) error
	Create(ctx context.Context, a *entity.Appointment) error
	Update(ctx context.Context, a *entity.Appointment) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// DoctorsClient represent doctor grpc client.
type DoctorsClient interface {
	GetByID(ctx context.Context, id uuid.UUID) (entity.Doctor, error)
}

//Producer represent Kafka producer.
type Producer interface {
	SendNotification(n interface{}) error
	SendAppointment(a *entity.Appointment) error
	SendBill(b entity.Bill) error
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
	// 	return fmt.Errorf("can't find a doctor with %v id, %w", a.DoctorID, err)
	// }
	a.To = a.From.Add(time.Minute * 30)
	// if a.To.After(d.EndAt) || a.From.Before(d.StartAt) {
	// 	return fmt.Errorf("doctor doesn't work %v - %v", a.From, a.To)
	// }
	a.ID = uuid.New()
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
func (uc *Usecases) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.ar.Delete(ctx, id)
}

// Create Add new appointment.
func (uc *Usecases) Update(ctx context.Context, a *entity.Appointment) error {
	a.To = a.From.Add(time.Minute * 30)
	if err := uc.ar.Update(ctx, a); err != nil {
		return fmt.Errorf("creating appointment:%w", err)
	}
	if err := uc.producer.SendNotification(a); err != nil {
		return fmt.Errorf("create mail event")
	}
	return nil
}

// Delete deletes a appointment from storage.
func (uc *Usecases) SendResult(ctx context.Context, v *entity.Visit) error {
	a, err := uc.ar.GetByID(ctx, v.AppointmentID)
	if err != nil {
		return customerrors.ErrNotFound
	}
	v.AppointmentID = a.ID
	v.DoctorID = a.DoctorID
	v.PatientID = a.PatientID
	if err := uc.ar.Delete(ctx, a.ID); err != nil {
		return fmt.Errorf("deleting appointment:%w", err)
	}
	b := entity.Bill{
		DoctorID:  v.DoctorID,
		PatientID: v.PatientID,
		Price:     v.Price,
	}
	if err := uc.producer.SendBill(b); err != nil {
		return fmt.Errorf("send bill events")
	}
	if err := uc.producer.SendNotification(v); err != nil {
		return fmt.Errorf("send bill events")
	}
	return nil
}

func (uc *Usecases) GetAll(ctx context.Context, al *entity.AppointmentList) error {
	return uc.ar.GetAll(ctx, al)
}

func (uc *Usecases) GetByDoctorID(ctx context.Context, id uuid.UUID, al *entity.AppointmentList) error {
	return uc.ar.GetByDoctorID(ctx, id, al)
}

func (uc *Usecases) GetByPatientID(ctx context.Context, id uuid.UUID, al *entity.AppointmentList) error {
	return uc.ar.GetByPatientID(ctx, id, al)
}

func (uc *Usecases) GetByID(ctx context.Context, id uuid.UUID) (entity.Appointment, error) {
	return uc.ar.GetByID(ctx, id)
}
