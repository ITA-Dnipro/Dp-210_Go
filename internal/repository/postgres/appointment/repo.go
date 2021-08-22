package appointment

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"

	usecases "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/appointment"
)

var _ usecases.AppointmentsRepository = (*Repository)(nil)

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

// Create Add new appointment.
func (r *Repository) Create(ctx context.Context, a *entity.Appointment) error {
	query := `INSERT INTO appointments (id , doctor_id,  patient_id , reason, time_range) 
              VALUES ( $1, $2, $3, $4, tstzrange($5, $6))`
	res, err := r.storage.ExecContext(ctx,
		query,
		a.ID,
		a.DoctorID,
		a.PatientID,
		a.Reason,
		a.From,
		a.To,
	)
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

// Delete deletes a appointment from storage.
func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM appointment WHERE id = $1`
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

func (r *Repository) fetch(ctx context.Context, query string, args ...interface{}) (res []entity.Appointment, err error) {
	rows, err := r.storage.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		a := entity.Appointment{}
		err = rows.Scan(
			&a.ID,
			&a.DoctorID,
			&a.PatientID,
			&a.From,
			&a.To,
		)
		if err != nil {
			return nil, fmt.Errorf("rows scan error: %w", err)
		}
		res = append(res, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("encountered during iteration %w", err)
	}
	return res, nil
}

//GetAll get all appointments.
func (r *Repository) GetAll(ctx context.Context) ([]entity.Appointment, error) {
	query := `SELECT id, doctor_id, patient_id, lower(time_range), upper(time_range) FROM appointments ORDER BY lower(time_range)`
	return r.fetch(ctx, query)
}

//GetByPatientID get all appointments by patient id.
func (r *Repository) GetByPatientID(ctx context.Context, id string) ([]entity.Appointment, error) {
	query := `SELECT id, doctor_id, patient_id, lower(time_range), upper(time_range) FROM appointments WHERE patient_id = $1 ORDER BY lower(time_range)`
	return r.fetch(ctx, query, id)
}

//GetByDoctorID get all appointments by doctor id.
func (r *Repository) GetByDoctorID(ctx context.Context, id string) ([]entity.Appointment, error) {
	query := `SELECT id, doctor_id, patient_id, lower(time_range), upper(time_range) FROM appointments WHERE doctor_id = $1 ORDER BY lower(time_range)`
	return r.fetch(ctx, query, id)
}

//GetByDoctorID get all appointments by user id.
func (r *Repository) GetByUserID(ctx context.Context, id string) ([]entity.Appointment, error) {
	query := `SELECT id, doctor_id, patient_id, lower(time_range), upper(time_range) FROM appointments WHERE doctor_id = $1 OR patient_id = $1 ORDER BY lower(time_range)`
	return r.fetch(ctx, query, id)
}
