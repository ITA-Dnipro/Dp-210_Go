package appointment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	usecases "github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/usecases/appointment"
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
	query := `
	INSERT INTO appointments 
	(id , doctor_id,  patient_id , reason, time_range) 
	VALUES 
	( $1, $2, $3, $4, tstzrange($5, $6))`
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ExclusionViolation {
			return fmt.Errorf("time is already taken")
		}
		return fmt.Errorf("create rows:%w", err)
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
	query := `DELETE FROM appointments WHERE id = $1`
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

//GetWithFilter.
func (r *Repository) GetWithFilter(ctx context.Context, filter entity.AppointmentFilter) ([]entity.Appointment, error) {
	var b strings.Builder
	b.WriteString(`
	SELECT id, doctor_id, patient_id, lower(time_range), upper(time_range) 
	FROM appointments 
	WHERE 1=1`)
	args := []interface{}{}
	if filter.DoctorID != nil {
		args = append(args, *filter.DoctorID)
		b.WriteString(fmt.Sprintf(` AND doctor_id = $%d`, len(args)))
	}
	if filter.PatientID != nil {
		args = append(args, *filter.PatientID)
		b.WriteString(fmt.Sprintf(` AND patient_id = $%d`, len(args)))
	}
	if filter.From != nil {
		args = append(args, *filter.From)
		b.WriteString(fmt.Sprintf(` AND lower(time_range) > $%d`, len(args)))
	}
	if filter.To != nil {
		args = append(args, *filter.To)
		b.WriteString(fmt.Sprintf(` AND upper(time_range) < $%d`, len(args)))
	}
	b.WriteString(` ORDER BY lower(time_range)`)
	query := b.String()
	return r.fetch(ctx, query, args...)
}
