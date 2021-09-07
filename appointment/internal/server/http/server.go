package http

import (
	"context"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	"github.com/go-chi/chi"
	"go.uber.org/zap"

	appointmentHandlers "github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/http/appointment"
)

type Usecase interface {
	GetWithFilter(ctx context.Context, filter entity.AppointmentFilter) ([]entity.Appointment, error)
	CreateRequest(ctx context.Context, a *entity.Appointment) error
	Delete(ctx context.Context, id string) error
}

// NewHTTPServer create new http server.
func NewHTTPServer(cfg config.Config, uc Usecase, logger *zap.Logger) *http.Server {
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/appointments", appointmentHandlers.NewHandlers(uc, logger))
	})

	return &http.Server{
		Addr:         cfg.APIHost,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}
