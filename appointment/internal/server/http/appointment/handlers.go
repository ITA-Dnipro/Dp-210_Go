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

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/http/middleware"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const idKey = "id"

// Usecase represent appointment usecase.
type Usecase interface {
	GetWithFilter(ctx context.Context, filter entity.AppointmentFilter) ([]entity.Appointment, error)
	CreateRequest(ctx context.Context, a *entity.Appointment) error
	Delete(ctx context.Context, id string) error
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
func NewHandlers(uc Usecase, logger *zap.Logger) *chi.Mux {
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
	a := entity.Appointment{
		PatientID: id,
		DoctorID:  chi.URLParam(r, idKey),
		From:      na.From,
	}

	if err := h.usecase.CreateRequest(r.Context(), &a); err != nil {
		h.logger.Error("can't create a appointment", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("appointment has been created", zap.String(idKey, a.ID))
	h.render(w, a)
}

func (h *Handlers) GetAppointments(w http.ResponseWriter, r *http.Request) {
	id, idOk := middleware.UserIDFromContext(r.Context())
	userRole, roleOk := middleware.UserRoleFromContext(r.Context())
	if !idOk || !roleOk {
		h.writeErrorResponse(http.StatusBadRequest, "context user data invalid", w)
		return
	}
	f := entity.AppointmentFilter{}
	if userRole == role.Doctor {
		f.DoctorID = &id
	}
	if userRole == role.Patient {
		f.PatientID = &id
	}
	a, err := h.usecase.GetWithFilter(r.Context(), f)
	if err != nil {
		h.logger.Error("can't get appointments", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("get all appointments by user id", zap.String(idKey, id))
	h.render(w, a)
}

// DeleteAppointment deletes a appointment from storage
func (h *Handlers) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idKey) // Gets params
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
