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
	Create(ctx context.Context, u entity.User) error
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
		PermissionRole: role.Viewer,
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
func (uc *Usecases) GetRoleByID(ctx context.Context, id string) (role.Role, error) {
	u, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("get role by id:%w", err)
	}
	return role.Role(u.PermissionRole), nil
}

// GetAll get all users.
func (uc *Usecases) GetAll(ctx context.Context) (res []entity.User, err error) {
	return uc.repo.GetAll(ctx)
}

// Authenticate user by email and password.
func (uc *Usecases) Authenticate(ctx context.Context, email, password string) (id string, err error) {
	u, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("authenticate get user by email:%w", err)
	}
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return "", fmt.Errorf("authentication failed:%w", err)
	}

	return u.ID, nil
}

func (uc *Usecases) ChangePassword(ctx context.Context, passw entity.UserNewPassword) error {

	u, err := uc.repo.GetByID(ctx, passw.UserID)
	if err != nil {
		return fmt.Errorf("change password userId: %v, %w", passw.UserID, err)
	}

	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(passw.OldPassword)); err != nil {
		return fmt.Errorf("wrong password")
	}

	if err := uc.repo.Update(ctx, &u); err != nil {
		return fmt.Errorf("change password: %w", err)
	}

	return nil
}
