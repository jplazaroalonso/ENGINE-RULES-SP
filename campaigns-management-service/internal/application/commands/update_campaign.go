package commands

import (
	"context"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/campaign"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// UpdateCampaignCommand represents the command to update an existing campaign
type UpdateCampaignCommand struct {
	CampaignID     string                   `json:"campaignId" validate:"required"`
	Name           *string                  `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description    *string                  `json:"description,omitempty" validate:"omitempty,max=1000"`
	TargetingRules []string                 `json:"targetingRules,omitempty" validate:"omitempty,min=1"`
	StartDate      *time.Time               `json:"startDate,omitempty"`
	EndDate        *time.Time               `json:"endDate,omitempty"`
	Budget         *BudgetRequest           `json:"budget,omitempty"`
	Settings       *CampaignSettingsRequest `json:"settings,omitempty"`
	UpdatedBy      string                   `json:"updatedBy" validate:"required"`
}

// UpdateCampaignResult represents the result of updating a campaign
type UpdateCampaignResult struct {
	CampaignID string    `json:"campaignId"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	Version    int       `json:"version"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// UpdateCampaignHandler handles campaign update commands
type UpdateCampaignHandler struct {
	campaignService *campaign.CampaignService
	validator       shared.Validator
}

// NewUpdateCampaignHandler creates a new UpdateCampaignHandler
func NewUpdateCampaignHandler(
	campaignService *campaign.CampaignService,
	validator shared.Validator,
) *UpdateCampaignHandler {
	return &UpdateCampaignHandler{
		campaignService: campaignService,
		validator:       validator,
	}
}

// Handle processes the update campaign command
func (h *UpdateCampaignHandler) Handle(ctx context.Context, cmd UpdateCampaignCommand) (*UpdateCampaignResult, error) {
	// Validate command input
	if err := h.validator.Validate(cmd); err != nil {
		return nil, shared.NewValidationError("invalid update campaign command", err)
	}

	// Parse campaign ID
	campaignID, err := campaign.NewCampaignIDFromString(cmd.CampaignID)
	if err != nil {
		return nil, shared.NewValidationError("invalid campaign ID", err)
	}

	// Get existing campaign
	existingCampaign, err := h.campaignService.GetCampaign(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	// Check if campaign can be updated
	if existingCampaign.Status() == campaign.CampaignStatusActive {
		return nil, shared.NewBusinessError("cannot update active campaign", nil)
	}

	// Create updated campaign with new values
	updatedCampaign := *existingCampaign

	// Update name if provided
	if cmd.Name != nil {
		// Check if new name already exists (if different from current)
		if *cmd.Name != existingCampaign.Name() {
			exists, err := h.campaignService.ExistsByName(ctx, *cmd.Name)
			if err != nil {
				return nil, shared.NewInfrastructureError("failed to check campaign name existence", err)
			}
			if exists {
				return nil, shared.NewBusinessError("campaign name already exists", nil)
			}
		}
		// Note: In a real implementation, you would need a method to update the name
		// For now, we'll assume the campaign has an UpdateName method
	}

	// Update description if provided
	if cmd.Description != nil {
		// Note: In a real implementation, you would need a method to update the description
	}

	// Update targeting rules if provided
	if cmd.TargetingRules != nil {
		targetingRules := make([]shared.RuleID, len(cmd.TargetingRules))
		for i, ruleID := range cmd.TargetingRules {
			parsedRuleID, err := shared.NewRuleIDFromString(ruleID)
			if err != nil {
				return nil, shared.NewValidationError("invalid rule ID", err)
			}
			targetingRules[i] = parsedRuleID
		}

		if err := updatedCampaign.UpdateTargetingRules(targetingRules); err != nil {
			return nil, err
		}
	}

	// Update budget if provided
	if cmd.Budget != nil {
		budget := &shared.Money{
			Amount:   cmd.Budget.Amount,
			Currency: cmd.Budget.Currency,
		}

		if err := updatedCampaign.UpdateBudget(budget); err != nil {
			return nil, err
		}
	}

	// Update settings if provided
	if cmd.Settings != nil {
		settings, err := h.parseCampaignSettings(*cmd.Settings)
		if err != nil {
			return nil, shared.NewValidationError("invalid campaign settings", err)
		}

		if err := updatedCampaign.UpdateSettings(settings); err != nil {
			return nil, err
		}
	}

	// Save updated campaign
	if err := h.campaignService.SaveCampaign(ctx, &updatedCampaign); err != nil {
		return nil, shared.NewInfrastructureError("failed to save updated campaign", err)
	}

	return &UpdateCampaignResult{
		CampaignID: updatedCampaign.ID().String(),
		Name:       updatedCampaign.Name(),
		Status:     updatedCampaign.Status().String(),
		Version:    updatedCampaign.Version(),
		UpdatedAt:  updatedCampaign.UpdatedAt(),
	}, nil
}

// parseCampaignSettings parses campaign settings from request
func (h *UpdateCampaignHandler) parseCampaignSettings(req CampaignSettingsRequest) (campaign.CampaignSettings, error) {
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
func (h *UpdateCampaignHandler) parseABTestSettings(req ABTestSettingsRequest) (*campaign.ABTestSettings, error) {
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
func (h *UpdateCampaignHandler) parseSchedulingRule(req SchedulingRuleRequest) (campaign.SchedulingRule, error) {
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
