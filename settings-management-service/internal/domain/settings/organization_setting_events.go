package settings

import (
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// OrganizationSettingCreatedEvent represents an organization setting created event
type OrganizationSettingCreatedEvent struct {
	shared.BaseDomainEvent
}

// NewOrganizationSettingCreatedEvent creates a new organization setting created event
func NewOrganizationSettingCreatedEvent(setting *OrganizationSetting) *OrganizationSettingCreatedEvent {
	eventData := map[string]interface{}{
		"settingId":      setting.GetID().String(),
		"organizationId": setting.GetOrganizationID().String(),
		"category":       setting.GetCategory(),
		"key":            setting.GetKey(),
		"value":          setting.GetValue(),
		"type":           setting.GetType().String(),
		"isInherited":    setting.IsInherited(),
		"parentId":       setting.GetParentID(),
		"createdBy":      setting.GetCreatedBy().String(),
		"createdAt":      setting.GetCreatedAt(),
		"version":        setting.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"OrganizationSettingCreated",
		setting.GetID().String(),
		"OrganizationSetting",
		eventData,
		setting.GetVersion(),
	)

	return &OrganizationSettingCreatedEvent{BaseDomainEvent: baseEvent}
}

// OrganizationSettingUpdatedEvent represents an organization setting updated event
type OrganizationSettingUpdatedEvent struct {
	shared.BaseDomainEvent
}

// NewOrganizationSettingUpdatedEvent creates a new organization setting updated event
func NewOrganizationSettingUpdatedEvent(setting *OrganizationSetting) *OrganizationSettingUpdatedEvent {
	eventData := map[string]interface{}{
		"settingId":      setting.GetID().String(),
		"organizationId": setting.GetOrganizationID().String(),
		"category":       setting.GetCategory(),
		"key":            setting.GetKey(),
		"value":          setting.GetValue(),
		"type":           setting.GetType().String(),
		"isInherited":    setting.IsInherited(),
		"parentId":       setting.GetParentID(),
		"updatedAt":      setting.GetUpdatedAt(),
		"version":        setting.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"OrganizationSettingUpdated",
		setting.GetID().String(),
		"OrganizationSetting",
		eventData,
		setting.GetVersion(),
	)

	return &OrganizationSettingUpdatedEvent{BaseDomainEvent: baseEvent}
}

// OrganizationSettingDeletedEvent represents an organization setting deleted event
type OrganizationSettingDeletedEvent struct {
	shared.BaseDomainEvent
}

// NewOrganizationSettingDeletedEvent creates a new organization setting deleted event
func NewOrganizationSettingDeletedEvent(setting *OrganizationSetting) *OrganizationSettingDeletedEvent {
	eventData := map[string]interface{}{
		"settingId":      setting.GetID().String(),
		"organizationId": setting.GetOrganizationID().String(),
		"category":       setting.GetCategory(),
		"key":            setting.GetKey(),
		"type":           setting.GetType().String(),
		"isInherited":    setting.IsInherited(),
		"parentId":       setting.GetParentID(),
		"deletedAt":      setting.GetUpdatedAt(),
		"version":        setting.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"OrganizationSettingDeleted",
		setting.GetID().String(),
		"OrganizationSetting",
		eventData,
		setting.GetVersion(),
	)

	return &OrganizationSettingDeletedEvent{BaseDomainEvent: baseEvent}
}
