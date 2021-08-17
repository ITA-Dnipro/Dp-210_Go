package router

import (
	"context"
	"database/sql"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	postgresUser "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/user"

	postgresDoctor "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/doctor"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware"

	handlersDoctor "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/doctor"
	handlersUser "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/user"

	usecasesDoctor "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/doctor"
	usecasesUser "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/user"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// NewRouter create http routes.
func NewRouter(db *sql.DB, logger *zap.Logger) chi.Router {
	repoU := postgresUser.NewRepository(db)
	usecaseU := usecasesUser.NewUsecases(repoU)
	hsU := handlersUser.NewHandlers(usecaseU, logger)
	mdU := &middleware.Middleware{Logger: logger, UserUC: usecaseU}

	repoD := postgresDoctor.NewRepository(db)
	usecaseD := usecasesDoctor.NewUsecases(repoD)
	hsD := handlersDoctor.NewHandlers(usecaseD, logger)

	// TODO remove. for testing purpose.
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.MinCost)
	repoU.Create(context.Background(), entity.User{
		ID:             "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		Name:           "admin",
		Email:          "admin@admin.com",
		PasswordHash:   hash,
		PermissionRole: role.Admin,
	})
	hash, _ = bcrypt.GenerateFromPassword([]byte("operator"), bcrypt.MinCost)
	repoU.Create(context.Background(), entity.User{
		ID:             "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		Name:           "operator",
		Email:          "operator@admin.com",
		PasswordHash:   hash,
		PermissionRole: role.Operator,
	})
	hash, _ = bcrypt.GenerateFromPassword([]byte("user"), bcrypt.MinCost)
	repoU.Create(context.Background(), entity.User{
		ID:             "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		Name:           "test",
		Email:          "test@admin.com",
		PasswordHash:   hash,
		PermissionRole: role.Viewer,
	})

	r := chi.NewRouter()
	r.Use(mdU.LoggingMiddleware)
	r.Route("/api/v1", func(r chi.Router) {
		//Anyone capabilities
		r.Post("/login", hsU.GetToken)    // POST /api/v1/login
		r.Post("/users", hsU.CreateUser)  // POST /api/v1/users
		r.Get("/doctors", hsD.GetDoctors) // GET    /api/v1/doctors

		//Tmp
		r.Post("/doctors", hsD.CreateDoctor) // POST	/api/v1/doctors

		r.Route("/", func(r chi.Router) { // route with permissions
			r.Use(mdU.AuthMiddleware)

			//O + A capabilities
			r.Group(func(r chi.Router) { // route with permissions
				r.Use(mdU.RoleOnly(role.Operator, role.Admin))

				r.Get("/users", hsU.GetUsers)     // GET /api/v1/users
				r.Get("/users/{id}", hsU.GetUser) // GET /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
			})

			//A only capabilities
			r.Group(func(r chi.Router) { // route with permission Admin.
				r.Use(mdU.RoleOnly(role.Admin))
				//User
				r.Put("/users/{id}", hsU.UpdateUser)    // PUT    /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
				r.Delete("/users/{id}", hsU.DeleteUser) // DELETE /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8

				//Doctor
				r.Put("/doctors/{id}", hsD.UpdateDoctor)    // PUT    /api/v1/doctors/<id>
				r.Delete("/doctors/{id}", hsD.DeleteDoctor) // DELETE /api/v1/doctors/<id>
				//r.Post("/doctors", hsD.CreateDoctor)			// POST	/api/v1/doctors
			})
		})
	})
	return r
}
