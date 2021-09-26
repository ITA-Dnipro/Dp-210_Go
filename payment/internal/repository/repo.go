package repository

import (
	"database/sql"
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/entity"
)

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		storage: db,
	}
}

type Repository struct {
	storage *sql.DB
}

func (repo *Repository) InsertDataToDb(bill *entity.Bill) error {
	query1 := `INSERT INTO patients(patient_id, patient_total) VALUES ($1, $2) 
				ON CONFLICT (patient_id) DO UPDATE 
				SET patient_total = patients.patient_total + $2;`
	query2 := `INSERT INTO doctors(doctor_id, doctor_total) VALUES ($1, $2) 
				ON CONFLICT (doctor_id) DO UPDATE 
				SET doctor_total = doctors.doctor_total + $2;`

	tx, err := repo.storage.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction error > %s", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer tx.Rollback()

	_, err = tx.Exec(query1, bill.PatientID, bill.Price)
	if err != nil {
		return fmt.Errorf("insert into users error > %s", err)
	}

	_, err = tx.Exec(query2, bill.DoctorID, bill.Price)
	if err != nil {
		return fmt.Errorf("insert into doctors error > %s", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx commit error > %s", err)
	}
	return nil
}

//goland:noinspection GoUnhandledErrorResult,SqlWithoutWhere
func (repo *Repository) GetAllDataFromDb() (*entity.DocPat, error) {
	var docPat entity.DocPat

	query1 := `SELECT * FROM patients;`
	query2 := `DELETE FROM patients;`
	query3 := `SELECT * FROM doctors;`
	query4 := `DELETE FROM doctors;`

	tx, err := repo.storage.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction error > %s", err)
	}
	defer tx.Rollback()

	// query1 := `SELECT * FROM patients;`
	patientsRows, err := tx.Query(query1)
	if err != nil {
		return nil, fmt.Errorf("query1 error > %s", err)
	}
	defer patientsRows.Close()
	for patientsRows.Next() {
		var p entity.Patient
		if err = patientsRows.Scan(
			&p.PatientId,
			&p.PatientTotal,
		); err != nil {
			return nil, fmt.Errorf("rows scan error > %s", err)
		}
		docPat.Patients = append(docPat.Patients, p)
	}

	// query2 := `DELETE FROM patients;`
	_, err = tx.Exec(query2)
	if err != nil {
		return nil, fmt.Errorf("query2 error > %s", err)
	}

	// query3 := `SELECT * FROM doctors;`
	doctorsRows, err := tx.Query(query3)
	if err != nil {
		return nil, fmt.Errorf("query1 error > %s", err)
	}
	defer doctorsRows.Close()
	for doctorsRows.Next() {
		var d entity.Doctor
		if err = doctorsRows.Scan(
			&d.DoctorId,
			&d.DoctorTotal,
		); err != nil {
			return nil, fmt.Errorf("rows scan error > %s", err)
		}
		docPat.Doctors = append(docPat.Doctors, d)
	}

	// query4 := `DELETE FROM doctors;`
	_, err = tx.Exec(query4)
	if err != nil {
		return nil, fmt.Errorf("query3 error > %s", err)
	}

	// tx commit!
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit error > %s", err)
	}

	return &docPat, nil
}
