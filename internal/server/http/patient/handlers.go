package patient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/customerrors"

	//Do i need auth here, probably yes
	//auth "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware/auth"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type PatientUsecases interface {
	Create(ctx context.Context, p *entity.Patient) error
	GetByID(ctx context.Context, id string) (entity.Patient, error)
	GetAll(ctx context.Context) ([]entity.Patient, error)
	Delete(ctx context.Context, id string) error
}

const idKey = "id"

type Handlers struct {
	usecases PatientUsecases
	logger   *zap.Logger
}

func NewHandlers(uc PatientUsecases, log *zap.Logger) *Handlers {
	return &Handlers{usecases: uc, logger: log}
}

func (h *Handlers) GetPatients(w http.ResponseWriter, r *http.Request) {
	patients, err := h.usecases.GetAll(r.Context())
	if err != nil {
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	h.logger.Info("ger all request succeeded")
	h.render(w, patients)
}

func (h *Handlers) GetPatient(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idKey) // Gets params
	patient, err := h.usecases.GetByID(r.Context(), id)
	if err == nil {
		h.render(w, patient)
		return
	}

	h.logger.Error("can't get a patient", zap.Error(err))

	// if you wan't you can set content type of the headers directly here
	w.Header().Set("Content-Type", "application/json")
	h.writeErrorResponse(http.StatusNotFound,
		fmt.Sprintf("can't find a patient with %v id", id), w)
}

func (h *Handlers) CreatePatient(w http.ResponseWriter, r *http.Request) {
	var patient entity.Patient
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse a patien", w)
		return
	}
	if ok := isRequestValid(&patient); !ok {
		h.writeErrorResponse(http.StatusBadRequest, "patient data invalid", w)
		return
	}
	if err := h.usecases.Create(r.Context(), &patient); err != nil {
		h.logger.Error("can't create a patient", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("patient has been created", zap.String(idKey, patient.ID))
	h.render(w, patient)
}

func (h *Handlers) DeletePatient(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idKey) // Gets params
	if err := h.usecases.Delete(r.Context(), id); err != nil {
		h.logger.Error("can't delete", zap.Error(err))
		if errors.Is(err, customerrors.NotFound) {
			h.writeErrorResponse(http.StatusNotFound,
				fmt.Sprintf("can't find a patient with %v id", id), w)
			return
		}

		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	h.render(w, Message{"deleted"})
}

func isRequestValid(nu *entity.Patient) bool {
	validate := validator.New()
	err := validate.Struct(nu)
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
