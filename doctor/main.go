package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"net/http"

	router "github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/server/http"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/repository/postgres"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

func main() {
	//Init logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	err = run(logger)
	if err != nil {
		logger.Error("doctor: error: ", zap.Error(err))
		os.Exit(1)
	}
	defer logger.Sync()
}
func run(logger *zap.Logger) error {
	var cfg config.Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return fmt.Errorf("read env: %w", err)
	}
	
	db, err := sql.Open("postgres", cfg.Postgres.String())
	defer func() {
		logger.Info("doctor: stopping database")
		db.Close()
	}()
	if err != nil {
		return fmt.Errorf("open db connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("database health check : %w", err)
	}

	err = postgres.MigrateUp("sql/migrations", cfg.Postgres)
	logger.Info("doctor: migrating database")
	if err != nil {
		return fmt.Errorf("migration : %w", err)
	}

	r := router.NewRouter(db, logger)
	return http.ListenAndServe(cfg.APIHost, r)
}
