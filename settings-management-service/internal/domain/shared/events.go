package shared

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// DomainEvent represents a domain event
type DomainEvent interface {
	EventID() string
	EventType() string
	AggregateID() string
	AggregateType() string
	OccurredAt() time.Time
	EventData() map[string]interface{}
	Version() int
}

// BaseDomainEvent provides a base implementation for domain events
type BaseDomainEvent struct {
	eventID       string
	eventType     string
	aggregateID   string
	aggregateType string
	occurredAt    time.Time
	eventData     map[string]interface{}
	version       int
}

// EventID returns the event ID
func (e *BaseDomainEvent) EventID() string {
	return e.eventID
}

// EventType returns the event type
func (e *BaseDomainEvent) EventType() string {
	return e.eventType
}

// AggregateID returns the aggregate ID
func (e *BaseDomainEvent) AggregateID() string {
	return e.aggregateID
}

// AggregateType returns the aggregate type
func (e *BaseDomainEvent) AggregateType() string {
	return e.aggregateType
}

// OccurredAt returns when the event occurred
func (e *BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// EventData returns the event data
func (e *BaseDomainEvent) EventData() map[string]interface{} {
	return e.eventData
}

// Version returns the event version
func (e *BaseDomainEvent) Version() int {
	return e.version
}

// NewBaseDomainEvent creates a new base domain event
func NewBaseDomainEvent(eventType, aggregateID, aggregateType string, eventData map[string]interface{}, version int) *BaseDomainEvent {
	return &BaseDomainEvent{
		eventID:       uuid.New().String(),
		eventType:     eventType,
		aggregateID:   aggregateID,
		aggregateType: aggregateType,
		occurredAt:    time.Now(),
		eventData:     eventData,
		version:       version,
	}
}

// EventBus defines the interface for publishing domain events
type EventBus interface {
	Publish(ctx context.Context, event DomainEvent) error
	PublishBatch(ctx context.Context, events []DomainEvent) error
	Subscribe(ctx context.Context, eventType string, handler EventHandler) error
}

// EventHandler defines the interface for handling domain events
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
	EventType() string
}

// EventStore defines the interface for storing and retrieving domain events
type EventStore interface {
	Save(ctx context.Context, events []DomainEvent) error
	GetEvents(ctx context.Context, aggregateID string, fromVersion int) ([]DomainEvent, error)
	GetEventsByType(ctx context.Context, eventType string, fromTime time.Time) ([]DomainEvent, error)
}

// EventPublisher defines the interface for publishing events to external systems
type EventPublisher interface {
	PublishEvent(ctx context.Context, event DomainEvent) error
	PublishEvents(ctx context.Context, events []DomainEvent) error
}
