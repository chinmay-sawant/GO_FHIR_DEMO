// Package fhirconv provides FHIR-specific conversion helpers.
package fhirconv

import (
	"encoding/json"
	"time"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// GenderPtr returns an AdministrativeGender pointer from string.
func GenderPtr(g string) *fhir.AdministrativeGender {
	var gender fhir.AdministrativeGender
	if err := gender.UnmarshalJSON([]byte(`"` + g + `"`)); err != nil {
		gender = fhir.AdministrativeGenderUnknown
	}
	return &gender
}

// SystemPtr returns a ContactPointSystem pointer from string.
func SystemPtr(s string) *fhir.ContactPointSystem {
	var sys fhir.ContactPointSystem
	if err := sys.UnmarshalJSON([]byte(`"` + s + `"`)); err != nil {
		sys = fhir.ContactPointSystemOther
	}
	return &sys
}

// UsePtr returns a ContactPointUse pointer from string.
func UsePtr(u string) *fhir.ContactPointUse {
	var use fhir.ContactPointUse
	if err := use.UnmarshalJSON([]byte(`"` + u + `"`)); err != nil {
		use = fhir.ContactPointUseHome
	}
	return &use
}

// NameUseOfficialPtr returns a pointer to the official NameUse.
func NameUseOfficialPtr() *fhir.NameUse {
	nameUse := fhir.NameUseOfficial
	return &nameUse
}

// CreateTimePtr returns a pointer to a time value parsed from string.
func CreateTimePtr(t string) *time.Time {
	var ft time.Time
	if err := ft.UnmarshalJSON([]byte(`"` + t + `"`)); err != nil {
		return nil
	}
	return &ft
}

// ConvertJSONToFHIRPatient converts raw JSON to a FHIR Patient resource.
func ConvertJSONToFHIRPatient(data json.RawMessage) (*fhir.Patient, error) {
	var patient fhir.Patient
	if err := json.Unmarshal(data, &patient); err != nil {
		return nil, err
	}
	return &patient, nil
}
