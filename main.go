package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres"
	router "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http"
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
	logger.Info("starting web server")
	r := router.NewRouter(db, logger)
	// Start server
	log.Fatal(http.ListenAndServe("localhost:8000", r))
}
