// Package consumer provides Kafka message consumers.
package consumer

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"go-fhir-demo/pkg/asyncdto"
	"go-fhir-demo/pkg/logger"

	"github.com/segmentio/kafka-go"
)

// StartAsyncConsumer starts a Kafka consumer in a background goroutine.
// It accepts a context for graceful shutdown and a WaitGroup for coordination.
func StartAsyncConsumer(ctx context.Context, wg *sync.WaitGroup, broker, topic, groupID string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if err := r.Close(); err != nil {
				logger.GetLogger().Errorf("Failed to close kafka reader: %v", err)
			}
		}()

		logger.GetLogger().Info("Async Kafka consumer started")

		retryTimer := time.NewTimer(0)
		<-retryTimer.C // drain immediately

		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				// Check if the error is due to context cancellation
				if ctx.Err() != nil {
					logger.GetLogger().Info("Async Kafka consumer shutting down")
					return
				}
				logger.GetLogger().Errorf("Failed to read message from Kafka: %v", err)

				retryTimer.Reset(2 * time.Second)
				select {
				case <-retryTimer.C:
					continue
				case <-ctx.Done():
					retryTimer.Stop()
					return
				}
			}

			var asyncData asyncdto.Message
			if err := json.Unmarshal(m.Value, &asyncData); err != nil {
				logger.GetLogger().Error("Failed to unmarshal async message: ", err)
				continue
			}
			logger.GetLogger().Infof("Consumed async message from Kafka: ID=%s, Data=%s", asyncData.ID, asyncData.Data)
		}
	}()
}
