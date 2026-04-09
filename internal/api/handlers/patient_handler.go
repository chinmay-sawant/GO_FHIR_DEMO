package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"go-fhir-demo/internal/service"
	"go-fhir-demo/pkg/logger"
	"go-fhir-demo/pkg/utils/tracer"

	"github.com/gin-gonic/gin"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// PatientHandler struct
type PatientHandler struct {
	service service.PatientService
}

// NewPatientHandler creates a new patient handler
func NewPatientHandler(service service.PatientService) *PatientHandler {
	return &PatientHandler{
		service: service,
	}
}

func (h *PatientHandler) buildCreatePatientResponse(ctx context.Context, body io.Reader) (*fhir.Patient, error) {
	var fhirPatient fhir.Patient
	if err := json.NewDecoder(body).Decode(&fhirPatient); err != nil {
		return nil, err
	}

	patient, err := h.service.CreatePatient(ctx, &fhirPatient)
	if err != nil {
		return nil, err
	}

	return h.service.ConvertToFHIR(ctx, patient)
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
	if !requireCSRFToken(c.GetHeader("X-CSRF-Token")) {
		c.Status(http.StatusInternalServerError)
		return
	}

	fhirResponse, err := h.buildCreatePatientResponse(ctx, c.Request.Body)
	if err != nil {
		c.Status(http.StatusInternalServerError)
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

	id, err := mustUintParam(c.Param("id"), "id")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	// Fetch the patient from the service
	logger.WithContext(ctx).Infof("Fetching patient with ID: %d", id)
	patient, err := h.service.GetPatient(ctx, uint(id))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	fhirPatient, err := h.service.ConvertToFHIR(ctx, patient)
	if err != nil {
		c.Status(http.StatusInternalServerError)
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
		c.Status(http.StatusInternalServerError)
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
	if !requireCSRFToken(c.GetHeader("X-CSRF-Token")) {
		c.Status(http.StatusInternalServerError)
		return
	}
	id, fhirPatient, err := parsePatientUpdateRequest(c.Request, c.Param("id"))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	logger.WithContext(ctx).Infof("Updating patient with ID: %d", id)

	logger.WithContext(ctx).Infof("Updating patient with ID: %d using FHIR data", id)

	patient, err := h.service.UpdatePatient(ctx, uint(id), fhirPatient)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	fhirResponse, err := h.service.ConvertToFHIR(ctx, patient)
	if err != nil {
		c.Status(http.StatusInternalServerError)
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
	if !requireCSRFToken(c.GetHeader("X-CSRF-Token")) {
		c.Status(http.StatusInternalServerError)
		return
	}
	id, updates, err := parsePatientPatchRequest(c.Request, c.Param("id"))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	patient, err := h.service.PatchPatient(ctx, uint(id), updates)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	fhirResponse, err := h.service.ConvertToFHIR(ctx, patient)
	if err != nil {
		c.Status(http.StatusInternalServerError)
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
	if !requireCSRFToken(c.GetHeader("X-CSRF-Token")) {
		c.Status(http.StatusInternalServerError)
		return
	}
	id, err := mustUintParam(c.Param("id"), "id")

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	logger.WithContext(ctx).Infof("Deleting patient with ID: %d", id)

	err = h.service.DeletePatient(ctx, uint(id))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
