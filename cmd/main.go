package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres"
	router "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http"
	"github.com/ITA-Dnipro/Dp-210_Go/visits/config"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

// Main function
func main() {
	logger, _ := zap.NewProduction()
	var cfg config.Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal(fmt.Errorf("parsing config: %w", err))
	}
	db, err := postgres.Open(cfg.Postgres)
	if err != nil {
		log.Fatal(fmt.Errorf("connecting to db: %w", err))
	}
	defer func() {
		log.Printf("visits: Database Stopping")
		db.Close()
	}()

	migrationsPath := "migrations"
	err = postgres.MigrateUp(migrationsPath, cfg.DatabaseStr())
	if err != nil {
		log.Fatal(fmt.Errorf("db migrations: %w", err))
	}

	logger.Info("starting web server")
	r := router.NewRouter(db, logger)
	// Start server
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", r))
}
