package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/domain/shared"
	"github.com/nats-io/nats.go"
)

// EventBus implements the EventBus interface using NATS JetStream
type EventBus struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

// NewEventBus creates a new NATS event bus
func NewEventBus(conn *nats.Conn) (*EventBus, error) {
	js, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	return &EventBus{
		conn: conn,
		js:   js,
	}, nil
}

// Publish publishes a single domain event
func (eb *EventBus) Publish(ctx context.Context, event shared.DomainEvent) error {
	subject := fmt.Sprintf("events.%s.%s", event.AggregateType(), event.EventType())

	eventData := map[string]interface{}{
		"eventId":       event.EventID(),
		"eventType":     event.EventType(),
		"aggregateId":   event.AggregateID(),
		"aggregateType": event.AggregateType(),
		"occurredAt":    event.OccurredAt(),
		"eventData":     event.EventData(),
		"version":       event.Version(),
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	_, err = eb.js.PublishAsync(subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

// PublishBatch publishes multiple domain events
func (eb *EventBus) PublishBatch(ctx context.Context, events []shared.DomainEvent) error {
	for _, event := range events {
		if err := eb.Publish(ctx, event); err != nil {
			return fmt.Errorf("failed to publish event %s: %w", event.EventID(), err)
		}
	}

	return nil
}

// Subscribe subscribes to domain events
func (eb *EventBus) Subscribe(ctx context.Context, eventType string, handler shared.EventHandler) error {
	subject := fmt.Sprintf("events.*.%s", eventType)

	_, err := eb.js.Subscribe(subject, func(msg *nats.Msg) {
		var eventData map[string]interface{}
		if err := json.Unmarshal(msg.Data, &eventData); err != nil {
			// Log error and continue
			return
		}

		// Create a domain event from the message
		event := &DomainEventWrapper{
			eventID:       getString(eventData, "eventId"),
			eventType:     getString(eventData, "eventType"),
			aggregateID:   getString(eventData, "aggregateId"),
			aggregateType: getString(eventData, "aggregateType"),
			occurredAt:    getTime(eventData, "occurredAt"),
			eventData:     getMap(eventData, "eventData"),
			version:       getInt(eventData, "version"),
		}

		// Handle the event
		if err := handler.Handle(ctx, event); err != nil {
			// Log error and continue
			return
		}

		// Acknowledge the message
		msg.Ack()
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to events: %w", err)
	}

	return nil
}

// DomainEventWrapper wraps event data for handling
type DomainEventWrapper struct {
	eventID       string
	eventType     string
	aggregateID   string
	aggregateType string
	occurredAt    time.Time
	eventData     map[string]interface{}
	version       int
}

// EventID returns the event ID
func (e *DomainEventWrapper) EventID() string {
	return e.eventID
}

// EventType returns the event type
func (e *DomainEventWrapper) EventType() string {
	return e.eventType
}

// AggregateID returns the aggregate ID
func (e *DomainEventWrapper) AggregateID() string {
	return e.aggregateID
}

// AggregateType returns the aggregate type
func (e *DomainEventWrapper) AggregateType() string {
	return e.aggregateType
}

// OccurredAt returns when the event occurred
func (e *DomainEventWrapper) OccurredAt() time.Time {
	return e.occurredAt
}

// EventData returns the event data
func (e *DomainEventWrapper) EventData() map[string]interface{} {
	return e.eventData
}

// Version returns the event version
func (e *DomainEventWrapper) Version() int {
	return e.version
}

// Helper functions for extracting data from maps
func getString(data map[string]interface{}, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}

func getInt(data map[string]interface{}, key string) int {
	if value, ok := data[key].(float64); ok {
		return int(value)
	}
	return 0
}

func getTime(data map[string]interface{}, key string) time.Time {
	if value, ok := data[key].(string); ok {
		if t, err := time.Parse(time.RFC3339, value); err == nil {
			return t
		}
	}
	return time.Now()
}

func getMap(data map[string]interface{}, key string) map[string]interface{} {
	if value, ok := data[key].(map[string]interface{}); ok {
		return value
	}
	return make(map[string]interface{})
}
