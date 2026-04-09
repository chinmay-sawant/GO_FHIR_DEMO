// Package bootstrap wires the application dependencies.
package bootstrap

import (
	"go-fhir-demo/config"
	"go-fhir-demo/internal/api/handlers"
	"go-fhir-demo/internal/api/handlers/cron"
	"go-fhir-demo/internal/api/routes"
	"go-fhir-demo/internal/repository"
	"go-fhir-demo/internal/service"
	"go-fhir-demo/pkg/cache"
	"go-fhir-demo/pkg/fhirclient"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/gorm"
)

// BuildRouter wires the application dependencies and returns the configured router.
func BuildRouter(cfg *config.Config, db *gorm.DB) (*gin.Engine, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	patientRepo := repository.NewPatientRepository(sqlDB)
	patientService := service.NewPatientService(patientRepo)

	fhirClient := fhirclient.NewClient(cfg.Server.ExternalFHIRServerBaseURL)

	var cacheService cache.RedisCache
	if cfg.Redis.Host != "" {
		cacheService = cache.NewRedisCache(cache.Config{
			Host:     cfg.Redis.Host,
			Port:     cfg.Redis.Port,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
	}

	externalPatientService := service.NewExternalPatientService(fhirClient, cacheService)
	patientHandler := handlers.NewPatientHandler(patientService)
	externalPatientHandler := handlers.NewExternalPatientHandler(externalPatientService)
	cronJobHandler := cron.NewJobHandler()
	consulHandler := handlers.NewConsulHandler(&cfg.Consul)
	vaultHandler := handlers.NewVaultHandler(&cfg.Vault)
	asyncHandler := handlers.NewAsyncHandler(cfg.Kafka.Broker, cfg.Kafka.Topic)

	router := routes.SetupRoutes(patientHandler, externalPatientHandler, vaultHandler, consulHandler)
	if cfg.Jaeger.Enabled {
		router.Use(otelgin.Middleware(cfg.Jaeger.ServiceName))
	}
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	cronRoutes := api.Group("/cron")
	cronRoutes.POST("/cleanup", cronJobHandler.TriggerCleanupJob)
	cronRoutes.POST("/sync", cronJobHandler.TriggerDataSyncJob)
	api.POST("/async/publish", asyncHandler.PublishAsync)

	return router, nil
}
