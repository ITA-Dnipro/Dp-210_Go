package appointment

import (
	"context"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	as "github.com/ITA-Dnipro/Dp-210_Go/appointment/proto/appointments"
	"go.uber.org/zap"
)

// UsersUsecases represent user usecases.
type Usecase interface {
	GetByFilter(ctx context.Context, filter entity.AppointmentFilter) ([]entity.Appointment, error)
}

// userService gRPC Service
type appointmentService struct {
	usecase Usecase
	logger  *zap.Logger
}

// NewUserService create new user grpc service.
func NewAppointmentService(uc Usecase, log *zap.Logger) *appointmentService {
	return &appointmentService{usecase: uc, logger: log}
}

// Get.
func (u *appointmentService) GetByDoctorID(ctx context.Context, req *as.GetByDoctrorIDReq) (*as.GetByDoctorIDRes, error) {
	doctorID, err := uuid.Parse(req.GetDoctorID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	from := req.GetFrom().AsTime()
	till := req.GetTill().AsTime()

	f := entity.AppointmentFilter{DoctorID: &doctorID, From: &from, To: &till}
	a, err := u.usecase.GetByFilter(ctx, f)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &as.GetByDoctorIDRes{Appointments: toProtoList(a)}, nil
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
