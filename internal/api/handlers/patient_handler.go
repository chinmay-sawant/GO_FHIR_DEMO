package handlers

import (
	"net/http"
	"strconv"

	"go-fhir-demo/internal/domain"
	"go-fhir-demo/pkg/logger"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	service domain.PatientService
}

// NewPatientHandler creates a new patient handler
func NewPatientHandler(service domain.PatientService) *PatientHandler {
	return &PatientHandler{
		service: service,
	}
}

// CreatePatient handles POST /patients
// @Summary Create a new Patient
// @Description Create a new FHIR Patient resource
// @Tags Patient
// @Accept json
// @Produce json
// @Param patient body domain.FHIRPatient true "FHIR Patient resource"
// @Success 201 {object} domain.FHIRPatient
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patients [post]
func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var fhirPatient domain.FHIRPatient

	if err := c.ShouldBindJSON(&fhirPatient); err != nil {
		logger.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON",
			"message": err.Error(),
		})
		return
	}

	// Create patient
	patient, err := h.service.CreatePatient(&fhirPatient)
	if err != nil {
		logger.Errorf("Failed to create patient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create patient",
			"message": err.Error(),
		})
		return
	}

	// Convert back to FHIR for response
	fhirResponse, err := h.service.ConvertToFHIR(patient)
	if err != nil {
		logger.Errorf("Failed to convert to FHIR: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to convert response",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, fhirResponse)
}

// GetPatient handles GET /patients/:id
// @Summary Get a Patient by ID
// @Description Get a FHIR Patient resource by its ID
// @Tags Patient
// @Produce json
// @Param id path int true "Patient ID"
// @Success 200 {object} domain.FHIRPatient
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patients/{id} [get]
func (h *PatientHandler) GetPatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid patient ID",
			"message": "Patient ID must be a valid number",
		})
		return
	}

	patient, err := h.service.GetPatient(uint(id))
	if err != nil {
		logger.Errorf("Failed to get patient: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Patient not found",
			"message": err.Error(),
		})
		return
	}

	// Convert to FHIR for response
	fhirResponse, err := h.service.ConvertToFHIR(patient)
	if err != nil {
		logger.Errorf("Failed to convert to FHIR: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to convert response",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, fhirResponse)
}

// GetPatients handles GET /patients
// @Summary List Patients
// @Description Get a list of FHIR Patient resources
// @Tags Patient
// @Produce json
// @Param _count query int false "Number of patients to return"
// @Param _offset query int false "Offset for pagination"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patients [get]
func (h *PatientHandler) GetPatients(c *gin.Context) {
	// Parse query parameters
	limitStr := c.DefaultQuery("_count", "10")
	offsetStr := c.DefaultQuery("_offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get patients
	patients, total, err := h.service.GetPatients(limit, offset)
	if err != nil {
		logger.Errorf("Failed to get patients: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve patients",
			"message": err.Error(),
		})
		return
	}

	// Convert to FHIR bundle response
	fhirPatients := make([]*domain.FHIRPatient, len(patients))
	for i, patient := range patients {
		fhirPatient, err := h.service.ConvertToFHIR(patient)
		if err != nil {
			logger.Errorf("Failed to convert patient to FHIR: %v", err)
			continue
		}
		fhirPatients[i] = fhirPatient
	}

	// FHIR Bundle format
	entries := make([]map[string]interface{}, len(fhirPatients))
	for i, fp := range fhirPatients {
		entries[i] = map[string]interface{}{
			"resource": fp,
		}
	}

	response := gin.H{
		"resourceType": "Bundle",
		"type":         "searchset",
		"total":        total,
		"entry":        entries,
	}

	c.JSON(http.StatusOK, response)
}

// UpdatePatient handles PUT /patients/:id
// @Summary Update a Patient
// @Description Update a FHIR Patient resource by its ID
// @Tags Patient
// @Accept json
// @Produce json
// @Param id path int true "Patient ID"
// @Param patient body domain.FHIRPatient true "FHIR Patient resource"
// @Success 200 {object} domain.FHIRPatient
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patients/{id} [put]
func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid patient ID",
			"message": "Patient ID must be a valid number",
		})
		return
	}

	var fhirPatient domain.FHIRPatient
	if err := c.ShouldBindJSON(&fhirPatient); err != nil {
		logger.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON",
			"message": err.Error(),
		})
		return
	}

	// Update patient
	patient, err := h.service.UpdatePatient(uint(id), &fhirPatient)
	if err != nil {
		logger.Errorf("Failed to update patient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update patient",
			"message": err.Error(),
		})
		return
	}

	// Convert back to FHIR for response
	fhirResponse, err := h.service.ConvertToFHIR(patient)
	if err != nil {
		logger.Errorf("Failed to convert to FHIR: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to convert response",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, fhirResponse)
}

// PatchPatient handles PATCH /patients/:id
// @Summary Partially update a Patient
// @Description Partially update a FHIR Patient resource by its ID
// @Tags Patient
// @Accept json
// @Produce json
// @Param id path int true "Patient ID"
// @Param patch body object true "Partial update fields"
// @Success 200 {object} domain.FHIRPatient
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patients/{id} [patch]
func (h *PatientHandler) PatchPatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid patient ID",
			"message": "Patient ID must be a valid number",
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		logger.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON",
			"message": err.Error(),
		})
		return
	}

	// Patch patient
	patient, err := h.service.PatchPatient(uint(id), updates)
	if err != nil {
		logger.Errorf("Failed to patch patient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to patch patient",
			"message": err.Error(),
		})
		return
	}

	// Convert back to FHIR for response
	fhirResponse, err := h.service.ConvertToFHIR(patient)
	if err != nil {
		logger.Errorf("Failed to convert to FHIR: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to convert response",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, fhirResponse)
}

// DeletePatient handles DELETE /patients/:id
// @Summary Delete a Patient
// @Description Delete a FHIR Patient resource by its ID
// @Tags Patient
// @Param id path int true "Patient ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patients/{id} [delete]
func (h *PatientHandler) DeletePatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid patient ID",
			"message": "Patient ID must be a valid number",
		})
		return
	}

	if err := h.service.DeletePatient(uint(id)); err != nil {
		logger.Errorf("Failed to delete patient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete patient",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
