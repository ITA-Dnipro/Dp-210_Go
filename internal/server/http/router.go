package router

import (
	"database/sql"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	postgres "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/user"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware"
	handlers "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/user"
	usecases "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/user"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// NewRouter create http routes.
func NewRouter(db *sql.DB, logger *zap.Logger) chi.Router {
	repo := postgres.NewRepository(db)
	usecase := usecases.NewUsecases(repo)
	hs := handlers.NewHandlers(usecase, logger)
	md := &middleware.Middleware{Logger: logger, UserUC: usecase}

	r := chi.NewRouter()
	r.Use(md.LoggingMiddleware)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/users", hs.CreateUser)
		r.Route("/users", func(r chi.Router) { // route with permission
			r.Use(md.RoleOnly(entity.Operator, entity.Admin))
			r.Get("/", hs.GetUsers)          // GET /api/v1/users
			r.Get("/{id}", hs.GetUser)       // GET /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
			r.Put("/{id}", hs.UpdateUser)    // PUT /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
			r.Delete("/{id}", hs.DeleteUser) // DELETE /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
		})
	})
	return r
}
