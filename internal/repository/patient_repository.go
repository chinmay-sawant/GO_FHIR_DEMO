package repository

import (
	"context"
	"go-fhir-demo/internal/domain"
	"go-fhir-demo/pkg/logger"
	"go-fhir-demo/pkg/utils/tracer"

	"gorm.io/gorm"
)

// PatientRepositoryInterface defines the contract for patient repository
type PatientRepositoryInterface interface {
	Create(ctx context.Context, patient *domain.Patient) error
	GetByID(ctx context.Context, id uint) (*domain.Patient, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Patient, error)
	Update(ctx context.Context, patient *domain.Patient) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
}

type patientRepository struct {
	db *gorm.DB
}

// NewPatientRepository creates a new patient repository
func NewPatientRepository(db *gorm.DB) PatientRepositoryInterface {
	return &patientRepository{
		db: db,
	}
}

// Create creates a new patient record
func (r *patientRepository) Create(ctx context.Context, patient *domain.Patient) error {
	ctx, span := tracer.StartSpan(ctx, "Create")
	defer span.End()
	// Ensure FHIRData is valid UTF-8 and does not contain null bytes
	if err := r.db.Create(patient).Error; err != nil {
		logger.WithContext(ctx).Errorf("Failed to create patient: %v", err)
		return err
	}
	logger.WithContext(ctx).Infof("Patient created successfully with ID: %d", patient.ID)
	return nil
}

// GetByID retrieves a patient by ID
func (r *patientRepository) GetByID(ctx context.Context, id uint) (*domain.Patient, error) {
	var patient domain.Patient
	if err := r.db.First(&patient, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithContext(ctx).Warnf("Patient not found with ID: %d", id)
			return nil, err
		}
		logger.WithContext(ctx).Errorf("Failed to get patient by ID %d: %v", id, err)
		return nil, err
	}
	return &patient, nil
}

// GetAll retrieves all patients with pagination
func (r *patientRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Patient, error) {
	var patients []*domain.Patient
	query := r.db.Limit(limit).Offset(offset)

	if err := query.Find(&patients).Error; err != nil {
		logger.WithContext(ctx).Errorf("Failed to get patients: %v", err)
		return nil, err
	}

	logger.WithContext(ctx).Infof("Retrieved %d patients", len(patients))
	return patients, nil
}

// Update updates an existing patient record
func (r *patientRepository) Update(ctx context.Context, patient *domain.Patient) error {
	if err := r.db.Save(patient).Error; err != nil {
		logger.WithContext(ctx).Errorf("Failed to update patient with ID %d: %v", patient.ID, err)
		return err
	}
	logger.WithContext(ctx).Infof("Patient updated successfully with ID: %d", patient.ID)
	return nil
}

// Delete soft deletes a patient record
func (r *patientRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.Delete(&domain.Patient{}, id).Error; err != nil {
		logger.WithContext(ctx).Errorf("Failed to delete patient with ID %d: %v", id, err)
		return err
	}
	logger.WithContext(ctx).Infof("Patient deleted successfully with ID: %d", id)
	return nil
}

// Count returns the total number of patients
func (r *patientRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.Model(&domain.Patient{}).Count(&count).Error; err != nil {
		logger.WithContext(ctx).Errorf("Failed to count patients: %v", err)
		return 0, err
	}
	return count, nil
}
