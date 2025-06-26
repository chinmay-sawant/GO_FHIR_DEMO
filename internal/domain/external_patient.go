package domain

import (
	"context"
	"time"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// ExternalPatientService defines the interface for interacting with external patient data.
type ExternalPatientService interface {
	GetExternalPatientByID(id string) (*fhir.Patient, error)
	SearchExternalPatients(params map[string]string) (*fhir.Bundle, error)
	CreateExternalPatient(patient *fhir.Patient) (*fhir.Patient, error)
	GetExternalPatientByIDDelayed(ctx context.Context, id string, timeout time.Duration) (*fhir.Patient, error)
}
