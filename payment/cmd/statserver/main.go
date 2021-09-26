package main

import (
	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/handlers/stathand"
	"log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/proto/statistics"
	server "github.com/ITA-Dnipro/Dp-210_Go/payment/internal/server/grpc"
)

func main() {
	logger, _ := zap.NewProduction()
	statHandler := stathand.NewHandler(logger)

	grpcServer := grpc.NewServer()
	srv := server.NewGRPCServer(statHandler)
	statistics.RegisterStatServer(grpcServer, srv)

	listener, err := net.Listen("tcp", ":1235")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
