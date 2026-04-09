// Package vaultdto provides Vault response DTOs.
package vaultdto

// Response represents the response from Vault API.
type Response struct {
	Data struct {
		Data map[string]string `json:"data"`
	} `json:"data"`
}
