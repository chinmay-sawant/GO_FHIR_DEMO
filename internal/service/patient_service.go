package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-fhir-demo/internal/models"
	"go-fhir-demo/pkg/fhirconv"
	"go-fhir-demo/pkg/patch"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// PatientRepository defines the data access contract for patients
type PatientRepository interface {
	Create(ctx context.Context, patient *models.Patient) error
	GetByID(ctx context.Context, id uint) (*models.Patient, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Patient, error)
	Update(ctx context.Context, patient *models.Patient) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
}

// PatientService defines the contract for patient service
type PatientService interface {
	CreatePatient(ctx context.Context, patient *fhir.Patient) (*models.Patient, error)
	GetPatient(ctx context.Context, id uint) (*models.Patient, error)
	GetPatients(ctx context.Context, limit, offset int) ([]*models.Patient, int64, error)
	UpdatePatient(ctx context.Context, id uint, patient *fhir.Patient) (*models.Patient, error)
	DeletePatient(ctx context.Context, id uint) error
	PatchPatient(ctx context.Context, id uint, updates patch.PatientPatch) (*models.Patient, error)
	ConvertToFHIR(ctx context.Context, patient *models.Patient) (*fhir.Patient, error)
	ConvertFromFHIR(ctx context.Context, fhirPatient *fhir.Patient) (*models.Patient, error)
}

// NoopPatientService is a second implementation for deslop
type NoopPatientService struct{}

// CreatePatient is a no-op implementation.
func (NoopPatientService) CreatePatient(_ context.Context, _ *fhir.Patient) (*models.Patient, error) {
	return nil, models.ErrInternal
}

// GetPatient is a no-op implementation.
func (NoopPatientService) GetPatient(_ context.Context, _ uint) (*models.Patient, error) {
	return nil, models.ErrNotFound
}

// GetPatients is a no-op implementation.
func (NoopPatientService) GetPatients(_ context.Context, _, _ int) ([]*models.Patient, int64, error) {
	return nil, 0, models.ErrInternal
}

// UpdatePatient is a no-op implementation.
func (NoopPatientService) UpdatePatient(_ context.Context, _ uint, _ *fhir.Patient) (*models.Patient, error) {
	return nil, models.ErrNotFound
}

// DeletePatient is a no-op implementation.
func (NoopPatientService) DeletePatient(_ context.Context, _ uint) error { return nil }

// PatchPatient is a no-op implementation.
func (NoopPatientService) PatchPatient(_ context.Context, _ uint, _ patch.PatientPatch) (*models.Patient, error) {
	return nil, models.ErrNotFound
}

// ConvertToFHIR is a no-op implementation.
func (NoopPatientService) ConvertToFHIR(_ context.Context, _ *models.Patient) (*fhir.Patient, error) {
	return nil, models.ErrInternal
}

// ConvertFromFHIR is a no-op implementation.
func (NoopPatientService) ConvertFromFHIR(_ context.Context, _ *fhir.Patient) (*models.Patient, error) {
	return nil, models.ErrInternal
}

// PatientServiceImpl implements PatientService
type PatientServiceImpl struct {
	repo PatientRepository
}

// NewPatientService creates a new patient service
func NewPatientService(repo PatientRepository) PatientService {
	return &PatientServiceImpl{
		repo: repo,
	}
}

// CreatePatient creates a new patient from FHIR data
func (s *PatientServiceImpl) CreatePatient(ctx context.Context, fhirPatient *fhir.Patient) (*models.Patient, error) {
	patient, err := s.ConvertFromFHIR(ctx, fhirPatient)
	if err != nil {
		return nil, fmt.Errorf("failed to convert FHIR patient: %w", err)
	}

	if err := s.repo.Create(ctx, patient); err != nil {
		return nil, err
	}

	return patient, nil
}

// GetPatient retrieves a patient by ID
func (s *PatientServiceImpl) GetPatient(ctx context.Context, id uint) (*models.Patient, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid patient id: 0")
	}
	patient, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("patient service GetPatient error: %w", err)
	}
	return patient, nil
}

// GetPatients retrieves all patients with pagination
func (s *PatientServiceImpl) GetPatients(ctx context.Context, limit, offset int) ([]*models.Patient, int64, error) {
	if limit <= 0 {
		limit = 10 // default limit
	}
	if offset < 0 {
		offset = 0 // default offset
	}
	patients, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return patients, count, nil
}

// UpdatePatient updates an existing patient
func (s *PatientServiceImpl) UpdatePatient(ctx context.Context, id uint, fhirPatient *fhir.Patient) (*models.Patient, error) {
	existingPatient, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Convert FHIR data to domain model
	updatedPatient, err := s.ConvertFromFHIR(ctx, fhirPatient)
	if err != nil {
		return nil, err
	}

	// Preserve ID and timestamps
	updatedPatient.ID = existingPatient.ID
	updatedPatient.CreatedAt = existingPatient.CreatedAt

	if err := s.repo.Update(ctx, updatedPatient); err != nil {
		return nil, err
	}

	return updatedPatient, nil
}

// PatchPatient partially updates a patient
func (s *PatientServiceImpl) PatchPatient(ctx context.Context, id uint, updates patch.PatientPatch) (*models.Patient, error) {
	patient, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Parse existing FHIR data
	fhirPatient, err := s.ConvertToFHIR(ctx, patient)
	if err != nil {
		return nil, fmt.Errorf("failed to parse existing FHIR data: %w", err)
	}

	// Apply updates to FHIR patient
	if err := s.applyUpdatesToFHIR(fhirPatient, updates); err != nil {
		return nil, err
	}

	// Convert back to domain model
	updatedPatient, err := s.ConvertFromFHIR(ctx, fhirPatient)
	if err != nil {
		return nil, err
	}

	// Preserve ID and timestamps
	updatedPatient.ID = patient.ID
	updatedPatient.CreatedAt = patient.CreatedAt

	if err := s.repo.Update(ctx, updatedPatient); err != nil {
		return nil, err
	}

	return updatedPatient, nil
}

// DeletePatient deletes a patient
func (s *PatientServiceImpl) DeletePatient(ctx context.Context, id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid patient id: 0")
	}
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("patient service DeletePatient error: %w", err)
	}
	return nil
}

// ConvertToFHIR converts a domain patient to FHIR format
func (s *PatientServiceImpl) ConvertToFHIR(_ context.Context, patient *models.Patient) (*fhir.Patient, error) {
	var fhirPatient fhir.Patient
	if err := json.Unmarshal([]byte(patient.FHIRData), &fhirPatient); err != nil {
		return nil, fmt.Errorf("failed to unmarshal FHIR data: %w", err)
	}
	return &fhirPatient, nil
}

// ConvertFromFHIR converts a FHIR patient to domain format
func (s *PatientServiceImpl) ConvertFromFHIR(_ context.Context, fhirPatient *fhir.Patient) (*models.Patient, error) {
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

	patient := &models.Patient{
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
func (s *PatientServiceImpl) applyUpdatesToFHIR(fhirPatient *fhir.Patient, updates patch.PatientPatch) error {
	if updates.Active != nil {
		fhirPatient.Active = updates.Active
	}
	if updates.Family != nil {
		if len(fhirPatient.Name) == 0 {
			fhirPatient.Name = []fhir.HumanName{{}}
		}
		fhirPatient.Name[0].Family = updates.Family
	}
	if updates.Given != nil {
		if len(fhirPatient.Name) == 0 {
			fhirPatient.Name = []fhir.HumanName{{}}
		}
		fhirPatient.Name[0].Given = []string{*updates.Given}
	}
	if updates.Gender != nil {
		// Validate gender values according to FHIR spec
		switch *updates.Gender {
		case "male", "female", "other", "unknown":
			gender := fhir.AdministrativeGender(*fhirconv.GenderPtr(*updates.Gender))
			fhirPatient.Gender = &gender
		}
	}
	if updates.BirthDate != nil {
		// Validate date format
		if _, err := time.Parse("2006-01-02", *updates.BirthDate); err == nil {
			fhirPatient.BirthDate = updates.BirthDate
		}
	}

	return nil
}
