package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// GetConsulKV fetches a key from Consul's KV store and returns the decoded value.
func GetConsulKV(consulAddr, key string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/v1/kv/%s?raw", strings.TrimRight(consulAddr, "/"), key)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to contact consul: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("consul returned status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read consul response: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// If not JSON, return as string
		return map[string]interface{}{"value": string(body)}, nil
	}
	return result, nil
}
