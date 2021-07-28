package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Metalscreame/go-training/day_6/networking-handlers/internal/entity"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/jarcoal/httpmock"
	"go.uber.org/zap"
)

func TestName(t *testing.T) {
	testCases := []struct {
		name         string
		wantError    bool
		mock         func(ctl *gomock.Controller) BookRepository
		bodyToSend   entity.Book
		expectedCode int
	}{
		{
			name: "testCase 1 - get all success",
			mock: func(ctl *gomock.Controller) BookRepository {
				r := NewMockBookRepository(ctl)
				r.EXPECT().GetAll(gomock.Any())
				return r
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			log, _ := zap.NewProduction()
			srv := NewServer(tc.mock(ctl), log)

			router := mux.NewRouter()
			router.HandleFunc("/books", srv.GetBooks).Methods(http.MethodPost)

			body, err := json.Marshal(tc.bodyToSend)
			if err != nil {
				t.Errorf("Can't encode body")
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			router.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Wanted code %v but got %v", tc.expectedCode, w.Code)
			}
		})
	}

}
