package patientdata

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"os"
)

// Repository is struct to hold db connection.
type Repository struct {
	storage *sql.DB
}

// NewRepository returning repo with db.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		storage: db,
	}
}

// CreateRecord for sending gotten data to db.
func (r *Repository) CreateRecord(ctx context.Context, pc *entity.Patient) error {
	query := `INSERT INTO patients(id, name, email, gender, age, phone, address, disability, reg_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
	res, err := r.storage.ExecContext(ctx, query,
		pc.Id, pc.Name, pc.Email, pc.Gender, pc.Age, pc.Phone, pc.Address, pc.Disability, pc.RegAt,
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
	_, _ = fmt.Fprintf(os.Stdout, "%#v\n", pc)

	return nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) *entity.Patient {
	query := `SELECT * FROM patients WHERE email = $1`
	var patient entity.Patient
	err := r.storage.QueryRowContext(ctx, query, email).Scan(
		&patient.Id,
		&patient.Name,
		&patient.Email,
		&patient.Gender,
		&patient.Age,
		&patient.Phone,
		&patient.Address,
		&patient.Disability,
		&patient.RegAt,
	)
	if err != nil {
		return nil
	}
	return &patient
}
