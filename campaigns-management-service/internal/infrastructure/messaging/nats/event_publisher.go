package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
	"github.com/nats-io/nats.go"
)

// EventPublisher implements the EventBus interface using NATS
type EventPublisher struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

// NewEventPublisher creates a new NATS event publisher
func NewEventPublisher(natsURL string) (*EventPublisher, error) {
	// Connect to NATS
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Create JetStream context
	js, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Create streams if they don't exist
	if err := createStreams(js); err != nil {
		return nil, fmt.Errorf("failed to create streams: %w", err)
	}

	return &EventPublisher{
		conn: conn,
		js:   js,
	}, nil
}

// Publish publishes a domain event to NATS
func (p *EventPublisher) Publish(ctx context.Context, event shared.DomainEvent) error {
	// Create event message
	eventMessage := EventMessage{
		EventID:      event.EventID(),
		EventType:    event.EventType(),
		AggregateID:  event.AggregateID(),
		OccurredAt:   event.OccurredAt(),
		EventVersion: event.EventVersion(),
		EventData:    event.EventData(),
		PublishedAt:  time.Now(),
	}

	// Serialize event message
	data, err := json.Marshal(eventMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal event message: %w", err)
	}

	// Determine subject based on event type
	subject := getSubjectForEventType(event.EventType())

	// Publish to JetStream
	_, err = p.js.PublishAsync(subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish event to NATS: %w", err)
	}

	return nil
}

// Subscribe subscribes to domain events
func (p *EventPublisher) Subscribe(eventType string, handler shared.EventHandler) error {
	subject := getSubjectForEventType(eventType)

	// Create subscription
	_, err := p.js.Subscribe(subject, func(msg *nats.Msg) {
		// Deserialize event message
		var eventMessage EventMessage
		if err := json.Unmarshal(msg.Data, &eventMessage); err != nil {
			fmt.Printf("Failed to unmarshal event message: %v\n", err)
			return
		}

		// Create domain event
		event := &DomainEventWrapper{
			eventMessage: eventMessage,
		}

		// Handle event
		if err := handler.Handle(context.Background(), event); err != nil {
			fmt.Printf("Failed to handle event: %v\n", err)
			// In production, you might want to implement retry logic or dead letter queue
		}

		// Acknowledge message
		msg.Ack()
	}, nats.Durable("campaigns-service"))

	if err != nil {
		return fmt.Errorf("failed to subscribe to events: %w", err)
	}

	return nil
}

// Close closes the NATS connection
func (p *EventPublisher) Close() error {
	if p.conn != nil {
		p.conn.Close()
	}
	return nil
}

// EventMessage represents a message sent to NATS
type EventMessage struct {
	EventID      string                 `json:"eventId"`
	EventType    string                 `json:"eventType"`
	AggregateID  string                 `json:"aggregateId"`
	OccurredAt   time.Time              `json:"occurredAt"`
	EventVersion int                    `json:"eventVersion"`
	EventData    map[string]interface{} `json:"eventData"`
	PublishedAt  time.Time              `json:"publishedAt"`
}

// DomainEventWrapper wraps EventMessage to implement DomainEvent interface
type DomainEventWrapper struct {
	eventMessage EventMessage
}

func (e *DomainEventWrapper) EventID() string {
	return e.eventMessage.EventID
}

func (e *DomainEventWrapper) EventType() string {
	return e.eventMessage.EventType
}

func (e *DomainEventWrapper) AggregateID() string {
	return e.eventMessage.AggregateID
}

func (e *DomainEventWrapper) OccurredAt() time.Time {
	return e.eventMessage.OccurredAt
}

func (e *DomainEventWrapper) EventVersion() int {
	return e.eventMessage.EventVersion
}

func (e *DomainEventWrapper) EventData() map[string]interface{} {
	return e.eventMessage.EventData
}

// createStreams creates necessary NATS streams
func createStreams(js nats.JetStreamContext) error {
	// Campaign events stream
	_, err := js.AddStream(&nats.StreamConfig{
		Name:        "CAMPAIGN_EVENTS",
		Description: "Stream for campaign domain events",
		Subjects:    []string{"campaigns.*"},
		Retention:   nats.LimitsPolicy,
		MaxAge:      24 * time.Hour * 30, // 30 days
		Storage:     nats.FileStorage,
		Replicas:    1,
	})
	if err != nil {
		return fmt.Errorf("failed to create campaign events stream: %w", err)
	}

	// Analytics events stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:        "ANALYTICS_EVENTS",
		Description: "Stream for analytics events",
		Subjects:    []string{"analytics.*"},
		Retention:   nats.LimitsPolicy,
		MaxAge:      24 * time.Hour * 90, // 90 days
		Storage:     nats.FileStorage,
		Replicas:    1,
	})
	if err != nil {
		return fmt.Errorf("failed to create analytics events stream: %w", err)
	}

	// Notification events stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:        "NOTIFICATION_EVENTS",
		Description: "Stream for notification events",
		Subjects:    []string{"notifications.*"},
		Retention:   nats.LimitsPolicy,
		MaxAge:      24 * time.Hour * 7, // 7 days
		Storage:     nats.FileStorage,
		Replicas:    1,
	})
	if err != nil {
		return fmt.Errorf("failed to create notification events stream: %w", err)
	}

	return nil
}

// getSubjectForEventType returns the NATS subject for a given event type
func getSubjectForEventType(eventType string) string {
	switch eventType {
	case "CampaignCreated":
		return "campaigns.created"
	case "CampaignActivated":
		return "campaigns.activated"
	case "CampaignPaused":
		return "campaigns.paused"
	case "CampaignResumed":
		return "campaigns.resumed"
	case "CampaignCompleted":
		return "campaigns.completed"
	case "CampaignCancelled":
		return "campaigns.cancelled"
	case "CampaignTargetingRulesUpdated":
		return "campaigns.targeting_updated"
	case "CampaignBudgetUpdated":
		return "campaigns.budget_updated"
	case "CampaignSettingsUpdated":
		return "campaigns.settings_updated"
	case "CampaignEventTracked":
		return "campaigns.event_tracked"
	case "CampaignMetricsUpdated":
		return "analytics.metrics_updated"
	default:
		return fmt.Sprintf("campaigns.%s", eventType)
	}
}

// NoOpEventPublisher is a no-op implementation for testing or when NATS is not available
type NoOpEventPublisher struct{}

// NewNoOpEventPublisher creates a new no-op event publisher
func NewNoOpEventPublisher() *NoOpEventPublisher {
	return &NoOpEventPublisher{}
}

// Publish does nothing (no-op)
func (p *NoOpEventPublisher) Publish(ctx context.Context, event shared.DomainEvent) error {
	// Log the event for debugging purposes
	fmt.Printf("NoOpEventPublisher: Would publish event %s for aggregate %s\n",
		event.EventType(), event.AggregateID())
	return nil
}

// Subscribe does nothing (no-op)
func (p *NoOpEventPublisher) Subscribe(eventType string, handler shared.EventHandler) error {
	fmt.Printf("NoOpEventPublisher: Would subscribe to event type %s\n", eventType)
	return nil
}

// Close does nothing (no-op)
func (p *NoOpEventPublisher) Close() error {
	return nil
}
