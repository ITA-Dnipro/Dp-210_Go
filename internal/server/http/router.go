package router

import (
	"database/sql"
	"time"

	restore "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/restore"
	postgres "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/user"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
	handlePasw "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/handlers/user/password"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware"
	handlers "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/user"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/service/auth"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/service/sender/mail"
	usecases "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/user"
	usecasesPasw "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/user/password"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Auth interface {
	CreateToken(user auth.UserAuth, lifetime time.Duration) (auth.JwtToken, error)
	ValidateToken(t auth.JwtToken) (auth.UserAuth, error)
}

// NewRouter create http routes.
func NewRouter(db *sql.DB, logger *zap.Logger, gmail *mail.GmailEmailSender, auth Auth) chi.Router {
	repo := postgres.NewRepository(db)
	usecase := usecases.NewUsecases(repo)

	mailSender := mail.NewPasswordCodeSender(gmail)

	paswCase := usecasesPasw.NewUsecases(mailSender, usecasesPasw.SixDigitGenerator{}, repo, restore.NewCodeRepo(db))

	hs := handlers.NewHandlers(usecase, logger, auth)
	paswHandler := handlePasw.NewHandler(paswCase, logger, auth)
	md := &middleware.Middleware{Logger: logger, UserUC: usecase}

	r := chi.NewRouter()
	r.Use(md.LoggingMiddleware)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", hs.GetToken) // POST /api/v1/login

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

		r.Post("/users", hs.CreateUser)   // POST /api/v1/users
		r.Route("/", func(r chi.Router) { // route with permissions
			r.Use(md.AuthMiddleware)

			r.Group(func(r chi.Router) { // route with permissions
				r.Use(md.RoleOnly(role.Operator, role.Admin))

				r.Get("/users", hs.GetUsers)     // GET /api/v1/users
				r.Get("/users/{id}", hs.GetUser) // GET /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
			})
			r.Group(func(r chi.Router) { // route with permission Admin.
				r.Use(md.RoleOnly(role.Admin))

				r.Put("/users/{id}", hs.UpdateUser)    // PUT /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
				r.Delete("/users/{id}", hs.DeleteUser) // DELETE /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
			})
		})
	})
	return r
}
