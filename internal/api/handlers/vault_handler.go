package handlers

import (
	"net/http"

	"go-fhir-demo/config"
	"go-fhir-demo/pkg/utils"

	"github.com/gin-gonic/gin"
)

// VaultHandlerInterface defines the contract for Vault handler

// VaultHandler handles requests for Vault secrets.
type VaultHandler struct {
	cfg *config.VaultConfig
}

// NewVaultHandler creates a new VaultHandler.
func NewVaultHandler(cfg *config.VaultConfig) *VaultHandler {
	return &VaultHandler{cfg: cfg}
}

// GetVaultSecret godoc
// @Summary Get secret from Vault KV
// @Description Fetches a secret from Vault Key Vault and returns it as JSON
// @Tags Vault
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/vault/secret [get]
func (h *VaultHandler) GetVaultSecret(c *gin.Context) {
	data, err := utils.GetVaultKV(c.Request.Context(), h.cfg.Address, h.cfg.Token, h.cfg.SecretPath)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}
