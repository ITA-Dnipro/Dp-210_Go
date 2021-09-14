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
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/kafka"

	appointmentRepo "github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/repository/postgres/appointment"
	appointmentUC "github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/usecases/appointment"
)

// Server users service.
func Serve(cfg config.Config, db *sql.DB, dc *doctor.Client, kf *kafka.Kafka, logger *zap.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ar := appointmentRepo.NewRepository(db)
	ac := appointmentUC.NewUsecases(ar, dc, kf)

	kf.OnAppointment(ac.CreateFromEvent)

	httpErrors := make(chan error, 1)
	httpServer := serverHTTP.NewHTTPServer(cfg, ac, logger)
	go func() { httpErrors <- httpServer.ListenAndServe() }()

	grpcErrors := make(chan error, 1)
	grpcServer := serverGRPC.NewGRPCServer(cfg, ac, logger)
	go func() { grpcErrors <- grpcServer.Serve() }()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-httpErrors:
		grpcServer.GracefulStop()
		return fmt.Errorf("http server error: %w", err)
	case err := <-grpcErrors:
		httpServer.GracefulShutdown()
		return fmt.Errorf("grpc server error: %w", err)
	case v := <-quit:
		logger.Info(fmt.Sprintf("signal.Notify: %v", v))
	case done := <-ctx.Done():
		logger.Info(fmt.Sprintf("ctx.Done: %v", done))
	}
	httpServer.GracefulShutdown()
	grpcServer.GracefulStop()
	return nil
}
