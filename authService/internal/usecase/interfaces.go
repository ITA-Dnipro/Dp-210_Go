package usecase

import (
	"context"

	"github.com/ITA-Dnipro/Dp-210_Go/authService/internal/entity"
)

type UsersRepository interface {
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	Update(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id string) (entity.User, error)
}

type Cache interface {
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

type EmailSender interface {
	Send(to, code string) error
}

type CodeGenerator interface {
	GenerateCode() (string, error)
}
