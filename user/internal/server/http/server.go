package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/config"
	handlers "github.com/ITA-Dnipro/Dp-210_Go/user/internal/server/http/user"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Server struct {
	srv    *http.Server
	cfg    config.Config
	logger *zap.Logger
}

// NewHTTPServer create new http server.
func NewHTTPServer(cfg config.Config, uc handlers.Usecase, logger *zap.Logger) *Server {
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/users", handlers.NewUserHandlers(uc, logger))
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
