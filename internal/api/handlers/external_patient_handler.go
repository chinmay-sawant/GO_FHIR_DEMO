package handlers

import (
	"net/http"

	"go-fhir-demo/internal/domain"
	"go-fhir-demo/pkg/logger"
	"go-fhir-demo/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ExternalPatientHandlerInterface defines the contract for external patient handlers
type ExternalPatientHandlerInterface interface {
	GetExternalPatientByID(c *gin.Context)
	SearchExternalPatients(c *gin.Context)
	CreateExternalPatient(c *gin.Context)
}

// ExternalPatientHandler handles requests for external patient data.
type ExternalPatientHandler struct {
	service domain.ExternalPatientService
}

// NewExternalPatientHandler creates a new ExternalPatientHandler.
func NewExternalPatientHandler(service domain.ExternalPatientService) ExternalPatientHandlerInterface {
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
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID is required"})
		return
	}

	patient, err := h.service.GetExternalPatientByID(id)
	if err != nil {
		// Basic error handling, can be improved to differentiate 404 from 500
		logger.Errorf("Failed to get external patient by ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve patient from external server", "details": err.Error()})
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
	queryParams := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0] // Taking the first value for simplicity
		}
	}

	bundle, err := h.service.SearchExternalPatients(queryParams)
	if err != nil {
		logger.Errorf("Failed to search external patients: %v", err)
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
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient data", "details": err.Error()})
		return
	}

	patient, err := utils.ConvertJsonToFHIRPatient(jsonData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to convert to FHIR Patient", "details": err.Error()})
		return
	}

	createdPatient, err := h.service.CreateExternalPatient(patient)
	if err != nil {
		logger.Errorf("Failed to create external patient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create patient on external server", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPatient)
}
