package patientdata

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	usecases "github.com/ITA-Dnipro/Dp-210_Go/internal/usecases/patientdata"
)

var _ usecases.PtDataRepository = (*Repository)(nil)

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
func (r *Repository) CreateRecord(ctx context.Context, pc *entity.PatientData) error {
	query := `INSERT INTO data_from_patients(id, first_name, last_name, email, gender, birthday_str, phone, 
                               address, job_info, disability, allergies, reg_date, patient_role) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);`
	res, err := r.storage.ExecContext(ctx, query,
		pc.Id, pc.FirstName, pc.LastName, pc.Email, pc.Gender, pc.BirthDayStr, pc.Phone, pc.Address,
		pc.JobInfo, pc.DisabilityB, pc.AllergiesB, pc.RegDayTime, pc.Role,
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
