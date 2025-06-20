package handlers

import (
	"net/http"

	"go-fhir-demo/config"
	"go-fhir-demo/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ConsulHandlerInterface defines the contract for Consul handler
type ConsulHandlerInterface interface {
	GetConsulSecret(c *gin.Context)
}

type ConsulHandler struct {
	cfg *config.ConsulConfig
}

func NewConsulHandler(cfg *config.ConsulConfig) ConsulHandlerInterface {
	return &ConsulHandler{cfg: cfg}
}

// GetConsulSecret godoc
// @Summary Get secret from Consul KV
// @Description Fetches a secret from Consul Key Vault and returns it as JSON
// @Tags Consul
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /consul/secret [get]
func (h *ConsulHandler) GetConsulSecret(c *gin.Context) {
	data, err := utils.GetConsulKV(h.cfg.Address, h.cfg.Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch from Consul",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, data)
}
