package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/client/grpc/doctor"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/repository/postgres"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/kafka"

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
		db.Close()
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
		return err
	}
	if err != nil {
		return fmt.Errorf("connecting to kafka: %w", err)
	}
	//defer k.Close()

	logger.Info("Initializing server")
	return server.Serve(cfg, db, dc, k, logger)
}
