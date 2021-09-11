package appointment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/role"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/http/customerrors"
	"github.com/google/uuid"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/http/middleware"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const idKey = "id"

// Usecase represent appointment usecase.
type Usecase interface {
	GetByFilter(ctx context.Context, filter entity.AppointmentFilter) ([]entity.Appointment, error)
	CreateRequest(ctx context.Context, a *entity.Appointment) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// Handlers represent a user handlers.
type Handlers struct {
	usecase Usecase
	logger  *zap.Logger
}

type NewAppointment struct {
	From time.Time `json:"from" validate:"required"`
}

// NewHandlers create new user handlers.
func NewHandlers(uc Usecase, logger *zap.Logger) http.Handler {
	h := &Handlers{usecase: uc, logger: logger}
	md := middleware.NewMiddleware(logger)
	r := chi.NewRouter()
	r.Use(md.AuthMiddleware)
	r.Group(func(r chi.Router) { // route with permissions
		r.Use(md.RoleOnly(role.Patient))
		r.Post("/doctors/{id}", h.CreateAppointment) // Post /api/v1/appointment
	})
	r.Group(func(r chi.Router) { // route with permissions
		r.Use(md.RoleOnly(role.Patient, role.Doctor))
		r.Get("/", h.GetAppointments) // GET /api/v1/appointments
	})
	r.Group(func(r chi.Router) { // route with permission Admin.
		r.Use(md.RoleOnly(role.Admin))
		r.Delete("/{id}", h.DeleteAppointment) // DELETE /api/v1/appointments/6ba7b810-9dad-11d1-80b4-00c04fd430c8
	})

	return r
}

// CreateUser Add new appointment.
func (h *Handlers) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var na NewAppointment
	if err := json.NewDecoder(r.Body).Decode(&na); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse a appointment", w)
		return
	}
	if err := validator.New().Struct(&na); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "appointment data invalid", w)
		return
	}
	id, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(http.StatusBadRequest, "context user data invalid", w)
		return
	}

	doctorID, err := uuid.Parse(chi.URLParam(r, idKey)) // Gets params
	if err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "id parse uuid", w)
		return
	}

	a := entity.Appointment{
		PatientID: id,
		DoctorID:  doctorID,
		From:      na.From,
	}

	if err := h.usecase.CreateRequest(r.Context(), &a); err != nil {
		h.logger.Error("can't create a appointment", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("appointment has been send", zap.String(idKey, a.ID.String()))
	h.render(w, a)
}

func (h *Handlers) GetAppointments(w http.ResponseWriter, r *http.Request) {
	id, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(http.StatusBadRequest, "context user data invalid", w)
		return
	}

	filter, err := getFilter(r)
	if err != nil {
		h.logger.Error("filter params", zap.Error(err))
		h.writeErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	userRole, ok := middleware.UserRoleFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(http.StatusBadRequest, "context user data invalid", w)
		return
	}

	if userRole == role.Doctor {
		filter.DoctorID = &id
	}

	if userRole == role.Patient {
		filter.PatientID = &id
	}

	a, err := h.usecase.GetByFilter(r.Context(), filter)
	if err != nil {
		h.logger.Error("can't get appointments", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("get all appointments by user id", zap.String(idKey, id.String()))
	h.render(w, a)
}

// DeleteAppointment deletes a appointment from storage
func (h *Handlers) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, idKey)) // Gets params
	if err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "id parse uuid", w)
		return
	}

	if err := h.usecase.Delete(r.Context(), id); err != nil {
		h.logger.Error("can't delete", zap.Error(err))
		if errors.Is(err, customerrors.NotFound) {
			h.writeErrorResponse(http.StatusNotFound,
				fmt.Sprintf("can't find a user with %v id", id), w)
			return
		}

		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.render(w, Message{"deleted"})
}

func getFilter(r *http.Request) (entity.AppointmentFilter, error) {
	filter := entity.AppointmentFilter{}
	const (
		fromKey      = "from"
		toKey        = "to"
		doctorIDKey  = "doctorID"
		patientIDKey = "patientID"
	)

	query := r.URL.Query()

	if query.Has(fromKey) {
		from, err := time.Parse(time.RFC3339, query.Get(fromKey))
		if err != nil {
			return filter, err
		}
		filter.From = &from
	}

	if query.Has(toKey) {
		to, err := time.Parse(time.RFC3339, query.Get(toKey))
		if err != nil {
			return filter, err
		}
		filter.To = &to
	}

	if query.Has(doctorIDKey) {
		doctorID, err := uuid.Parse(query.Get(doctorIDKey))
		if err != nil {
			return filter, err
		}
		filter.DoctorID = &doctorID
	}

	return filter, nil
}

// Message represent error message.
type Message struct {
	Msg string
}

func (*Handlers) writeErrorResponse(code int, msg string, w http.ResponseWriter) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Message{msg})
}

func (h *Handlers) render(w http.ResponseWriter, data interface{}) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(data)
}
