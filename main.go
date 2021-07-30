package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/middlware"
	postgres "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/user"
	handlers "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/user"
	usecases "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/user"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

// Main function
func main() {
	logger, _ := zap.NewProduction()
	dsn := "postgres://postgres:secret@0.0.0.0:5432/test?sslmode=disable&timezone=utc"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(fmt.Errorf("creating db: %w", err))
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatal(fmt.Errorf("ping db %s : %w", dsn, err))
	}

	migrationsPath := "migrations"
	err = postgres.MigrateUp(migrationsPath, dsn)
	if err != nil {
		log.Fatal(fmt.Errorf("db migrations: %w", err))
	}
	repo := postgres.NewRepository(db)
	usecase := usecases.NewUsecases(repo)
	hs := handlers.NewHandlers(usecase, logger)

	md := &middlware.Middleware{Logger: logger}
	// Init router
	r := mux.NewRouter()

	// type Handler interface {
	//    ServeHTTP(ResponseWriter, *Request)
	//}
	//http.HandleFunc("/", h1)
	//	http.HandleFunc("/endpoint", h2)
	//https://golang.org/pkg/net/http/#HandleFunc

	// we can also use middleware
	r.Use(md.LoggingMiddleware)

	// Route handles & endpoints
	r.HandleFunc("/users", hs.GetUsers).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", hs.GetUser).Methods("GET")
	r.HandleFunc("/users", hs.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", hs.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", hs.DeleteUser).Methods("DELETE")
	logger.Info("starting web server")

	// Start server
	log.Fatal(http.ListenAndServe("localhost:8000", r))
}
