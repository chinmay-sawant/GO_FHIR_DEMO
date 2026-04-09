// Package patch provides request patch DTOs shared across layers.
package patch

// PatientPatch captures supported partial patient updates.
type PatientPatch struct {
	Active    *bool   `json:"active,omitempty"`
	Family    *string `json:"family,omitempty"`
	Given     *string `json:"given,omitempty"`
	Gender    *string `json:"gender,omitempty"`
	BirthDate *string `json:"birthDate,omitempty"`
}
