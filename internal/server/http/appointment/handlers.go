package appointment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/customerrors"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const idKey = "id"

// UsersUsecase represent appointment usecases.
type UsersUsecase interface {
	GetByUser(ctx context.Context, userID string) ([]entity.Appointment, error)
	GetByPatientID(ctx context.Context, id string) ([]entity.Appointment, error)
	GetByDoctorID(ctx context.Context, id string) ([]entity.Appointment, error)
	//GetAll(ctx context.Context) (res []entity.Appointment, err error)
	CreateRequest(ctx context.Context, a *entity.Appointment) error
	Delete(ctx context.Context, id string) error
}

// Handlers represent a user handlers.
type Handlers struct {
	usecase UsersUsecase
	logger  *zap.Logger
}

// NewHandlers create new user handlers.
func NewHandlers(uc UsersUsecase, log *zap.Logger) *Handlers {
	return &Handlers{usecase: uc, logger: log}
}

// CreateUser Add new appointment.
func (h *Handlers) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var a entity.Appointment
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse a appointment", w)
		return
	}
	if ok := isRequestValid(&a); !ok {
		h.writeErrorResponse(http.StatusBadRequest, "appointment data invalid", w)
		return
	}
	id, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(http.StatusBadRequest, "context user data invalid", w)
		return
	}
	a.PatientID = id
	a.DoctorID = chi.URLParam(r, idKey) // Gets params
	if err := h.usecase.CreateRequest(r.Context(), &a); err != nil {
		h.logger.Error("can't create a appointment", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("user has been created", zap.String(idKey, a.ID))
	h.render(w, a)
}

func (h *Handlers) GetAppointments(w http.ResponseWriter, r *http.Request) {
	id, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(http.StatusBadRequest, "context user data invalid", w)
		return
	}
	a, err := h.usecase.GetByUser(r.Context(), id)
	if err != nil {
		h.logger.Error("can't create a appointment", zap.Error(err))
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

func isRequestValid(a *entity.Appointment) bool {
	validate := validator.New()
	err := validate.Struct(a)
	return err == nil
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
