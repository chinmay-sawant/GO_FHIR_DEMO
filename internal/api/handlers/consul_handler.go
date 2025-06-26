package handlers

import (
	"net/http"

	"go-fhir-demo/config"
	"go-fhir-demo/pkg/logger"
	"go-fhir-demo/pkg/utils"
	"go-fhir-demo/pkg/utils/tracer"

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
// @Router /api/v1/consul/secret [get]
func (h *ConsulHandler) GetConsulSecret(c *gin.Context) {
	// Start a child span for the background job
	ctx, span := tracer.StartSpan(c.Request.Context(), "GetConsulSecret")
	defer span.End()

	logger.WithContext(ctx).Infof("Fetching secret from Consul KV at %s with key %s", h.cfg.Address, h.cfg.Key)
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
