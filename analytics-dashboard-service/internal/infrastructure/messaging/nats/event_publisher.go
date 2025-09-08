package nats

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/infrastructure/config"
	"github.com/nats-io/nats.go"
)

// EventPublisher implements the EventBus interface using NATS
type EventPublisher struct {
	conn *nats.Conn
}

// NewEventPublisher creates a new NATS event publisher
func NewEventPublisher(cfg config.NATSConfig) (*EventPublisher, error) {
	opts := []nats.Option{
		nats.Name(cfg.ClientID),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS reconnected to %v", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Printf("NATS connection closed")
		}),
	}

	conn, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return &EventPublisher{conn: conn}, nil
}

// Publish publishes a domain event to NATS
func (p *EventPublisher) Publish(event *shared.DomainEvent) error {
	subject := fmt.Sprintf("analytics.events.%s", event.Type)

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if err := p.conn.Publish(subject, data); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Published event %s to subject %s", event.Type, subject)
	return nil
}

// Subscribe subscribes to domain events
func (p *EventPublisher) Subscribe(eventType string, handler func(*shared.DomainEvent) error) error {
	subject := fmt.Sprintf("analytics.events.%s", eventType)

	_, err := p.conn.Subscribe(subject, func(m *nats.Msg) {
		var event shared.DomainEvent
		if err := json.Unmarshal(m.Data, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			return
		}

		if err := handler(&event); err != nil {
			log.Printf("Failed to handle event: %v", err)
		}
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", subject, err)
	}

	log.Printf("Subscribed to events on subject %s", subject)
	return nil
}

// Close closes the NATS connection
func (p *EventPublisher) Close() error {
	if p.conn != nil {
		p.conn.Close()
	}
	return nil
}

// NoOpEventPublisher is a no-op implementation of EventBus for testing
type NoOpEventPublisher struct{}

// Publish does nothing
func (p *NoOpEventPublisher) Publish(event *shared.DomainEvent) error {
	return nil
}

// Subscribe does nothing
func (p *NoOpEventPublisher) Subscribe(eventType string, handler func(*shared.DomainEvent) error) error {
	return nil
}
