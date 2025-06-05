package routes

import (
	"go-fhir-demo/internal/api/handlers"
	"go-fhir-demo/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(patientHandler *handlers.PatientHandler) *gin.Engine {
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
