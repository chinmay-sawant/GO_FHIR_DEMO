package handlers

import (
	"net/http"

	"encoding/json"
	"go-fhir-demo/internal/domain"
	"go-fhir-demo/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

type AsyncHandlerInterface interface {
	PublishAsync(c *gin.Context)
}

type AsyncHandler struct {
	KafkaWriter *kafka.Writer
	Topic       string
}

func NewAsyncHandler(broker, topic string) AsyncHandlerInterface {
	return &AsyncHandler{
		KafkaWriter: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
		Topic: topic,
	}
}

// PublishAsync godoc
// @Summary Publish async data to Kafka
// @Description Publishes async data to Kafka topic
// @Tags Async
// @Accept json
// @Produce json
// @Param async body domain.Async true "Async data"
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /async/publish [post]
func (h *AsyncHandler) PublishAsync(c *gin.Context) {
	var asyncData domain.Async
	if err := c.ShouldBindJSON(&asyncData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "message": err.Error()})
		return
	}
	value, err := json.Marshal(asyncData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}
	msg := kafka.Message{
		Value: value,
	}
	if err := h.KafkaWriter.WriteMessages(c, msg); err != nil {
		logger.Error("Failed to publish to Kafka: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish to Kafka"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "published"})
}
