package patientdata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"

	"github.com/go-playground/validator/v10"
)

const idKey = "id"

type PatientUsecases interface {
	CreateRecord(ctx context.Context, u *entity.Patient) error
	GetByEmail(ctx context.Context, email string) *entity.Patient
	CreateUserFromPatient(ctx context.Context, u entity.NewUser, p *entity.Patient) error
}

type Handlers struct {
	patientUsecases PatientUsecases
	logger          *zap.Logger
}

func NewHandlers(pdUcs PatientUsecases, log *zap.Logger) *Handlers {
	return &Handlers{patientUsecases: pdUcs, logger: log}
}

func (h *Handlers) CreatePatientDataFromJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	patient := new(entity.Patient)
	patient.AdditionalInfo()
	if err := json.NewDecoder(r.Body).Decode(patient); err != nil {
		h.logger.Error("can't create a patient", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	if err := isValid(patient); err != nil {
		h.logger.Error("can't validate the patient", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	if patientFromDb := h.patientUsecases.GetByEmail(ctx, patient.Email); patientFromDb != nil {
		h.logger.Error("db error", zap.Error(errors.New(fmt.Sprintf("patient with email {%s} already exist", patient.Email))))
		h.writeErrorResponse(http.StatusInternalServerError, fmt.Sprintf("patient with email {%s} already exist", patient.Email), w)
		return
	}

	if err := h.patientUsecases.CreateUserFromPatient(ctx, entity.NewUser{Password: "dp210go"}, patient); err != nil {
		h.logger.Error("can't create a user from patient", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	if err := h.patientUsecases.CreateRecord(ctx, patient); err != nil {
		h.logger.Error("can't create a patient", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("patient and user have been created", zap.String(idKey, patient.Id.String()))
	h.render(w, patient)
}

func (h *Handlers) CreatePatientDataFromCSV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bytes, _ := io.ReadAll(r.Body)
	csvStrData := strings.Split(strings.TrimSpace(string(bytes)), "\n")[5]
	csvArr := strings.Split(csvStrData, ",")
	patient, err := entity.NewPatient(csvArr)
	if err != nil {
		h.logger.Error("can't create a patient", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	if err = isValid(patient); err != nil {
		h.logger.Error("can't validate the patient", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	if patientFromDb := h.patientUsecases.GetByEmail(ctx, patient.Email); patientFromDb != nil {
		h.logger.Error("db error", zap.Error(errors.New(fmt.Sprintf("patient with email {%s} already exist", patient.Email))))
		h.writeErrorResponse(http.StatusInternalServerError, fmt.Sprintf("patient with email {%s} already exist", patient.Email), w)
		return
	}

	if err := h.patientUsecases.CreateUserFromPatient(ctx, entity.NewUser{Password: "dp210go"}, patient); err != nil {
		h.logger.Error("can't create a user from patient", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	if err = h.patientUsecases.CreateRecord(r.Context(), patient); err != nil {
		h.logger.Error("can't create a patient", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("patient has been created", zap.String(idKey, patient.Id.String()))
	h.render(w, patient)
}

type Message struct {
	Msg string
}

func (*Handlers) writeErrorResponse(code int, msg string, w http.ResponseWriter) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(Message{msg})
}

func (h *Handlers) render(w http.ResponseWriter, data interface{}) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	_ = enc.Encode(data)
}

func isValid(p *entity.Patient) error {
	validate := validator.New()
	return validate.Struct(*p)
}
