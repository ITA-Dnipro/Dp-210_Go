package doctor

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	usecases "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/doctor"
)

var _ usecases.DoctorsRepository = (*Repository)(nil)

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
func (r *Repository) Create(ctx context.Context, d entity.Doctor) error {
	query := `INSERT INTO doctors (id, user_id, first_name, last_name,
			  speciality, schedule_id) VALUES ($1, $2, $3, $4, $5, $6)`
	res, err := r.storage.ExecContext(ctx,
		query,
		d.ID,
		d.UserID,
		d.FirstName,
		d.LastName,
		d.Speciality,
		d.ScheduleId)
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

// Update updates a doctor
func (r *Repository) Update(ctx context.Context, d *entity.Doctor) error {
	query := `UPDATE doctors SET id = $1, user_id = $2, first_name = $3, 
			  last_name = $4, speciality = $5, schedule_id = $6 `
	res, err := r.storage.ExecContext(ctx,
		query,
		&d.ID,
		&d.UserID,
		&d.FirstName,
		&d.LastName,
		&d.Speciality,
		&d.ScheduleId)

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

// Delete deletes a doctor from storage
func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM doctors WHERE id = $1`
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
func (r *Repository) GetByID(ctx context.Context, id string) (entity.Doctor, error) {
	query := `SELECT id, user_id, first_name, last_name, speciality,
			  schedule_id FROM doctors WHERE id = $1`
	d := entity.Doctor{}
	d.ID = id
	err := r.storage.QueryRowContext(ctx, query, id).Scan(&d.ID,
														  &d.UserID,
														  &d.FirstName,
														  &d.LastName,
														  &d.Speciality,
														  &d.ScheduleId)

	if err != nil {
		return entity.Doctor{}, fmt.Errorf("there is no doctors with %s id", id)
	}
	return d, nil
}

// GetAll get all doctors.
func (r *Repository) GetAll(ctx context.Context) (res []entity.Doctor, err error) {
	query := `SELECT id, user_id, first_name, last_name, speciality,
			  schedule_id FROM doctors ORDER BY id`
	rows, err := r.storage.QueryContext(ctx, query)                   
	if err != nil {                                                   
	    return nil, fmt.Errorf("query error: %w", err)                
	}                                                                 
	defer rows.Close()                                                
	for rows.Next() {                                                 
	    d := entity.Doctor{}                                            
	    err = rows.Scan(&d.ID,
						&d.UserID,
						&d.FirstName,
						&d.LastName,
						&d.Speciality,
						&d.ScheduleId)
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

