package user

import (
	"context"
	"errors"

	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/role"
	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/server/http/customerrors"
	"github.com/go-playground/validator/v10"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	us "github.com/ITA-Dnipro/Dp-210_Go/user/proto/user"
	"go.uber.org/zap"
)

// UsersUsecases represent user usecases.
type Usecase interface {
	GetByID(ctx context.Context, id string) (entity.User, error)
	Create(ctx context.Context, u *entity.User) error
	Update(ctx context.Context, u *entity.User) error
}

// userService gRPC Service
type userServiceServer struct {
	usecase Usecase
	logger  *zap.Logger
}

// NewUserService create new user grpc service.
func RegisterUserServiceServer(s *grpc.Server, uc Usecase, log *zap.Logger) {
	uss := &userServiceServer{usecase: uc, logger: log}
	us.RegisterUsersServiceServer(s, uss)
}

// Create create new user.
func (u *userServiceServer) Create(ctx context.Context, req *us.CreateReq) (*us.CreateRes, error) {
	user := entity.User{
		Name:           req.GetName(),
		Email:          req.GetEmail(),
		PermissionRole: role.Role(req.GetRole()),
		PasswordHash:   []byte(req.GetPassword()),
	}

	if err := validator.New().Struct(user); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := u.usecase.Create(ctx, &user); err != nil {
		if errors.Is(err, customerrors.ErrDublication) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &us.CreateRes{User: &us.User{UserID: user.ID, Name: user.Name, Email: user.Email, Role: string(user.PermissionRole)}}, nil
}

// Update update user.
func (u *userServiceServer) Update(ctx context.Context, req *us.UpdateReq) (*us.UpdateRes, error) {
	user := entity.User{
		ID:             req.GetUserID(),
		Name:           req.GetName(),
		Email:          req.GetEmail(),
		PermissionRole: role.Role(req.GetRole()),
	}

	if err := validator.New().Struct(user); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := u.usecase.Update(ctx, &user); err != nil {
		if errors.Is(err, customerrors.ErrDublication) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, customerrors.ErrForeignKey) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &us.UpdateRes{User: &us.User{UserID: user.ID, Name: user.Name, Email: user.Email, Role: string(user.PermissionRole)}}, nil
}

// Update update user.
func (u *userServiceServer) GetByID(ctx context.Context, req *us.GetByIDReq) (*us.GetByIDRes, error) {
	user, err := u.usecase.GetByID(ctx, req.GetUserID())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &us.GetByIDRes{User: &us.User{UserID: user.ID, Name: user.Name, Email: user.Email, Role: string(user.PermissionRole)}}, nil
}
