package code

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
)

type Repo struct {
	storage *sql.DB
}

func NewCodeRepo(db *sql.DB) *Repo {
	return &Repo{storage: db}
}

func (r *Repo) Set(ctx context.Context, key, value string) error {
	q := `INSERT INTO password_codes (email, code) VALUES ($1, $2)`

	res, err := r.storage.ExecContext(ctx, q,
		key,
		value)

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

func (r *Repo) Get(ctx context.Context, key string) (string, error) {
	q := `SELECT code FROM password_codes WHERE email = $1`

	c := entity.PasswordCode{Email: key}
	if err := r.storage.QueryRowContext(ctx, q, key).Scan(&c.Code); err != nil {
		return "", fmt.Errorf("no such record with email %v", key)
	}
	return c.Code, nil
}

func (r *Repo) Del(ctx context.Context, key string) error {
	q := `DELETE FROM password_codes WHERE email = $1`
	res, err := r.storage.ExecContext(ctx, q, key)
	if err != nil {
		return fmt.Errorf("no such record with %v", key)
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
