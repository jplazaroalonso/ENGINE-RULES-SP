package commands

import (
	"context"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/campaign"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// CreateCampaignCommand represents the command to create a new campaign
type CreateCampaignCommand struct {
	Name           string                  `json:"name" validate:"required,min=3,max=255"`
	Description    string                  `json:"description" validate:"max=1000"`
	CampaignType   string                  `json:"campaignType" validate:"required,oneof=PROMOTION LOYALTY COUPON SEGMENTATION RETARGETING"`
	TargetingRules []string                `json:"targetingRules" validate:"required,min=1"`
	StartDate      time.Time               `json:"startDate" validate:"required"`
	EndDate        *time.Time              `json:"endDate,omitempty"`
	Budget         *BudgetRequest          `json:"budget,omitempty"`
	Settings       CampaignSettingsRequest `json:"settings" validate:"required"`
	CreatedBy      string                  `json:"createdBy" validate:"required"`
}

// BudgetRequest represents budget information in the request
type BudgetRequest struct {
	Amount   float64 `json:"amount" validate:"required,gt=0"`
	Currency string  `json:"currency" validate:"required,len=3"`
}

// CampaignSettingsRequest represents campaign settings in the request
type CampaignSettingsRequest struct {
	TargetAudience  []string                     `json:"targetAudience" validate:"required,min=1"`
	Channels        []string                     `json:"channels" validate:"required,min=1"`
	Frequency       string                       `json:"frequency" validate:"required,oneof=ONCE DAILY WEEKLY MONTHLY"`
	MaxImpressions  *int                         `json:"maxImpressions,omitempty" validate:"omitempty,gt=0"`
	BudgetLimit     *BudgetRequest               `json:"budgetLimit,omitempty"`
	ABTestSettings  *ABTestSettingsRequest       `json:"abTestSettings,omitempty"`
	SchedulingRules []SchedulingRuleRequest      `json:"schedulingRules,omitempty"`
	Personalization PersonalizationConfigRequest `json:"personalization" validate:"required"`
}

// ABTestSettingsRequest represents A/B test settings in the request
type ABTestSettingsRequest struct {
	Enabled       bool             `json:"enabled"`
	Variants      []VariantRequest `json:"variants" validate:"omitempty,min=2,max=10"`
	TrafficSplit  float64          `json:"trafficSplit" validate:"omitempty,min=0,max=1"`
	SuccessMetric string           `json:"successMetric" validate:"omitempty,min=1"`
	Duration      int              `json:"duration" validate:"omitempty,gt=0"`
}

// VariantRequest represents a variant in the request
type VariantRequest struct {
	ID          string                 `json:"id" validate:"required"`
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description"`
	Settings    map[string]interface{} `json:"settings"`
	Weight      float64                `json:"weight" validate:"required,min=0,max=1"`
}

// SchedulingRuleRequest represents a scheduling rule in the request
type SchedulingRuleRequest struct {
	ID          string                       `json:"id" validate:"required"`
	Name        string                       `json:"name" validate:"required"`
	Description string                       `json:"description"`
	Conditions  []SchedulingConditionRequest `json:"conditions" validate:"required,min=1"`
	Actions     []SchedulingActionRequest    `json:"actions" validate:"required,min=1"`
	IsActive    bool                         `json:"isActive"`
}

// SchedulingConditionRequest represents a scheduling condition in the request
type SchedulingConditionRequest struct {
	Type     string                 `json:"type" validate:"required,oneof=time date event metric"`
	Operator string                 `json:"operator" validate:"required,oneof=equals greater_than less_than contains not_equals"`
	Value    interface{}            `json:"value" validate:"required"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// SchedulingActionRequest represents a scheduling action in the request
type SchedulingActionRequest struct {
	Type       string                 `json:"type" validate:"required,oneof=activate pause stop update_settings send_notification"`
	Parameters map[string]interface{} `json:"parameters"`
}

// PersonalizationConfigRequest represents personalization config in the request
type PersonalizationConfigRequest struct {
	Enabled     bool     `json:"enabled"`
	Rules       []string `json:"rules"`
	Fallback    string   `json:"fallback"`
	MaxVariants int      `json:"maxVariants" validate:"omitempty,min=1,max=100"`
}

// CreateCampaignResult represents the result of creating a campaign
type CreateCampaignResult struct {
	CampaignID   string    `json:"campaignId"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	CampaignType string    `json:"campaignType"`
	Version      int       `json:"version"`
	CreatedAt    time.Time `json:"createdAt"`
}

// CreateCampaignHandler handles campaign creation commands
type CreateCampaignHandler struct {
	campaignService *campaign.CampaignService
	validator       shared.Validator
}

// NewCreateCampaignHandler creates a new CreateCampaignHandler
func NewCreateCampaignHandler(
	campaignService *campaign.CampaignService,
	validator shared.Validator,
) *CreateCampaignHandler {
	return &CreateCampaignHandler{
		campaignService: campaignService,
		validator:       validator,
	}
}

// Handle processes the create campaign command
func (h *CreateCampaignHandler) Handle(ctx context.Context, cmd CreateCampaignCommand) (*CreateCampaignResult, error) {
	// Validate command input
	if err := h.validator.Validate(cmd); err != nil {
		return nil, shared.NewValidationError("invalid create campaign command", err)
	}

	// Parse campaign type
	campaignType, err := campaign.ParseCampaignType(cmd.CampaignType)
	if err != nil {
		return nil, shared.NewValidationError("invalid campaign type", err)
	}

	// Parse targeting rules
	targetingRules := make([]shared.RuleID, len(cmd.TargetingRules))
	for i, ruleID := range cmd.TargetingRules {
		parsedRuleID, err := shared.NewRuleIDFromString(ruleID)
		if err != nil {
			return nil, shared.NewValidationError("invalid rule ID", err)
		}
		targetingRules[i] = parsedRuleID
	}

	// Parse budget
	var budget *shared.Money
	if cmd.Budget != nil {
		budget = &shared.Money{
			Amount:   cmd.Budget.Amount,
			Currency: cmd.Budget.Currency,
		}
	}

	// Parse settings
	settings, err := h.parseCampaignSettings(cmd.Settings)
	if err != nil {
		return nil, shared.NewValidationError("invalid campaign settings", err)
	}

	// Parse user ID
	userID, err := shared.NewUserIDFromString(cmd.CreatedBy)
	if err != nil {
		return nil, shared.NewValidationError("invalid user ID", err)
	}

	// Create campaign using domain service
	createdCampaign, err := h.campaignService.CreateCampaign(
		ctx,
		cmd.Name,
		cmd.Description,
		campaignType,
		targetingRules,
		cmd.StartDate,
		cmd.EndDate,
		budget,
		userID,
		settings,
	)
	if err != nil {
		return nil, err
	}

	return &CreateCampaignResult{
		CampaignID:   createdCampaign.ID().String(),
		Name:         createdCampaign.Name(),
		Status:       createdCampaign.Status().String(),
		CampaignType: createdCampaign.CampaignType().String(),
		Version:      createdCampaign.Version(),
		CreatedAt:    createdCampaign.CreatedAt(),
	}, nil
}

// parseCampaignSettings parses campaign settings from request
func (h *CreateCampaignHandler) parseCampaignSettings(req CampaignSettingsRequest) (campaign.CampaignSettings, error) {
	// Parse channels
	channels := make([]campaign.Channel, len(req.Channels))
	for i, channelStr := range req.Channels {
		channel, err := campaign.ParseChannel(channelStr)
		if err != nil {
			return campaign.CampaignSettings{}, err
		}
		channels[i] = channel
	}

	// Parse frequency
	frequency, err := campaign.ParseFrequency(req.Frequency)
	if err != nil {
		return campaign.CampaignSettings{}, err
	}

	// Parse budget limit
	var budgetLimit *shared.Money
	if req.BudgetLimit != nil {
		budgetLimit = &shared.Money{
			Amount:   req.BudgetLimit.Amount,
			Currency: req.BudgetLimit.Currency,
		}
	}

	// Parse A/B test settings
	var abTestSettings *campaign.ABTestSettings
	if req.ABTestSettings != nil {
		abTestSettings, err = h.parseABTestSettings(*req.ABTestSettings)
		if err != nil {
			return campaign.CampaignSettings{}, err
		}
	}

	// Parse scheduling rules
	schedulingRules := make([]campaign.SchedulingRule, len(req.SchedulingRules))
	for i, ruleReq := range req.SchedulingRules {
		rule, err := h.parseSchedulingRule(ruleReq)
		if err != nil {
			return campaign.CampaignSettings{}, err
		}
		schedulingRules[i] = rule
	}

	// Parse personalization config
	personalization := campaign.PersonalizationConfig{
		Enabled:     req.Personalization.Enabled,
		Rules:       req.Personalization.Rules,
		Fallback:    req.Personalization.Fallback,
		MaxVariants: req.Personalization.MaxVariants,
	}

	return campaign.NewCampaignSettings(
		req.TargetAudience,
		channels,
		frequency,
		req.MaxImpressions,
		budgetLimit,
		abTestSettings,
		schedulingRules,
		personalization,
	)
}

// parseABTestSettings parses A/B test settings from request
func (h *CreateCampaignHandler) parseABTestSettings(req ABTestSettingsRequest) (*campaign.ABTestSettings, error) {
	variants := make([]campaign.Variant, len(req.Variants))
	for i, variantReq := range req.Variants {
		variants[i] = campaign.Variant{
			ID:          variantReq.ID,
			Name:        variantReq.Name,
			Description: variantReq.Description,
			Settings:    variantReq.Settings,
			Weight:      variantReq.Weight,
		}
	}

	return &campaign.ABTestSettings{
		Enabled:       req.Enabled,
		Variants:      variants,
		TrafficSplit:  req.TrafficSplit,
		SuccessMetric: req.SuccessMetric,
		Duration:      req.Duration,
	}, nil
}

// parseSchedulingRule parses a scheduling rule from request
func (h *CreateCampaignHandler) parseSchedulingRule(req SchedulingRuleRequest) (campaign.SchedulingRule, error) {
	conditions := make([]campaign.SchedulingCondition, len(req.Conditions))
	for i, conditionReq := range req.Conditions {
		conditions[i] = campaign.SchedulingCondition{
			Type:     conditionReq.Type,
			Operator: conditionReq.Operator,
			Value:    conditionReq.Value,
			Metadata: conditionReq.Metadata,
		}
	}

	actions := make([]campaign.SchedulingAction, len(req.Actions))
	for i, actionReq := range req.Actions {
		actions[i] = campaign.SchedulingAction{
			Type:       actionReq.Type,
			Parameters: actionReq.Parameters,
		}
	}

	return campaign.SchedulingRule{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Conditions:  conditions,
		Actions:     actions,
		IsActive:    req.IsActive,
	}, nil
}
