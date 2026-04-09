package domain

import (
	"time"
)

// Patient represents a FHIR Patient resource in the domain
type Patient struct {
	ID        uint
	FHIRData  []byte
	Active    *bool
	Family    string
	Given     string
	Gender    string
	BirthDate *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
