// Package domain contains the domain models for the FHIR demo.
package domain

// Async represents the data structure for async Kafka messages.
type Async struct {
	ID   string
	Data string
}
