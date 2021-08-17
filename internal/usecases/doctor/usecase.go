package doctor

import (
	"context"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
)

// DoctorsRepository represent doctor repository.
type DoctorsRepository interface {
	Create(ctx context.Context, u entity.Doctor) error
	Update(ctx context.Context, u *entity.Doctor) error
	GetByID(ctx context.Context, id string) (entity.Doctor, error)
	GetAll(ctx context.Context) ([]entity.Doctor, error)
	Delete(ctx context.Context, id string) error
}

// NewUsecases create new doctor usecases.
func NewUsecases(r DoctorsRepository) *Usecases {
	return &Usecases{
		repo: r,
	}
}

// Usecases represent a doctor usecases.
type Usecases struct {
	repo DoctorsRepository
}

// Create Add new doctor
func (uc *Usecases) Create(ctx context.Context, nd entity.Doctor) (string, error) {
	return nd.ID, uc.repo.Create(ctx, nd)
}

// Update updates a doctor
func (uc *Usecases) Update(ctx context.Context, nd entity.Doctor) (entity.Doctor, error) {
	return nd, uc.repo.Update(ctx, &nd)
}

// Delete deletes a doctor from storage
func (uc *Usecases) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

// GetByID get single doctor by id.
func (uc *Usecases) GetByID(ctx context.Context, id string) (entity.Doctor, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetRoleByID get doctor permission role.
//func (uc *Usecases) GetRoleByID(ctx context.Context, id string) (role.Role, error) {
//	u, err := uc.repo.GetByID(ctx, id)
//	if err != nil {
//		return "", fmt.Errorf("get role by id:%w", err)
//	}
//	return role.Role(u.PermissionRole), nil
//}

// GetAll get all doctors.
func (uc *Usecases) GetAll(ctx context.Context) (res []entity.Doctor, err error) {
	return uc.repo.GetAll(ctx)
}
