package password

import (
	"context"
	"fmt"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
)

type Usecases struct {
	sender   EmailSender
	codeGen  CodeGenerator
	userRepo UsersRepository
	cache    Cache
}

type UsersRepository interface {
	UserExists(ctx context.Context, email string) bool
	GetByEmail(ctx context.Context, email string) (entity.User, error)
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

func NewUsecases(es EmailSender, cg CodeGenerator, ur UsersRepository, cache Cache) *Usecases {
	return &Usecases{
		sender:   es,
		codeGen:  cg,
		userRepo: ur,
		cache:    cache,
	}
}

func (uc *Usecases) SendRestorePasswordCode(ctx context.Context, email string) (code string, err error) {
	if !uc.userRepo.UserExists(ctx, email) {
		return "", fmt.Errorf("no such user with email: %v", email)
	}

	code, err = uc.codeGen.GenerateCode()
	if err != nil {
		return code, fmt.Errorf("send restore passw code: %w", err)
	}

	if err := uc.cache.Set(ctx, email, code); err != nil {
		return code, fmt.Errorf("save restore code: %w", err)
	}

	if err = uc.sender.Send(email, code); err != nil {
		uc.cache.Del(ctx, email)
		return code, fmt.Errorf("send restore code: %w", err)
	}

	return code, nil
}

func (uc *Usecases) Authenticate(ctx context.Context, pc entity.PasswordCode) (string, error) {
	ent, err := uc.cache.Get(ctx, pc.Email)
	if err != nil || ent != pc.Code {
		return "", fmt.Errorf("no such code found: %v", pc)
	}

	user, err := uc.userRepo.GetByEmail(ctx, pc.Email)
	if err != nil {
		return "", fmt.Errorf("auth via passw code, get user: %w", err)
	}

	return user.ID, nil
}

func (uc *Usecases) DeleteCode(ctx context.Context, email string) error {
	return fmt.Errorf("delete passw code: %w", uc.cache.Del(ctx, email))
}
