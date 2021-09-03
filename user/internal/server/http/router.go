package http

import (
	"database/sql"

	postgres "github.com/ITA-Dnipro/Dp-210_Go/user/internal/repository/postgres/user"
	handlers "github.com/ITA-Dnipro/Dp-210_Go/user/internal/server/http/user"
	usecases "github.com/ITA-Dnipro/Dp-210_Go/user/internal/usecases/user"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// NewRouter create http routes.
func NewRouter(db *sql.DB, logger *zap.Logger) chi.Router {
	repo := postgres.NewRepository(db)
	usecase := usecases.NewUsecases(repo)
	hs := handlers.NewHandlers(usecase, logger)
	r := chi.NewRouter()
	//r.Use(md.LoggingMiddleware)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/users", hs.CreateUser) // POST /api/v1/users

		// TODO: remove for test purpose.
		r.Get("/users", hs.GetUsers)           // GET /api/v1/users
		r.Get("/users/{id}", hs.GetUser)       // GET /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
		r.Put("/users/{id}", hs.UpdateUser)    // PUT /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
		r.Delete("/users/{id}", hs.DeleteUser) // DELETE /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8

		// r.Route("/", func(r chi.Router) { // route with permissions
		// 	r.Use(md.AuthMiddleware)
		// 	r.Group(func(r chi.Router) { // route with permissions
		// 		r.Use(md.RoleOnly(role.Operator, role.Admin))
		// 		r.Get("/users", hs.GetUsers)     // GET /api/v1/users
		// 		r.Get("/users/{id}", hs.GetUser) // GET /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
		// 	})
		// 	r.Group(func(r chi.Router) { // route with permission Admin.
		// 		r.Use(md.RoleOnly(role.Admin))
		// 		r.Put("/users/{id}", hs.UpdateUser)    // PUT /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
		// 		r.Delete("/users/{id}", hs.DeleteUser) // DELETE /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
		// 	})
		// })
	})
	return r
}
