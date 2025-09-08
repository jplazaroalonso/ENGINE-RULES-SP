package campaign

import (
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// CampaignCreatedEvent represents a campaign creation event
type CampaignCreatedEvent struct {
	shared.BaseDomainEvent
}

func NewCampaignCreatedEvent(campaignID CampaignID, name string, campaignType CampaignType) CampaignCreatedEvent {
	eventData := map[string]interface{}{
		"name":         name,
		"campaignType": campaignType.String(),
	}

	return CampaignCreatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent("CampaignCreated", campaignID.String(), eventData),
	}
}

// CampaignActivatedEvent represents a campaign activation event
type CampaignActivatedEvent struct {
	shared.BaseDomainEvent
}

func NewCampaignActivatedEvent(campaignID CampaignID, name string) CampaignActivatedEvent {
	eventData := map[string]interface{}{
		"name": name,
	}

	return CampaignActivatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent("CampaignActivated", campaignID.String(), eventData),
	}
}

// CampaignPausedEvent represents a campaign pause event
type CampaignPausedEvent struct {
	shared.BaseDomainEvent
}

func NewCampaignPausedEvent(campaignID CampaignID, name string) CampaignPausedEvent {
	eventData := map[string]interface{}{
		"name": name,
	}

	return CampaignPausedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent("CampaignPaused", campaignID.String(), eventData),
	}
}

// CampaignResumedEvent represents a campaign resume event
type CampaignResumedEvent struct {
	shared.BaseDomainEvent
}

func NewCampaignResumedEvent(campaignID CampaignID, name string) CampaignResumedEvent {
	eventData := map[string]interface{}{
		"name": name,
	}

	return CampaignResumedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent("CampaignResumed", campaignID.String(), eventData),
	}
}

// CampaignCompletedEvent represents a campaign completion event
type CampaignCompletedEvent struct {
	shared.BaseDomainEvent
}

func NewCampaignCompletedEvent(campaignID CampaignID, name string) CampaignCompletedEvent {
	eventData := map[string]interface{}{
		"name": name,
	}

	return CampaignCompletedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent("CampaignCompleted", campaignID.String(), eventData),
	}
}

// CampaignCancelledEvent represents a campaign cancellation event
type CampaignCancelledEvent struct {
	shared.BaseDomainEvent
}

func NewCampaignCancelledEvent(campaignID CampaignID, name, reason string) CampaignCancelledEvent {
	eventData := map[string]interface{}{
		"name":   name,
		"reason": reason,
	}

	return CampaignCancelledEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent("CampaignCancelled", campaignID.String(), eventData),
	}
}

// CampaignTargetingRulesUpdatedEvent represents a targeting rules update event
type CampaignTargetingRulesUpdatedEvent struct {
	shared.BaseDomainEvent
}

func NewCampaignTargetingRulesUpdatedEvent(campaignID CampaignID, name string, rules []shared.RuleID) CampaignTargetingRulesUpdatedEvent {
	ruleStrings := make([]string, len(rules))
	for i, rule := range rules {
		ruleStrings[i] = rule.String()
	}

	eventData := map[string]interface{}{
		"name":  name,
		"rules": ruleStrings,
	}

	return CampaignTargetingRulesUpdatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent("CampaignTargetingRulesUpdated", campaignID.String(), eventData),
	}
}

// CampaignBudgetUpdatedEvent represents a budget update event
type CampaignBudgetUpdatedEvent struct {
	shared.BaseDomainEvent
}

func NewCampaignBudgetUpdatedEvent(campaignID CampaignID, name string, budget *shared.Money) CampaignBudgetUpdatedEvent {
	eventData := map[string]interface{}{
		"name": name,
	}

	if budget != nil {
		eventData["budget"] = map[string]interface{}{
			"amount":   budget.Amount,
			"currency": budget.Currency,
		}
	} else {
		eventData["budget"] = nil
	}

	return CampaignBudgetUpdatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent("CampaignBudgetUpdated", campaignID.String(), eventData),
	}
}

// CampaignSettingsUpdatedEvent represents a settings update event
type CampaignSettingsUpdatedEvent struct {
	shared.BaseDomainEvent
}

func NewCampaignSettingsUpdatedEvent(campaignID CampaignID, name string, settings CampaignSettings) CampaignSettingsUpdatedEvent {
	eventData := map[string]interface{}{
		"name":     name,
		"settings": settings,
	}

	return CampaignSettingsUpdatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent("CampaignSettingsUpdated", campaignID.String(), eventData),
	}
}

// CampaignEventTrackedEvent represents a campaign event tracking event
type CampaignEventTrackedEvent struct {
	shared.BaseDomainEvent
}

func NewCampaignEventTrackedEvent(campaignID CampaignID, eventType CampaignEventType, customerID *shared.CustomerID) CampaignEventTrackedEvent {
	eventData := map[string]interface{}{
		"eventType": eventType.String(),
	}

	if customerID != nil {
		eventData["customerId"] = customerID.String()
	}

	return CampaignEventTrackedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent("CampaignEventTracked", campaignID.String(), eventData),
	}
}

// CampaignMetricsUpdatedEvent represents a metrics update event
type CampaignMetricsUpdatedEvent struct {
	shared.BaseDomainEvent
}

func NewCampaignMetricsUpdatedEvent(campaignID CampaignID, metrics CampaignMetrics) CampaignMetricsUpdatedEvent {
	eventData := map[string]interface{}{
		"metrics": map[string]interface{}{
			"impressions":       metrics.Impressions,
			"clicks":            metrics.Clicks,
			"conversions":       metrics.Conversions,
			"revenue":           metrics.Revenue,
			"cost":              metrics.Cost,
			"ctr":               metrics.CTR,
			"conversionRate":    metrics.ConversionRate,
			"costPerClick":      metrics.CostPerClick,
			"costPerConversion": metrics.CostPerConversion,
			"roas":              metrics.ROAS,
			"roi":               metrics.ROI,
		},
	}

	return CampaignMetricsUpdatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent("CampaignMetricsUpdated", campaignID.String(), eventData),
	}
}
