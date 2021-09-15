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
	GetByPatientID(ctx context.Context, id uuid.UUID, al *entity.AppointmentList) error
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
	al := entity.AppointmentList{
		From: req.GetFrom().AsTime(),
		To:   req.GetTill().AsTime(),
	}
	if err := u.usecase.GetByPatientID(ctx, doctorID, &al); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &as.GetByDoctorIDRes{Appointments: toProtoList(al.Appointments)}, nil
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
