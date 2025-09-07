package shared

// DomainEvent represents an event that occurred in the domain.
type DomainEvent interface {
	EventType() string
}

// EventBus defines the interface for an event bus.
type EventBus interface {
	Publish(event DomainEvent) error
}
