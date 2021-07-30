package usecases

import (
	"context"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
)

type UsersRepository interface {
	Create(ctx context.Context, u entity.User) (string, error)
	Update(ctx context.Context, u entity.User) error
	GetByID(ctx context.Context, id string) (entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Delete(ctx context.Context, id string) error
}

func NewUsecases(r UsersRepository) *Usecases {
	return &Usecases{
		repo: r,
	}
}

type Usecases struct {
	repo UsersRepository
}

func (uc *Usecases) Create(ctx context.Context, u entity.User) (string, error) {
	return uc.repo.Create(ctx, u)
}

func (uc *Usecases) Update(ctx context.Context, u entity.User) error {
	return uc.repo.Update(ctx, u)
}

func (uc *Usecases) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *Usecases) GetByID(ctx context.Context, id string) (entity.User, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *Usecases) GetAll(ctx context.Context) (res []entity.User, err error) {
	return uc.repo.GetAll(ctx)
}
