package customer

import (
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/domain/shared"
)

// CustomerCreatedEvent represents a customer created event
type CustomerCreatedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerCreatedEvent creates a new customer created event
func NewCustomerCreatedEvent(customer *Customer) *CustomerCreatedEvent {
	eventData := map[string]interface{}{
		"customerId":  customer.GetID().String(),
		"email":       customer.GetEmail().String(),
		"name":        customer.GetName(),
		"status":      customer.GetStatus().String(),
		"createdAt":   customer.GetCreatedAt(),
		"preferences": customer.GetPreferences(),
		"metadata":    customer.GetMetadata(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerCreated",
		customer.GetID().String(),
		"Customer",
		eventData,
		customer.GetVersion(),
	)

	return &CustomerCreatedEvent{BaseDomainEvent: baseEvent}
}

// CustomerUpdatedEvent represents a customer updated event
type CustomerUpdatedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerUpdatedEvent creates a new customer updated event
func NewCustomerUpdatedEvent(customer *Customer) *CustomerUpdatedEvent {
	eventData := map[string]interface{}{
		"customerId": customer.GetID().String(),
		"email":      customer.GetEmail().String(),
		"name":       customer.GetName(),
		"status":     customer.GetStatus().String(),
		"updatedAt":  customer.GetUpdatedAt(),
		"version":    customer.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerUpdated",
		customer.GetID().String(),
		"Customer",
		eventData,
		customer.GetVersion(),
	)

	return &CustomerUpdatedEvent{BaseDomainEvent: baseEvent}
}

// CustomerEmailUpdatedEvent represents a customer email updated event
type CustomerEmailUpdatedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerEmailUpdatedEvent creates a new customer email updated event
func NewCustomerEmailUpdatedEvent(customer *Customer) *CustomerEmailUpdatedEvent {
	eventData := map[string]interface{}{
		"customerId": customer.GetID().String(),
		"oldEmail":   "", // This would need to be tracked separately
		"newEmail":   customer.GetEmail().String(),
		"updatedAt":  customer.GetUpdatedAt(),
		"version":    customer.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerEmailUpdated",
		customer.GetID().String(),
		"Customer",
		eventData,
		customer.GetVersion(),
	)

	return &CustomerEmailUpdatedEvent{BaseDomainEvent: baseEvent}
}

// CustomerPreferencesUpdatedEvent represents a customer preferences updated event
type CustomerPreferencesUpdatedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerPreferencesUpdatedEvent creates a new customer preferences updated event
func NewCustomerPreferencesUpdatedEvent(customer *Customer) *CustomerPreferencesUpdatedEvent {
	eventData := map[string]interface{}{
		"customerId":  customer.GetID().String(),
		"preferences": customer.GetPreferences(),
		"updatedAt":   customer.GetUpdatedAt(),
		"version":     customer.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerPreferencesUpdated",
		customer.GetID().String(),
		"Customer",
		eventData,
		customer.GetVersion(),
	)

	return &CustomerPreferencesUpdatedEvent{BaseDomainEvent: baseEvent}
}

// CustomerActivatedEvent represents a customer activated event
type CustomerActivatedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerActivatedEvent creates a new customer activated event
func NewCustomerActivatedEvent(customer *Customer) *CustomerActivatedEvent {
	eventData := map[string]interface{}{
		"customerId":  customer.GetID().String(),
		"email":       customer.GetEmail().String(),
		"activatedAt": customer.GetUpdatedAt(),
		"version":     customer.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerActivated",
		customer.GetID().String(),
		"Customer",
		eventData,
		customer.GetVersion(),
	)

	return &CustomerActivatedEvent{BaseDomainEvent: baseEvent}
}

// CustomerDeactivatedEvent represents a customer deactivated event
type CustomerDeactivatedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerDeactivatedEvent creates a new customer deactivated event
func NewCustomerDeactivatedEvent(customer *Customer) *CustomerDeactivatedEvent {
	eventData := map[string]interface{}{
		"customerId":    customer.GetID().String(),
		"email":         customer.GetEmail().String(),
		"deactivatedAt": customer.GetUpdatedAt(),
		"version":       customer.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerDeactivated",
		customer.GetID().String(),
		"Customer",
		eventData,
		customer.GetVersion(),
	)

	return &CustomerDeactivatedEvent{BaseDomainEvent: baseEvent}
}

// CustomerSuspendedEvent represents a customer suspended event
type CustomerSuspendedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerSuspendedEvent creates a new customer suspended event
func NewCustomerSuspendedEvent(customer *Customer) *CustomerSuspendedEvent {
	eventData := map[string]interface{}{
		"customerId":  customer.GetID().String(),
		"email":       customer.GetEmail().String(),
		"suspendedAt": customer.GetUpdatedAt(),
		"version":     customer.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerSuspended",
		customer.GetID().String(),
		"Customer",
		eventData,
		customer.GetVersion(),
	)

	return &CustomerSuspendedEvent{BaseDomainEvent: baseEvent}
}

// CustomerDeletedEvent represents a customer deleted event
type CustomerDeletedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerDeletedEvent creates a new customer deleted event
func NewCustomerDeletedEvent(customer *Customer) *CustomerDeletedEvent {
	eventData := map[string]interface{}{
		"customerId": customer.GetID().String(),
		"email":      customer.GetEmail().String(),
		"deletedAt":  customer.GetUpdatedAt(),
		"version":    customer.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerDeleted",
		customer.GetID().String(),
		"Customer",
		eventData,
		customer.GetVersion(),
	)

	return &CustomerDeletedEvent{BaseDomainEvent: baseEvent}
}

// CustomerSegmentJoinedEvent represents a customer joined segment event
type CustomerSegmentJoinedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerSegmentJoinedEvent creates a new customer joined segment event
func NewCustomerSegmentJoinedEvent(customer *Customer, segmentID SegmentID) *CustomerSegmentJoinedEvent {
	eventData := map[string]interface{}{
		"customerId": customer.GetID().String(),
		"segmentId":  segmentID.String(),
		"joinedAt":   time.Now(),
		"version":    customer.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerSegmentJoined",
		customer.GetID().String(),
		"Customer",
		eventData,
		customer.GetVersion(),
	)

	return &CustomerSegmentJoinedEvent{BaseDomainEvent: baseEvent}
}

// CustomerSegmentLeftEvent represents a customer left segment event
type CustomerSegmentLeftEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerSegmentLeftEvent creates a new customer left segment event
func NewCustomerSegmentLeftEvent(customer *Customer, segmentID SegmentID) *CustomerSegmentLeftEvent {
	eventData := map[string]interface{}{
		"customerId": customer.GetID().String(),
		"segmentId":  segmentID.String(),
		"leftAt":     time.Now(),
		"version":    customer.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerSegmentLeft",
		customer.GetID().String(),
		"Customer",
		eventData,
		customer.GetVersion(),
	)

	return &CustomerSegmentLeftEvent{BaseDomainEvent: baseEvent}
}

// CustomerLoggedInEvent represents a customer logged in event
type CustomerLoggedInEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerLoggedInEvent creates a new customer logged in event
func NewCustomerLoggedInEvent(customer *Customer) *CustomerLoggedInEvent {
	eventData := map[string]interface{}{
		"customerId": customer.GetID().String(),
		"email":      customer.GetEmail().String(),
		"loggedInAt": time.Now(),
		"loginCount": customer.GetMetadata().LoginCount,
		"version":    customer.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerLoggedIn",
		customer.GetID().String(),
		"Customer",
		eventData,
		customer.GetVersion(),
	)

	return &CustomerLoggedInEvent{BaseDomainEvent: baseEvent}
}

// CustomerPurchaseEvent represents a customer purchase event
type CustomerPurchaseEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerPurchaseEvent creates a new customer purchase event
func NewCustomerPurchaseEvent(customer *Customer, purchase PurchaseRecord) *CustomerPurchaseEvent {
	eventData := map[string]interface{}{
		"customerId":   customer.GetID().String(),
		"purchaseId":   purchase.ID,
		"amount":       purchase.Amount,
		"product":      purchase.Product,
		"category":     purchase.Category,
		"purchaseDate": purchase.PurchaseDate,
		"channel":      purchase.Channel,
		"version":      customer.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerPurchase",
		customer.GetID().String(),
		"Customer",
		eventData,
		customer.GetVersion(),
	)

	return &CustomerPurchaseEvent{BaseDomainEvent: baseEvent}
}

// CustomerInteractionEvent represents a customer interaction event
type CustomerInteractionEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerInteractionEvent creates a new customer interaction event
func NewCustomerInteractionEvent(customer *Customer, interaction InteractionRecord) *CustomerInteractionEvent {
	eventData := map[string]interface{}{
		"customerId":    customer.GetID().String(),
		"interactionId": interaction.ID,
		"type":          interaction.Type,
		"channel":       interaction.Channel,
		"action":        interaction.Action,
		"timestamp":     interaction.Timestamp,
		"duration":      interaction.Duration,
		"outcome":       interaction.Outcome,
		"version":       customer.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerInteraction",
		customer.GetID().String(),
		"Customer",
		eventData,
		customer.GetVersion(),
	)

	return &CustomerInteractionEvent{BaseDomainEvent: baseEvent}
}

// CustomerDeviceAddedEvent represents a customer device added event
type CustomerDeviceAddedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerDeviceAddedEvent creates a new customer device added event
func NewCustomerDeviceAddedEvent(customer *Customer, device DeviceInfo) *CustomerDeviceAddedEvent {
	eventData := map[string]interface{}{
		"customerId": customer.GetID().String(),
		"deviceType": device.Type,
		"os":         device.OS,
		"browser":    device.Browser,
		"firstSeen":  device.FirstSeen,
		"isActive":   device.IsActive,
		"version":    customer.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerDeviceAdded",
		customer.GetID().String(),
		"Customer",
		eventData,
		customer.GetVersion(),
	)

	return &CustomerDeviceAddedEvent{BaseDomainEvent: baseEvent}
}

// CustomerDeviceUpdatedEvent represents a customer device updated event
type CustomerDeviceUpdatedEvent struct {
	*shared.BaseDomainEvent
}

// NewCustomerDeviceUpdatedEvent creates a new customer device updated event
func NewCustomerDeviceUpdatedEvent(customer *Customer, device DeviceInfo) *CustomerDeviceUpdatedEvent {
	eventData := map[string]interface{}{
		"customerId": customer.GetID().String(),
		"deviceType": device.Type,
		"os":         device.OS,
		"browser":    device.Browser,
		"lastSeen":   device.LastSeen,
		"isActive":   device.IsActive,
		"version":    customer.GetVersion(),
	}

	baseEvent := shared.NewBaseDomainEvent(
		"CustomerDeviceUpdated",
		customer.GetID().String(),
		"Customer",
		eventData,
		customer.GetVersion(),
	)

	return &CustomerDeviceUpdatedEvent{BaseDomainEvent: baseEvent}
}
