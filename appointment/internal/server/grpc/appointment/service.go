package appointment

import (
	"context"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	"github.com/google/uuid"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	as "github.com/ITA-Dnipro/Dp-210_Go/appointment/proto/appointments"
	"go.uber.org/zap"
)

// UsersUsecases represent user usecases.
type Usecase interface {
	GetByDoctorID(ctx context.Context, id uuid.UUID, p *entity.AppointmentsParam) ([]entity.Appointment, string, error)
}

// userService gRPC Service
type appointmentServiceServer struct {
	usecase Usecase
	logger  *zap.Logger
}

// NewUserService create new user grpc service.
func RegisterAppointmentServiceServer(s *grpc.Server, uc Usecase, log *zap.Logger) {
	ass := &appointmentServiceServer{usecase: uc, logger: log}
	as.RegisterAppointmentServiceServer(s, ass)
}

// GetByDoctorID get appointments by doctor id.
func (u *appointmentServiceServer) GetByDoctorID(ctx context.Context, req *as.GetByDoctrorIDReq) (*as.GetByDoctorIDRes, error) {
	doctorID, err := uuid.Parse(req.GetDoctorID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	p := entity.AppointmentsParam{
		From: req.GetFrom().AsTime(),
		To:   req.GetTill().AsTime(),
	}
	resp, _, err := u.usecase.GetByDoctorID(ctx, doctorID, &p)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &as.GetByDoctorIDRes{Appointments: toProtoList(resp)}, nil
}
func toProto(a entity.Appointment) *as.Appointment {
	return &as.Appointment{
		AppointmentID: a.ID.String(),
		DoctorID:      a.DoctorID.String(),
		PatientID:     a.PatientID.String(),
		Reason:        a.Reason,
		From:          timestamppb.New(a.From),
		To:            timestamppb.New(a.To),
	}
}

func toProtoList(al []entity.Appointment) []*as.Appointment {
	pl := make([]*as.Appointment, 0, len(al))
	for _, a := range al {
		pl = append(pl, toProto(a))
	}
	return pl
}
