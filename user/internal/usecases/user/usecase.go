package user

import (
	"context"
	"fmt"

	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/role"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UsersRepository represent user repository.
type UsersRepository interface {
	Create(ctx context.Context, u *entity.User) error
	Update(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id string) (entity.User, error)
	GetAll(ctx context.Context, ul *entity.UserList) error
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
	hash, err := bcrypt.GenerateFromPassword(u.PasswordHash, bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generate password hash:%w", err)
	}
	u.ID = uuid.New().String()
	u.PermissionRole = role.Viewer
	u.PasswordHash = hash
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

// GetRoleByID get user permission role.
func (uc *Usecases) GetRoleByID(ctx context.Context, id string) (role.Role, error) {
	u, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("get role by id:%w", err)
	}
	return role.Role(u.PermissionRole), nil
}

// GetAll get all users.
func (uc *Usecases) GetAll(ctx context.Context, ul *entity.UserList) error {
	return uc.repo.GetAll(ctx, ul)
}

