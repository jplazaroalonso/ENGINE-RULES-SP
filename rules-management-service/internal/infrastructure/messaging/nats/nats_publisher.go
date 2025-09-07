package nats

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"

	"rules-management-service/internal/domain/shared"
	"rules-management-service/internal/infrastructure/config"
)

// EventPublisher is a NATS-based event publisher.
type EventPublisher struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

// NewEventPublisher creates a new NATS event publisher.
func NewEventPublisher(cfg config.NATSConfig) (*EventPublisher, error) {
	conn, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to get JetStream context: %w", err)
	}

	// This is where you would configure your streams, e.g.,
	// js.AddStream(&nats.StreamConfig{...})
	// For simplicity, we assume streams are configured externally.

	return &EventPublisher{conn: conn, js: js}, nil
}

// Publish publishes a domain event to a NATS subject.
func (p *EventPublisher) Publish(event shared.DomainEvent) error {
	subject := fmt.Sprintf("rules.%s", event.EventType())
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	log.Printf("Publishing event to subject %s", subject)
	_, err = p.js.Publish(subject, eventData)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

// Close closes the NATS connection.
func (p *EventPublisher) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}
