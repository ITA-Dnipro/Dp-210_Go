package appointment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/customerrors"
	"github.com/google/uuid"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const (
	idKey = "id"
)

type AppointmentsResponse struct {
	Appointments []entity.Appointment `json:"data"`
	Cursor       string               `json:"cursor"`
}

// Usecase represent appointment usecase.
type Usecase interface {
	GetByPatientID(ctx context.Context, id uuid.UUID, p *entity.AppointmentsParam) ([]entity.Appointment, string, error)
	GetByDoctorID(ctx context.Context, id uuid.UUID, p *entity.AppointmentsParam) ([]entity.Appointment, string, error)
	GetAll(ctx context.Context, p *entity.AppointmentsParam) ([]entity.Appointment, string, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.Appointment, error)
	CreateRequest(ctx context.Context, a *entity.Appointment) error
	Create(ctx context.Context, a *entity.Appointment) error
	Update(ctx context.Context, a *entity.Appointment) error
	Delete(ctx context.Context, id uuid.UUID) error
	SendResult(ctx context.Context, v *entity.Visit) error
}

// Handlers represent a appointment handlers.
type Handlers struct {
	usecase Usecase
	logger  *zap.Logger
}

// NewHandlers create new appointment handlers.
func NewHandlers(uc Usecase, logger *zap.Logger) http.Handler {
	h := &Handlers{usecase: uc, logger: logger}
	r := chi.NewRouter()
	r.Post("/", h.Create)                     // Post /api/v1/appointment
	r.Get("/", h.GetAll)                      // GET /api/v1/appointments
	r.Put("/{id}", h.Update)                  // PUT /api/v1/appointments/6ba7b810-9dad-11d1-80b4-00c04fd430c8
	r.Get("/{id}", h.GetByID)                 // GET /api/v1/appointments/6ba7b810-9dad-11d1-80b4-00c04fd430c8
	r.Delete("/{id}", h.Delete)               // DELETE /api/v1/appointments/6ba7b810-9dad-11d1-80b4-00c04fd430c8
	r.Post("/{id}/result", h.SendResult)      // DELETE /api/v1/appointments/6ba7b810-9dad-11d1-80b4-00c04fd430c8
	r.Get("/doctors/{id}", h.GetByDoctorID)   // GET /api/v1/appointments/doctors/6ba7b810-9dad-11d1-80b4-00c04fd430c8
	r.Get("/patients/{id}", h.GetByPatientID) // GET /api/v1/appointments/patient/6ba7b810-9dad-11d1-80b4-00c04fd430c8
	return r
}

// CreateUser Add new appointment.
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	var a entity.Appointment
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse a appointment", w)
		return
	}

	if err := validator.New().Struct(&a); err != nil {
		fmt.Println(err)
		h.writeErrorResponse(http.StatusBadRequest, "appointment data invalid", w)
		return
	}

	if err := h.usecase.CreateRequest(r.Context(), &a); err != nil {
		if errors.Is(err, customerrors.ErrBadParamInput) {
			h.writeErrorResponse(http.StatusBadRequest, err.Error(), w)
			return
		}
		h.logger.Error("can't create a appointment", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("appointment has been send", zap.String(idKey, a.ID.String()))
	h.render(w, a)
}

//Get all get all appointments.
func (h *Handlers) GetAll(w http.ResponseWriter, r *http.Request) {
	var p entity.AppointmentsParam
	if err := queryParameters(&p, r); err != nil {
		h.logger.Error("can't get appointments", zap.Error(err))
		h.writeErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}
	ap, cursor, err := h.usecase.GetAll(r.Context(), &p)
	if err != nil {
		h.logger.Error("can't get appointments", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("get all appointments")
	h.render(w, AppointmentsResponse{Appointments: ap, Cursor: cursor})
}

//GetByID get appointments by id.
func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, idKey)) // Gets params
	if err != nil {
		h.writeErrorResponse(http.StatusNotFound, "appointment id invalid", w)
		return
	}
	a, err := h.usecase.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("can't get appointments", zap.Error(err))
		h.writeErrorResponse(http.StatusNotFound, err.Error(), w)
		return
	}
	h.logger.Info("get all appointments")
	h.render(w, a)
}

//GetByDoctorID get appointments by doctor id.
func (h *Handlers) GetByDoctorID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, idKey)) // Gets params
	if err != nil {
		h.writeErrorResponse(http.StatusNotFound, "doctor id invalid", w)
		return
	}
	p := entity.AppointmentsParam{}
	if err := queryParameters(&p, r); err != nil {
		h.logger.Error("can't get appointments", zap.Error(err))
		h.writeErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}
	ap, cursor, err := h.usecase.GetByDoctorID(r.Context(), id, &p)
	if err != nil {
		h.logger.Error("can't get appointments", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("get all appointments")
	h.render(w, AppointmentsResponse{Appointments: ap, Cursor: cursor})
}

//GetByPatientID get appointments by doctor id.
func (h *Handlers) GetByPatientID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, idKey)) // Gets params
	if err != nil {
		h.writeErrorResponse(http.StatusNotFound, "patient id invalid", w)
		return
	}
	p := entity.AppointmentsParam{}
	if err := queryParameters(&p, r); err != nil {
		h.logger.Error("can't get appointments", zap.Error(err))
		h.writeErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}
	ap, cursor, err := h.usecase.GetByPatientID(r.Context(), id, &p)
	if err != nil {
		h.logger.Error("can't get appointments", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("get all appointments")
	h.render(w, AppointmentsResponse{Appointments: ap, Cursor: cursor})
}

// DeleteAppointment deletes a appointment from storage
func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, idKey)) // Gets params
	if err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "id parse uuid", w)
		return
	}

	if err := h.usecase.Delete(r.Context(), id); err != nil {
		h.logger.Error("can't delete", zap.Error(err))
		if errors.Is(err, customerrors.ErrNotFound) {
			h.writeErrorResponse(http.StatusNotFound,
				fmt.Sprintf("can't find a appointment with %v id", id), w)
			return
		}

		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.render(w, Message{"deleted"})
}

// Update update appointment.
func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, idKey)) // Gets params
	if err != nil {
		h.writeErrorResponse(http.StatusNotFound, "id parse uuid", w)
		return
	}

	a := entity.Appointment{ID: id}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse a appointment", w)
		return
	}

	if err := validator.New().Struct(&a); err != nil {
		fmt.Println(err)
		h.writeErrorResponse(http.StatusBadRequest, "appointment data invalid", w)
		return
	}

	if err := h.usecase.Update(r.Context(), &a); err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			h.writeErrorResponse(http.StatusNotFound,
				fmt.Sprintf("can't find a appointment with %v id", a.ID), w)
			return
		}
		h.logger.Error("can't update a appointment", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("appointment has been send", zap.String(idKey, a.ID.String()))
	h.render(w, a)
}

func (h *Handlers) SendResult(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, idKey)) // Gets params
	if err != nil {
		h.writeErrorResponse(http.StatusNotFound, "id parse uuid", w)
		return
	}
	var v entity.Visit
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse result", w)
		return
	}

	if err := validator.New().Struct(&v); err != nil {
		fmt.Println(err)
		h.writeErrorResponse(http.StatusBadRequest, "result data invalid", w)
		return
	}
	v.AppointmentID = id
	if err := h.usecase.SendResult(r.Context(), &v); err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			h.writeErrorResponse(http.StatusBadRequest, err.Error(), w)
			return
		}
		h.logger.Error("can't send result", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("result has been send",
		zap.String("appointment id", v.AppointmentID.String()),
	)
	h.render(w, v)
}

func queryParameters(p *entity.AppointmentsParam, r *http.Request) (err error) {
	layout := "2006-01-02T15:04:05"
	query := r.URL.Query()
	p.Cursor = query.Get("cursor")
	if query.Has("limit") {
		p.Limits, err = strconv.Atoi(query.Get("limit"))
		if err != nil {
			return
		}
	}
	if query.Has("from") {
		p.From, err = time.Parse(layout, query.Get("from"))
		if err != nil {
			return
		}
	}
	if query.Has("to") {
		p.To, err = time.Parse(layout, query.Get("to"))
		if err != nil {
			return
		}
	}
	return nil
}

// Message represent error message.
type Message struct {
	Msg string
}

func (h *Handlers) writeErrorResponse(code int, msg string, w http.ResponseWriter) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(Message{msg}); err != nil {
		h.logger.Error("can't encode error data", zap.Error(err))
	}
}

func (h *Handlers) render(w http.ResponseWriter, data interface{}) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	if err := enc.Encode(data); err != nil {
		h.logger.Error("can't encode data", zap.Error(err))
	}
}
