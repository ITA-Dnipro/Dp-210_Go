package server

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/config"
	"go.uber.org/zap"

	serverGRPC "github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/grpc"
	serverHTTP "github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/http"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/client/grpc/doctor"
	k "github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/kafka"

	appointmentRepo "github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/repository/postgres/appointment"
	appointmentUC "github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/usecases/appointment"
)

// Server users service.
func Serve(cfg config.Config, db *sql.DB, dc *doctor.Client, kafka *k.Kafka, logger *zap.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ar := appointmentRepo.NewRepository(db)
	ac := appointmentUC.NewUsecases(ar, dc, kafka)

	serverErrors := make(chan error, 1)

	sh := serverHTTP.NewHTTPServer(cfg, ac, logger)
	logger.Info(fmt.Sprintf("startup http server:%s", cfg.APIHost))
	go func() { serverErrors <- sh.ListenAndServe() }()

	logger.Info(fmt.Sprintf("startup grpc server:%s", cfg.GRPCHost))
	sg := serverGRPC.NewGRPCServer(cfg, ac, logger)
	go func() { serverErrors <- sg.Serve() }()

	kafka.OnAppointment(ac.CreateFromEvent)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case v := <-quit:
		logger.Info(fmt.Sprintf("signal.Notify: %v", v))
	case done := <-ctx.Done():
		logger.Info(fmt.Sprintf("ctx.Done: %v", done))
	}
	sg.GracefulStop()
	if err := sh.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
