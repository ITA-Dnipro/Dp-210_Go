package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/zap"

	appointmentHandlers "github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/http/appointment"
)

type Usecase interface {
	GetByFilter(ctx context.Context, filter entity.AppointmentFilter) ([]entity.Appointment, error)
	CreateRequest(ctx context.Context, a *entity.Appointment) error
	Delete(ctx context.Context, id uuid.UUID) error
}
type Server struct {
	srv    *http.Server
	cfg    config.Config
	logger *zap.Logger
}

// NewHTTPServer create new http server.
func NewHTTPServer(cfg config.Config, uc Usecase, logger *zap.Logger) *Server {
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/appointments", appointmentHandlers.NewHandlers(uc, logger))
	})

	return &Server{srv: &http.Server{
		Addr:         cfg.APIHost,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}, cfg: cfg, logger: logger}
}

func (s *Server) ListenAndServe() error {
	s.logger.Info(fmt.Sprintf("startup http server:%s", s.cfg.APIHost))
	return s.srv.ListenAndServe()
}

func (s *Server) GracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		s.srv.Close()
		s.logger.Error("could not stop server gracefully: %w", zap.Error(err))
	}
	s.logger.Info("http server shutdown")
}
