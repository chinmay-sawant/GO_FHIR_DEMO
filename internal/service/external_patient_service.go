package service

import (
	"go-fhir-demo/pkg/fhirclient"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// ExternalPatientServiceInterface defines the contract for external patient service
type ExternalPatientServiceInterface interface {
	GetExternalPatientByID(id string) (*fhir.Patient, error)
	SearchExternalPatients(params map[string]string) (*fhir.Bundle, error)
}

type externalPatientService struct {
	client fhirclient.ClientInterface
}

// NewExternalPatientService creates a new ExternalPatientService.
func NewExternalPatientService(client fhirclient.ClientInterface) ExternalPatientServiceInterface {
	return &externalPatientService{
		client: client,
	}
}

// GetExternalPatientByID retrieves a patient from the external FHIR server by ID.
func (s *externalPatientService) GetExternalPatientByID(id string) (*fhir.Patient, error) {
	return s.client.GetPatientByID(id)
}

// SearchExternalPatients searches for patients on the external FHIR server.
func (s *externalPatientService) SearchExternalPatients(params map[string]string) (*fhir.Bundle, error) {
	return s.client.SearchPatients(params)
}
