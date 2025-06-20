package consul

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// RegisterWithConsul registers this service with Consul agent
func RegisterWithConsul(consulAddr, serviceName, serviceID, serviceHost, servicePort string) error {
	portInt, err := strconv.Atoi(servicePort)
	if err != nil {
		return fmt.Errorf("invalid service port: %w", err)
	}
	reg := map[string]interface{}{
		"Name":    serviceName,
		"ID":      serviceID,
		"Address": serviceHost,
		"Port":    portInt,
		"Check": map[string]interface{}{
			"HTTP":     fmt.Sprintf("http://%s:%d/health", serviceHost, portInt),
			"Interval": "10s",
			"Timeout":  "5s",
		},
	}
	body, _ := json.Marshal(reg)
	url := fmt.Sprintf("%s/v1/agent/service/register", consulAddr)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to register with consul: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("consul registration failed: %s", resp.Status)
	}
	return nil
}
