package routes

import (
	"go-fhir-demo/internal/api/handlers"
	"go-fhir-demo/internal/api/handlers/cron"
	"go-fhir-demo/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RouteSetupInterface defines the contract for route setup
type RouteSetupInterface interface {
	SetupRoutes(patientHandler handlers.PatientHandlerInterface, externalPatientHandler handlers.ExternalPatientHandlerInterface, cronJobHandler cron.CronJobHandlerInterface, consulHandler ...handlers.ConsulHandlerInterface) *gin.Engine
}

// RouteSetup implements RouteSetupInterface
type RouteSetup struct{}

// NewRouteSetup creates a new RouteSetup instance
func NewRouteSetup() RouteSetupInterface {
	return &RouteSetup{}
}

// Legacy function for backward compatibility
func SetupRoutes(patientHandler handlers.PatientHandlerInterface, externalPatientHandler handlers.ExternalPatientHandlerInterface, cronJobHandler cron.CronJobHandlerInterface, consulHandler ...handlers.ConsulHandlerInterface) *gin.Engine {
	routeSetup := NewRouteSetup()
	return routeSetup.SetupRoutes(patientHandler, externalPatientHandler, cronJobHandler, consulHandler...)
}

// SetupRoutes configures all the routes for the application
func (r *RouteSetup) SetupRoutes(
	patientHandler handlers.PatientHandlerInterface,
	externalPatientHandler handlers.ExternalPatientHandlerInterface,
	cronJobHandler cron.CronJobHandlerInterface,
	// Add optional handlers
	consulHandler ...handlers.ConsulHandlerInterface,
) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(middleware.RequestTracker())
	router.Use(middleware.RequestTimer())
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "FHIR Patient API",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Patient routes
		patients := v1.Group("/patients")
		{
			patients.GET("", patientHandler.GetPatients)
			patients.POST("", patientHandler.CreatePatient)
			patients.GET("/:id", patientHandler.GetPatient)
			patients.PUT("/:id", patientHandler.UpdatePatient)
			patients.PATCH("/:id", patientHandler.PatchPatient)
			patients.DELETE("/:id", patientHandler.DeletePatient)
		}

		// External Patient routes
		externalPatients := v1.Group("/external-patients")
		{
			externalPatients.GET("/:id", externalPatientHandler.GetExternalPatientByID)
			externalPatients.GET("/:id/cached", externalPatientHandler.GetExternalPatientByIDCached)
			externalPatients.GET("/:id/delayed", externalPatientHandler.GetExternalPatientByIDDelayed)
			externalPatients.GET("", externalPatientHandler.SearchExternalPatients)
			externalPatients.POST("", externalPatientHandler.CreateExternalPatient)
		}

		// Cron job routes
		if cronJobHandler != nil {
			cronJobs := v1.Group("/cron")
			{
				cronJobs.POST("/cleanup", cronJobHandler.TriggerCleanupJob)
				cronJobs.POST("/sync", cronJobHandler.TriggerDataSyncJob)
			}
		}
		// Consul secret endpoint
		if len(consulHandler) > 0 && consulHandler[0] != nil {
			v1.GET("/consul/secret", consulHandler[0].GetConsulSecret)
		}
	}

	// FHIR metadata endpoint
	router.GET("/metadata", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"resourceType": "CapabilityStatement",
			"status":       "active",
			"date":         "2025-06-05",
			"publisher":    "FHIR Demo",
			"kind":         "instance",
			"software": gin.H{
				"name":    "FHIR Patient API",
				"version": "1.0.0",
			},
			"fhirVersion": "4.0.1",
			"format":      []string{"json"},
			"rest": []gin.H{
				{
					"mode": "server",
					"resource": []gin.H{
						{
							"type": "Patient",
							"interaction": []gin.H{
								{"code": "read"},
								{"code": "create"},
								{"code": "update"},
								{"code": "patch"},
								{"code": "delete"},
								{"code": "search-type"},
							},
						},
					},
				},
			},
		})
	})

	return router
}
