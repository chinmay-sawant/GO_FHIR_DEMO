package fhirclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// Client is a client for interacting with a FHIR server.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new FHIR client.
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetPatientByID fetches a Patient resource by its ID.
func (c *Client) GetPatientByID(id string) (*fhir.Patient, error) {
	reqURL := fmt.Sprintf("%s/Patient/%s", c.BaseURL, id)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/fhir+json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fhir server returned non-OK status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var patient fhir.Patient
	if err := json.NewDecoder(resp.Body).Decode(&patient); err != nil {
		return nil, fmt.Errorf("failed to decode patient response: %w", err)
	}

	return &patient, nil
}

// SearchPatients searches for Patient resources based on query parameters.
// It returns a FHIR Bundle containing the search results.
func (c *Client) SearchPatients(queryParams map[string]string) (*fhir.Bundle, error) {
	baseURL, err := url.Parse(fmt.Sprintf("%s/Patient", c.BaseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	params := url.Values{}
	for k, v := range queryParams {
		params.Add(k, v)
	}
	baseURL.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}
	req.Header.Set("Accept", "application/fhir+json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fhir server returned non-OK status for search %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var bundle fhir.Bundle
	if err := json.NewDecoder(resp.Body).Decode(&bundle); err != nil {
		return nil, fmt.Errorf("failed to decode bundle response: %w", err)
	}

	return &bundle, nil
}
