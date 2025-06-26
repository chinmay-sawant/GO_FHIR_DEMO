package handlers

import (
	"net/http"

	"go-fhir-demo/config"
	"go-fhir-demo/pkg/utils"

	"github.com/gin-gonic/gin"
)

// VaultHandlerInterface defines the contract for Vault handler
type VaultHandlerInterface interface {
	GetVaultSecret(c *gin.Context)
}

type VaultHandler struct {
	cfg *config.VaultConfig
}

func NewVaultHandler(cfg *config.VaultConfig) VaultHandlerInterface {
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
	data, err := utils.GetVaultKV(h.cfg.Address, h.cfg.Token, h.cfg.SecretPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch from Vault",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, data)
}
