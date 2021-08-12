package code

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	usecases "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/user/password"
)

var _ usecases.PasswordCodeRepository = (*Repository)(nil)

type Repository struct {
	storage *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{storage: db}
}

func (r *Repository) Create(ctx context.Context, c entity.PasswordCode) error {
	q := `INSERT INTO password_codes (email, code) VALUES ($1, $2)`

	res, err := r.storage.ExecContext(ctx, q,
		c.Email,
		c.Code)

	if err != nil {
		return fmt.Errorf("create passw code: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("create rows affected: %w", err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("create total affected: %v", rowsAffected)
	}
	return nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (entity.PasswordCode, error) {
	q := `SELECT code FROM password_codes WHERE email = $1`

	c := entity.PasswordCode{Email: email}
	if err := r.storage.QueryRowContext(ctx, q, email).Scan(&c.Code); err != nil {
		return entity.PasswordCode{}, fmt.Errorf("no such record with email %v", email)
	}
	return c, nil
}

func (r *Repository) Delete(ctx context.Context, email string) error {
	q := `DELETE FROM password_codes WHERE email = $1`
	res, err := r.storage.ExecContext(ctx, q, email)
	if err != nil {
		return fmt.Errorf("no such record with %v", email)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete rows affected: %w", err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("delete total affected: %v", rowsAffected)
	}

	return nil
}
