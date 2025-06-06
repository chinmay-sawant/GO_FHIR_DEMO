package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-fhir-demo/config"
	"go-fhir-demo/internal/api/handlers"
	"go-fhir-demo/internal/api/routes"
	"go-fhir-demo/internal/domain"
	"go-fhir-demo/internal/repository"
	"go-fhir-demo/internal/service"
	"go-fhir-demo/pkg/database"
	"go-fhir-demo/pkg/fhirclient" // Import the new fhirclient package
	"go-fhir-demo/pkg/logger"
	"go-fhir-demo/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"

	// Swagger imports
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// docs is generated by Swag CLI, you have to import it.
	_ "go-fhir-demo/docs"
)

// @title Go FHIR Demo API
// @version 1.0
// @description This is a sample FHIR Patient API server in Go using Gin.
// @BasePath /api/v1

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Initialize(cfg.Logging.Level, cfg.Logging.Format, cfg.Logging.File); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	logger.Info("Starting FHIR Patient API server...")

	// Initialize database
	if err := database.Initialize(&cfg.Database); err != nil {
		logger.Errorf("Failed to initialize database: %v", err)
		os.Exit(1)
	}
	defer database.Close()

	// Auto-migrate the database schema
	db := database.GetDB()
	if err := db.AutoMigrate(&domain.Patient{}); err != nil {
		logger.Errorf("Failed to migrate database: %v", err)
		os.Exit(1)
	}

	// Initialize repositories
	patientRepo := repository.NewPatientRepository(db)

	// Initialize services
	patientService := service.NewPatientService(patientRepo)

	// Initialize FHIR client
	fhirClient := fhirclient.NewClient(cfg.Server.ExternalFHIRServerBaseURL)

	// Initialize external patient service
	externalPatientService := service.NewExternalPatientService(fhirClient)

	// --- Seed dummy patients if not present ---
	dummyPatients := []fhir.Patient{
		{
			Active: utils.CreateBoolPtr(true),
			Name: []fhir.HumanName{
				{
					Use:    utils.NameUseOfficialPtr(),
					Family: utils.CreateStringPtr("Doe"),
					Given:  []string{"John"},
				},
			},
			Gender:    utils.GenderPtr("male"),
			BirthDate: utils.CreateStringPtr("1980-01-01"),
			Telecom: []fhir.ContactPoint{
				{
					System: utils.SystemPtr("phone"),
					Value:  utils.CreateStringPtr("1234567890"),
					Use:    utils.UsePtr("mobile"),
				},
			},
			Address: []fhir.Address{
				{
					Line:       []string{"123 Main St"},
					City:       utils.CreateStringPtr("Metropolis"),
					State:      utils.CreateStringPtr("NY"),
					PostalCode: utils.CreateStringPtr("12345"),
					Country:    utils.CreateStringPtr("USA"),
				},
			},
		},
		{
			Active: utils.CreateBoolPtr(true),
			Name: []fhir.HumanName{
				{
					Use:    utils.NameUseOfficialPtr(),
					Family: utils.CreateStringPtr("Smith"),
					Given:  []string{"Jane"},
				},
			},
			Gender:    utils.GenderPtr("female"),
			BirthDate: utils.CreateStringPtr("1990-05-15"),
			Telecom: []fhir.ContactPoint{
				{
					System: utils.SystemPtr("email"),
					Value:  utils.CreateStringPtr("jane.smith@example.com"),
					Use:    utils.UsePtr("home"),
				},
			},
			Address: []fhir.Address{
				{
					Line:       []string{"456 Oak Ave"},
					City:       utils.CreateStringPtr("Gotham"),
					State:      utils.CreateStringPtr("CA"),
					PostalCode: utils.CreateStringPtr("67890"),
					Country:    utils.CreateStringPtr("USA"),
				},
			},
		},
		{
			Active: utils.CreateBoolPtr(false),
			Name: []fhir.HumanName{
				{
					Use:    utils.NameUseOfficialPtr(),
					Family: utils.CreateStringPtr("Brown"),
					Given:  []string{"Charlie"},
				},
			},
			Gender:    utils.GenderPtr("other"),
			BirthDate: utils.CreateStringPtr("2000-12-31"),
			Telecom: []fhir.ContactPoint{
				{
					System: utils.SystemPtr("email"),
					Value:  utils.CreateStringPtr("charlie.brown@example.com"),
					Use:    utils.UsePtr("work"),
				},
			},
			Address: []fhir.Address{
				{
					Line:       []string{"789 Pine Rd"},
					City:       utils.CreateStringPtr("Star City"),
					State:      utils.CreateStringPtr("WA"),
					PostalCode: utils.CreateStringPtr("24680"),
					Country:    utils.CreateStringPtr("USA"),
				},
			},
		},
	}

	for _, dummy := range dummyPatients {
		// Check if patient exists by unique fields (family, given, birthdate, gender)
		var count int64
		db.Model(&domain.Patient{}).
			Where("family = ? AND given = ? AND gender = ? AND birth_date = ?",
				*dummy.Name[0].Family,
				func() string {
					if len(dummy.Name[0].Given) > 0 {
						return dummy.Name[0].Given[0]
					} else {
						return ""
					}
				}(),
				dummy.Gender,
				*dummy.BirthDate,
			).Count(&count)
		if count == 0 {
			_, err := patientService.CreatePatient(&dummy)
			if err != nil {
				logger.Warnf("Failed to seed dummy patient: %v", err)
			} else {
				logger.Infof("Seeded dummy patient: %s %s", dummy.Name[0].Given[0], *dummy.Name[0].Family)
			}
		}
	}

	// --- End seed logic ---

	// Initialize handlers
	patientHandler := handlers.NewPatientHandler(patientService)
	externalPatientHandler := handlers.NewExternalPatientHandler(externalPatientService) // Initialize new handler

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Setup routes
	router := routes.SetupRoutes(patientHandler, externalPatientHandler) // Pass new handler to SetupRoutes

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Configure server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Gracefully shutdown the server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	logger.Info("Server exited")
}
