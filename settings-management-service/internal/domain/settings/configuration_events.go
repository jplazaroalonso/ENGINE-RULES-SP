package settings

import (
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// ConfigurationCreatedEvent represents a configuration created event
type ConfigurationCreatedEvent struct {
	shared.BaseDomainEvent
}

// NewConfigurationCreatedEvent creates a new configuration created event
func NewConfigurationCreatedEvent(configuration *Configuration) *ConfigurationCreatedEvent {
	eventData := map[string]interface{}{
		"configurationId": configuration.GetID().String(),
		"key":             configuration.GetKey(),
		"type":            configuration.GetType().String(),
		"category":        configuration.GetCategory(),
		"description":     configuration.GetDescription(),
		"environment":     configuration.GetEnvironment().String(),
		"organizationId":  configuration.GetOrganizationID(),
		"service":         configuration.GetService(),
		"isEncrypted":     configuration.IsEncrypted(),
		"isSensitive":     configuration.IsSensitive(),
		"createdBy":       configuration.GetCreatedBy().String(),
		"createdAt":       configuration.GetCreatedAt(),
		"version":         configuration.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"ConfigurationCreated",
		configuration.GetID().String(),
		"Configuration",
		eventData,
		configuration.GetVersion(),
	)

	return &ConfigurationCreatedEvent{BaseDomainEvent: baseEvent}
}

// ConfigurationUpdatedEvent represents a configuration updated event
type ConfigurationUpdatedEvent struct {
	shared.BaseDomainEvent
}

// NewConfigurationUpdatedEvent creates a new configuration updated event
func NewConfigurationUpdatedEvent(configuration *Configuration) *ConfigurationUpdatedEvent {
	eventData := map[string]interface{}{
		"configurationId": configuration.GetID().String(),
		"key":             configuration.GetKey(),
		"type":            configuration.GetType().String(),
		"category":        configuration.GetCategory(),
		"description":     configuration.GetDescription(),
		"environment":     configuration.GetEnvironment().String(),
		"organizationId":  configuration.GetOrganizationID(),
		"service":         configuration.GetService(),
		"isEncrypted":     configuration.IsEncrypted(),
		"isSensitive":     configuration.IsSensitive(),
		"updatedAt":       configuration.GetUpdatedAt(),
		"version":         configuration.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"ConfigurationUpdated",
		configuration.GetID().String(),
		"Configuration",
		eventData,
		configuration.GetVersion(),
	)

	return &ConfigurationUpdatedEvent{BaseDomainEvent: baseEvent}
}

// ConfigurationDeletedEvent represents a configuration deleted event
type ConfigurationDeletedEvent struct {
	shared.BaseDomainEvent
}

// NewConfigurationDeletedEvent creates a new configuration deleted event
func NewConfigurationDeletedEvent(configuration *Configuration) *ConfigurationDeletedEvent {
	eventData := map[string]interface{}{
		"configurationId": configuration.GetID().String(),
		"key":             configuration.GetKey(),
		"type":            configuration.GetType().String(),
		"category":        configuration.GetCategory(),
		"environment":     configuration.GetEnvironment().String(),
		"organizationId":  configuration.GetOrganizationID(),
		"service":         configuration.GetService(),
		"deletedAt":       configuration.GetUpdatedAt(),
		"version":         configuration.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"ConfigurationDeleted",
		configuration.GetID().String(),
		"Configuration",
		eventData,
		configuration.GetVersion(),
	)

	return &ConfigurationDeletedEvent{BaseDomainEvent: baseEvent}
}
