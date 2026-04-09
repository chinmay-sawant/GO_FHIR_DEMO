// Package models contains the domain models for the FHIR demo.
package models

// Async represents the data structure for async Kafka messages.
type Async struct {
	ID   string
	Data string
}
