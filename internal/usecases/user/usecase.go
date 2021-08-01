package user

import (
	"context"
	"fmt"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UsersRepository represent user repository.
type UsersRepository interface {
	Create(ctx context.Context, u entity.User) error
	Update(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id string) (entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Delete(ctx context.Context, id string) error
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
func (uc *Usecases) Create(ctx context.Context, nu entity.NewUser) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate password hash:%w", err)
	}

	id := uuid.New().String()

	u := entity.User{
		ID:             id,
		Name:           nu.Name,
		Email:          nu.Email,
		PermissionRole: entity.Viewer,
		PasswordHash:   hash,
	}
	return id, uc.repo.Create(ctx, u)
}

// Update updates a user
func (uc *Usecases) Update(ctx context.Context, nu entity.NewUser) (entity.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, fmt.Errorf("generate password hash:%w", err)
	}
	u := entity.User{
		ID:           nu.ID,
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
	}
	return u, uc.repo.Update(ctx, &u)
}

// Delete deletes a user from storage
func (uc *Usecases) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

// GetByID get single user by id.
func (uc *Usecases) GetByID(ctx context.Context, id string) (entity.User, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetRoleByID get user permission role.
func (uc *Usecases) GetRoleByID(ctx context.Context, id string) (entity.Role, error) {
	u, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("get role by id:%w", err)
	}
	return entity.Role(u.PermissionRole), nil
}

// GetAll get all users.
func (uc *Usecases) GetAll(ctx context.Context) (res []entity.User, err error) {
	return uc.repo.GetAll(ctx)
}
