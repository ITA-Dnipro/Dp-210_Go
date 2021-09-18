package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/usecase"
)

var _ usecase.UsersRepository = (*Repository)(nil)

// NewRepository create new user repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		storage: db,
	}
}

// Repository represent a user repository.
type Repository struct {
	storage *sql.DB
}

func (r *Repository) Update(ctx context.Context, u *entity.User) error {
	query := `UPDATE users SET email=$2, password_hash=$3 WHERE id=$1`
	res, err := r.storage.ExecContext(ctx, query, &u.ID, &u.Email, &u.PasswordHash)
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update rows afected:%w", err)
	}

	if rowsAfected != 1 {
		return fmt.Errorf("update affected: %d", rowsAfected)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (entity.User, error) {
	query := `SELECT email, role, password_hash FROM users WHERE id = $1`
	u := entity.User{}
	u.ID = id
	err := r.storage.QueryRowContext(ctx, query, id).Scan(&u.Email, &u.PermissionRole, &u.PasswordHash)
	if err != nil {
		return entity.User{}, fmt.Errorf("there is no users with %s id", id)
	}
	return u, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	query := `SELECT id, role, password_hash FROM users WHERE email = $1`
	u := entity.User{}
	u.Email = email
	err := r.storage.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.PermissionRole, &u.PasswordHash)
	if err != nil {
		return entity.User{}, fmt.Errorf("there is no users with %s email: (%w)", email, err)
	}
	return u, nil
}
