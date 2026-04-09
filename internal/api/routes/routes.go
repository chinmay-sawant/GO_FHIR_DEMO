// Package routes provides routing configuration for the application.
package routes

import (
	"go-fhir-demo/internal/api/handlers"
	"go-fhir-demo/internal/models"
	"go-fhir-demo/internal/middleware"

	"github.com/gin-gonic/gin"
)

const (
	healthPath   = "/health"
	metadataPath = "/metadata"
	patientID    = "/:id"
)

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

type capabilityStatement struct {
	ResourceType string `json:"resourceType"`
	Status       string `json:"status"`
	Date         string `json:"date"`
	Publisher    string `json:"publisher"`
	Kind         string `json:"kind"`
	Software     struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"software"`
	FHIRVersion string   `json:"fhirVersion"`
	Format      []string `json:"format"`
	Rest        []struct {
		Mode     string `json:"mode"`
		Resource []struct {
			Type        string `json:"type"`
			Interaction []struct {
				Code string `json:"code"`
			} `json:"interaction"`
		} `json:"resource"`
	} `json:"rest"`
}

var metadataPayload = capabilityStatement{
	ResourceType: models.CapabilityResource,
	Status:       models.CapabilityStatus,
	Date:         models.CapabilityDate,
	Publisher:    "FHIR Demo",
	Kind:         models.CapabilityKind,
	FHIRVersion:  models.CapabilityVersion,
	Format:       []string{models.CapabilityFormat},
}

// RouteSetup implements *RouteSetup
type RouteSetup struct{}

// NewRouteSetup creates a new RouteSetup instance
func NewRouteSetup() *RouteSetup {
	return &RouteSetup{}
}

// SetupRoutes is a legacy function for backward compatibility.
func SetupRoutes(patientHandler *handlers.PatientHandler, externalPatientHandler *handlers.ExternalPatientHandler, vaultHandler *handlers.VaultHandler, consulHandler ...*handlers.ConsulHandler) *gin.Engine {
	routeSetup := NewRouteSetup()
	return routeSetup.SetupRoutes(patientHandler, externalPatientHandler, vaultHandler, consulHandler...)
}

// SetupRoutes configures all the routes for the application
func (r *RouteSetup) SetupRoutes(
	patientHandler *handlers.PatientHandler,
	externalPatientHandler *handlers.ExternalPatientHandler,
	vaultHandler *handlers.VaultHandler,
	consulHandler ...*handlers.ConsulHandler,
) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(middleware.RequestTracker())
	router.Use(middleware.RequestTimer())
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET(healthPath, healthHandler)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Patient routes
		patients := v1.Group("/patients")
		{
			patients.GET("", patientHandler.GetPatients)
			patients.POST("", middleware.CSRFProtection(), patientHandler.CreatePatient)
			patients.GET(patientID, patientHandler.GetPatient)
			patients.PUT(patientID, middleware.CSRFProtection(), patientHandler.UpdatePatient)
			patients.PATCH(patientID, middleware.CSRFProtection(), patientHandler.PatchPatient)
			patients.DELETE(patientID, middleware.CSRFProtection(), patientHandler.DeletePatient)
		}

		// External Patient routes
		externalPatients := v1.Group("/external-patients")
		{
			externalPatients.GET(patientID, externalPatientHandler.GetExternalPatientByID)
			externalPatients.GET(patientID+"/cached", externalPatientHandler.GetPatientCached)
			externalPatients.GET(patientID+"/delayed", externalPatientHandler.GetPatientDelayed)
			externalPatients.GET("", externalPatientHandler.SearchExternalPatients)
			externalPatients.POST("", middleware.CSRFProtection(), externalPatientHandler.CreateExternalPatient)
		}

		// Consul routes
		if len(consulHandler) > 0 && consulHandler[0] != nil {
			v1.GET("/consul/secret", consulHandler[0].GetConsulSecret)
		}

		// Vault routes
		v1.GET("/vault/secret", vaultHandler.GetVaultSecret)

	}

	router.GET(metadataPath, metadataHandler)

	return router
}

func healthHandler(c *gin.Context) {
	c.JSON(200, healthResponse{
		Status:  "active",
		Service: "go-fhir-demo",
		Version: "1.0",
	})
}

func metadataHandler(c *gin.Context) {
	payload := metadataPayload
	payload.Software.Name = "go-fhir-demo"
	payload.Software.Version = "1.0"

	c.JSON(200, payload)
}
