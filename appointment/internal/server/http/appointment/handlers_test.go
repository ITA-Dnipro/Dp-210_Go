package appointment

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/customerrors"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/http/appointment/mocks"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandlers_Create(t *testing.T) {
	tests := []struct {
		name         string
		wantError    bool
		mock         func(ctl *gomock.Controller) Usecase
		bodyToSend   entity.Appointment
		expectedCode int
	}{
		{
			name: "testCase 1 - create success",
			mock: func(ctl *gomock.Controller) Usecase {
				r := mocks.NewMockUsecase(ctl)
				r.EXPECT().CreateRequest(gomock.Any(), gomock.Any())
				return r
			},
			bodyToSend: entity.Appointment{
				DoctorID: uuid.New(), PatientID: uuid.New(), From: time.Now().Add(time.Minute),
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "testCase 2 - not valid appointment",
			mock: func(ctl *gomock.Controller) Usecase {
				r := mocks.NewMockUsecase(ctl)
				return r
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "testCase 3 - create appointment internal error",
			mock: func(ctl *gomock.Controller) Usecase {
				r := mocks.NewMockUsecase(ctl)
				r.EXPECT().CreateRequest(gomock.Any(), gomock.Any()).Return(errors.New("test"))
				return r
			},
			bodyToSend: entity.Appointment{
				DoctorID: uuid.New(), PatientID: uuid.New(), From: time.Now().Add(time.Minute),
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "testCase 4 - create appointment bad request error",
			mock: func(ctl *gomock.Controller) Usecase {
				r := mocks.NewMockUsecase(ctl)
				r.EXPECT().CreateRequest(gomock.Any(), gomock.Any()).Return(customerrors.ErrBadParamInput)
				return r
			},
			bodyToSend: entity.Appointment{
				DoctorID: uuid.New(), PatientID: uuid.New(), From: time.Now().Add(time.Minute),
			},
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			log, _ := zap.NewProduction()
			router := chi.NewRouter()
			router.Mount("/appointments", NewHandlers(tt.mock(ctl), log))

			body, err := json.Marshal(tt.bodyToSend)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/appointments", bytes.NewReader(body))

			router.ServeHTTP(w, r)
			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestHandlers_GetAll(t *testing.T) {
	type fields struct {
		usecase Usecase
		logger  *zap.Logger
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				usecase: tt.fields.usecase,
				logger:  tt.fields.logger,
			}
			h.GetAll(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlers_GetByID(t *testing.T) {
	type fields struct {
		usecase Usecase
		logger  *zap.Logger
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				usecase: tt.fields.usecase,
				logger:  tt.fields.logger,
			}
			h.GetByID(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlers_GetByDoctorID(t *testing.T) {
	type fields struct {
		usecase Usecase
		logger  *zap.Logger
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				usecase: tt.fields.usecase,
				logger:  tt.fields.logger,
			}
			h.GetByDoctorID(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlers_GetByPatientID(t *testing.T) {
	type fields struct {
		usecase Usecase
		logger  *zap.Logger
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				usecase: tt.fields.usecase,
				logger:  tt.fields.logger,
			}
			h.GetByPatientID(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlers_Delete(t *testing.T) {
	type fields struct {
		usecase Usecase
		logger  *zap.Logger
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				usecase: tt.fields.usecase,
				logger:  tt.fields.logger,
			}
			h.Delete(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlers_Update(t *testing.T) {
	type fields struct {
		usecase Usecase
		logger  *zap.Logger
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				usecase: tt.fields.usecase,
				logger:  tt.fields.logger,
			}
			h.Update(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlers_SendResult(t *testing.T) {
	type fields struct {
		usecase Usecase
		logger  *zap.Logger
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				usecase: tt.fields.usecase,
				logger:  tt.fields.logger,
			}
			h.SendResult(tt.args.w, tt.args.r)
		})
	}
}

func Test_queryParameters(t *testing.T) {
	type args struct {
		p *entity.AppointmentsParam
		r *http.Request
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, queryParameters(tt.args.p, tt.args.r))
		})
	}
}
