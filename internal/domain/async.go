package domain

// Async represents the data structure for async Kafka messages.
type Async struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}
