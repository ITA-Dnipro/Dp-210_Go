package user

import (
	"context"
	"fmt"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UsersRepository represent user repository.
type UsersRepository interface {
	Create(ctx context.Context, u *entity.User) error
	Update(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id string) (entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Delete(ctx context.Context, id string) error
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

// NewUsecases create new user usecases.
func NewUsecases(r UsersRepository) *Usecases {
	return &Usecases{
		repo: r,
	}
}

// Usecases represent a user usecases.
type Usecases struct {
	repo UsersRepository
}

// Create Add new user
func (uc *Usecases) Create(ctx context.Context, u *entity.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate password hash:%w", err)
	}

	u.ID = uuid.New().String()
	u.PasswordHash = string(hash)
	u.PermissionRole = role.Viewer

	return uc.repo.Create(ctx, u)
}

// Update updates a user
func (uc *Usecases) Update(ctx context.Context, u *entity.User) error {
	return uc.repo.Update(ctx, u)
}

// Delete deletes a user from storage
func (uc *Usecases) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

// GetByID get single user by id.
func (uc *Usecases) GetByID(ctx context.Context, id string) (entity.User, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetAll get all users.
func (uc *Usecases) GetAll(ctx context.Context) (res []entity.User, err error) {
	return uc.repo.GetAll(ctx)
}

func (uc *Usecases) ChangePassword(ctx context.Context, passw entity.UserNewPassword) error {

	u, err := uc.userRepo.GetByID(ctx, passw.UserID)
	if err != nil {
		return fmt.Errorf("change password userId: %v, %w", passw.UserID, err)
	}

	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(passw.OldPassword)); err != nil {
		return fmt.Errorf("wrong password")
	}

	u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(passw.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generate password hash:%w", err)
	}

	if err := uc.userRepo.Update(ctx, &u); err != nil {
		return fmt.Errorf("change password: %w", err)
	}

	return nil
}
