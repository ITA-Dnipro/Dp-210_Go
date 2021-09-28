package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/http/middleware"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	ah "github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/http/appointment"
)

type Server struct {
	srv    *http.Server
	cfg    config.Config
	logger *zap.Logger
}

// NewHTTPServer create new http server.
func NewHTTPServer(cfg config.Config, uc ah.Usecase, logger *zap.Logger) *Server {
	r := chi.NewRouter()
	r.Mount("/metrics", promhttp.Handler())
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.MeasureResponseDuration)
		r.Mount("/appointments", ah.NewHandlers(uc, logger))
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
		s.logger.Error("could not stop server gracefully: %w", zap.Error(err))
		if err := s.srv.Close(); err != nil {
			s.logger.Error("could not close server: %w", zap.Error(err))
		}
	}
	s.logger.Info("http server shutdown")
}
