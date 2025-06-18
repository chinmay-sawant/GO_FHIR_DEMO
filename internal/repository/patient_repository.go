package repository

import (
	"go-fhir-demo/internal/domain"
	"go-fhir-demo/pkg/logger"

	"gorm.io/gorm"
)

// PatientRepositoryInterface defines the contract for patient repository
type PatientRepositoryInterface interface {
	Create(patient *domain.Patient) error
	GetByID(id uint) (*domain.Patient, error)
	GetAll(limit, offset int) ([]*domain.Patient, error)
	Update(patient *domain.Patient) error
	Delete(id uint) error
	Count() (int64, error)
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
func (r *patientRepository) Create(patient *domain.Patient) error {
	// Ensure FHIRData is valid UTF-8 and does not contain null bytes

	if err := r.db.Create(patient).Error; err != nil {
		logger.Errorf("Failed to create patient: %v", err)
		return err
	}
	logger.Infof("Patient created successfully with ID: %d", patient.ID)
	return nil
}

// GetByID retrieves a patient by ID
func (r *patientRepository) GetByID(id uint) (*domain.Patient, error) {
	var patient domain.Patient
	if err := r.db.First(&patient, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warnf("Patient not found with ID: %d", id)
			return nil, err
		}
		logger.Errorf("Failed to get patient by ID %d: %v", id, err)
		return nil, err
	}
	return &patient, nil
}

// GetAll retrieves all patients with pagination
func (r *patientRepository) GetAll(limit, offset int) ([]*domain.Patient, error) {
	var patients []*domain.Patient
	query := r.db.Limit(limit).Offset(offset)

	if err := query.Find(&patients).Error; err != nil {
		logger.Errorf("Failed to get patients: %v", err)
		return nil, err
	}

	logger.Infof("Retrieved %d patients", len(patients))
	return patients, nil
}

// Update updates an existing patient record
func (r *patientRepository) Update(patient *domain.Patient) error {
	if err := r.db.Save(patient).Error; err != nil {
		logger.Errorf("Failed to update patient with ID %d: %v", patient.ID, err)
		return err
	}
	logger.Infof("Patient updated successfully with ID: %d", patient.ID)
	return nil
}

// Delete soft deletes a patient record
func (r *patientRepository) Delete(id uint) error {
	if err := r.db.Delete(&domain.Patient{}, id).Error; err != nil {
		logger.Errorf("Failed to delete patient with ID %d: %v", id, err)
		return err
	}
	logger.Infof("Patient deleted successfully with ID: %d", id)
	return nil
}

// Count returns the total number of patients
func (r *patientRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&domain.Patient{}).Count(&count).Error; err != nil {
		logger.Errorf("Failed to count patients: %v", err)
		return 0, err
	}
	return count, nil
}
