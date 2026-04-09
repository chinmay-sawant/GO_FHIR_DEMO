package bootstrap

import (
	"context"
	"errors"

	"go-fhir-demo/internal/models"
	"gorm.io/gorm"
)

// Migrate applies the database schema for the application.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Patient{})
}

// SeedDummyPatients inserts sample patients when they are missing.
func SeedDummyPatients(ctx context.Context, db *gorm.DB) error {
	var patient models.Patient
	err := db.WithContext(ctx).Select("id").First(&patient).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	dummy := models.Patient{
		Family: "Doe",
		Given:  "John",
		Gender: "male",
	}
	return db.WithContext(ctx).Create(&dummy).Error
}
