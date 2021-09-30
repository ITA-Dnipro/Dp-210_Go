package appointment

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/server/customerrors"
	"github.com/google/uuid"
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
			return customerrors.ErrTime
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

// Update updates a appointment.
func (r *Repository) Update(ctx context.Context, a *entity.Appointment) error {
	query := `
	UPDATE appointments 
	SET doctor_id=$2,  patient_id=$3, reason=$4, time_range=tstzrange($5, $6)
	WHERE id=$1`
	res, err := r.storage.ExecContext(
		ctx,
		query,
		&a.ID,
		&a.DoctorID,
		&a.PatientID,
		&a.Reason,
		&a.From,
		&a.To,
	)
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update rows afected:%w", err)
	}

	if rowsAfected == 0 {
		return customerrors.ErrNotFound
	}

	return nil
}

// Delete deletes a appointment from storage.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM appointments WHERE id = $1`
	res, err := r.storage.ExecContext(ctx, query, id)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return customerrors.ErrNotFound
		}
		return fmt.Errorf("no doctor with %s id", id)
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete rows afected: %w", err)
	}

	if rowsAfected == 0 {
		return customerrors.ErrNotFound
	}
	return nil
}

// GetByID get single doctor by id.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (entity.Appointment, error) {
	query := `SELECT doctor_id, patient_id, lower(time_range), upper(time_range) 
	FROM appointments 
	WHERE id = $1`
	a := entity.Appointment{}
	a.ID = id
	err := r.storage.QueryRowContext(ctx, query, id).Scan(
		&a.DoctorID,
		&a.PatientID,
		&a.From,
		&a.To,
	)

	if err != nil {
		return entity.Appointment{}, fmt.Errorf("there is no appointments with %s id", id)
	}
	return a, nil
}

func (r *Repository) fetch(ctx context.Context, query string, args ...interface{}) (res []entity.Appointment, err error) {
	rows, err := r.storage.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	// nolint:errcheck
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
		a.From = a.From.UTC()
		a.To = a.To.UTC()
		res = append(res, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("encountered during iteration %w", err)
	}
	return res, nil
}

//GetAll get all appointments.
func (r *Repository) GetAll(ctx context.Context, p *entity.AppointmentsParam) (res []entity.Appointment, cursor string, err error) {
	query := `
	SELECT id, doctor_id, patient_id, lower(time_range), upper(time_range) 
	FROM appointments
	WHERE 1=1`
	q, a, err := addParam(p, query)
	if err != nil {
		return
	}
	res, err = r.fetch(ctx, q, a...)
	if err != nil {
		return
	}
	if len(res) != 0 && len(res) == p.Limits {
		next := res[len(res)-1]
		cursor = encodeCursor(next.From, next.ID)
		return
	}
	return res, cursor, nil
}

//GetByPatientID get all appointments by patient id.
func (r *Repository) GetByPatientID(ctx context.Context, id uuid.UUID, p *entity.AppointmentsParam) (res []entity.Appointment, cursor string, err error) {
	query := `
	SELECT id, doctor_id, patient_id, lower(time_range), upper(time_range) 
	FROM appointments 
	WHERE patient_id = $1`
	q, a, err := addParam(p, query, id)
	if err != nil {
		return
	}
	res, err = r.fetch(ctx, q, a...)
	if err != nil {
		return
	}
	if len(res) != 0 && len(res) == p.Limits {
		next := res[len(res)-1]
		cursor = encodeCursor(next.From, next.ID)
		return
	}
	return res, cursor, nil
}

//GetByDoctorID get all appointments by doctor id.
func (r *Repository) GetByDoctorID(ctx context.Context, id uuid.UUID, p *entity.AppointmentsParam) (res []entity.Appointment, cursor string, err error) {
	query := `
	SELECT id, doctor_id, patient_id, lower(time_range), upper(time_range) 
	FROM appointments 
	WHERE doctor_id = $1`
	q, a, err := addParam(p, query, id)
	if err != nil {
		return
	}
	res, err = r.fetch(ctx, q, a...)
	if err != nil {
		return
	}
	if len(res) != 0 && len(res) == p.Limits {
		next := res[len(res)-1]
		cursor = encodeCursor(next.From, next.ID)
		return
	}
	return res, cursor, nil
}

func addParam(p *entity.AppointmentsParam, query string, args ...interface{}) (string, []interface{}, error) {
	if p.Cursor != "" {
		ct, id, err := decodeCursor(p.Cursor)
		if err != nil {
			return "", nil, customerrors.ErrBadParamInput
		}
		args = append(args, id)
		query += ` AND id < $` + strconv.Itoa(len(args))
		if p.From.Before(ct) {
			p.From = ct
		}
	}
	if !p.From.IsZero() {
		args = append(args, p.From)
		query += ` AND upper(time_range) > $` + strconv.Itoa(len(args))
	}
	if !p.To.IsZero() {
		args = append(args, p.To)
		query += ` AND lower(time_range) < $` + strconv.Itoa(len(args))
	}
	if p.Limits == 0 {
		p.Limits = 10
	}
	args = append(args, p.Limits)
	query += ` ORDER BY lower(time_range) LIMIT $` + strconv.Itoa(len(args))
	return query, args, nil
}

func decodeCursor(encodedCursor string) (res time.Time, id uuid.UUID, err error) {
	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return
	}

	arrStr := strings.Split(string(byt), ",")
	if len(arrStr) != 2 {
		err = errors.New("cursor is invalid")
		return
	}

	res, err = time.Parse(time.RFC3339Nano, arrStr[0])
	if err != nil {
		return
	}
	id, err = uuid.Parse(arrStr[1])
	return
}

func encodeCursor(t time.Time, uuid uuid.UUID) string {
	key := fmt.Sprintf("%s,%s", t.Format(time.RFC3339Nano), uuid)
	return base64.StdEncoding.EncodeToString([]byte(key))
}
