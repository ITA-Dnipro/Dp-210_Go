package doctor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/server/http/customerrors"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// DoctorsUsecases represent doctor usecases.
type DoctorsUsecases interface {
	GetByID(ctx context.Context, id string) (entity.Doctor, error)
	Update(ctx context.Context, u *entity.Doctor) error
	Create(ctx context.Context, u *entity.Doctor) error
	GetAll(ctx context.Context) ([]entity.Doctor, error)
	Delete(ctx context.Context, id string) error
}

const idKey = "id"

// Handlers represent a doctor handlers.
type Handlers struct {
	usecases DoctorsUsecases
	logger   *zap.Logger
}

// NewHandlers create new doctor handlers.
func NewHandlers(uc DoctorsUsecases, log *zap.Logger) *Handlers {
	return &Handlers{usecases: uc, logger: log}
}

// GetDoctors Get all doctors.
func (h *Handlers) GetDoctors(w http.ResponseWriter, r *http.Request) {
	doctors, err := h.usecases.GetAll(r.Context())
	if err != nil {
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	h.logger.Info("ger all request succeeded")
	h.render(w, doctors)
}

// GetDoctor Get single doctor by id.
func (h *Handlers) GetDoctor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idKey) // Gets params
	doctor, err := h.usecases.GetByID(r.Context(), id)
	if err == nil {
		h.render(w, doctor)
		return
	}

	h.logger.Error("can't get a doctor", zap.Error(err))

	// if you wan't you can set content type of the headers directly here
	w.Header().Set("Content-Type", "application/json")
	h.writeErrorResponse(http.StatusNotFound,
		fmt.Sprintf("can't find a doctor with %v id", id), w)
}

// CreateDoctor Add new doctor
func (h *Handlers) CreateDoctor(w http.ResponseWriter, r *http.Request) {
	var d entity.Doctor
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse a doctor", w)
		return
	}
	if ok := isRequestValid(&d); !ok {
		h.writeErrorResponse(http.StatusBadRequest, "doctor data invalid", w)
		return
	}
	if err := h.usecases.Create(r.Context(), &d); err != nil {
		h.logger.Error("can't create a doctor", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("doctor has been created", zap.String(idKey, d.ID))
	h.render(w, d)
}

// UpdateDoctor updates a doctor
func (h *Handlers) UpdateDoctor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idKey) // Gets params
	var d entity.Doctor
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse a doctor", w)
		return
	}
	d.ID = id
	if ok := isRequestValid(&d); !ok {
		h.writeErrorResponse(http.StatusBadRequest, "doctor data invalid", w)
		return
	}
	if err := h.usecases.Update(r.Context(), &d); err != nil {
		h.logger.Error("can't update a doctor", zap.Error(err))
		if errors.Is(err, customerrors.NotFound) {
			h.writeErrorResponse(http.StatusNotFound,
				fmt.Sprintf("can't find a doctor with %v id", d.ID), w)
			return
		}

		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	h.render(w, d)
}

// DeleteDoctor deletes a doctor from storage
func (h *Handlers) DeleteDoctor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idKey) // Gets params
	if err := h.usecases.Delete(r.Context(), id); err != nil {
		h.logger.Error("can't delete", zap.Error(err))
		if errors.Is(err, customerrors.NotFound) {
			h.writeErrorResponse(http.StatusNotFound,
				fmt.Sprintf("can't find a doctor with %v id", id), w)
			return
		}

		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	h.render(w, Message{"deleted"})
}

func isRequestValid(nu *entity.Doctor) bool {
	validate := validator.New()
	err := validate.Struct(nu)
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
