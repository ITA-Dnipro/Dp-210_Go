package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	usecases "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/user"
)

var _ usecases.UsersRepository = (*Repository)(nil)

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

// Create Add new user
func (r *Repository) Create(ctx context.Context, u entity.User) error {
	query := `INSERT INTO users (id, name, email, role, password_hash) 
              VALUES ( $1, $2, $3, $4, $5)`
	res, err := r.storage.ExecContext(ctx,
		query,
		u.ID,
		u.Name,
		u.Email,
		u.PermissionRole,
		u.PasswordHash)
	if err != nil {
		return fmt.Errorf("store error: %w", err)
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("create rows afected:%w", err)
	}

	if rowsAfected != 1 {
		return fmt.Errorf("create affected: %d", rowsAfected)
	}

	return nil
}

// Update updates a user
func (r *Repository) Update(ctx context.Context, u *entity.User) error {
	query := `UPDATE users SET  name=$2, email=$3, password_hash=$4 WHERE id=$1`
	res, err := r.storage.ExecContext(ctx, query, &u.ID, &u.Name, &u.Email, &u.PasswordHash)
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

// Delete deletes a user from storage
func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	res, err := r.storage.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("no user with %s id", id)
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete rows afected:%w", err)
	}

	if rowsAfected != 1 {
		return fmt.Errorf("delete total affected: %d", rowsAfected)
	}

	return nil
}

// GetByID get single user by id.
func (r *Repository) GetByID(ctx context.Context, id string) (entity.User, error) {
	query := `SELECT name, email, role, password_hash FROM users WHERE id = $1`
	u := entity.User{}
	u.ID = id
	err := r.storage.QueryRowContext(ctx, query, id).Scan(&u.Name, &u.Email, &u.PermissionRole, &u.PasswordHash)
	if err != nil {
		return entity.User{}, fmt.Errorf("there is no users with %s id", id)
	}
	return u, nil
}

// GetByEmail get single user by id.
func (r *Repository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	query := `SELECT id, name, role, password_hash FROM users WHERE email = $1`
	u := entity.User{}
	u.Email = email
	err := r.storage.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Name, &u.PermissionRole, &u.PasswordHash)
	if err != nil {
		return entity.User{}, fmt.Errorf("there is no users with %s id", email)
	}
	return u, nil
}

// GetAll get all users.
func (r *Repository) GetAll(ctx context.Context) (res []entity.User, err error) {
	query := `SELECT id, name, email, role FROM users ORDER BY name`
	rows, err := r.storage.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		u := entity.User{}
		err = rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.PermissionRole,
		)
		if err != nil {
			return nil, fmt.Errorf("rows scan error: %w", err)
		}
		res = append(res, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("encountered during iteration %w", err)
	}
	return
}
