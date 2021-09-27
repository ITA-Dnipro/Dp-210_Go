package grpc

import (
	"context"
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/handlers/stathand"
	stat "github.com/ITA-Dnipro/Dp-210_Go/payment/internal/proto/statistics"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
)

func NewGRPCServer(h *stathand.Handler) *GRPCServer {
	return &GRPCServer{
		h: h,
	}
}

//goland:noinspection GoNameStartsWithPackageName
type GRPCServer struct {
	h *stathand.Handler
}

func (s *GRPCServer) DocStat(ctx context.Context, req *stat.DocRequest) (*empty.Empty, error) {
	// 1) Достать врачей.
	doctorsArr, err := s.h.DoctorsUnmarshal(req.DocsBytesArr)
	if err != nil {
		s.h.Logger.Error("error in DocStat", zap.String("server", err.Error()))
		return &empty.Empty{}, nil
	}
	// 2) Получить лучшего.
	bestDoctorsArr := s.h.GetBest(doctorsArr)
	// 3) Форматированный вывод.
	var out string
	for _, doctor := range bestDoctorsArr {
		out += fmt.Sprintf("| Top doctor for the last period = Doctor id > [%s] with result > [%d]|", doctor.DoctorId.String(), doctor.DoctorTotal)
	}
	s.h.Logger.Info(out)
	return &empty.Empty{}, nil
}
