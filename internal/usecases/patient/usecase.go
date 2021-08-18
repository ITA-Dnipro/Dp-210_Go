package doctor

import (
	"context"
	"fmt"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
)

// UsersRepository represent user repository.
type UsersRepository interface {
	Update(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id string) (entity.User, error)
}

type PatientsRepository interface {
	Create(ctx context.Context, p *entity.Patient) error
	GetByID(ctx context.Context, id string) (entity.Patient, error)
	GetAll(ctx context.Context) ([]entity.Patient, error)
	Delete(ctx context.Context, id string) error
}

func NewUsecases(pr PatientsRepository, ur UsersRepository) *Usecases {
	return &Usecases{
		pr: pr,
		ur: ur,
	}
}

type Usecases struct {
	pr PatientsRepository
	ur UsersRepository
}

func (uc *Usecases) Create(ctx context.Context, p *entity.Patient) error {
	user, err := uc.ur.GetByID(ctx, p.ID)
	if err != nil {
		return fmt.Errorf("get user by %s id: %w", p.ID, err)
	}
	if user.PermissionRole != role.Viewer {
		return fmt.Errorf("user alredy registered as %s", user.PermissionRole)
	}
	user.PermissionRole = role.Patient
	if err := uc.ur.Update(ctx, &user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return uc.pr.Create(ctx, p)
}

func (uc *Usecases) Delete(ctx context.Context, id string) error {
	return uc.pr.Delete(ctx, id)
}

func (uc *Usecases) GetByID(ctx context.Context, id string) (entity.Patient, error) {
	return uc.pr.GetByID(ctx, id)
}

func (uc *Usecases) GetAll(ctx context.Context) (res []entity.Patient, err error) {
	return uc.pr.GetAll(ctx)
}
