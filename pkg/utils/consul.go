// Package utils provides miscellaneous utilities for the FHIR demo.
package utils //nolint:revive

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// GetConsulKV fetches a key from Consul's KV store and returns the raw secret value.
func GetConsulKV(ctx context.Context, consulAddr, key string) (string, error) {
	addr := strings.TrimRight(consulAddr, "/")
	reqURL := addr + "/v1/kv/" + key + "?raw"
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create search request: %w", err)
	}

	resp, err := SharedClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to contact consul: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("consul returned status %d", resp.StatusCode)
	}

	var builder strings.Builder
	if _, err := io.Copy(&builder, resp.Body); err != nil {
		return "", fmt.Errorf("failed to read consul response: %w", err)
	}
	return builder.String(), nil
}
