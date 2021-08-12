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
	codeRepo PasswordCodeRepository
}

type PasswordCodeRepository interface {
	Create(ctx context.Context, c entity.PasswordCode) error
	GetByEmail(ctx context.Context, email string) (entity.PasswordCode, error)
	Delete(ctx context.Context, email string) error
}

// UsersRepository represent user repository.
type UsersRepository interface {
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

type Cache interface {
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
}

type EmailSender interface {
	Send(to, code string) error
}

type CodeGenerator interface {
	GenerateCode() (string, error)
}

// NewUsecases create new user usecases.
func NewUsecases(es EmailSender, cg CodeGenerator, ur UsersRepository, cr PasswordCodeRepository) *Usecases {
	return &Usecases{
		sender:   es,
		codeGen:  cg,
		userRepo: ur,
		codeRepo: cr,
	}
}

func (uc *Usecases) SendRestorePasswordCode(ctx context.Context, email string) (code string, err error) {
	if _, err := uc.userRepo.GetByEmail(ctx, email); err != nil {
		return "", fmt.Errorf("no such user with email %v; %w", email, err)
	}

	code, err = uc.codeGen.GenerateCode()
	if err != nil {
		return code, fmt.Errorf("send restore passw code: %w", err)
	}

	if err = uc.sender.Send(email, code); err != nil {
		return code, fmt.Errorf("send restore code: %w", err)
	}

	pc := entity.PasswordCode{Email: email, Code: code}
	if err := uc.codeRepo.Create(ctx, pc); err != nil {
		return code, fmt.Errorf("save restore code: %w", err)
	}

	return code, nil
}

func (uc *Usecases) VerifyCode(ctx context.Context, pc entity.PasswordCode) error {
	ent, err := uc.codeRepo.GetByEmail(ctx, pc.Email)
	if err != nil {
		return fmt.Errorf("passw code verify: %w", err)
	}

	if ent != pc {
		return fmt.Errorf("no such code found: %v", pc)
	}

	return nil
}
