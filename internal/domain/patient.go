package domain

import (
	"context"
	"time"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
	"gorm.io/gorm"
)

// Patient represents a FHIR Patient resource in the database
type Patient struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	FHIRData  []byte         `json:"fhir_data" gorm:"type:jsonb;not null"` // Store FHIR JSON as bytes to avoid invalid UTF-8
	Active    *bool          `json:"active" gorm:"index"`
	Family    string         `json:"family" gorm:"index"`
	Given     string         `json:"given" gorm:"index"`
	Gender    string         `json:"gender" gorm:"type:varchar(20);index"`
	BirthDate *time.Time     `json:"birth_date" gorm:"index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// PatientRepository defines the interface for patient data operations
type PatientRepository interface {
	Create(ctx context.Context, patient *Patient) error
	GetByID(ctx context.Context, id uint) (*Patient, error)
	GetAll(ctx context.Context, limit, offset int) ([]*Patient, error)
	Update(ctx context.Context, patient *Patient) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
}

// PatientService defines the interface for patient business logic
type PatientService interface {
	CreatePatient(ctx context.Context, fhirPatient *fhir.Patient) (*Patient, error)
	GetPatient(ctx context.Context, id uint) (*Patient, error)
	GetPatients(ctx context.Context, limit, offset int) ([]*Patient, int64, error)
	UpdatePatient(ctx context.Context, id uint, fhirPatient *fhir.Patient) (*Patient, error)
	PatchPatient(ctx context.Context, id uint, updates map[string]interface{}) (*Patient, error)
	DeletePatient(ctx context.Context, id uint) error
	ConvertToFHIR(ctx context.Context, patient *Patient) (*fhir.Patient, error)
	ConvertFromFHIR(ctx context.Context, fhirPatient *fhir.Patient) (*Patient, error)
}

// TableName specifies the table name for Patient model
func (Patient) TableName() string {
	return "patients"
}
