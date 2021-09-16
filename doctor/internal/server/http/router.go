package http

import (
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/repository/postgres/doctor"
	handlers "github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/server/http/doctor"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/server/http/middleware"

	//"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/service/auth"
	agc "github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/client/grpc/appointments"
	usecases "github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/usecases/doctor"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

//~ type Auth interface {
//~ CreateToken(user auth.UserAuth) (auth.JwtToken, error)
//~ ValidateToken(t auth.JwtToken) (auth.UserAuth, error)
//~ InvalidateToken(userId string) error
//~ }

// NewRouter create http routes.
func NewRouter(repo *doctor.Repository, usecases *usecases.Usecases, logger *zap.Logger, md *middleware.Middleware, agc *agc.Client) chi.Router {
	r := chi.NewRouter()
	hs := handlers.NewHandlers(usecases, logger, agc)
	//repo := postgres.NewRepository(db)
	//usecase := usecases.NewUsecases(repo)
	//md := &middleware.Middleware{Logger: logger, UserUC: usecase, Auth: auth}

	//r.Use(md.LoggingMiddleware)
	r.Route("/api/v1", func(r chi.Router) { // perm All

		r.Get("/doctors/appointment/{id}", hs.GetAppointment) // GET /api/v1/doctors/ppointment/<id>
		r.Get("/doctors", hs.GetDoctors)                      // GET    /api/v1/doctors
		r.Get("/doctors/{id}", hs.GetDoctor)                  // GET    /api/v1/doctors/<id>
		r.Post("/doctors", hs.CreateDoctor)                   // POST	/api/v1/doctors
		r.Put("/doctors/{id}", hs.UpdateDoctor)               // PUT    /api/v1/doctors/<id>
		r.Delete("/doctors/{id}", hs.DeleteDoctor)            // DELETE /api/v1/doctors/<id>
	})

	//r.Use(md.LoggingMiddleware)
	//r.Route("/api/v1", func(r chi.Router) {
	//r.Post("/users", hs.CreateUser) // POST /api/v1/users

	// TODO: remove for test purpose.

	//r.Get("/users", hs.GetUsers)           // GET /api/v1/users
	//r.Get("/users/{id}", hs.GetUser)       // GET /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
	//r.Put("/users/{id}", hs.UpdateUser)    // PUT /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
	//r.Delete("/users/{id}", hs.DeleteUser) // DELETE /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8

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
	//})
	return r
}
