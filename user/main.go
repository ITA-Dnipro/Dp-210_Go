package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/repository/postgres"
	router "github.com/ITA-Dnipro/Dp-210_Go/user/internal/server/http"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

const (
	migrationsPath = "sql/migrations"
)

func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("building logger", err)
	}
	if err := run(zapLogger); err != nil {
		zapLogger.Error("user: error:", zap.Error(err))
		os.Exit(1)
	}
}

func run(logger *zap.Logger) error {
	var cfg config.Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("read env: %w", err))
	}
	logger.Info("user: Initializing database support")
	db, err := postgres.Open(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		log.Printf("user: Database Stopping")
		db.Close()
	}()
	err = postgres.MigrateUp(migrationsPath, cfg.Postgres)
	if err != nil {
		log.Fatal(fmt.Errorf("db migrations: %w", err))
		return fmt.Errorf("migrations db: %w", err)
	}
	r := router.NewRouter(db, logger)
	// Start server
	logger.Info("user: Initializing API support", zap.String("host", cfg.APIHost))
	return http.ListenAndServe(cfg.APIHost, r)
}
