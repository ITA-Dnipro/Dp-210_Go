package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/client/grpc/doctor"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/client/grpc/user"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/repository/postgres"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/kafka"

	serverGRPC "github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/grpc"
	serverHTTP "github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/http"

	appointmentRepo "github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/repository/postgres/appointment"
	appointmentUC "github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/usecases/appointment"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

const (
	migrationsPath = "migrations"
)

func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("building logger", err)
	}
	if err := run(zapLogger); err != nil {
		zapLogger.Error("error:", zap.Error(err))
		os.Exit(1)
	}
}

func run(logger *zap.Logger) error {
	var cfg config.Config
	if err := envconfig.Process("appointment", &cfg); err != nil {
		return fmt.Errorf("read env: %w", err)
	}
	logger.Info("Initializing database")
	db, err := postgres.Open(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("close database", zap.Error(err))
		}
		logger.Info("Database Stopping")
	}()
	if err = postgres.MigrateUp(migrationsPath, cfg.Postgres); err != nil {
		return fmt.Errorf("migrations db: %w", err)
	}
	logger.Info("Initializing kafka")
	k, err := kafka.NewKafka(cfg.KafkaBrokers, logger)
	if err != nil {
		return fmt.Errorf("connecting to kafka: %w", err)
	}
	defer k.Close()
	logger.Info("Initializing doctor client")
	dc, err := doctor.NewDoctorClient(cfg.DocrotGRPCHost)
	if err != nil {
		return fmt.Errorf("connecting to doctor client: %w", err)
	}
	logger.Info("Initializing user client")
	uc, err := user.NewUserClient(cfg.UserGRPCHost)
	if err != nil {
		return fmt.Errorf("connecting to user client: %w", err)
	}

	ar := appointmentRepo.NewRepository(db)
	ac := appointmentUC.NewUsecases(ar, uc, dc, k)

	if err := k.OnAppointment(ac.CreateFromEvent); err != nil {
		return fmt.Errorf("kafka consumer: %w", err)
	}

	errors := make(chan error, 1)
	httpServer := serverHTTP.NewHTTPServer(cfg, ac, logger)
	go func() { errors <- httpServer.ListenAndServe() }()
	grpcServer := serverGRPC.NewGRPCServer(cfg, ac, logger)
	go func() { errors <- grpcServer.Serve() }()

	return <-errors
}
