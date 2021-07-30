package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/user/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/user/usecases"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
)

var _ usecases.UsersRepository = (*Repository)(nil)

// suggest to have it as repository method
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		storage: db,
	}
}

type Repository struct {
	storage *sql.DB
}

func (r *Repository) Create(ctx context.Context, u entity.User) (string, error) {
	u.ID = uuid.New().String()
	query := `INSERT INTO users (id, first_name, last_name) VALUES ( $1, $2, $3) RETURNING id`
	row := r.storage.QueryRowContext(ctx, query, u.ID, u.Firstname, u.Lastname)
	err := row.Scan(&u.ID)
	if err != nil {
		return "", fmt.Errorf("store error: %w", err)
	}

	return u.ID, nil
}

func (r *Repository) Update(ctx context.Context, u entity.User) error {
	query := `UPDATE users SET  first_name=$2, last_name=$3 WHERE id = $1`
	_, err := r.storage.ExecContext(ctx, query, u.ID, u.Firstname, u.Lastname)
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.storage.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("no user with %s id", id)
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (entity.User, error) {
	query := `SELECT first_name, last_name FROM users WHERE id =$1`
	u := entity.User{}
	u.ID = id
	err := r.storage.QueryRowContext(ctx, query, id).Scan(&u.Firstname, &u.Lastname)
	if err != nil {
		return entity.User{}, fmt.Errorf("there is no users with %s id", id)
	}
	return u, nil
}

func (r *Repository) GetAll(ctx context.Context) (res []entity.User, err error) {
	query := `SELECT id,first_name,last_name FROM users ORDER BY first_name`
	rows, err := r.storage.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		u := entity.User{}
		err = rows.Scan(
			&u.ID,
			&u.Firstname,
			&u.Lastname,
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

// MigrateUp runs migration and applies everything new to the DB provided in dsn string
func MigrateUp(migrationsPath, dsn string) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dsn)
	if err != nil {
		return fmt.Errorf("migration failed, %v", err)
	}

	if err := m.Up(); err != nil {
		if err.Error() != "no change" {
			return fmt.Errorf("migration failed, %v", err)
		}
	}
	return nil
}
