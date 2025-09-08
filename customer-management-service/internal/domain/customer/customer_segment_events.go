package customer

import (
	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/domain/shared"
)

// CustomerSegmentCreatedEvent represents a customer segment created event
type CustomerSegmentCreatedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerSegmentCreatedEvent creates a new customer segment created event
func NewCustomerSegmentCreatedEvent(segment *CustomerSegment) *CustomerSegmentCreatedEvent {
	eventData := map[string]interface{}{
		"segmentId":   segment.GetID().String(),
		"name":        segment.GetName(),
		"description": segment.GetDescription(),
		"ruleId":      segment.GetRuleID().String(),
		"criteria":    segment.GetCriteria(),
		"status":      segment.GetStatus().String(),
		"createdBy":   segment.GetCreatedBy().String(),
		"createdAt":   segment.GetCreatedAt(),
		"version":     segment.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerSegmentCreated",
		segment.GetID().String(),
		"CustomerSegment",
		eventData,
		segment.GetVersion(),
	)

	return &CustomerSegmentCreatedEvent{BaseDomainEvent: baseEvent}
}

// CustomerSegmentUpdatedEvent represents a customer segment updated event
type CustomerSegmentUpdatedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerSegmentUpdatedEvent creates a new customer segment updated event
func NewCustomerSegmentUpdatedEvent(segment *CustomerSegment) *CustomerSegmentUpdatedEvent {
	eventData := map[string]interface{}{
		"segmentId":     segment.GetID().String(),
		"name":          segment.GetName(),
		"description":   segment.GetDescription(),
		"customerCount": segment.GetCustomerCount(),
		"status":        segment.GetStatus().String(),
		"updatedAt":     segment.GetUpdatedAt(),
		"version":       segment.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerSegmentUpdated",
		segment.GetID().String(),
		"CustomerSegment",
		eventData,
		segment.GetVersion(),
	)

	return &CustomerSegmentUpdatedEvent{BaseDomainEvent: baseEvent}
}

// CustomerSegmentCriteriaUpdatedEvent represents a customer segment criteria updated event
type CustomerSegmentCriteriaUpdatedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerSegmentCriteriaUpdatedEvent creates a new customer segment criteria updated event
func NewCustomerSegmentCriteriaUpdatedEvent(segment *CustomerSegment) *CustomerSegmentCriteriaUpdatedEvent {
	eventData := map[string]interface{}{
		"segmentId": segment.GetID().String(),
		"criteria":  segment.GetCriteria(),
		"updatedAt": segment.GetUpdatedAt(),
		"version":   segment.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerSegmentCriteriaUpdated",
		segment.GetID().String(),
		"CustomerSegment",
		eventData,
		segment.GetVersion(),
	)

	return &CustomerSegmentCriteriaUpdatedEvent{BaseDomainEvent: baseEvent}
}

// CustomerSegmentActivatedEvent represents a customer segment activated event
type CustomerSegmentActivatedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerSegmentActivatedEvent creates a new customer segment activated event
func NewCustomerSegmentActivatedEvent(segment *CustomerSegment) *CustomerSegmentActivatedEvent {
	eventData := map[string]interface{}{
		"segmentId":   segment.GetID().String(),
		"name":        segment.GetName(),
		"activatedAt": segment.GetUpdatedAt(),
		"version":     segment.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerSegmentActivated",
		segment.GetID().String(),
		"CustomerSegment",
		eventData,
		segment.GetVersion(),
	)

	return &CustomerSegmentActivatedEvent{BaseDomainEvent: baseEvent}
}

// CustomerSegmentDeactivatedEvent represents a customer segment deactivated event
type CustomerSegmentDeactivatedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerSegmentDeactivatedEvent creates a new customer segment deactivated event
func NewCustomerSegmentDeactivatedEvent(segment *CustomerSegment) *CustomerSegmentDeactivatedEvent {
	eventData := map[string]interface{}{
		"segmentId":     segment.GetID().String(),
		"name":          segment.GetName(),
		"deactivatedAt": segment.GetUpdatedAt(),
		"version":       segment.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerSegmentDeactivated",
		segment.GetID().String(),
		"CustomerSegment",
		eventData,
		segment.GetVersion(),
	)

	return &CustomerSegmentDeactivatedEvent{BaseDomainEvent: baseEvent}
}

// CustomerSegmentCalculationStartedEvent represents a customer segment calculation started event
type CustomerSegmentCalculationStartedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerSegmentCalculationStartedEvent creates a new customer segment calculation started event
func NewCustomerSegmentCalculationStartedEvent(segment *CustomerSegment) *CustomerSegmentCalculationStartedEvent {
	eventData := map[string]interface{}{
		"segmentId": segment.GetID().String(),
		"name":      segment.GetName(),
		"startedAt": segment.GetUpdatedAt(),
		"version":   segment.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerSegmentCalculationStarted",
		segment.GetID().String(),
		"CustomerSegment",
		eventData,
		segment.GetVersion(),
	)

	return &CustomerSegmentCalculationStartedEvent{BaseDomainEvent: baseEvent}
}

// CustomerSegmentCalculationCompletedEvent represents a customer segment calculation completed event
type CustomerSegmentCalculationCompletedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerSegmentCalculationCompletedEvent creates a new customer segment calculation completed event
func NewCustomerSegmentCalculationCompletedEvent(segment *CustomerSegment, customerCount int) *CustomerSegmentCalculationCompletedEvent {
	eventData := map[string]interface{}{
		"segmentId":      segment.GetID().String(),
		"name":           segment.GetName(),
		"customerCount":  customerCount,
		"completedAt":    segment.GetUpdatedAt(),
		"lastCalculated": segment.GetLastCalculated(),
		"version":        segment.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerSegmentCalculationCompleted",
		segment.GetID().String(),
		"CustomerSegment",
		eventData,
		segment.GetVersion(),
	)

	return &CustomerSegmentCalculationCompletedEvent{BaseDomainEvent: baseEvent}
}

// CustomerSegmentCalculationFailedEvent represents a customer segment calculation failed event
type CustomerSegmentCalculationFailedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerSegmentCalculationFailedEvent creates a new customer segment calculation failed event
func NewCustomerSegmentCalculationFailedEvent(segment *CustomerSegment, errorMessage string) *CustomerSegmentCalculationFailedEvent {
	eventData := map[string]interface{}{
		"segmentId":    segment.GetID().String(),
		"name":         segment.GetName(),
		"errorMessage": errorMessage,
		"failedAt":     segment.GetUpdatedAt(),
		"version":      segment.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerSegmentCalculationFailed",
		segment.GetID().String(),
		"CustomerSegment",
		eventData,
		segment.GetVersion(),
	)

	return &CustomerSegmentCalculationFailedEvent{BaseDomainEvent: baseEvent}
}

// CustomerSegmentDeletedEvent represents a customer segment deleted event
type CustomerSegmentDeletedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerSegmentDeletedEvent creates a new customer segment deleted event
func NewCustomerSegmentDeletedEvent(segment *CustomerSegment) *CustomerSegmentDeletedEvent {
	eventData := map[string]interface{}{
		"segmentId": segment.GetID().String(),
		"name":      segment.GetName(),
		"deletedAt": segment.GetUpdatedAt(),
		"version":   segment.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerSegmentDeleted",
		segment.GetID().String(),
		"CustomerSegment",
		eventData,
		segment.GetVersion(),
	)

	return &CustomerSegmentDeletedEvent{BaseDomainEvent: baseEvent}
}
