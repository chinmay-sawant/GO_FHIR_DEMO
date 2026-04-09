// Package main is the entry point for the FHIR Patient API server.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-fhir-demo/config"
	"go-fhir-demo/internal/bootstrap"
	"go-fhir-demo/pkg/database"
	"go-fhir-demo/pkg/utils/consul"
	"go-fhir-demo/pkg/utils/tracer"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

import _ "go-fhir-demo/docs"

// @title Go FHIR Demo API
// @version 1.0
// @description This is a sample FHIR Patient API server in Go using Gin.
// @BasePath /api/v1

const localhost = "localhost"

func main() {
	if err := run(); err != nil {
		log.Printf("Fatal error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx := context.Background()

	tracerProvider, err := tracer.InitJaeger(tracer.Config{
		Endpoint:    cfg.Jaeger.Endpoint,
		ServiceName: cfg.Jaeger.ServiceName,
		Environment: cfg.Jaeger.Environment,
		Enabled:     cfg.Jaeger.Enabled,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize Jaeger: %w", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tracerProvider.Shutdown(shutdownCtx); err != nil {
			log.Printf("Failed to shutdown tracer: %v", err)
		}
	}()

	db, err := database.Initialize(&cfg.DB)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	if err := bootstrap.Migrate(db); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	if err := bootstrap.SeedDummyPatients(ctx, db); err != nil {
		log.Printf("Failed to seed dummy patients: %v", err)
	}

	gin.SetMode(cfg.Server.Mode)
	router, err := bootstrap.BuildRouter(cfg, db)
	if err != nil {
		return fmt.Errorf("failed to build router: %w", err)
	}



	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	appName := "go-fhir-demo"
	appID := appName + "-" + cfg.Server.Port
	var appHost, checkHost string
	if cfg.Server.DevMode {
		appHost = localhost
		checkHost = "host.docker.internal"
	} else {
		appHost = getLocalIP()
		checkHost = appHost
	}
	if err := consul.RegisterWithConsul(ctx, cfg.Consul.Address, appName, appID, appHost, cfg.Server.Port, checkHost); err != nil {
		return fmt.Errorf("consul registration failed: %w", err)
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("failed to start server: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-quit:
			log.Printf("Received signal: %v", sig)
		case <-gCtx.Done():
			return gCtx.Err()
		}

		log.Println("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil && err != context.Canceled {
		return err
	}

	log.Printf("Server exited")
	return nil
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return localhost
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ipnet.IP = ipnet.IP.To4()
			if ipnet.IP[0] == 127 && ipnet.IP[1] == 0 && ipnet.IP[2] == 0 && ipnet.IP[3] == 1 {
				return localhost
			}
			return ipnet.IP.String()
		}
	}
	return localhost
}
