package grpc

import (
	"context"
	"net"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/grpc/appointment"
	appointmentsService "github.com/ITA-Dnipro/Dp-210_Go/appointment/proto/appointments"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Usecase interface {
	GetWithFilter(ctx context.Context, filter entity.AppointmentFilter) ([]entity.Appointment, error)
}
type Server struct {
	srv  *grpc.Server
	host string
}

// NewGRPCServer create grpc server.
func NewGRPCServer(cfg config.Config, uc Usecase, logger *zap.Logger) *Server {
	grpcServer := grpc.NewServer()
	as := appointment.NewAppointmentService(uc, logger)
	appointmentsService.RegisterUsersServiceServer(grpcServer, as)
	return &Server{srv: grpcServer, host: cfg.GRPCHost}
}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", s.host)
	if err != nil {
		return err
	}
	defer l.Close()
	return s.srv.Serve(l)
}

func (s *Server) GracefulStop() {
	s.srv.GracefulStop()
}
