package router

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"net/http"
	"os"
	"strings"

	postgres "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/user"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
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
		r.Post("/login", hs.GetToken)   // POST /api/v1/login
		r.Post("/users", hs.CreateUser) // POST /api/v1/users
		r.Route("/patient", func(r chi.Router) {
			r.Post("/card", func(w http.ResponseWriter, r *http.Request) {
				records := make([]string, 0)
				scanner := bufio.NewScanner(r.Body)
				for scanner.Scan() {
					records = append(records, scanner.Text())
				}
				if scanner.Err() != nil {
					fmt.Println(scanner.Err())
				}

				if len(records) != 7 {
					_, _ = fmt.Fprintln(os.Stderr, "no data or can't proceed")
					return
				}

				patientCard := entity.New(strings.Split(records[5], ";"))
				fmt.Printf("%#v\n", *patientCard)
				// Send patientCard to accept by Operator.
			})
			r.Post("/info", func(w http.ResponseWriter, r *http.Request) {
				var patientCard entity.PatientCard
				if err := json.NewDecoder(r.Body).Decode(&patientCard); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
					return
				}
				patientCard.CorrectData()
				fmt.Printf("%#v\n", patientCard)
				// Send patientCard to accept by Operator.
			})
		})
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
