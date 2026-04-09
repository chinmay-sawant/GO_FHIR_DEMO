// Package asyncdto defines payloads for asynchronous messages.
package asyncdto

// Message is the transport payload used for async Kafka messages.
type Message struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}
