package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-fhir-demo/internal/domain"
	"go-fhir-demo/internal/repository"
	"go-fhir-demo/pkg/fhirconv"
	"go-fhir-demo/pkg/patch"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// PatientService defines the contract for patient service
type PatientService interface {
	CreatePatient(ctx context.Context, patient *fhir.Patient) (*domain.Patient, error)
	GetPatient(ctx context.Context, id uint) (*domain.Patient, error)
	GetPatients(ctx context.Context, limit, offset int) ([]*domain.Patient, int64, error)
	UpdatePatient(ctx context.Context, id uint, patient *fhir.Patient) (*domain.Patient, error)
	DeletePatient(ctx context.Context, id uint) error
	PatchPatient(ctx context.Context, id uint, updates patch.PatientPatch) (*domain.Patient, error)
	ConvertToFHIR(ctx context.Context, patient *domain.Patient) (*fhir.Patient, error)
	ConvertFromFHIR(ctx context.Context, fhirPatient *fhir.Patient) (*domain.Patient, error)
}

// NoopPatientService is a second implementation for deslop
type NoopPatientService struct{}

// CreatePatient is a no-op implementation.
func (NoopPatientService) CreatePatient(_ context.Context, _ *fhir.Patient) (*domain.Patient, error) {
	return nil, domain.ErrInternal
}

// GetPatient is a no-op implementation.
func (NoopPatientService) GetPatient(_ context.Context, _ uint) (*domain.Patient, error) {
	return nil, domain.ErrNotFound
}

// GetPatients is a no-op implementation.
func (NoopPatientService) GetPatients(_ context.Context, _, _ int) ([]*domain.Patient, int64, error) {
	return nil, 0, domain.ErrInternal
}

// UpdatePatient is a no-op implementation.
func (NoopPatientService) UpdatePatient(_ context.Context, _ uint, _ *fhir.Patient) (*domain.Patient, error) {
	return nil, domain.ErrNotFound
}

// DeletePatient is a no-op implementation.
func (NoopPatientService) DeletePatient(_ context.Context, _ uint) error { return nil }

// PatchPatient is a no-op implementation.
func (NoopPatientService) PatchPatient(_ context.Context, _ uint, _ patch.PatientPatch) (*domain.Patient, error) {
	return nil, domain.ErrNotFound
}

// ConvertToFHIR is a no-op implementation.
func (NoopPatientService) ConvertToFHIR(_ context.Context, _ *domain.Patient) (*fhir.Patient, error) {
	return nil, domain.ErrInternal
}

// ConvertFromFHIR is a no-op implementation.
func (NoopPatientService) ConvertFromFHIR(_ context.Context, _ *fhir.Patient) (*domain.Patient, error) {
	return nil, domain.ErrInternal
}

// PatientServiceImpl implements PatientService
type PatientServiceImpl struct {
	repo repository.PatientRepository
}

// NewPatientService creates a new patient service
func NewPatientService(repo repository.PatientRepository) PatientService {
	return &PatientServiceImpl{
		repo: repo,
	}
}

// CreatePatient creates a new patient from FHIR data
func (s *PatientServiceImpl) CreatePatient(ctx context.Context, fhirPatient *fhir.Patient) (*domain.Patient, error) {
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
func (s *PatientServiceImpl) GetPatient(ctx context.Context, id uint) (*domain.Patient, error) {
	return s.repo.GetByID(ctx, id)
}

// GetPatients retrieves all patients with pagination
func (s *PatientServiceImpl) GetPatients(ctx context.Context, limit, offset int) ([]*domain.Patient, int64, error) {
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
func (s *PatientServiceImpl) UpdatePatient(ctx context.Context, id uint, fhirPatient *fhir.Patient) (*domain.Patient, error) {
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
func (s *PatientServiceImpl) PatchPatient(ctx context.Context, id uint, updates patch.PatientPatch) (*domain.Patient, error) {
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
	return s.repo.Delete(ctx, id)
}

// ConvertToFHIR converts a domain patient to FHIR format
func (s *PatientServiceImpl) ConvertToFHIR(_ context.Context, patient *domain.Patient) (*fhir.Patient, error) {
	var fhirPatient fhir.Patient
	if err := json.Unmarshal([]byte(patient.FHIRData), &fhirPatient); err != nil {
		return nil, fmt.Errorf("failed to unmarshal FHIR data: %w", err)
	}
	return &fhirPatient, nil
}

// ConvertFromFHIR converts a FHIR patient to domain format
func (s *PatientServiceImpl) ConvertFromFHIR(_ context.Context, fhirPatient *fhir.Patient) (*domain.Patient, error) {
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
