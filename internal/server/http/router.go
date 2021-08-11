package router

import (
	"context"
	"database/sql"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	postgres "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/user"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware"
	handlers "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/user"
	usecases "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/user"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// NewRouter create http routes.
func NewRouter(db *sql.DB, logger *zap.Logger) chi.Router {
	repo := postgres.NewRepository(db)
	usecase := usecases.NewUsecases(repo)
	hs := handlers.NewHandlers(usecase, logger)
	md := &middleware.Middleware{Logger: logger, UserUC: usecase}
	// TODO remove. for testing purpose.
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.MinCost)
	repo.Create(context.Background(), entity.User{
		ID:             "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		Name:           "admin",
		Email:          "admin@admin.com",
		PasswordHash:   hash,
		PermissionRole: role.Admin,
	})
	hash, _ = bcrypt.GenerateFromPassword([]byte("operator"), bcrypt.MinCost)
	repo.Create(context.Background(), entity.User{
		ID:             "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		Name:           "operator",
		Email:          "operator@admin.com",
		PasswordHash:   hash,
		PermissionRole: role.Operator,
	})
	hash, _ = bcrypt.GenerateFromPassword([]byte("user"), bcrypt.MinCost)
	repo.Create(context.Background(), entity.User{
		ID:             "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		Name:           "test",
		Email:          "test@admin.com",
		PasswordHash:   hash,
		PermissionRole: role.Viewer,
	})

	r := chi.NewRouter()
	r.Use(md.LoggingMiddleware)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", hs.GetToken) // POST /api/v1/login

		r.Route("/password", func(r chi.Router) {
			r.Post("/restore", hs.RestorePassword)
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
