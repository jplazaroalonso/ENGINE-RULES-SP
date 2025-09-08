package campaign

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// CampaignID represents a campaign identifier
type CampaignID struct {
	value uuid.UUID
}

func NewCampaignID() CampaignID {
	return CampaignID{value: uuid.New()}
}

func NewCampaignIDFromString(value string) (CampaignID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return CampaignID{}, shared.NewValidationError("invalid campaign ID format", err)
	}
	return CampaignID{value: id}, nil
}

func (c CampaignID) String() string {
	return c.value.String()
}

// CampaignType represents the type of campaign
type CampaignType string

const (
	CampaignTypePromotion    CampaignType = "PROMOTION"
	CampaignTypeLoyalty      CampaignType = "LOYALTY"
	CampaignTypeCoupon       CampaignType = "COUPON"
	CampaignTypeSegmentation CampaignType = "SEGMENTATION"
	CampaignTypeRetargeting  CampaignType = "RETARGETING"
)

func (ct CampaignType) String() string {
	return string(ct)
}

func ParseCampaignType(s string) (CampaignType, error) {
	switch CampaignType(s) {
	case CampaignTypePromotion, CampaignTypeLoyalty, CampaignTypeCoupon, CampaignTypeSegmentation, CampaignTypeRetargeting:
		return CampaignType(s), nil
	default:
		return "", shared.NewValidationError("invalid campaign type", fmt.Errorf("unknown type: %s", s))
	}
}

// CampaignStatus represents the status of a campaign
type CampaignStatus string

const (
	CampaignStatusDraft     CampaignStatus = "DRAFT"
	CampaignStatusActive    CampaignStatus = "ACTIVE"
	CampaignStatusPaused    CampaignStatus = "PAUSED"
	CampaignStatusCompleted CampaignStatus = "COMPLETED"
	CampaignStatusCancelled CampaignStatus = "CANCELLED"
)

func (cs CampaignStatus) String() string {
	return string(cs)
}

func ParseCampaignStatus(s string) (CampaignStatus, error) {
	switch CampaignStatus(s) {
	case CampaignStatusDraft, CampaignStatusActive, CampaignStatusPaused, CampaignStatusCompleted, CampaignStatusCancelled:
		return CampaignStatus(s), nil
	default:
		return "", shared.NewValidationError("invalid campaign status", fmt.Errorf("unknown status: %s", s))
	}
}

// Campaign represents a marketing campaign aggregate
type Campaign struct {
	id             CampaignID
	name           string
	description    string
	status         CampaignStatus
	campaignType   CampaignType
	targetingRules []shared.RuleID
	startDate      time.Time
	endDate        *time.Time
	budget         *shared.Money
	createdBy      shared.UserID
	createdAt      time.Time
	updatedAt      time.Time
	settings       CampaignSettings
	metrics        CampaignMetrics
	version        int
	events         []shared.DomainEvent
}

// NewCampaign creates a new campaign with validation
func NewCampaign(
	name, description string,
	campaignType CampaignType,
	targetingRules []shared.RuleID,
	startDate time.Time,
	endDate *time.Time,
	budget *shared.Money,
	createdBy shared.UserID,
	settings CampaignSettings,
) (*Campaign, error) {
	// Validate required fields
	if err := shared.ValidateRequired(name, "campaign name"); err != nil {
		return nil, shared.NewValidationError("invalid campaign name", err)
	}

	if err := shared.ValidateMaxLength(name, 255, "campaign name"); err != nil {
		return nil, shared.NewValidationError("invalid campaign name", err)
	}

	if err := shared.ValidateMaxLength(description, 1000, "campaign description"); err != nil {
		return nil, shared.NewValidationError("invalid campaign description", err)
	}

	// Validate dates
	if startDate.IsZero() {
		return nil, shared.NewValidationError("start date is required", nil)
	}

	if endDate != nil && endDate.Before(startDate) {
		return nil, shared.NewValidationError("end date must be after start date", nil)
	}

	// Validate budget
	if budget != nil && !budget.IsPositive() {
		return nil, shared.NewValidationError("budget must be positive", nil)
	}

	// Validate targeting rules
	if len(targetingRules) == 0 {
		return nil, shared.NewValidationError("at least one targeting rule is required", nil)
	}

	// Validate settings
	if err := settings.Validate(); err != nil {
		return nil, shared.NewValidationError("invalid campaign settings", err)
	}

	now := time.Now()
	campaign := &Campaign{
		id:             NewCampaignID(),
		name:           name,
		description:    description,
		status:         CampaignStatusDraft,
		campaignType:   campaignType,
		targetingRules: targetingRules,
		startDate:      startDate,
		endDate:        endDate,
		budget:         budget,
		createdBy:      createdBy,
		createdAt:      now,
		updatedAt:      now,
		settings:       settings,
		metrics:        NewCampaignMetrics(),
		version:        1,
		events:         make([]shared.DomainEvent, 0),
	}

	// Raise domain event
	campaign.addEvent(NewCampaignCreatedEvent(campaign.id, campaign.name, campaign.campaignType))

	return campaign, nil
}

// Getters
func (c *Campaign) ID() CampaignID {
	return c.id
}

func (c *Campaign) Name() string {
	return c.name
}

func (c *Campaign) Description() string {
	return c.description
}

func (c *Campaign) Status() CampaignStatus {
	return c.status
}

func (c *Campaign) CampaignType() CampaignType {
	return c.campaignType
}

func (c *Campaign) TargetingRules() []shared.RuleID {
	return c.targetingRules
}

func (c *Campaign) StartDate() time.Time {
	return c.startDate
}

func (c *Campaign) EndDate() *time.Time {
	return c.endDate
}

func (c *Campaign) Budget() *shared.Money {
	return c.budget
}

func (c *Campaign) CreatedBy() shared.UserID {
	return c.createdBy
}

func (c *Campaign) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Campaign) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *Campaign) Settings() CampaignSettings {
	return c.settings
}

func (c *Campaign) Metrics() CampaignMetrics {
	return c.metrics
}

func (c *Campaign) Version() int {
	return c.version
}

// Business methods
func (c *Campaign) Activate() error {
	if c.status != CampaignStatusDraft {
		return shared.NewBusinessError("only draft campaigns can be activated", nil)
	}

	// Check if start date is in the future
	if c.startDate.After(time.Now()) {
		return shared.NewBusinessError("cannot activate campaign with future start date", nil)
	}

	c.status = CampaignStatusActive
	c.updatedAt = time.Now()
	c.addEvent(NewCampaignActivatedEvent(c.id, c.name))

	return nil
}

func (c *Campaign) Pause() error {
	if c.status != CampaignStatusActive {
		return shared.NewBusinessError("only active campaigns can be paused", nil)
	}

	c.status = CampaignStatusPaused
	c.updatedAt = time.Now()
	c.addEvent(NewCampaignPausedEvent(c.id, c.name))

	return nil
}

func (c *Campaign) Resume() error {
	if c.status != CampaignStatusPaused {
		return shared.NewBusinessError("only paused campaigns can be resumed", nil)
	}

	// Check if campaign is still within date range
	now := time.Now()
	if c.endDate != nil && now.After(*c.endDate) {
		return shared.NewBusinessError("cannot resume campaign that has ended", nil)
	}

	c.status = CampaignStatusActive
	c.updatedAt = time.Now()
	c.addEvent(NewCampaignResumedEvent(c.id, c.name))

	return nil
}

func (c *Campaign) Complete() error {
	if c.status != CampaignStatusActive && c.status != CampaignStatusPaused {
		return shared.NewBusinessError("only active or paused campaigns can be completed", nil)
	}

	c.status = CampaignStatusCompleted
	c.updatedAt = time.Now()
	c.addEvent(NewCampaignCompletedEvent(c.id, c.name))

	return nil
}

func (c *Campaign) Cancel(reason string) error {
	if c.status == CampaignStatusCompleted {
		return shared.NewBusinessError("cannot cancel completed campaign", nil)
	}

	if c.status == CampaignStatusCancelled {
		return shared.NewBusinessError("campaign is already cancelled", nil)
	}

	c.status = CampaignStatusCancelled
	c.updatedAt = time.Now()
	c.addEvent(NewCampaignCancelledEvent(c.id, c.name, reason))

	return nil
}

func (c *Campaign) UpdateTargetingRules(rules []shared.RuleID) error {
	if c.status == CampaignStatusActive {
		return shared.NewBusinessError("cannot update targeting rules for active campaign", nil)
	}

	if len(rules) == 0 {
		return shared.NewValidationError("at least one targeting rule is required", nil)
	}

	c.targetingRules = rules
	c.updatedAt = time.Now()
	c.addEvent(NewCampaignTargetingRulesUpdatedEvent(c.id, c.name, rules))

	return nil
}

func (c *Campaign) UpdateBudget(budget *shared.Money) error {
	if budget != nil && !budget.IsPositive() {
		return shared.NewValidationError("budget must be positive", nil)
	}

	c.budget = budget
	c.updatedAt = time.Now()
	c.addEvent(NewCampaignBudgetUpdatedEvent(c.id, c.name, budget))

	return nil
}

func (c *Campaign) UpdateSettings(settings CampaignSettings) error {
	if err := settings.Validate(); err != nil {
		return shared.NewValidationError("invalid campaign settings", err)
	}

	c.settings = settings
	c.updatedAt = time.Now()
	c.addEvent(NewCampaignSettingsUpdatedEvent(c.id, c.name, settings))

	return nil
}

func (c *Campaign) TrackEvent(eventType CampaignEventType, customerID *shared.CustomerID, eventData map[string]interface{}) error {
	event := NewCampaignEvent(c.id, eventType, customerID, eventData)
	c.metrics.UpdateFromEvent(event)
	c.updatedAt = time.Now()
	c.addEvent(NewCampaignEventTrackedEvent(c.id, eventType, customerID))

	return nil
}

func (c *Campaign) GetEvents() []shared.DomainEvent {
	events := c.events
	c.events = make([]shared.DomainEvent, 0)
	return events
}

func (c *Campaign) addEvent(event shared.DomainEvent) {
	c.events = append(c.events, event)
}

// Check if campaign is currently active
func (c *Campaign) IsActive() bool {
	if c.status != CampaignStatusActive {
		return false
	}

	now := time.Now()
	if now.Before(c.startDate) {
		return false
	}

	if c.endDate != nil && now.After(*c.endDate) {
		return false
	}

	return true
}

// Check if campaign has exceeded budget
func (c *Campaign) HasExceededBudget() bool {
	if c.budget == nil {
		return false
	}

	return c.metrics.Cost.Amount >= c.budget.Amount
}

// Check if campaign is approaching budget limit (90%)
func (c *Campaign) IsApproachingBudgetLimit() bool {
	if c.budget == nil {
		return false
	}

	threshold := c.budget.Amount * 0.9
	return c.metrics.Cost.Amount >= threshold
}
