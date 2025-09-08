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
	OccurredAt() time.Time
	EventVersion() int
	EventData() map[string]interface{}
}

// BaseDomainEvent provides common fields for domain events
type BaseDomainEvent struct {
	eventID      string
	eventType    string
	aggregateID  string
	occurredAt   time.Time
	eventVersion int
	eventData    map[string]interface{}
}

func NewBaseDomainEvent(eventType, aggregateID string, eventData map[string]interface{}) BaseDomainEvent {
	return BaseDomainEvent{
		eventID:      uuid.New().String(),
		eventType:    eventType,
		aggregateID:  aggregateID,
		occurredAt:   time.Now(),
		eventVersion: 1,
		eventData:    eventData,
	}
}

func (e BaseDomainEvent) EventID() string {
	return e.eventID
}

func (e BaseDomainEvent) EventType() string {
	return e.eventType
}

func (e BaseDomainEvent) AggregateID() string {
	return e.aggregateID
}

func (e BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e BaseDomainEvent) EventVersion() int {
	return e.eventVersion
}

func (e BaseDomainEvent) EventData() map[string]interface{} {
	return e.eventData
}

// EventBus defines the interface for publishing domain events
type EventBus interface {
	Publish(ctx context.Context, event DomainEvent) error
	Subscribe(eventType string, handler EventHandler) error
}

// EventHandler defines the interface for handling domain events
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
}

// EventStore defines the interface for storing domain events
type EventStore interface {
	Save(ctx context.Context, events []DomainEvent) error
	GetEvents(ctx context.Context, aggregateID string) ([]DomainEvent, error)
	GetEventsByType(ctx context.Context, eventType string) ([]DomainEvent, error)
}

// Validator defines the interface for input validation
type Validator interface {
	Validate(value interface{}) error
}
