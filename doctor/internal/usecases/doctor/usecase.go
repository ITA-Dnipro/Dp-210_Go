package doctor

import (
	"context"
	"fmt"
	_ "fmt"

	ugc "github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/client/grpc/users"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/role"
	_ "github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/role"
)

// UsersRepository represent user repository.
type UsersRepository interface {
	Update(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id string) (entity.User, error)
}

// DoctorsRepository represent doctor repository.
type DoctorsRepository interface {
	Create(ctx context.Context, d *entity.Doctor) error
	Update(ctx context.Context, d *entity.Doctor) error
	GetByID(ctx context.Context, id string) (entity.Doctor, error)
	GetAll(ctx context.Context) ([]entity.Doctor, error)
	Delete(ctx context.Context, id string) error
}

// NewUsecases create new doctor usecases.
func NewUsecases(dr DoctorsRepository, uc *ugc.Client) *Usecases {
	return &Usecases{
		dr:  dr,
		ugc: uc,
		//	ur: ur,
	}
}

// Usecases represent a doctor usecases.
type Usecases struct {
	dr  DoctorsRepository
	ugc *ugc.Client
}

// Create Add new doctor
func (uc *Usecases) Create(ctx context.Context, d *entity.Doctor) error {

	user, err := uc.ugc.GetByID(ctx, d.ID) //Call client
	if err != nil {
		return fmt.Errorf("get user by %s id: %w", d.ID, err)
	}

	if user.PermissionRole != role.Viewer {
		return fmt.Errorf("user alredy registered as %s", user.PermissionRole)
	}
	user.PermissionRole = role.Doctor
	if err := uc.ugc.Update(ctx, &user); err != nil { //Call client
		return fmt.Errorf("update user: %w", err)
	}
	return uc.dr.Create(ctx, d)
}

// Update updates a doctor
func (uc *Usecases) Update(ctx context.Context, d *entity.Doctor) error {
	return uc.dr.Update(ctx, d)
}

// Delete deletes a doctor from storage
func (uc *Usecases) Delete(ctx context.Context, id string) error {
	return uc.dr.Delete(ctx, id)
}

// GetByID get single doctor by id.
func (uc *Usecases) GetByID(ctx context.Context, id string) (entity.Doctor, error) {
	return uc.dr.GetByID(ctx, id)
}

// GetAll get all doctors.
func (uc *Usecases) GetAll(ctx context.Context) (res []entity.Doctor, err error) {
	return uc.dr.GetAll(ctx)
}
