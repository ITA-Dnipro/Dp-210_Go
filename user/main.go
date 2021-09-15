package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/repository/postgres"
	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/server"

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
		return fmt.Errorf("read env: %w", err)
	}
	logger.Info("user: Initializing database support")
	db, err := postgres.Open(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		logger.Info("user: Database Stopping")
		db.Close()
	}()
	err = postgres.MigrateUp(migrationsPath, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("migrations db: %w", err)
	}
	return server.Serve(cfg, db, logger)
}
