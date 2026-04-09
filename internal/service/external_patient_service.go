// Package service provides business logic for the application.
package service

import (
	"context"
	"time"

	"go-fhir-demo/internal/models"
	"go-fhir-demo/pkg/cache"
	"go-fhir-demo/pkg/fhirclient"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// ExternalPatientService defines the contract for external patient service
type ExternalPatientService interface {
	GetExternalPatientByID(ctx context.Context, id string) (*fhir.Patient, error)
	SearchExternalPatients(ctx context.Context, params map[string]string) (*fhir.Bundle, error)
	CreateExternalPatient(ctx context.Context, patient *fhir.Patient) (*fhir.Patient, error)
	GetPatientCached(ctx context.Context, id string) (*fhir.Patient, error)
	GetPatientDelayed(ctx context.Context, id string, timeout time.Duration) (*fhir.Patient, error)
}

// NoopExternalPatientService is a second implementation for deslop
type NoopExternalPatientService struct{}

// GetExternalPatientByID is a no-op implementation.
func (NoopExternalPatientService) GetExternalPatientByID(_ context.Context, _ string) (*fhir.Patient, error) {
	return nil, models.ErrNotFound
}

// SearchExternalPatients is a no-op implementation.
func (NoopExternalPatientService) SearchExternalPatients(_ context.Context, _ map[string]string) (*fhir.Bundle, error) {
	return nil, models.ErrInternal
}

// CreateExternalPatient is a no-op implementation.
func (NoopExternalPatientService) CreateExternalPatient(_ context.Context, _ *fhir.Patient) (*fhir.Patient, error) {
	return nil, models.ErrInternal
}

// GetPatientCached is a no-op implementation.
func (NoopExternalPatientService) GetPatientCached(_ context.Context, _ string) (*fhir.Patient, error) {
	return nil, models.ErrNotFound
}

// GetPatientDelayed is a no-op implementation.
func (NoopExternalPatientService) GetPatientDelayed(_ context.Context, _ string, _ time.Duration) (*fhir.Patient, error) {
	return nil, models.ErrNotFound
}

// ExternalPatientServiceImpl implements ExternalPatientService.
type ExternalPatientServiceImpl struct {
	client fhirclient.Client
	cache  cache.RedisCache
}

// NewExternalPatientService creates a new ExternalPatientService.
func NewExternalPatientService(client fhirclient.Client, cache cache.RedisCache) ExternalPatientService {
	return &ExternalPatientServiceImpl{
		client: client,
		cache:  cache,
	}
}

// GetExternalPatientByID retrieves a patient from the external FHIR server by ID.
func (s *ExternalPatientServiceImpl) GetExternalPatientByID(ctx context.Context, id string) (*fhir.Patient, error) {
	return s.client.GetPatientByID(ctx, id)
}

// SearchExternalPatients searches for patients on the external FHIR server.
func (s *ExternalPatientServiceImpl) SearchExternalPatients(ctx context.Context, params map[string]string) (*fhir.Bundle, error) {
	return s.client.SearchPatients(ctx, params)
}

// CreateExternalPatient creates a patient on the external FHIR server.
func (s *ExternalPatientServiceImpl) CreateExternalPatient(ctx context.Context, patient *fhir.Patient) (*fhir.Patient, error) {
	return s.client.CreatePatient(ctx, patient)
}

// GetPatientCached retrieves a patient with Redis caching
func (s *ExternalPatientServiceImpl) GetPatientCached(ctx context.Context, id string) (*fhir.Patient, error) {
	// Try to get from cache first
	if s.cache != nil {
		cachedPatient, err := s.cache.GetPatient(ctx, id)
		if err == nil && cachedPatient != nil {
			return cachedPatient, nil
		}
	}

	// Cache miss or error, fetch from external FHIR server
	patient, err := s.client.GetPatientByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests (expire after 1 hour)
	if s.cache != nil {
		_ = s.cache.SetPatient(ctx, id, patient, time.Hour)
	}

	return patient, nil
}

// GetPatientDelayed retrieves a patient with timeout logic
func (s *ExternalPatientServiceImpl) GetPatientDelayed(ctx context.Context, id string, timeout time.Duration) (*fhir.Patient, error) {
	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Call GetPatientByID with timeout context
	patient, err := s.client.GetPatientByID(timeoutCtx, id)
	if err != nil {
		return nil, err
	}

	return patient, nil
}
