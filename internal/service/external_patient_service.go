package service

import (
	"context"
	"time"

	"go-fhir-demo/pkg/cache"
	"go-fhir-demo/pkg/fhirclient"
	"go-fhir-demo/pkg/logger"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// ExternalPatientServiceInterface defines the contract for external patient service
type ExternalPatientServiceInterface interface {
	GetExternalPatientByID(id string) (*fhir.Patient, error)
	GetExternalPatientByIDCached(ctx context.Context, id string) (*fhir.Patient, error)
	GetExternalPatientByIDDelayed(ctx context.Context, id string, timeout time.Duration) (*fhir.Patient, error)
	SearchExternalPatients(params map[string]string) (*fhir.Bundle, error)
	CreateExternalPatient(patient *fhir.Patient) (*fhir.Patient, error)
}

type externalPatientService struct {
	client fhirclient.ClientInterface
	cache  cache.CacheInterface
}

// NewExternalPatientService creates a new ExternalPatientService.
func NewExternalPatientService(client fhirclient.ClientInterface, cache cache.CacheInterface) ExternalPatientServiceInterface {
	return &externalPatientService{
		client: client,
		cache:  cache,
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

// CreateExternalPatient creates a patient on the external FHIR server.
func (s *externalPatientService) CreateExternalPatient(patient *fhir.Patient) (*fhir.Patient, error) {
	return s.client.CreatePatient(patient)
}

// GetExternalPatientByIDCached retrieves a patient with Redis caching
func (s *externalPatientService) GetExternalPatientByIDCached(ctx context.Context, id string) (*fhir.Patient, error) {
	// Try to get from cache first
	if s.cache != nil {
		cachedPatient, err := s.cache.GetPatient(ctx, id)
		if err != nil {
			logger.Warnf("Failed to get patient from cache: %v", err)
		} else if cachedPatient != nil {
			logger.Infof("Patient %s retrieved from cache", id)
			return cachedPatient, nil
		}
	}

	// Cache miss or error, fetch from external FHIR server
	logger.Infof("Cache miss for patient %s, fetching from external server", id)
	patient, err := s.client.GetPatientByID(id)
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests (expire after 1 hour)
	if s.cache != nil {
		if err := s.cache.SetPatient(ctx, id, patient, time.Hour); err != nil {
			logger.Warnf("Failed to cache patient %s: %v", id, err)
		} else {
			logger.Infof("Patient %s cached successfully", id)
		}
	}

	return patient, nil
}

// GetExternalPatientByIDDelayed retrieves a patient with timeout logic
func (s *externalPatientService) GetExternalPatientByIDDelayed(ctx context.Context, id string, timeout time.Duration) (*fhir.Patient, error) {
	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Channel to receive the result
	resultChan := make(chan struct {
		patient *fhir.Patient
		err     error
	}, 1)

	// Start the API call in a goroutine
	go func() {
		patient, err := s.client.GetPatientByID(id)
		resultChan <- struct {
			patient *fhir.Patient
			err     error
		}{patient, err}
	}()

	// Wait for either the result or timeout
	select {
	case result := <-resultChan:
		if result.err != nil {
			logger.Errorf("Failed to get patient %s from external server: %v", id, result.err)
			return nil, result.err
		}
		logger.Infof("Patient %s retrieved from external server within timeout", id)
		return result.patient, nil
	case <-timeoutCtx.Done():
		logger.Errorf("Timeout occurred while fetching patient %s from external server", id)
		return nil, timeoutCtx.Err()
	}
}
