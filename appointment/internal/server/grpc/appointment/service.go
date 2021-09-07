package appointment

import (
	"context"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"

	as "github.com/ITA-Dnipro/Dp-210_Go/appointment/proto/appointments"
	"go.uber.org/zap"
)

// UsersUsecases represent user usecases.
type Usecase interface {
	GetWithFilter(ctx context.Context, filter entity.AppointmentFilter) ([]entity.Appointment, error)
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

// Update update user.
func (u *appointmentService) GetByID(ctx context.Context, req *as.GetByIDReq) (*as.GetByIDRes, error) {
	// user, err := u.usecase.GetByDoctorID(ctx, req.GetUserID())
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, err.Error())
	// }
	// return &us.GetByIDRes{User: &us.User{UserID: user.ID, Name: user.Name, Email: user.Email, Role: string(user.PermissionRole)}}, nil
	return nil, nil
}
