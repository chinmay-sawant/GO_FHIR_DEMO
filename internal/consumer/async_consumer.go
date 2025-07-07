package consumer

import (
	"context"
	"encoding/json"
	"time"

	"go-fhir-demo/internal/domain"
	"go-fhir-demo/pkg/logger"

	"github.com/segmentio/kafka-go"
)

func StartAsyncConsumer(broker, topic, groupID string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			m, err := r.ReadMessage(ctx)
			cancel()
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}
			var asyncData domain.Async
			if err := json.Unmarshal(m.Value, &asyncData); err != nil {
				logger.Error("Failed to unmarshal async message: ", err)
				continue
			}
			logger.Infof("Consumed async message from Kafka: ID=%s, Data=%s", asyncData.ID, asyncData.Data)
		}
	}()
}
