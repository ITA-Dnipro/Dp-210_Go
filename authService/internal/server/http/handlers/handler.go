package handlers

import (
	"context"
	"github.com/ITA-Dnipro/Dp-210_Go/authService/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/authService/internal/usecase"
	"go.uber.org/zap"
)

type Auth interface {
	CreateToken(user usecase.UserAuth) (usecase.JwtToken, error)
	InvalidateToken(userId string) error
	ValidateToken(t usecase.JwtToken) (usecase.UserAuth, error)
}

type PasswordUsecases interface {
	SendRestorePasswordCode(ctx context.Context, email string) (code string, err error)
	Authenticate(ctx context.Context, pc entity.PasswordCode) (entity.User, error)
	Auth(ctx context.Context, email, password string) (u entity.User, err error)
	DeleteCode(ctx context.Context, email string) error
	ChangePassword(ctx context.Context, passw entity.UserNewPassword) error
	SetNewPassword(ctx context.Context, password string, user *entity.User) error
}

type Handlers struct {
	paswCases PasswordUsecases
	logger    *zap.Logger
	auth      Auth
}

func NewHandler(paswCases PasswordUsecases, logger *zap.Logger, auth Auth) *Handlers {
	return &Handlers{
		paswCases: paswCases,
		logger:    logger,
		auth:      auth,
	}
}
