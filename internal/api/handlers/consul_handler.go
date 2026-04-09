package handlers

import (
	"net/http"
	"strings"

	"go-fhir-demo/config"
	"go-fhir-demo/pkg/utils"
	"go-fhir-demo/pkg/utils/tracer"

	"github.com/gin-gonic/gin"
)

// ConsulHandlerInterface defines the contract for Consul handler

// ConsulHandler handles requests for Consul KV secrets.
type ConsulHandler struct {
	cfg *config.ConsulConfig
}

// NewConsulHandler creates a new ConsulHandler.
func NewConsulHandler(cfg *config.ConsulConfig) *ConsulHandler {
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

	secret, err := utils.GetConsulKV(ctx, h.cfg.Address, h.cfg.Key)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"secret": maskSecret(secret)})
}

func maskSecret(value string) string {
	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}
	return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}
