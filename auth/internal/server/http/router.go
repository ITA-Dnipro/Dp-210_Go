package http

import (
	"database/sql"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/auth"
	cache "github.com/ITA-Dnipro/Dp-210_Go/auth/internal/cache/redis"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/usecase"

	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/repository/postgres"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/server/http/handlers"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/server/http/middleware"
	mail "github.com/ITA-Dnipro/Dp-210_Go/email_sender"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Auth interface {
	CreateToken(user auth.UserAuth) (auth.JwtToken, error)
	ValidateToken(t auth.JwtToken) (auth.UserAuth, error)
	InvalidateToken(userId string) error
}

// NewRouter create http routes.
func NewRouter(db *sql.DB, logger *zap.Logger, gmail *mail.GmailEmailSender, rdb *redis.Client) (chi.Router, error) {
	repo := postgres.NewRepository(db)
	expire := time.Minute * 15
	jwt, err := auth.NewJwtAuth(cache.NewSessionCache(rdb, expire, "jwtToken"), expire)
	if err != nil {
		return nil, err
	}
	mailSender := mail.NewPasswordCodeSender(gmail)

	md := &middleware.Middleware{Logger: logger}

	paswCase := usecase.NewUsecases(
		mailSender,
		usecase.SixDigitGenerator{},
		repo,
		cache.NewRestoreCodeCache(rdb, time.Minute*5, "restore"),
	)

	hs := handlers.NewHandler(paswCase, logger, jwt)

	r := chi.NewRouter()
	r.Use(md.LoggingMiddleware)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", hs.LogIn)

		r.Group(func(r chi.Router) {
			r.Use(md.AuthMiddleware)
			r.Post("/logout", hs.LogOut)
		})

		r.Route("/usecase", func(r chi.Router) {
			r.Route("/restore", func(r chi.Router) {
				r.Post("/code/send", hs.SendRestorePasswordCode)
				r.Post("/code/check", hs.RestorePassword)
			})

			r.Group(func(r chi.Router) {
				r.Use(md.AuthMiddleware)
				r.Post("/change", hs.ChangePassword)
			})
		})
	})
	return r, nil
}
