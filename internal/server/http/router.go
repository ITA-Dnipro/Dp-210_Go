package http

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"net/http"
	"os"
	"strings"
	"time"

	cache "github.com/ITA-Dnipro/Dp-210_Go/internal/cache/redis"
	postgres "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/user"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
	handlers "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/handlers/user"
	handlePasw "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/handlers/user/password"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/service/auth"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/service/sender/mail"
	usecases "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/user"
	usecasesPasw "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/user/password"
	"github.com/go-redis/redis/v8"

	postgresDoctor "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/doctor"
	handlersDoctor "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/handlers/doctor"
	usecasesDoctor "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/doctor"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Auth interface {
	CreateToken(user auth.UserAuth) (auth.JwtToken, error)
	ValidateToken(t auth.JwtToken) (auth.UserAuth, error)
	InvalidateToken(userId string) error
}

// NewRouter create http routes.
func NewRouter(db *sql.DB, logger *zap.Logger, gmail *mail.GmailEmailSender, auth Auth, rdb *redis.Client) chi.Router {
	repo := postgres.NewRepository(db)
	repoD := postgresDoctor.NewRepository(db)

	usecase := usecases.NewUsecases(repo)

	mailSender := mail.NewPasswordCodeSender(gmail)

	paswCase := usecasesPasw.NewUsecases(
		mailSender,
		usecasesPasw.SixDigitGenerator{},
		repo,
		cache.NewRestoreCodeCache(rdb, time.Minute*5, "restore"),
	)

	md := &middleware.Middleware{Logger: logger, UserUC: usecase, Auth: auth}
	hs := handlers.NewHandlers(usecase, logger, auth)

	paswHandler := handlePasw.NewHandler(paswCase, logger, auth)

	usecaseD := usecasesDoctor.NewUsecases(repoD, repo)
	hsD := handlersDoctor.NewHandlers(usecaseD, logger)

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
				patientCard.SendToDb()
			})
			r.Post("/info", func(w http.ResponseWriter, r *http.Request) {
				var patientCard entity.PatientCard
				if err := json.NewDecoder(r.Body).Decode(&patientCard); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
					return
				}
				patientCard.CorrectData()
				patientCard.SendToDb()
			})
		})
		r.Route("/", func(r chi.Router) { // route with permissions
			r.Use(md.AuthMiddleware)

			r.Group(func(r chi.Router) { // route with permissions
				r.Use(md.RoleOnly(role.Operator, role.Admin))

				r.Get("/users", hs.GetUsers)     // GET /api/v1/users
				r.Get("/users/{id}", hs.GetUser) // GET /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8

				r.Post("/doctors", hsD.CreateDoctor) 		// POST	/api/v1/doctors
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
