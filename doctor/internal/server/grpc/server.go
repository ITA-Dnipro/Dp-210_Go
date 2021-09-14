package grpc

import (
	"net"

	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/usecases/doctor"
	doctorsService "github.com/ITA-Dnipro/Dp-210_Go/doctor/proto/doctors"

	handlers "github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/server/grpc/doctor"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
	host   string
}

func NewGRPCServer(cfg config.Config, uc *doctor.Usecases, logger *zap.Logger) *Server {
	grpcServer := grpc.NewServer()
	dh := handlers.NewHandlers(uc, logger)
	doctorsService.RegisterDoctorsServiceServer(grpcServer, dh)
	return &Server{server: grpcServer, host: cfg.GRPCHost}
}

func (s *Server) Serve() error {
	listener, err := net.Listen("tcp", s.host)
	if err != nil {
		return err
	}
	defer listener.Close()
	return s.server.Serve(listener)
}
