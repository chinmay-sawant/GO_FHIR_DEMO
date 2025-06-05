package domain

import (
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
	Create(patient *Patient) error
	GetByID(id uint) (*Patient, error)
	GetAll(limit, offset int) ([]*Patient, error)
	Update(patient *Patient) error
	Delete(id uint) error
	Count() (int64, error)
}

// PatientService defines the interface for patient business logic
type PatientService interface {
	CreatePatient(fhirPatient *fhir.Patient) (*Patient, error)
	GetPatient(id uint) (*Patient, error)
	GetPatients(limit, offset int) ([]*Patient, int64, error)
	UpdatePatient(id uint, fhirPatient *fhir.Patient) (*Patient, error)
	PatchPatient(id uint, updates map[string]interface{}) (*Patient, error)
	DeletePatient(id uint) error
	ConvertToFHIR(patient *Patient) (*fhir.Patient, error)
	ConvertFromFHIR(fhirPatient *fhir.Patient) (*Patient, error)
}

// TableName specifies the table name for Patient model
func (Patient) TableName() string {
	return "patients"
}
