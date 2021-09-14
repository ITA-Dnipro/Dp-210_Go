package server

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/config"
	postgres "github.com/ITA-Dnipro/Dp-210_Go/user/internal/repository/postgres/user"
	serverGRPC "github.com/ITA-Dnipro/Dp-210_Go/user/internal/server/grpc"
	serverHTTP "github.com/ITA-Dnipro/Dp-210_Go/user/internal/server/http"
	usecases "github.com/ITA-Dnipro/Dp-210_Go/user/internal/usecases/user"

	"go.uber.org/zap"
)

// Server users service.
func Serve(cfg config.Config, db *sql.DB, logger *zap.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := postgres.NewRepository(db)
	usecase := usecases.NewUsecases(repo)

	httpErrors := make(chan error, 1)
	httpServer := serverHTTP.NewHTTPServer(cfg, usecase, logger)
	go func() { httpErrors <- httpServer.ListenAndServe() }()

	grpcErrors := make(chan error, 1)
	grpcServer := serverGRPC.NewGRPCServer(cfg, usecase, logger)
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
