package models

// Shared application and FHIR catalog values.
const (
	HealthStatus         = "healthy"
	ServiceName          = "FHIR Patient API"
	ServiceVersion       = "1.0.0"
	CapabilityResource   = "CapabilityStatement"
	CapabilityStatus     = "active"
	CapabilityDate       = "2025-06-05"
	CapabilityKind       = "instance"
	CapabilityVersion    = "4.0.1"
	CapabilityFormat     = "json"
	CapabilityModeServer = "server"
	ResourceTypePatient  = "Patient"

	InteractionRead       = "read"
	InteractionCreate     = "create"
	InteractionUpdate     = "update"
	InteractionPatch      = "patch"
	InteractionDelete     = "delete"
	InteractionSearchType = "search-type"
)
