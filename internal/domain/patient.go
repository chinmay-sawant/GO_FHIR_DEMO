package domain

import (
	"time"

	"gorm.io/gorm"
)

// FHIRPatient represents a simplified FHIR Patient resource
type FHIRPatient struct {
	ResourceType string             `json:"resourceType"`
	ID           string             `json:"id,omitempty"`
	Active       *bool              `json:"active,omitempty"`
	Name         []FHIRHumanName    `json:"name,omitempty"`
	Gender       string             `json:"gender,omitempty"`
	BirthDate    string             `json:"birthDate,omitempty"`
	Telecom      []FHIRContactPoint `json:"telecom,omitempty"`
	Address      []FHIRAddress      `json:"address,omitempty"`
}

type FHIRHumanName struct {
	Use    string   `json:"use,omitempty"`
	Family string   `json:"family,omitempty"`
	Given  []string `json:"given,omitempty"`
}

type FHIRContactPoint struct {
	System string `json:"system,omitempty"`
	Value  string `json:"value,omitempty"`
	Use    string `json:"use,omitempty"`
}

type FHIRAddress struct {
	Use        string   `json:"use,omitempty"`
	Line       []string `json:"line,omitempty"`
	City       string   `json:"city,omitempty"`
	State      string   `json:"state,omitempty"`
	PostalCode string   `json:"postalCode,omitempty"`
	Country    string   `json:"country,omitempty"`
}

// Patient represents a FHIR Patient resource in the database
type Patient struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	FHIRData  string         `json:"fhir_data" gorm:"type:jsonb;not null"` // Store FHIR JSON
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
	CreatePatient(fhirPatient *FHIRPatient) (*Patient, error)
	GetPatient(id uint) (*Patient, error)
	GetPatients(limit, offset int) ([]*Patient, int64, error)
	UpdatePatient(id uint, fhirPatient *FHIRPatient) (*Patient, error)
	PatchPatient(id uint, updates map[string]interface{}) (*Patient, error)
	DeletePatient(id uint) error
	ConvertToFHIR(patient *Patient) (*FHIRPatient, error)
	ConvertFromFHIR(fhirPatient *FHIRPatient) (*Patient, error)
}

// TableName specifies the table name for Patient model
func (Patient) TableName() string {
	return "patients"
}
