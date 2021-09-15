package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/server/grpc/user"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Usecase interface {
	Create(ctx context.Context, u *entity.User) error
	Update(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id string) (entity.User, error)
}

type Server struct {
	srv    *grpc.Server
	cfg    config.Config
	logger *zap.Logger
}

// NewGRPCServer create grpc server.
func NewGRPCServer(cfg config.Config, uc Usecase, logger *zap.Logger) *Server {
	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, uc, logger)
	return &Server{srv: grpcServer, cfg: cfg, logger: logger}
}

func (s *Server) Serve() error {
	s.logger.Info(fmt.Sprintf("startup grpc server:%s", s.cfg.GRPCHost))
	l, err := net.Listen("tcp", s.cfg.GRPCHost)
	if err != nil {
		return err
	}
	return s.srv.Serve(l)
}

func (s *Server) GracefulStop() {
	s.srv.GracefulStop()
	s.logger.Info("grpc server shutdown")
}
