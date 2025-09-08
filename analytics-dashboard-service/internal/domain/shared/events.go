package shared

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent represents a domain event
type DomainEvent struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	AggregateID string                 `json:"aggregateId"`
	Version     int                    `json:"version"`
	Data        map[string]interface{} `json:"data"`
	Metadata    map[string]interface{} `json:"metadata"`
	Timestamp   time.Time              `json:"timestamp"`
}

// NewDomainEvent creates a new domain event
func NewDomainEvent(eventType, aggregateID string, version int, data map[string]interface{}) *DomainEvent {
	return &DomainEvent{
		ID:          uuid.New().String(),
		Type:        eventType,
		AggregateID: aggregateID,
		Version:     version,
		Data:        data,
		Metadata:    make(map[string]interface{}),
		Timestamp:   time.Now(),
	}
}

// EventBus interface for publishing domain events
type EventBus interface {
	Publish(event *DomainEvent) error
	Subscribe(eventType string, handler func(*DomainEvent) error) error
}

// EventStore interface for storing domain events
type EventStore interface {
	Save(event *DomainEvent) error
	GetEvents(aggregateID string) ([]*DomainEvent, error)
	GetEventsByType(eventType string) ([]*DomainEvent, error)
}
