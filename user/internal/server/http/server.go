package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/entity"
	handlers "github.com/ITA-Dnipro/Dp-210_Go/user/internal/server/http/user"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// UsersUsecases represent user usecases.
type Usecase interface {
	Create(ctx context.Context, u *entity.User) error
	Update(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id string) (entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Delete(ctx context.Context, id string) error
	Authenticate(ctx context.Context, email, password string) (entity.User, error)
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
