package service

import (
	"encoding/json"
	"fmt"
	"time"

	"go-fhir-demo/internal/domain"
	"go-fhir-demo/pkg/logger"
	"go-fhir-demo/pkg/utils"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// PatientServiceInterface defines the contract for patient service
type PatientServiceInterface interface {
	CreatePatient(fhirPatient *fhir.Patient) (*domain.Patient, error)
	GetPatient(id uint) (*domain.Patient, error)
	GetPatients(limit, offset int) ([]*domain.Patient, int64, error)
	UpdatePatient(id uint, fhirPatient *fhir.Patient) (*domain.Patient, error)
	PatchPatient(id uint, updates map[string]interface{}) (*domain.Patient, error)
	DeletePatient(id uint) error
	ConvertToFHIR(patient *domain.Patient) (*fhir.Patient, error)
	ConvertFromFHIR(fhirPatient *fhir.Patient) (*domain.Patient, error)
}

type patientService struct {
	repo domain.PatientRepository
}

// NewPatientService creates a new patient service
func NewPatientService(repo domain.PatientRepository) PatientServiceInterface {
	return &patientService{
		repo: repo,
	}
}

// CreatePatient creates a new patient from FHIR data
func (s *patientService) CreatePatient(fhirPatient *fhir.Patient) (*domain.Patient, error) {
	patient, err := s.ConvertFromFHIR(fhirPatient)
	if err != nil {
		logger.Errorf("Failed to convert FHIR patient: %v", err)
		return nil, err
	}

	if err := s.repo.Create(patient); err != nil {
		return nil, err
	}

	return patient, nil
}

// GetPatient retrieves a patient by ID
func (s *patientService) GetPatient(id uint) (*domain.Patient, error) {
	return s.repo.GetByID(id)
}

// GetPatients retrieves all patients with pagination
func (s *patientService) GetPatients(limit, offset int) ([]*domain.Patient, int64, error) {
	patients, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return patients, count, nil
}

// UpdatePatient updates an existing patient
func (s *patientService) UpdatePatient(id uint, fhirPatient *fhir.Patient) (*domain.Patient, error) {
	existingPatient, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Convert FHIR data to domain model
	updatedPatient, err := s.ConvertFromFHIR(fhirPatient)
	if err != nil {
		return nil, err
	}

	// Preserve ID and timestamps
	updatedPatient.ID = existingPatient.ID
	updatedPatient.CreatedAt = existingPatient.CreatedAt

	if err := s.repo.Update(updatedPatient); err != nil {
		return nil, err
	}

	return updatedPatient, nil
}

// PatchPatient partially updates a patient
func (s *patientService) PatchPatient(id uint, updates map[string]interface{}) (*domain.Patient, error) {
	patient, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Parse existing FHIR data
	fhirPatient, err := s.ConvertToFHIR(patient)
	if err != nil {
		return nil, fmt.Errorf("failed to parse existing FHIR data: %w", err)
	}

	// Apply updates to FHIR patient
	if err := s.applyUpdatesToFHIR(fhirPatient, updates); err != nil {
		return nil, err
	}

	// Convert back to domain model
	updatedPatient, err := s.ConvertFromFHIR(fhirPatient)
	if err != nil {
		return nil, err
	}

	// Preserve ID and timestamps
	updatedPatient.ID = patient.ID
	updatedPatient.CreatedAt = patient.CreatedAt

	if err := s.repo.Update(updatedPatient); err != nil {
		return nil, err
	}

	return updatedPatient, nil
}

// DeletePatient deletes a patient
func (s *patientService) DeletePatient(id uint) error {
	return s.repo.Delete(id)
}

// ConvertToFHIR converts a domain patient to FHIR format
func (s *patientService) ConvertToFHIR(patient *domain.Patient) (*fhir.Patient, error) {
	var fhirPatient fhir.Patient
	if err := json.Unmarshal([]byte(patient.FHIRData), &fhirPatient); err != nil {
		return nil, fmt.Errorf("failed to unmarshal FHIR data: %w", err)
	}
	return &fhirPatient, nil
}

// ConvertFromFHIR converts a FHIR patient to domain format
func (s *patientService) ConvertFromFHIR(fhirPatient *fhir.Patient) (*domain.Patient, error) {
	// Marshal FHIR patient to JSON
	fhirJSON, err := json.Marshal(fhirPatient)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal FHIR patient: %w", err)
	}

	// Remove any invalid UTF-8 byte 0x00 from the JSON
	cleaned := make([]byte, 0, len(fhirJSON))
	for _, b := range fhirJSON {
		if b != 0x00 {
			cleaned = append(cleaned, b)
		}
	}

	patient := &domain.Patient{
		FHIRData: cleaned,
	}

	// Extract searchable fields
	if fhirPatient.Active != nil {
		patient.Active = fhirPatient.Active
	}

	// Extract name information
	if len(fhirPatient.Name) > 0 {
		name := fhirPatient.Name[0]
		if name.Family != nil {
			patient.Family = *name.Family
		}
		if len(name.Given) > 0 {
			patient.Given = name.Given[0]
		}
	}

	// Extract gender
	if fhirPatient.Gender != nil {
		patient.Gender = fhirPatient.Gender.String()
	}

	// Extract birth date
	if fhirPatient.BirthDate != nil {
		if birthDate, err := time.Parse("2006-01-02", *fhirPatient.BirthDate); err == nil {
			patient.BirthDate = &birthDate
		}
	}

	return patient, nil
}

// applyUpdatesToFHIR applies partial updates to a FHIR patient
func (s *patientService) applyUpdatesToFHIR(fhirPatient *fhir.Patient, updates map[string]interface{}) error {
	for key, value := range updates {
		switch key {
		case "active":
			if active, ok := value.(bool); ok {
				fhirPatient.Active = &active
			}
		case "family":
			if family, ok := value.(string); ok {
				if len(fhirPatient.Name) == 0 {
					fhirPatient.Name = []fhir.HumanName{{}}
				}
				fhirPatient.Name[0].Family = &family
			}
		case "given":
			if given, ok := value.(string); ok {
				if len(fhirPatient.Name) == 0 {
					fhirPatient.Name = []fhir.HumanName{{}}
				}
				if fhirPatient.Name[0].Given == nil {
					fhirPatient.Name[0].Given = []string{given}
				} else {
					fhirPatient.Name[0].Given = []string{given}
				}
			}
		case "gender":
			if gender, ok := value.(string); ok {
				// Validate gender values according to FHIR spec
				if gender == "male" || gender == "female" || gender == "other" || gender == "unknown" {
					gender := fhir.AdministrativeGender(*utils.GenderPtr(gender))
					fhirPatient.Gender = &gender
				}
			}
		case "birthDate":
			if birthDate, ok := value.(string); ok {
				// Validate date format
				if _, err := time.Parse("2006-01-02", birthDate); err == nil {
					fhirPatient.BirthDate = &birthDate
				}
			}
		}
	}

	return nil
}
