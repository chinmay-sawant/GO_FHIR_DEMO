package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"go-fhir-demo/internal/service"
	"go-fhir-demo/pkg/logger"
	"go-fhir-demo/pkg/utils"
	"go-fhir-demo/pkg/utils/tracer"

	"github.com/gin-gonic/gin"
)

// ExternalPatientHandlerInterface defines the contract for external patient handlers
type ExternalPatientHandlerInterface interface {
	GetExternalPatientByID(c *gin.Context)
	GetExternalPatientByIDCached(c *gin.Context)
	GetExternalPatientByIDDelayed(c *gin.Context)
	SearchExternalPatients(c *gin.Context)
	CreateExternalPatient(c *gin.Context)
}

// ExternalPatientHandler handles requests for external patient data.
type ExternalPatientHandler struct {
	service service.ExternalPatientServiceInterface
}

// NewExternalPatientHandler creates a new ExternalPatientHandler.
func NewExternalPatientHandler(service service.ExternalPatientServiceInterface) ExternalPatientHandlerInterface {
	return &ExternalPatientHandler{
		service: service,
	}
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
func (h *ExternalPatientHandler) GetExternalPatientByID(c *gin.Context) {
	ctx, span := tracer.StartSpan(c.Request.Context(), "GetExternalPatientByID")
	defer span.End()
	logger.WithContext(ctx).Infof("Fetching external patient by ID: %s", c.Param("id"))
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID is required"})
		return
	}

	patient, err := h.service.GetExternalPatientByID(ctx, id)
	if err != nil {
		// Basic error handling, can be improved to differentiate 404 from 500
		logger.WithContext(ctx).Errorf("Failed to get external patient by ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve patient from external server", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, patient)
}

// GetExternalPatientByIDCached godoc
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
func (h *ExternalPatientHandler) GetExternalPatientByIDCached(c *gin.Context) {
	ctx, span := tracer.StartSpan(c.Request.Context(), "GetExternalPatientByIDCached")
	defer span.End()

	id := c.Param("id")
	logger.WithContext(ctx).Infof("Fetching cached external patient by ID: %s", id)

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID is required"})
		return
	}

	// Attempt to retrieve patient from cache or external server
	logger.WithContext(ctx).Infof("Attempting to get cached external patient by ID: %s", id)

	patient, err := h.service.GetExternalPatientByIDCached(ctx, id)
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to get cached external patient by ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve patient from external server or cache",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, patient)
}

// GetExternalPatientByIDDelayed godoc
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
func (h *ExternalPatientHandler) GetExternalPatientByIDDelayed(c *gin.Context) {
	ctx, span := tracer.StartSpan(c.Request.Context(), "GetExternalPatientByIDDelayed")
	defer span.End()

	id := c.Param("id")
	logger.WithContext(ctx).Infof("Fetching delayed external patient by ID: %s", id)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID is required"})
		return
	}

	// Parse timeout parameter (default: 10 seconds)
	timeoutStr := c.DefaultQuery("timeout", "10")
	timeoutSeconds, err := strconv.Atoi(timeoutStr)
	if err != nil || timeoutSeconds <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timeout parameter"})
		return
	}

	timeout := time.Duration(timeoutSeconds) * time.Second

	logger.WithContext(ctx).Infof("Using timeout of %d seconds for patient %s", timeoutSeconds, id)

	patient, err := h.service.GetExternalPatientByIDDelayed(ctx, id, timeout)
	if err != nil {
		if err == context.DeadlineExceeded {
			logger.WithContext(ctx).Errorf("Timeout occurred for patient %s: %v", id, err)
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error":   "Request timeout",
				"details": "The external FHIR server did not respond within the specified timeout",
			})
			return
		}
		logger.WithContext(ctx).Errorf("Failed to get delayed external patient by ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve patient from external server",
			"details": err.Error(),
		})
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
	queryParams := make(map[string]string)

	logger.WithContext(ctx).Infof("Searching external patients with query parameters: %v", c.Request.URL.Query())
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0] // Taking the first value for simplicity
		}
	}

	bundle, err := h.service.SearchExternalPatients(ctx, queryParams)
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to search external patients: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search patients on external server", "details": err.Error()})
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

	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient data", "details": err.Error()})
		return
	}

	logger.WithContext(ctx).Infof("Creating external patient with data: %v", jsonData)

	patient, err := utils.ConvertJsonToFHIRPatient(jsonData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to convert to FHIR Patient", "details": err.Error()})
		return
	}

	createdPatient, err := h.service.CreateExternalPatient(ctx, patient)
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to create external patient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create patient on external server", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPatient)
}
