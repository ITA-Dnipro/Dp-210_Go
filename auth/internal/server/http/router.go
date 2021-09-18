package http

import (
	"database/sql"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/config"
	"time"

	cache "github.com/ITA-Dnipro/Dp-210_Go/auth/internal/cache/redis"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/usecase"

	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/repository/postgres"
	mail "github.com/ITA-Dnipro/Dp-210_Go/auth/internal/sender"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/server/http/handlers"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/server/http/middleware"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Auth interface {
	CreateToken(user usecase.UserAuth) (usecase.JwtToken, error)
	InvalidateToken(userId string) error
	ValidateToken(t usecase.JwtToken) (usecase.UserAuth, error)
}

func NewRouter(db *sql.DB, logger *zap.Logger, gmail *mail.GmailEmailSender, rdb *redis.Client, auth Auth) (chi.Router, error) {
	repo := postgres.NewRepository(db)

	mailSender := mail.NewPasswordCodeSender(gmail)

	md := &middleware.Middleware{Logger: logger, Auth: auth}

	cfg := config.GetConfig()
	exp := time.Duration(cfg.RestoreCodeExpirationMillis) * time.Millisecond

	paswCase := usecase.NewUsecases(
		mailSender,
		usecase.SixDigitGenerator{},
		repo,
		cache.NewRestoreCodeCache(rdb, exp, cfg.RestoreCodeType),
	)

	hs := handlers.NewHandler(paswCase, logger, auth)

	r := chi.NewRouter()
	r.Use(md.LoggingMiddleware)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", hs.LogIn)

		r.Group(func(r chi.Router) {
			r.Use(md.AuthMiddleware)
			r.Post("/logout", hs.LogOut)
		})

		r.Route("/password", func(r chi.Router) {
			r.Route("/restore", func(r chi.Router) {
				r.Post("/", hs.RestorePassword)
				r.Post("/code", hs.SendRestorePasswordCode)
			})

			r.Group(func(r chi.Router) {
				r.Use(md.AuthMiddleware)
				r.Post("/change", hs.ChangePassword)
			})
		})
	})

	return r, nil
}
