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
}

// UsersRepository represent user repository.
type UsersRepository interface {
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

type EmailSender interface {
	Send(to, code string) error
}

type CodeGenerator interface {
	GenerateCode() (string, error)
}

// NewUsecases create new user usecases.
func NewUsecases(es EmailSender, cg CodeGenerator, ur UsersRepository) *Usecases {
	return &Usecases{
		sender:   es,
		codeGen:  cg,
		userRepo: ur,
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

	return code, nil
}
