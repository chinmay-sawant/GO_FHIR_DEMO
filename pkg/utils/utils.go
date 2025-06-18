package utils

import (
	"time"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// PtrBool returns a pointer to the given bool value.
func PtrBool(b bool) *bool {
	return &b
}

// Helper function to create bool pointers
func CreateBoolPtr(b bool) *bool {
	return &b
}

// Helper function to create gender pointer using UnmarshalJSON
func GenderPtr(g string) *fhir.AdministrativeGender {
	var gender fhir.AdministrativeGender
	if err := gender.UnmarshalJSON([]byte(`"` + g + `"`)); err != nil {
		// If unknown or error, default to AdministrativeGenderUnknown
		gender = fhir.AdministrativeGenderUnknown
	}
	return &gender
}

// Helper function to create system pointer
func SystemPtr(s string) *fhir.ContactPointSystem {
	var sys fhir.ContactPointSystem
	// Use UnmarshalJSON to leverage the built-in conversion logic
	if err := sys.UnmarshalJSON([]byte(`"` + s + `"`)); err != nil {
		// If unknown, default to ContactPointSystemOther
		sys = fhir.ContactPointSystemOther
	}
	return &sys
}

// Helper function to create use pointer
func UsePtr(u string) *fhir.ContactPointUse {
	var use fhir.ContactPointUse
	if err := use.UnmarshalJSON([]byte(`"` + u + `"`)); err != nil {
		// If unknown, default to ContactPointUseHome or another sensible default
		use = fhir.ContactPointUseHome
	}
	return &use
}

// Helper function to create string pointers
func CreateStringPtr(s string) *string {
	return &s
}

// Helper function to create NameUse pointer
func NameUseOfficialPtr() *fhir.NameUse {
	nameUse := fhir.NameUseOfficial
	return &nameUse
}

// Helper function to create time pointer
func CreateTimePtr(t string) *time.Time {
	var ft time.Time
	if err := ft.UnmarshalJSON([]byte(`"` + t + `"`)); err != nil {
		return nil
	}
	return &ft
}

// Helper function to create date pointer
