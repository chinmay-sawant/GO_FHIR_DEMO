package handlers

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"

	"go-fhir-demo/pkg/patch"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

func mustUintParam(value, name string) (uint64, error) {
	if value == "" {
		return 0, fmt.Errorf("%s is required", name)
	}
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s", name)
	}
	return id, nil
}

func mustStringParam(value, name string) (string, error) {
	if value == "" {
		return "", fmt.Errorf("%s is required", name)
	}
	return value, nil
}

func mustRawJSONBody(body io.Reader) (json.RawMessage, error) {
	var raw json.RawMessage
	if err := json.NewDecoder(body).Decode(&raw); err != nil {
		return nil, err
	}
	if len(raw) == 0 || !json.Valid(raw) {
		return nil, fmt.Errorf("invalid json")
	}
	return raw, nil
}

func parseFHIRPatientRequest(req *http.Request) (*fhir.Patient, error) {
	raw, err := mustRawJSONBody(req.Body)
	if err != nil {
		return nil, err
	}
	patient, err := fhirconvToPatient(raw)
	if err != nil {
		return nil, err
	}
	return patient, nil
}

func parsePatientUpdateRequest(req *http.Request, id string) (uint64, *fhir.Patient, error) {
	patientID, err := mustUintParam(id, "id")
	if err != nil {
		return 0, nil, err
	}
	patient, err := parseFHIRPatientRequest(req)
	if err != nil {
		return 0, nil, err
	}
	return patientID, patient, nil
}

func parsePatientPatchRequest(req *http.Request, id string) (uint64, patch.PatientPatch, error) {
	patientID, err := mustUintParam(id, "id")
	if err != nil {
		return 0, patch.PatientPatch{}, err
	}
	var updates patch.PatientPatch
	if err := json.NewDecoder(req.Body).Decode(&updates); err != nil {
		return 0, patch.PatientPatch{}, err
	}
	return patientID, updates, nil
}

func parseRawQueryParams(rawQuery string) map[string]string {
	params := make(map[string]string, strings.Count(rawQuery, "&")+1)
	if rawQuery == "" {
		return params
	}
	for _, part := range strings.Split(rawQuery, "&") {
		if part == "" {
			continue
		}
		key, value, found := strings.Cut(part, "=")
		if !found {
			continue
		}
		decodedKey, err := neturl.QueryUnescape(key)
		if err != nil {
			continue
		}
		decodedValue, err := neturl.QueryUnescape(value)
		if err != nil {
			continue
		}
		params[decodedKey] = decodedValue
	}
	return params
}

func fhirconvToPatient(raw json.RawMessage) (*fhir.Patient, error) {
	var patient fhir.Patient
	if err := json.Unmarshal(raw, &patient); err != nil {
		return nil, err
	}
	return &patient, nil
}

func requireCSRFToken(token string) bool {
	return subtle.ConstantTimeCompare([]byte(token), []byte("test-token")) == 1
}
