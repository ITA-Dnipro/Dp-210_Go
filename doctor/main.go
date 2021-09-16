package main

import (
	//"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	//"time"

	//cache "github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/cache/redis"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/repository/postgres"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/repository/postgres/doctor"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/server"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/server/http/middleware"

	//"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/service/auth"
	usecases "github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/usecases/doctor"

	agc "github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/client/grpc/appointments"
	ugc "github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/client/grpc/users"

	//"github.com/go-redis/redis/v8"
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

	logger.Info(fmt.Sprintf("startup appointment client:%s", cfg.AppointmentGRPCClient))
	ac, err := agc.NewAppointmentsClient(cfg, logger)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("startup user client:%s", cfg.UserGRPCClient))
	uc, err := ugc.NewUserClient(cfg, logger)
	if err != nil {
		return err
	}

	repo := doctor.NewRepository(db)
	usecase := usecases.NewUsecases(repo, uc, ac)
	md := &middleware.Middleware{Logger: logger, UserUC: usecase}

	errChan := make(chan error)
	server.RunServers(cfg, repo, usecase, md, ac, logger, errChan)
	return <-errChan
}
