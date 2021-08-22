package appointment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const idKey = "id"

// UsersUsecases represent appointment usecases.
type UsersUsecases interface {
	GetByUser(ctx context.Context, userID string) ([]entity.Appointment, error)
	GetByPatientID(ctx context.Context, id string) ([]entity.Appointment, error)
	GetByDoctorID(ctx context.Context, id string) ([]entity.Appointment, error)
	GetAll(ctx context.Context) (res []entity.Appointment, err error)
	Create(ctx context.Context, a *entity.Appointment) error
	Delete(ctx context.Context, id string) error
}

// Handlers represent a user handlers.
type Handlers struct {
	usecases UsersUsecases
	logger   *zap.Logger
}

// NewHandlers create new user handlers.
func NewHandlers(uc UsersUsecases, log *zap.Logger) *Handlers {
	return &Handlers{usecases: uc, logger: log}
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
	if err := h.usecases.Create(r.Context(), &a); err != nil {
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
	a, err := h.usecases.GetByUser(r.Context(), id)
	if err != nil {
		h.logger.Error("can't create a appointment", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("get all appointments by user id", zap.String(idKey, id))
	h.render(w, a)
}

func isRequestValid(a *entity.Appointment) bool {
	validate := validator.New()
	err := validate.Struct(a)
	fmt.Println(err)
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
