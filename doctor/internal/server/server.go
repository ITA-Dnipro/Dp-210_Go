package server

import (
	"fmt"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/repository/postgres/doctor"
	"go.uber.org/zap"

	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/server/grpc"
	router "github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/server/http"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/server/http/middleware"
	usecases "github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/usecases/doctor"
)

func RunServers(cfg config.Config, repo *doctor.Repository, usecase *usecases.Usecases, middleware *middleware.Middleware, logger *zap.Logger, errChan chan error) {
	r := router.NewRouter(repo, usecase, logger, middleware)
	//Create instance of grpc
	grpcServer := grpc.NewGRPCServer(cfg, usecase, logger)

	go func() {
		logger.Info(fmt.Sprintf("startup http server:%s", cfg.APIHost))
		errChan <- http.ListenAndServe(cfg.APIHost, r)
	}()
	go func() {
		logger.Info(fmt.Sprintf("startup grpc server:%s", cfg.GRPCHost))
		errChan <- grpcServer.Serve()
	}()
}
