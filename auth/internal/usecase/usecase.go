package usecase

import (
	"context"
	"fmt"

	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/entity"

	"golang.org/x/crypto/bcrypt"
)

type Usecases struct {
	sender   EmailSender
	codeGen  CodeGenerator
	userRepo UsersRepository
	cache    Cache
}

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

func NewUsecases(es EmailSender, cg CodeGenerator, ur UsersRepository, cache Cache) *Usecases {
	return &Usecases{
		sender:   es,
		codeGen:  cg,
		userRepo: ur,
		cache:    cache,
	}
}

func (uc *Usecases) SendRestorePasswordCode(ctx context.Context, email string) (code string, err error) {
	u, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil || u.Email != email {
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

func (uc *Usecases) Auth(ctx context.Context, email, password string) (u entity.User, err error) {
	u, err = uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return entity.User{}, fmt.Errorf("authenticate get user by email:%w", err)
	}
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return entity.User{}, fmt.Errorf("authentication failed:%w", err)
	}

	return u, nil
}

func (uc *Usecases) Authenticate(ctx context.Context, pc entity.PasswordCode) (entity.User, error) {
	ent, err := uc.cache.Get(ctx, pc.Email)
	if err != nil || ent != pc.Code {
		return entity.User{}, fmt.Errorf("no such code found: %v", pc)
	}

	user, err := uc.userRepo.GetByEmail(ctx, pc.Email)
	if err != nil {
		return entity.User{}, fmt.Errorf("auth via passw code, get user: %w", err)
	}

	return user, nil
}

func (uc *Usecases) SetNewPassword(ctx context.Context, password string, user *entity.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("set new usecase: %w", err)
	}

	user.PasswordHash = hash
	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("set new usecase: %w", err)
	}

	return nil
}

func (uc *Usecases) DeleteCode(ctx context.Context, email string) error {
	return fmt.Errorf("delete passw code: %w", uc.cache.Del(ctx, email))
}

func (uc *Usecases) ChangePassword(ctx context.Context, passw entity.UserNewPassword) error {

	u, err := uc.userRepo.GetByID(ctx, passw.UserID)
	if err != nil {
		return fmt.Errorf("change usecase userId: %v, %w", passw.UserID, err)
	}

	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(passw.OldPassword)); err != nil {
		return fmt.Errorf("wrong usecase")
	}

	u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(passw.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generate usecase hash:%w", err)
	}

	if err := uc.userRepo.Update(ctx, &u); err != nil {
		return fmt.Errorf("change usecase: %w", err)
	}

	return nil
}
