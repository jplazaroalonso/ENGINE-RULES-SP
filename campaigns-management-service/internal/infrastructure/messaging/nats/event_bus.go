package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
	"github.com/nats-io/nats.go"
)

// NATSEventBus implements the shared.EventBus interface using NATS
type NATSEventBus struct {
	conn *nats.Conn
}

// NewNATSEventBus creates a new NATS event bus
func NewNATSEventBus(conn *nats.Conn) *NATSEventBus {
	return &NATSEventBus{conn: conn}
}

// Publish publishes a domain event to NATS
func (e *NATSEventBus) Publish(ctx context.Context, event shared.DomainEvent) error {
	// Create event envelope
	envelope := EventEnvelope{
		EventID:       event.EventID(),
		AggregateID:   event.AggregateID(),
		AggregateType: "Campaign", // Default aggregate type for campaign events
		EventType:     event.EventType(),
		OccurredAt:    event.OccurredAt(),
		EventData:     event.EventData(),
		PublishedAt:   time.Now(),
	}

	// Marshal event to JSON
	data, err := json.Marshal(envelope)
	if err != nil {
		return shared.NewInfrastructureError("failed to marshal event", err)
	}

	// Create subject based on event type
	subject := fmt.Sprintf("campaigns.events.%s", event.EventType())

	// Publish with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return shared.NewInfrastructureError("publish timeout", ctx.Err())
	default:
		err := e.conn.Publish(subject, data)
		if err != nil {
			return shared.NewInfrastructureError("failed to publish event", err)
		}
	}

	log.Printf("Published event %s for aggregate %s to subject %s", event.EventType(), event.AggregateID(), subject)
	return nil
}

// Subscribe subscribes to domain events from NATS
func (e *NATSEventBus) Subscribe(eventType string, handler shared.EventHandler) error {
	subject := fmt.Sprintf("campaign.events.%s", eventType)

	_, err := e.conn.Subscribe(subject, func(msg *nats.Msg) {
		// Parse event envelope
		var envelope EventEnvelope
		if err := json.Unmarshal(msg.Data, &envelope); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			return
		}

		// Create domain event from envelope
		eventData, ok := envelope.EventData.(map[string]interface{})
		if !ok {
			log.Printf("Failed to convert event data to map[string]interface{}")
			return
		}
		event := shared.NewBaseDomainEvent(envelope.EventType, envelope.AggregateID, eventData)

		// Handle event
		if err := handler.Handle(context.Background(), &event); err != nil {
			log.Printf("Failed to handle event %s: %v", eventType, err)
		}
	})

	if err != nil {
		return shared.NewInfrastructureError("failed to subscribe to events", err)
	}

	log.Printf("Subscribed to events of type %s on subject %s", eventType, subject)
	return nil
}

// EventEnvelope wraps domain events for NATS transport
type EventEnvelope struct {
	EventID       string      `json:"eventId"`
	AggregateID   string      `json:"aggregateId"`
	AggregateType string      `json:"aggregateType"`
	EventType     string      `json:"eventType"`
	OccurredAt    time.Time   `json:"occurredAt"`
	EventData     interface{} `json:"eventData"`
	PublishedAt   time.Time   `json:"publishedAt"`
}
