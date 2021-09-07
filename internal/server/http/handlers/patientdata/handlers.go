package patientdata

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"
)

const idKey = "id"

type PatDataUsecases interface {
	CreateRecord(ctx context.Context, u *entity.PatientData) error
}

type Handlers struct {
	usecases PatDataUsecases
	logger   *zap.Logger
}

func NewHandlers(pdUcs PatDataUsecases, log *zap.Logger) *Handlers {
	return &Handlers{usecases: pdUcs, logger: log}
}

func (h *Handlers) CreatePatientDataFromCSV(w http.ResponseWriter, r *http.Request) {
	records := make([]string, 0)
	scanner := bufio.NewScanner(r.Body)
	for scanner.Scan() {
		records = append(records, scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Println(scanner.Err())
	}

	if len(records) != 7 {
		_, _ = fmt.Fprintln(os.Stderr, "no data or can't proceed")
		return
	}

	patientData := entity.NewPatientData(strings.Split(records[5], ";"))

	if err := h.usecases.CreateRecord(r.Context(), patientData); err != nil {
		h.logger.Error("can't create a patientData", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("patientData has been created", zap.String(idKey, patientData.Id.String()))
	h.render(w, patientData)
}

func (h *Handlers) CreatePatientDataFromJSON(w http.ResponseWriter, r *http.Request) {
	patientData := new(entity.PatientData)
	if err := json.NewDecoder(r.Body).Decode(&patientData); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	patientData.CorrectData()

	if err := h.usecases.CreateRecord(r.Context(), patientData); err != nil {
		h.logger.Error("can't create a patientData", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("patientData has been created", zap.String(idKey, patientData.Id.String()))
	h.render(w, patientData)
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
