package service

import (
	"encoding/json"
	"fmt"
	"time"

	"go-fhir-demo/internal/domain"
	"go-fhir-demo/pkg/logger"
	"strconv"
)

type patientService struct {
	repo domain.PatientRepository
}

// NewPatientService creates a new patient service
func NewPatientService(repo domain.PatientRepository) domain.PatientService {
	return &patientService{
		repo: repo,
	}
}

// CreatePatient creates a new patient from FHIR data
func (s *patientService) CreatePatient(fhirPatient *domain.FHIRPatient) (*domain.Patient, error) {
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
func (s *patientService) UpdatePatient(id uint, fhirPatient *domain.FHIRPatient) (*domain.Patient, error) {
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
	var fhirPatient domain.FHIRPatient
	if err := json.Unmarshal([]byte(patient.FHIRData), &fhirPatient); err != nil {
		return nil, fmt.Errorf("failed to parse existing FHIR data: %w", err)
	}

	// Apply updates to FHIR patient
	if err := s.applyUpdatesToFHIR(&fhirPatient, updates); err != nil {
		return nil, err
	}

	// Convert back to domain model
	updatedPatient, err := s.ConvertFromFHIR(&fhirPatient)
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
func (s *patientService) ConvertToFHIR(patient *domain.Patient) (*domain.FHIRPatient, error) {
	var fhirPatient domain.FHIRPatient
	if err := json.Unmarshal([]byte(patient.FHIRData), &fhirPatient); err != nil {
		return nil, fmt.Errorf("failed to unmarshal FHIR data: %w", err)
	}
	return &fhirPatient, nil
}

// ConvertFromFHIR converts a FHIR patient to domain format
func (s *patientService) ConvertFromFHIR(fhirPatient *domain.FHIRPatient) (*domain.Patient, error) {
	// Ensure resource type is set
	fhirPatient.ResourceType = "Patient"

	// Marshal FHIR patient to JSON
	fhirJSON, err := json.Marshal(fhirPatient)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal FHIR patient: %w", err)
	}

	patient := &domain.Patient{
		FHIRData: string(fhirJSON),
	}

	// Extract searchable fields
	if fhirPatient.Active != nil {
		patient.Active = fhirPatient.Active
	}

	// Extract name information
	if len(fhirPatient.Name) > 0 {
		name := fhirPatient.Name[0]
		patient.Family = name.Family
		if len(name.Given) > 0 {
			patient.Given = name.Given[0]
		}
	}

	// Extract gender
	patient.Gender = fhirPatient.Gender

	// Extract birth date
	if fhirPatient.BirthDate != "" {
		if birthDate, err := time.Parse("2006-01-02", fhirPatient.BirthDate); err == nil {
			patient.BirthDate = &birthDate
		}
	}

	return patient, nil
}

// applyUpdatesToFHIR applies partial updates to a FHIR patient
func (s *patientService) applyUpdatesToFHIR(fhirPatient *domain.FHIRPatient, updates map[string]interface{}) error {
	for key, value := range updates {
		switch key {
		case "active":
			if active, ok := value.(bool); ok {
				fhirPatient.Active = &active
			}
		case "family":
			if family, ok := value.(string); ok {
				if len(fhirPatient.Name) == 0 {
					fhirPatient.Name = []domain.FHIRHumanName{{}}
				}
				fhirPatient.Name[0].Family = family
			}
		case "given":
			if given, ok := value.(string); ok {
				if len(fhirPatient.Name) == 0 {
					fhirPatient.Name = []domain.FHIRHumanName{{}}
				}
				fhirPatient.Name[0].Given = []string{given}
			}
		case "gender":
			if gender, ok := value.(string); ok {
				// Validate gender values according to FHIR spec
				if gender == "male" || gender == "female" || gender == "other" || gender == "unknown" {
					fhirPatient.Gender = gender
				}
			}
		case "birthDate":
			if birthDate, ok := value.(string); ok {
				// Validate date format
				if _, err := time.Parse("2006-01-02", birthDate); err == nil {
					fhirPatient.BirthDate = birthDate
				}
			}
		}
	}

	// Convert back to domain model
	updatedPatient, err := s.ConvertFromFHIR(fhirPatient)
	if err != nil {
		return err
	}

	//Preserve ID and timestamps
	// Convert the ID to uint
	idUint, err := strconv.ParseInt(fhirPatient.ID, 10, 64)
	if err != nil {
		return err
	}
	updatedPatient.ID = uint(idUint)

	// updatedPatient.CreatedAt = fhirPatient.

	if err := s.repo.Update(updatedPatient); err != nil {
		return err
	}

	return nil
}
