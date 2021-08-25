package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/kafka"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware"
	"github.com/ITA-Dnipro/Dp-210_Go/visits/config"
	"github.com/ITA-Dnipro/Dp-210_Go/visits/store/postgres"
	"github.com/go-chi/chi"

	appointmentRepo "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/appointment"
	doctorRepo "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/doctor"
	patientRepo "github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/patient"
	appointmentHandlers "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/appointment"
	appointmentUsecases "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/appointment"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

func main() {
	ZapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("building logger", err)
	}
	if err := run(ZapLogger); err != nil {
		ZapLogger.Error("visits: error:", zap.Error(err))
		os.Exit(1)
	}
}
func run(logger *zap.Logger) error {
	var cfg config.Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	logger.Info("visits: Initializing database support")
	db, err := postgres.Open(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		log.Printf("visits: Database Stopping")
		db.Close()
	}()
	events, err := kafka.NewEvents([]string{"localhost:9092"})
	if err != nil {
		return fmt.Errorf("connecting to kafka: %w", err)
	}

	events.On(kafka.MailTopic, func(payload []byte) error {
		fmt.Println(string(payload))
		return nil
	})
	events.On(kafka.BillTopic, func(payload []byte) error {
		fmt.Println(string(payload))
		return nil
	})

	dr := doctorRepo.NewRepository(db)
	pr := patientRepo.NewRepository(db)
	ar := appointmentRepo.NewRepository(db)
	ac := appointmentUsecases.NewUsecases(ar, dr, pr, events)
	ah := appointmentHandlers.NewHandlers(ac, logger)
	md := &middleware.Middleware{Logger: logger}

	r := chi.NewRouter()
	r.Use(md.LoggingMiddleware)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(md.AuthMiddleware)
		r.Group(func(r chi.Router) { // route with permissions
			r.Use(md.RoleOnly(role.Patient))
			r.Post("/appointments", ah.CreateAppointment) // Post /api/v1/appointment
		})
		r.Group(func(r chi.Router) { // route with permissions
			r.Use(md.RoleOnly(role.Patient, role.Doctor, role.Admin, role.Operator))
			r.Get("/appointments", ah.GetAppointments) // GET /api/v1/appointments
		})
	})
	logger.Info("visits: Initializing API support")
	api := http.Server{
		Addr:         cfg.APIHost,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	serverErrors := make(chan error, 1)
	go func() {
		logger.Sugar().Infof("visits: API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				ac.DeleteOld(ctx, time.Now())
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}(ctx)

	events.On(kafka.AppoinmentTopic, func(payload []byte) error {
		if err := ac.CreateEvent(payload); err != nil {
			logger.Error("visits: create error:", zap.Error(err))
		}
		return nil
	})
	_ = cancel
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select { // Blocking visits and waiting for shutdown.
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		logger.Sugar().Infof("visits: %v : Start shutdown", sig)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("visits could not stop server gracefully: %w", err)
		}
	}
	return nil
}
