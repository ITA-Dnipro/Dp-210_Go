package http

import (
	"database/sql"
	"google.golang.org/grpc"

	postgres "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/user"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
	handlers "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/handlers/user"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware"
	usecases "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/user"

	postgresDoctor "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/doctor"
	handlersDoctor "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/handlers/doctor"
	usecasesDoctor "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/doctor"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// NewRouter create http routes.
func NewRouter(db *sql.DB, logger *zap.Logger, conn *grpc.ClientConn) chi.Router {
	repo := postgres.NewRepository(db)
	repoD := postgresDoctor.NewRepository(db)

	usecase := usecases.NewUsecases(repo)

	md := middleware.New(logger, conn)
	hs := handlers.NewHandlers(usecase, logger)

	usecaseD := usecasesDoctor.NewUsecases(repoD, repo)
	hsD := handlersDoctor.NewHandlers(usecaseD, logger)

	r := chi.NewRouter()
	r.Use(md.LoggingMiddleware)
	r.Route("/api/v1", func(r chi.Router) {

		r.Post("/users", hs.CreateUser)       // POST /api/v1/users
		r.Get("/doctors", hsD.GetDoctors)     // GET    /api/v1/doctors
		r.Get("/doctors/{id}", hsD.GetDoctor) // GET /api/v1/doctors/<id>
		r.Route("/", func(r chi.Router) {     // route with permissions
			r.Use(md.AuthMiddleware)

			r.Group(func(r chi.Router) { // route with permissions
				r.Use(md.RoleOnly(role.Operator, role.Admin))

				r.Get("/users", hs.GetUsers)     // GET /api/v1/users
				r.Get("/users/{id}", hs.GetUser) // GET /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8

				r.Post("/doctors", hsD.CreateDoctor) // POST	/api/v1/doctors
			})
			r.Group(func(r chi.Router) { // route with permission Admin.
				r.Use(md.RoleOnly(role.Admin))

				r.Put("/users/{id}", hs.UpdateUser)    // PUT /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
				r.Delete("/users/{id}", hs.DeleteUser) // DELETE /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8

				r.Put("/doctors/{id}", hsD.UpdateDoctor)    // PUT    /api/v1/doctors/<id>
				r.Delete("/doctors/{id}", hsD.DeleteDoctor) // DELETE /api/v1/doctors/<id>
			})
		})
	})
	return r
}
