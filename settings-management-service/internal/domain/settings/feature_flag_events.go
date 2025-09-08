package settings

import (
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// FeatureFlagCreatedEvent represents a feature flag created event
type FeatureFlagCreatedEvent struct {
	shared.BaseDomainEvent
}

// NewFeatureFlagCreatedEvent creates a new feature flag created event
func NewFeatureFlagCreatedEvent(featureFlag *FeatureFlag) *FeatureFlagCreatedEvent {
	eventData := map[string]interface{}{
		"featureFlagId":   featureFlag.GetID().String(),
		"name":            featureFlag.GetName(),
		"description":     featureFlag.GetDescription(),
		"key":             featureFlag.GetKey(),
		"isEnabled":       featureFlag.IsEnabled(),
		"rolloutStrategy": featureFlag.GetRolloutStrategy().String(),
		"targetAudience":  featureFlag.GetTargetAudience(),
		"variants":        featureFlag.GetVariants(),
		"rules":           featureFlag.GetRules(),
		"environment":     featureFlag.GetEnvironment().String(),
		"organizationId":  featureFlag.GetOrganizationID(),
		"service":         featureFlag.GetService(),
		"createdBy":       featureFlag.GetCreatedBy().String(),
		"createdAt":       featureFlag.GetCreatedAt(),
		"version":         featureFlag.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"FeatureFlagCreated",
		featureFlag.GetID().String(),
		"FeatureFlag",
		eventData,
		featureFlag.GetVersion(),
	)

	return &FeatureFlagCreatedEvent{BaseDomainEvent: baseEvent}
}

// FeatureFlagUpdatedEvent represents a feature flag updated event
type FeatureFlagUpdatedEvent struct {
	shared.BaseDomainEvent
}

// NewFeatureFlagUpdatedEvent creates a new feature flag updated event
func NewFeatureFlagUpdatedEvent(featureFlag *FeatureFlag) *FeatureFlagUpdatedEvent {
	eventData := map[string]interface{}{
		"featureFlagId":   featureFlag.GetID().String(),
		"name":            featureFlag.GetName(),
		"key":             featureFlag.GetKey(),
		"isEnabled":       featureFlag.IsEnabled(),
		"rolloutStrategy": featureFlag.GetRolloutStrategy().String(),
		"targetAudience":  featureFlag.GetTargetAudience(),
		"variants":        featureFlag.GetVariants(),
		"rules":           featureFlag.GetRules(),
		"environment":     featureFlag.GetEnvironment().String(),
		"organizationId":  featureFlag.GetOrganizationID(),
		"service":         featureFlag.GetService(),
		"updatedAt":       featureFlag.GetUpdatedAt(),
		"version":         featureFlag.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"FeatureFlagUpdated",
		featureFlag.GetID().String(),
		"FeatureFlag",
		eventData,
		featureFlag.GetVersion(),
	)

	return &FeatureFlagUpdatedEvent{BaseDomainEvent: baseEvent}
}

// FeatureFlagEnabledEvent represents a feature flag enabled event
type FeatureFlagEnabledEvent struct {
	*shared.BaseDomainEvent
}

// NewFeatureFlagEnabledEvent creates a new feature flag enabled event
func NewFeatureFlagEnabledEvent(featureFlag *FeatureFlag) *FeatureFlagEnabledEvent {
	eventData := map[string]interface{}{
		"featureFlagId": featureFlag.GetID().String(),
		"name":          featureFlag.GetName(),
		"key":           featureFlag.GetKey(),
		"enabledAt":     featureFlag.GetUpdatedAt(),
		"version":       featureFlag.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"FeatureFlagEnabled",
		featureFlag.GetID().String(),
		"FeatureFlag",
		eventData,
		featureFlag.GetVersion(),
	)

	return &FeatureFlagEnabledEvent{BaseDomainEvent: baseEvent}
}

// FeatureFlagDisabledEvent represents a feature flag disabled event
type FeatureFlagDisabledEvent struct {
	*shared.BaseDomainEvent
}

// NewFeatureFlagDisabledEvent creates a new feature flag disabled event
func NewFeatureFlagDisabledEvent(featureFlag *FeatureFlag) *FeatureFlagDisabledEvent {
	eventData := map[string]interface{}{
		"featureFlagId": featureFlag.GetID().String(),
		"name":          featureFlag.GetName(),
		"key":           featureFlag.GetKey(),
		"disabledAt":    featureFlag.GetUpdatedAt(),
		"version":       featureFlag.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"FeatureFlagDisabled",
		featureFlag.GetID().String(),
		"FeatureFlag",
		eventData,
		featureFlag.GetVersion(),
	)

	return &FeatureFlagDisabledEvent{BaseDomainEvent: baseEvent}
}

// FeatureFlagDeletedEvent represents a feature flag deleted event
type FeatureFlagDeletedEvent struct {
	shared.BaseDomainEvent
}

// NewFeatureFlagDeletedEvent creates a new feature flag deleted event
func NewFeatureFlagDeletedEvent(featureFlag *FeatureFlag) *FeatureFlagDeletedEvent {
	eventData := map[string]interface{}{
		"featureFlagId":  featureFlag.GetID().String(),
		"name":           featureFlag.GetName(),
		"key":            featureFlag.GetKey(),
		"environment":    featureFlag.GetEnvironment().String(),
		"organizationId": featureFlag.GetOrganizationID(),
		"service":        featureFlag.GetService(),
		"deletedAt":      featureFlag.GetUpdatedAt(),
		"version":        featureFlag.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"FeatureFlagDeleted",
		featureFlag.GetID().String(),
		"FeatureFlag",
		eventData,
		featureFlag.GetVersion(),
	)

	return &FeatureFlagDeletedEvent{BaseDomainEvent: baseEvent}
}
