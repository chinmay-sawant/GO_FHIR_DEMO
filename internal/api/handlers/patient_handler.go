package handlers

import (
	"net/http"
	"strconv"

	"go-fhir-demo/internal/domain"
	"go-fhir-demo/pkg/logger"
	"go-fhir-demo/pkg/utils/tracer"

	"github.com/gin-gonic/gin"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// PatientHandlerInterface defines the contract for patient handlers
type PatientHandlerInterface interface {
	CreatePatient(c *gin.Context)
	GetPatient(c *gin.Context)
	GetPatients(c *gin.Context)
	UpdatePatient(c *gin.Context)
	PatchPatient(c *gin.Context)
	DeletePatient(c *gin.Context)
}

// PatientHandler struct
type PatientHandler struct {
	service domain.PatientService
}

// NewPatientHandler creates a new patient handler
func NewPatientHandler(service domain.PatientService) PatientHandlerInterface {
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
// @Param patient body fhir.Patient true "FHIR Patient resource"
// @Success 201 {object} fhir.Patient
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patients [post]
func (h *PatientHandler) CreatePatient(c *gin.Context) {
	ctx, span := tracer.StartSpan(c.Request.Context(), "CreatePatient")
	defer span.End()

	var fhirPatient fhir.Patient

	if err := c.ShouldBindJSON(&fhirPatient); err != nil {
		logger.WithContext(ctx).Errorf("Failed to bind JSON: %v", err)
		// Return a 400 Bad Request if JSON binding fails
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON",
			"message": err.Error(),
		})
		return
	}

	// Create patient
	patient, err := h.service.CreatePatient(ctx, &fhirPatient)
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to create patient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create patient",
			"message": err.Error(),
		})
		return
	}

	// Convert back to FHIR for response
	fhirResponse, err := h.service.ConvertToFHIR(ctx, patient)
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to convert to FHIR: %v", err)
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
// @Success 200 {object} fhir.Patient
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patients/{id} [get]
func (h *PatientHandler) GetPatient(c *gin.Context) {
	ctx, span := tracer.StartSpan(c.Request.Context(), "GetPatient")
	defer span.End()

	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid patient ID",
			"message": "Patient ID is required",
		})
		return
	}
	// Parse the ID from the URL parameter
	id, err := strconv.ParseUint(idStr, 10, 64)
	// If parsing fails, return a 400 Bad Request
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid patient ID",
			"message": "Patient ID must be a valid number",
		})
		return
	}
	// Fetch the patient from the service
	logger.WithContext(ctx).Infof("Fetching patient with ID: %d", id)
	patient, err := h.service.GetPatient(ctx, uint(id))
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to get patient: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Patient not found",
			"message": err.Error(),
		})
		return
	}

	fhirPatient, err := h.service.ConvertToFHIR(ctx, patient)
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to convert to FHIR: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to convert patient data",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, fhirPatient)
}

// GetPatients handles GET /patients
// @Summary Get all Patients
// @Description Get all FHIR Patient resources with pagination
// @Tags Patient
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patients [get]
func (h *PatientHandler) GetPatients(c *gin.Context) {
	ctx, span := tracer.StartSpan(c.Request.Context(), "GetPatients")
	defer span.End()

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	logger.WithContext(ctx).Infof("Fetching patients with limit %d and offset %d", limit, offset)

	patients, total, err := h.service.GetPatients(ctx, limit, offset)
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to get patients: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get patients",
			"message": err.Error(),
		})
		return
	}

	// Convert patients to FHIR format
	fhirPatients := make([]*fhir.Patient, 0, len(patients))
	for _, patient := range patients {
		fhirPatient, err := h.service.ConvertToFHIR(ctx, patient)
		if err != nil {
			logger.WithContext(ctx).Warnf("Failed to convert patient %d to FHIR: %v", patient.ID, err)
			continue
		}
		fhirPatients = append(fhirPatients, fhirPatient)
	}

	c.JSON(http.StatusOK, gin.H{
		"patients": fhirPatients,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

// UpdatePatient handles PUT /patients/:id
// @Summary Update a Patient
// @Description Update an existing FHIR Patient resource
// @Tags Patient
// @Accept json
// @Produce json
// @Param id path int true "Patient ID"
// @Param patient body fhir.Patient true "FHIR Patient resource"
// @Success 200 {object} fhir.Patient
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patients/{id} [put]
func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	ctx, span := tracer.StartSpan(c.Request.Context(), "UpdatePatient")
	defer span.End()
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid patient ID",
			"message": "Patient ID is required",
		})
		return
	}
	// Parse the ID from the URL parameter
	id, err := strconv.ParseUint(idStr, 10, 64)

	// If parsing fails, return a 400 Bad Request
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid patient ID",
			"message": "Patient ID must be a valid number",
		})
		return
	}

	logger.WithContext(ctx).Infof("Updating patient with ID: %d", id)

	var fhirPatient fhir.Patient
	if err := c.ShouldBindJSON(&fhirPatient); err != nil {
		logger.WithContext(ctx).Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON",
			"message": err.Error(),
		})
		return
	}

	logger.WithContext(ctx).Infof("Updating patient with ID: %d using FHIR data", id)

	patient, err := h.service.UpdatePatient(ctx, uint(id), &fhirPatient)
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to update patient %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update patient",
			"message": err.Error(),
		})
		return
	}

	fhirResponse, err := h.service.ConvertToFHIR(ctx, patient)
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to convert to FHIR: %v", err)
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
// @Description Partially update an existing FHIR Patient resource
// @Tags Patient
// @Accept json
// @Produce json
// @Param id path int true "Patient ID"
// @Param patches body map[string]interface{} true "Partial updates"
// @Success 200 {object} fhir.Patient
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patients/{id} [patch]
func (h *PatientHandler) PatchPatient(c *gin.Context) {
	ctx, span := tracer.StartSpan(c.Request.Context(), "PatchPatient")
	defer span.End()

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
		logger.WithContext(ctx).Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON",
			"message": err.Error(),
		})
		return
	}

	patient, err := h.service.PatchPatient(ctx, uint(id), updates)
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to patch patient %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to patch patient",
			"message": err.Error(),
		})
		return
	}

	fhirResponse, err := h.service.ConvertToFHIR(ctx, patient)
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to convert to FHIR: %v", err)
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
// @Description Delete an existing FHIR Patient resource
// @Tags Patient
// @Produce json
// @Param id path int true "Patient ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patients/{id} [delete]
func (h *PatientHandler) DeletePatient(c *gin.Context) {
	ctx, span := tracer.StartSpan(c.Request.Context(), "DeletePatient")
	defer span.End()
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid patient ID",
			"message": "Patient ID must be a valid number",
		})
		return
	}
	logger.WithContext(ctx).Infof("Deleting patient with ID: %d", id)

	err = h.service.DeletePatient(ctx, uint(id))
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to delete patient %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete patient",
			"message": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
