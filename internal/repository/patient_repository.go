// Package repository provides data access for patients.
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"go-fhir-demo/internal/models"
	"go-fhir-demo/pkg/logger"
	"go-fhir-demo/pkg/utils/tracer"
)



// NewPatientRepository creates a new patient repository
func NewPatientRepository(db *sql.DB) *PatientRepositoryImpl {
	return &PatientRepositoryImpl{db: db}
}

// PatientRepositoryImpl is the SQL implementation.
type PatientRepositoryImpl struct {
	db *sql.DB
}

// Create creates a new patient record.
func (r *PatientRepositoryImpl) Create(ctx context.Context, patient *models.Patient) error {
	ctx, span := tracer.StartSpan(ctx, "Create")
	defer span.End()

	var active any
	if patient.Active != nil {
		active = *patient.Active
	}

	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO patients (fhir_data, active, family, given, gender, birth_date, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		 RETURNING id`,
		patient.FHIRData,
		active,
		patient.Family,
		patient.Given,
		patient.Gender,
		patient.BirthDate,
	).Scan(&patient.ID)
	if err != nil {
		return fmt.Errorf("failed to create patient: %w", err)
	}

	logger.WithContext(ctx).Infof("Patient created successfully with ID: %d", patient.ID)
	return nil
}

// GetByID retrieves a patient by ID.
func (r *PatientRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Patient, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, fhir_data, active, family, given, gender, birth_date, created_at, updated_at
		 FROM patients
		 WHERE id = $1`,
		id,
	)

	var patient models.Patient
	var active sql.NullBool
	var birthDate sql.NullTime
	if err := row.Scan(
		&patient.ID,
		&patient.FHIRData,
		&active,
		&patient.Family,
		&patient.Given,
		&patient.Gender,
		&birthDate,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("failed to get patient by ID %d: no rows found", id)
		}
		return nil, fmt.Errorf("failed to get patient by ID %d: %w", id, err)
	}

	if active.Valid {
		patient.Active = &active.Bool
	}
	if birthDate.Valid {
		patient.BirthDate = &birthDate.Time
	}

	return &patient, nil
}

// GetAll retrieves all patients with pagination.
func (r *PatientRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]*models.Patient, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, fhir_data, active, family, given, gender, birth_date, created_at, updated_at
		 FROM patients
		 ORDER BY id
		 LIMIT $1 OFFSET $2`,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get patients: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.WithContext(ctx).Warnf("failed to close patient rows: %v", err)
		}
	}()

	patients := make([]*models.Patient, 0, limit)
	for rows.Next() {
		var patient models.Patient
		var active sql.NullBool
		var birthDate sql.NullTime
		if err := rows.Scan(
			&patient.ID,
			&patient.FHIRData,
			&active,
			&patient.Family,
			&patient.Given,
			&patient.Gender,
			&birthDate,
			&patient.CreatedAt,
			&patient.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to get patients: %w", err)
		}
		if active.Valid {
			patient.Active = &active.Bool
		}
		if birthDate.Valid {
			patient.BirthDate = &birthDate.Time
		}
		patients = append(patients, &patient)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get patients: %w", err)
	}

	logger.WithContext(ctx).Infof("Retrieved %d patients", len(patients))
	return patients, nil
}

// Update updates an existing patient record.
func (r *PatientRepositoryImpl) Update(ctx context.Context, patient *models.Patient) error {
	var active any
	if patient.Active != nil {
		active = *patient.Active
	}

	if _, err := r.db.ExecContext(
		ctx,
		`UPDATE patients
		 SET fhir_data = $1, active = $2, family = $3, given = $4, gender = $5, birth_date = $6, updated_at = CURRENT_TIMESTAMP
		 WHERE id = $7`,
		patient.FHIRData,
		active,
		patient.Family,
		patient.Given,
		patient.Gender,
		patient.BirthDate,
		patient.ID,
	); err != nil {
		return fmt.Errorf("failed to update patient with ID %d: %w", patient.ID, err)
	}

	logger.WithContext(ctx).Infof("Patient updated successfully with ID: %d", patient.ID)
	return nil
}

// Delete soft deletes a patient record.
func (r *PatientRepositoryImpl) Delete(ctx context.Context, id uint) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM patients WHERE id = $1", id); err != nil {
		return fmt.Errorf("failed to delete patient with ID %d: %w", id, err)
	}
	logger.WithContext(ctx).Infof("Patient deleted successfully with ID: %d", id)
	return nil
}

// Count returns the total number of patients.
func (r *PatientRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM patients").Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count patients: %w", err)
	}
	return count, nil
}
