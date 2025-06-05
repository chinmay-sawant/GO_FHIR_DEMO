package service

import (
	"go-fhir-demo/internal/domain"
	"go-fhir-demo/pkg/fhirclient"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

type externalPatientService struct {
	client *fhirclient.Client
}

// NewExternalPatientService creates a new ExternalPatientService.
func NewExternalPatientService(client *fhirclient.Client) domain.ExternalPatientService {
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
