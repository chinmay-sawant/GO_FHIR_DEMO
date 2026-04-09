package handlers

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"go-fhir-demo/internal/service"
	"go-fhir-demo/pkg/fhirconv"
	"go-fhir-demo/pkg/logger"
	"go-fhir-demo/pkg/utils/tracer"

	"github.com/gin-gonic/gin"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// ExternalPatientHandler handles requests for external patient data.
type ExternalPatientHandler struct {
	service service.ExternalPatientService
}

// NewExternalPatientHandler creates a new ExternalPatientHandler.
func NewExternalPatientHandler(service service.ExternalPatientService) *ExternalPatientHandler {
	return &ExternalPatientHandler{
		service: service,
	}
}

func (h *ExternalPatientHandler) createExternalPatientResponse(ctx context.Context, body io.Reader) (*fhir.Patient, error) {
	raw, err := mustRawJSONBody(body)
	if err != nil {
		return nil, err
	}
	patient, err := fhirconv.ConvertJSONToFHIRPatient(raw)
	if err != nil {
		return nil, err
	}
	return h.service.CreateExternalPatient(ctx, patient)
}

// GetExternalPatientByID godoc
// @Summary Get an external patient by ID
// @Description Retrieves a patient resource from an external FHIR server by its ID
// @Tags ExternalPatients
// @Produce json
// @Param id path string true "Patient ID"
// @Success 200 {object} fhir.Patient "Successfully retrieved patient"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Patient not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /external-patients/{id} [get]
// @Router /external-patients/{id} [get]
func (h *ExternalPatientHandler) GetExternalPatientByID(c *gin.Context) {
	ctx, span := tracer.StartSpan(c.Request.Context(), "GetExternalPatientByID")
	defer span.End()
	id, err := mustStringParam(c.Param("id"), "id")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	logger.WithContext(ctx).Infof("Fetching external patient by ID: %s", id)

	patient, err := h.service.GetExternalPatientByID(ctx, id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, patient)
}

// GetPatientCached godoc
// @Summary Get an external patient by ID with caching
// @Description Retrieves a patient resource from an external FHIR server by its ID with Redis caching
// @Tags ExternalPatients
// @Produce json
// @Param id path string true "Patient ID"
// @Success 200 {object} fhir.Patient "Successfully retrieved patient"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Patient not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /external-patients/{id}/cached [get]
func (h *ExternalPatientHandler) GetPatientCached(c *gin.Context) {
	ctx, span := tracer.StartSpan(c.Request.Context(), "GetPatientCached")
	defer span.End()

	id, err := mustStringParam(c.Param("id"), "id")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	logger.WithContext(ctx).Infof("Fetching cached external patient by ID: %s", id)

	// Attempt to retrieve patient from cache or external server
	logger.WithContext(ctx).Infof("Attempting to get cached external patient by ID: %s", id)

	patient, err := h.service.GetPatientCached(ctx, id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, patient)
}

// GetPatientDelayed godoc
// @Summary Get an external patient by ID with timeout
// @Description Retrieves a patient resource from an external FHIR server by its ID with configurable timeout
// @Tags ExternalPatients
// @Produce json
// @Param id path string true "Patient ID"
// @Param timeout query int false "Timeout in seconds (default: 10)"
// @Success 200 {object} fhir.Patient "Successfully retrieved patient"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Patient not found"
// @Failure 408 {object} map[string]string "Request timeout"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /external-patients/{id}/delayed [get]
func (h *ExternalPatientHandler) GetPatientDelayed(c *gin.Context) {
	ctx, span := tracer.StartSpan(c.Request.Context(), "GetPatientDelayed")
	defer span.End()

	id, err := mustStringParam(c.Param("id"), "id")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	logger.WithContext(ctx).Infof("Fetching delayed external patient by ID: %s", id)

	// Parse timeout parameter (default: 10 seconds)
	timeoutStr := c.DefaultQuery("timeout", "10")
	timeoutSeconds, err := strconv.Atoi(timeoutStr)
	if err != nil || timeoutSeconds <= 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	timeout := time.Duration(timeoutSeconds) * time.Second

	logger.WithContext(ctx).Infof("Using timeout of %d seconds for patient %s", timeoutSeconds, id)

	patient, err := h.service.GetPatientDelayed(ctx, id, timeout)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, patient)
}

// SearchExternalPatients godoc
// @Summary Search for external patients
// @Description Searches for patient resources on an external FHIR server based on query parameters
// @Tags ExternalPatients
// @Produce json
// @Param _query query string false "FHIR search parameters (e.g., name=John,birthdate=1990-01-01)"
// @Success 200 {object} fhir.Bundle "Successfully retrieved search results"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /external-patients [get]
func (h *ExternalPatientHandler) SearchExternalPatients(c *gin.Context) {
	ctx, span := tracer.StartSpan(c.Request.Context(), "SearchExternalPatients")
	defer span.End()

	filters := parseRawQueryParams(c.Request.URL.RawQuery)
	logger.WithContext(ctx).Infof("Searching external patients with query parameters: %v", filters)

	bundle, err := h.service.SearchExternalPatients(ctx, filters)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, bundle)
}

// CreateExternalPatient godoc
// @Summary Create an external patient
// @Description Creates a new patient resource on an external FHIR server
// @Tags ExternalPatients
// @Accept json
// @Produce json
// @Param patient body object true "Patient resource to create (FHIR-compliant JSON)"
// @Success 201 {object} fhir.Patient "Successfully created patient"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /external-patients [post]
func (h *ExternalPatientHandler) CreateExternalPatient(c *gin.Context) {
	ctx, span := tracer.StartSpan(c.Request.Context(), "CreateExternalPatient")
	defer span.End()
	if !requireCSRFToken(c.GetHeader("X-CSRF-Token")) {
		c.Status(http.StatusInternalServerError)
		return
	}

	createdPatient, err := h.createExternalPatientResponse(ctx, c.Request.Body)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, createdPatient)
}
