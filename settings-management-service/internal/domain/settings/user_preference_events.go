package settings

import (
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// UserPreferenceCreatedEvent represents a user preference created event
type UserPreferenceCreatedEvent struct {
	shared.BaseDomainEvent
}

// NewUserPreferenceCreatedEvent creates a new user preference created event
func NewUserPreferenceCreatedEvent(preference *UserPreference) *UserPreferenceCreatedEvent {
	eventData := map[string]interface{}{
		"preferenceId":   preference.GetID().String(),
		"userId":         preference.GetUserID().String(),
		"organizationId": preference.GetOrganizationID(),
		"category":       preference.GetCategory(),
		"key":            preference.GetKey(),
		"value":          preference.GetValue(),
		"type":           preference.GetType().String(),
		"isDefault":      preference.IsDefault(),
		"createdAt":      preference.GetCreatedAt(),
		"version":        preference.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"UserPreferenceCreated",
		preference.GetID().String(),
		"UserPreference",
		eventData,
		preference.GetVersion(),
	)

	return &UserPreferenceCreatedEvent{BaseDomainEvent: baseEvent}
}

// UserPreferenceUpdatedEvent represents a user preference updated event
type UserPreferenceUpdatedEvent struct {
	shared.BaseDomainEvent
}

// NewUserPreferenceUpdatedEvent creates a new user preference updated event
func NewUserPreferenceUpdatedEvent(preference *UserPreference) *UserPreferenceUpdatedEvent {
	eventData := map[string]interface{}{
		"preferenceId":   preference.GetID().String(),
		"userId":         preference.GetUserID().String(),
		"organizationId": preference.GetOrganizationID(),
		"category":       preference.GetCategory(),
		"key":            preference.GetKey(),
		"value":          preference.GetValue(),
		"type":           preference.GetType().String(),
		"isDefault":      preference.IsDefault(),
		"updatedAt":      preference.GetUpdatedAt(),
		"version":        preference.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"UserPreferenceUpdated",
		preference.GetID().String(),
		"UserPreference",
		eventData,
		preference.GetVersion(),
	)

	return &UserPreferenceUpdatedEvent{BaseDomainEvent: baseEvent}
}

// UserPreferenceDeletedEvent represents a user preference deleted event
type UserPreferenceDeletedEvent struct {
	shared.BaseDomainEvent
}

// NewUserPreferenceDeletedEvent creates a new user preference deleted event
func NewUserPreferenceDeletedEvent(preference *UserPreference) *UserPreferenceDeletedEvent {
	eventData := map[string]interface{}{
		"preferenceId":   preference.GetID().String(),
		"userId":         preference.GetUserID().String(),
		"organizationId": preference.GetOrganizationID(),
		"category":       preference.GetCategory(),
		"key":            preference.GetKey(),
		"type":           preference.GetType().String(),
		"deletedAt":      preference.GetUpdatedAt(),
		"version":        preference.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"UserPreferenceDeleted",
		preference.GetID().String(),
		"UserPreference",
		eventData,
		preference.GetVersion(),
	)

	return &UserPreferenceDeletedEvent{BaseDomainEvent: baseEvent}
}
