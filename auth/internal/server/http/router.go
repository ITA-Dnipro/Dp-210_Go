package http

import (
	"database/sql"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/service/auth"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Auth interface {
	CreateToken(user auth.UserAuth) (auth.JwtToken, error)
	ValidateToken(t auth.JwtToken) (auth.UserAuth, error)
	InvalidateToken(userId string) error
}

// NewRouter create http routes.
func NewRouter(db *sql.DB, logger *zap.Logger) chi.Router {
	// repo := postgres.NewRepository(db)
	// usecase := usecases.NewUsecases(repo)

	// mailSender := mail.NewPasswordCodeSender(gmail)

	// paswCase := usecasesPasw.NewUsecases(mailSender, usecasesPasw.SixDigitGenerator{}, repo, restore.NewCodeRepo(db))

	md := &middleware.Middleware{Logger: logger, Auth: auth}
	// hs := handlers.NewHandlers(usecase, logger, auth)

	// paswHandler := handlePasw.NewHandler(paswCase, logger, auth)

	r := chi.NewRouter()
	r.Use(md.LoggingMiddleware)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", hs.GetToken) // POST /api/v1/login

		r.Group(func(r chi.Router) {
			r.Use(md.AuthMiddleware)
			r.Post("/logout", hs.LogOut)
		})

		r.Route("/password", func(r chi.Router) {
			r.Route("/restore", func(r chi.Router) {
				r.Post("/code/send", paswHandler.SendRestorePasswordCode)
				r.Post("/code/check", paswHandler.CheckPasswordCode)
			})

			r.Group(func(r chi.Router) {
				r.Use(md.AuthMiddleware)
				r.Post("/change", paswHandler.ChangePassword)
			})
		})
	})
	return r
}
