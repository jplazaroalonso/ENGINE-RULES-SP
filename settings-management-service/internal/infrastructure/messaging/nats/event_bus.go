package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// EventBus implements the shared.EventBus interface using NATS JetStream
type EventBus struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

// NewEventBus creates a new NATS-based event bus
func NewEventBus(conn *nats.Conn) (*EventBus, error) {
	js, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to get JetStream context: %w", err)
	}

	// Create streams for different event types
	if err := createStreams(js); err != nil {
		return nil, fmt.Errorf("failed to create streams: %w", err)
	}

	return &EventBus{
		conn: conn,
		js:   js,
	}, nil
}

// Publish publishes a domain event to the event bus
func (eb *EventBus) Publish(ctx context.Context, event shared.DomainEvent) error {
	// Serialize the event
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Determine the subject based on event type
	subject := getSubjectForEventType(event.EventType())

	// Publish the event
	_, err = eb.js.PublishAsync(subject, eventData)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Published event %s with ID %s to subject %s", event.EventType(), event.EventID(), subject)
	return nil
}

// Subscribe subscribes to events of a specific type
func (eb *EventBus) Subscribe(ctx context.Context, eventType string, handler func(context.Context, shared.DomainEvent) error) error {
	subject := getSubjectForEventType(eventType)
	streamName := getStreamNameForEventType(eventType)

	// Create a consumer for the subscription
	consumerName := fmt.Sprintf("settings-service-%s-consumer", eventType)

	// Subscribe to the subject
	sub, err := eb.js.PullSubscribe(subject, consumerName, nats.Durable(consumerName))
	if err != nil {
		return fmt.Errorf("failed to create pull subscription: %w", err)
	}

	// Start processing messages in a goroutine
	go func() {
		defer sub.Unsubscribe()

		for {
			select {
			case <-ctx.Done():
				log.Printf("Stopping subscription for event type %s", eventType)
				return
			default:
				// Fetch messages
				msgs, err := sub.Fetch(1, nats.MaxWait(1*time.Second))
				if err != nil {
					if err == nats.ErrTimeout {
						continue
					}
					log.Printf("Error fetching messages for event type %s: %v", eventType, err)
					continue
				}

				// Process each message
				for _, msg := range msgs {
					if err := eb.processMessage(ctx, msg, handler); err != nil {
						log.Printf("Error processing message for event type %s: %v", eventType, err)
						// Acknowledge the message even if processing failed to avoid reprocessing
						msg.Ack()
					} else {
						msg.Ack()
					}
				}
			}
		}
	}()

	log.Printf("Subscribed to event type %s on subject %s", eventType, subject)
	return nil
}

// processMessage processes a single NATS message
func (eb *EventBus) processMessage(ctx context.Context, msg *nats.Msg, handler func(context.Context, shared.DomainEvent) error) error {
	// Deserialize the event
	var event shared.BaseDomainEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	// Call the handler
	return handler(ctx, &event)
}

// createStreams creates the necessary JetStream streams
func createStreams(js nats.JetStreamContext) error {
	streams := []struct {
		name     string
		subjects []string
	}{
		{
			name:     "SETTINGS_CONFIGURATION_EVENTS",
			subjects: []string{"settings.configuration.*"},
		},
		{
			name:     "SETTINGS_FEATURE_FLAG_EVENTS",
			subjects: []string{"settings.feature-flag.*"},
		},
		{
			name:     "SETTINGS_USER_PREFERENCE_EVENTS",
			subjects: []string{"settings.user-preference.*"},
		},
		{
			name:     "SETTINGS_ORGANIZATION_SETTING_EVENTS",
			subjects: []string{"settings.organization-setting.*"},
		},
		{
			name:     "SETTINGS_CACHE_EVENTS",
			subjects: []string{"settings.cache.*"},
		},
		{
			name:     "SETTINGS_COMPLIANCE_EVENTS",
			subjects: []string{"settings.compliance.*"},
		},
	}

	for _, stream := range streams {
		// Check if stream already exists
		_, err := js.StreamInfo(stream.name)
		if err == nil {
			// Stream already exists, skip
			continue
		}

		// Create the stream
		_, err = js.AddStream(&nats.StreamConfig{
			Name:      stream.name,
			Subjects:  stream.subjects,
			Retention: nats.WorkQueuePolicy,
			MaxAge:    24 * time.Hour, // Keep events for 24 hours
			Storage:   nats.FileStorage,
			Replicas:  1,
		})
		if err != nil {
			return fmt.Errorf("failed to create stream %s: %w", stream.name, err)
		}

		log.Printf("Created JetStream stream: %s", stream.name)
	}

	return nil
}

// getSubjectForEventType returns the NATS subject for a given event type
func getSubjectForEventType(eventType string) string {
	switch eventType {
	case "ConfigurationCreated", "ConfigurationUpdated", "ConfigurationDeleted":
		return fmt.Sprintf("settings.configuration.%s", eventType)
	case "FeatureFlagCreated", "FeatureFlagUpdated", "FeatureFlagDeleted", "FeatureFlagEvaluated":
		return fmt.Sprintf("settings.feature-flag.%s", eventType)
	case "UserPreferenceCreated", "UserPreferenceUpdated", "UserPreferenceDeleted":
		return fmt.Sprintf("settings.user-preference.%s", eventType)
	case "OrganizationSettingCreated", "OrganizationSettingUpdated", "OrganizationSettingDeleted":
		return fmt.Sprintf("settings.organization-setting.%s", eventType)
	case "SettingsCacheInvalidated":
		return fmt.Sprintf("settings.cache.%s", eventType)
	case "ComplianceReportGenerated", "ComplianceValidationCompleted":
		return fmt.Sprintf("settings.compliance.%s", eventType)
	default:
		return fmt.Sprintf("settings.unknown.%s", eventType)
	}
}

// getStreamNameForEventType returns the JetStream stream name for a given event type
func getStreamNameForEventType(eventType string) string {
	switch eventType {
	case "ConfigurationCreated", "ConfigurationUpdated", "ConfigurationDeleted":
		return "SETTINGS_CONFIGURATION_EVENTS"
	case "FeatureFlagCreated", "FeatureFlagUpdated", "FeatureFlagDeleted", "FeatureFlagEvaluated":
		return "SETTINGS_FEATURE_FLAG_EVENTS"
	case "UserPreferenceCreated", "UserPreferenceUpdated", "UserPreferenceDeleted":
		return "SETTINGS_USER_PREFERENCE_EVENTS"
	case "OrganizationSettingCreated", "OrganizationSettingUpdated", "OrganizationSettingDeleted":
		return "SETTINGS_ORGANIZATION_SETTING_EVENTS"
	case "SettingsCacheInvalidated":
		return "SETTINGS_CACHE_EVENTS"
	case "ComplianceReportGenerated", "ComplianceValidationCompleted":
		return "SETTINGS_COMPLIANCE_EVENTS"
	default:
		return "SETTINGS_UNKNOWN_EVENTS"
	}
}

// Close closes the event bus connection
func (eb *EventBus) Close() error {
	if eb.conn != nil {
		eb.conn.Close()
	}
	return nil
}

// HealthCheck performs a health check on the NATS connection
func (eb *EventBus) HealthCheck() error {
	if eb.conn == nil {
		return fmt.Errorf("NATS connection is nil")
	}

	if !eb.conn.IsConnected() {
		return fmt.Errorf("NATS connection is not connected")
	}

	// Check JetStream availability
	if eb.js == nil {
		return fmt.Errorf("JetStream context is nil")
	}

	return nil
}
