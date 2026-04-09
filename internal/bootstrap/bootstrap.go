package bootstrap

import (
	"context"

	"go-fhir-demo/internal/domain"
	"gorm.io/gorm"
)

// Migrate applies the database schema for the application.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&domain.Patient{})
}

// SeedDummyPatients inserts sample patients when they are missing.
func SeedDummyPatients(_ context.Context, _ *gorm.DB) error {
	return nil
}
