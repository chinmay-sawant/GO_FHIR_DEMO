// Package handlers provides Gin handlers for the application.
package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"go-fhir-demo/pkg/asyncdto"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

// AsyncHandler handles requests for asynchronous Kafka publishing.
type AsyncHandler struct {
	KafkaWriter *kafka.Writer
	Topic       string
}

func (h *AsyncHandler) publishAsyncMessage(ctx context.Context, asyncData asyncdto.Message) error {
	value, err := json.Marshal(asyncData)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Value: value,
	}

	return h.KafkaWriter.WriteMessages(ctx, msg)
}

// NewAsyncHandler creates a new AsyncHandler.
func NewAsyncHandler(broker, topic string) *AsyncHandler {
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
// @Param async body asyncdto.Message true "Async data"
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /async/publish [post]
func (h *AsyncHandler) PublishAsync(c *gin.Context) {
	ctx := context.Background()
	var asyncData asyncdto.Message
	if err := c.ShouldBindJSON(&asyncData); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if err := h.publishAsyncMessage(ctx, asyncData); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusAccepted, StatusResponse{Status: "published"})
}

// StatusResponse represents an asynchronous operation status.
type StatusResponse struct {
	Status string `json:"status"`
}
