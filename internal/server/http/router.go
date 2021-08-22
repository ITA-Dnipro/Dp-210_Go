package router

import (
	"database/sql"

	userRepo "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/user"
	userHandlers "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/user"
	userUsecases "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/user"

	patientRepo "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/patient"
	patientHandlers "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/patient"
	patientUsecases "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/patient"

	doctorRepo "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/doctor"
	doctorHandlers "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/doctor"
	doctorUsecases "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/doctor"

	appointmentRepo "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/appointment"
	appointmentHandlers "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/appointment"
	appointmentUsecases "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/appointment"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// NewRouter create http routes.
func NewRouter(db *sql.DB, logger *zap.Logger) chi.Router {
	ur := userRepo.NewRepository(db)
	dr := doctorRepo.NewRepository(db)
	pr := patientRepo.NewRepository(db)
	ar := appointmentRepo.NewRepository(db)

	uc := userUsecases.NewUsecases(ur)
	dc := doctorUsecases.NewUsecases(dr, ur)
	pc := patientUsecases.NewUsecases(pr, ur)
	ac := appointmentUsecases.NewUsecases(ar, dr, pr)

	uh := userHandlers.NewHandlers(uc, logger)
	ph := patientHandlers.NewHandlers(pc, logger)
	dh := doctorHandlers.NewHandlers(dc, logger)
	ah := appointmentHandlers.NewHandlers(ac, logger)

	md := &middleware.Middleware{Logger: logger}

	r := chi.NewRouter()
	r.Use(md.LoggingMiddleware)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", uh.GetToken)     // POST /api/v1/login
		r.Post("/users", uh.CreateUser)   // POST /api/v1/users
		r.Route("/", func(r chi.Router) { // route with permissions
			r.Use(md.AuthMiddleware)
			r.Group(func(r chi.Router) { // route with permissions
				r.Use(md.RoleOnly(role.Patient, role.Admin))
				r.Get("/doctors", dh.GetDoctors)     // GET /api/v1/doctors
				r.Get("/doctors/{id}", dh.GetDoctor) // GET /api/v1/doctors/6ba7b810-9dad-11d1-80b4-00c04fd430c8
			})
			r.Group(func(r chi.Router) { // route with permissions
				r.Use(md.RoleOnly(role.Patient))
				r.Post("/appointments", ah.CreateAppointment) // Post /api/v1/appointment
			})
			r.Group(func(r chi.Router) { // route with permissions
				r.Use(md.RoleOnly(role.Doctor, role.Admin))
				r.Get("/patients", ph.GetPatients)     // GET /api/v1/patients
				r.Get("/patients/{id}", ph.GetPatient) // GET /api/v1/patients/6ba7b810-9dad-11d1-80b4-00c04fd430c8
			})
			r.Group(func(r chi.Router) { // route with permissions
				r.Use(md.RoleOnly(role.Patient, role.Doctor, role.Admin))
				r.Get("/appointments", ah.GetAppointments) // GET /api/v1/appointments
			})
			r.Group(func(r chi.Router) { // route with permissions
				r.Use(md.RoleOnly(role.Operator, role.Admin))
				r.Get("/users", uh.GetUsers)                 // GET /api/v1/users
				r.Get("/users/{id}", uh.GetUser)             // GET /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
				r.Post("/doctors", dh.CreateDoctor)          // POST /api/v1/doctors
				r.Put("/doctors/{id}", dh.UpdateDoctor)      // PUT /api/v1/doctors//6ba7b810-9dad-11d1-80b4-00c04fd430c8
				r.Delete("/doctors/{id}", dh.DeleteDoctor)   // DELET /api/v1/doctors/6ba7b810-9dad-11d1-80b4-00c04fd430c8
				r.Post("/patients", ph.CreatePatient)        // POST /api/v1/patients
				r.Delete("/patients/{id}", ph.DeletePatient) // DELET /api/v1/patients//6ba7b810-9dad-11d1-80b4-00c04fd430c8
				//r.Get("appointments/doctors/{id}", ah.GetAppointmentsByDoctorID)   // Post /api/v1/appointment/doctors/6ba7b810-9dad-11d1-80b4-00c04fd430c8
				//r.Get("appointments/patients/{id}", ah.GetAppointmentsByPatientID) // Post /api/v1/appointment/doctors/6ba7b810-9dad-11d1-80b4-00c04fd430c8
			})
			r.Group(func(r chi.Router) { // route with permission Admin.
				r.Use(md.RoleOnly(role.Admin))
				r.Put("/users/{id}", uh.UpdateUser)    // PUT /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
				r.Delete("/users/{id}", uh.DeleteUser) // DELETE /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
			})
		})
	})
	return r
}
