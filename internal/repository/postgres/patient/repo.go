package doctor

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
)

// NewRepository create new doctor repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		storage: db,
	}
}

// Repository represent a doctor repository.
type Repository struct {
	storage *sql.DB
}

// Create Add new doctor
func (r *Repository) Create(ctx context.Context, p *entity.Patient) error {
	query := `INSERT INTO patients (id, first_name, last_name) VALUES ($1, $2, $3)`
	res, err := r.storage.ExecContext(ctx,
		query,
		&p.ID,
		&p.FirstName,
		&p.LastName,
	)
	if err != nil {
		return fmt.Errorf("store error: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("create rows affected: %w", err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("create affected: %d", rowsAffected)
	}

	return nil
}

// Delete deletes a doctor from storage
func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM patients WHERE id = $1`
	res, err := r.storage.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("no doctor with %s id", id)
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete rows afected: %w", err)
	}

	if rowsAfected != 1 {
		return fmt.Errorf("delete total affected: %d", rowsAfected)
	}
	return nil
}

// GetByID get single doctor by id.
func (r *Repository) GetByID(ctx context.Context, id string) (entity.Patient, error) {
	query := `SELECT first_name, last_name FROM patients WHERE id = $1`
	d := entity.Patient{}
	d.ID = id
	err := r.storage.QueryRowContext(ctx, query, id).Scan(
		&d.FirstName,
		&d.LastName,
	)

	if err != nil {
		return entity.Patient{}, fmt.Errorf("there is no doctors with %s id", id)
	}
	return d, nil
}

// GetAll get all doctors.
func (r *Repository) GetAll(ctx context.Context) (res []entity.Patient, err error) {
	query := `SELECT id, first_name, last_name FROM patients ORDER BY id`
	rows, err := r.storage.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		d := entity.Patient{}
		err = rows.Scan(
			&d.ID,
			&d.FirstName,
			&d.LastName,
		)
		if err != nil {
			return nil, fmt.Errorf("rows scan error: %w", err)
		}
		res = append(res, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("encountered during iteration %w", err)
	}
	return res, nil
}
