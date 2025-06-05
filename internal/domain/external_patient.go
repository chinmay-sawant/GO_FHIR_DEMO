package domain

import (
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// ExternalPatientService defines the interface for interacting with external patient data.
type ExternalPatientService interface {
	GetExternalPatientByID(id string) (*fhir.Patient, error)
	SearchExternalPatients(params map[string]string) (*fhir.Bundle, error)
}
