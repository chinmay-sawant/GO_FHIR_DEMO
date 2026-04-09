package utils //nolint:revive

import (
	"context"
	"encoding/json"
	"fmt"
	"go-fhir-demo/pkg/vaultdto"
	"net/http"
	"time"
)

// GetVaultKV fetches a key-value secret from Vault.
func GetVaultKV(ctx context.Context, vaultAddr, token, secretPath string) (map[string]string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	reqURL := vaultAddr + "/v1/" + secretPath

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Vault-Token", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Vault: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vault returned status %d", resp.StatusCode)
	}

	var vaultResp vaultdto.Response
	if err := json.NewDecoder(resp.Body).Decode(&vaultResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return vaultResp.Data.Data, nil
}
