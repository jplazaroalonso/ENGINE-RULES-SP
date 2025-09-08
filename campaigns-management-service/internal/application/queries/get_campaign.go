package queries

import (
	"context"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/campaign"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// GetCampaignQuery represents the query to get a campaign by ID
type GetCampaignQuery struct {
	CampaignID string `json:"campaignId" validate:"required"`
}

// GetCampaignResult represents the result of getting a campaign
type GetCampaignResult struct {
	CampaignID     string                   `json:"campaignId"`
	Name           string                   `json:"name"`
	Description    string                   `json:"description"`
	Status         string                   `json:"status"`
	CampaignType   string                   `json:"campaignType"`
	TargetingRules []string                 `json:"targetingRules"`
	StartDate      string                   `json:"startDate"`
	EndDate        *string                  `json:"endDate,omitempty"`
	Budget         *BudgetResponse          `json:"budget,omitempty"`
	CreatedBy      string                   `json:"createdBy"`
	CreatedAt      string                   `json:"createdAt"`
	UpdatedAt      string                   `json:"updatedAt"`
	Settings       CampaignSettingsResponse `json:"settings"`
	Metrics        CampaignMetricsResponse  `json:"metrics"`
	Version        int                      `json:"version"`
}

// BudgetResponse represents budget information in the response
type BudgetResponse struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// CampaignSettingsResponse represents campaign settings in the response
type CampaignSettingsResponse struct {
	TargetAudience  []string                      `json:"targetAudience"`
	Channels        []string                      `json:"channels"`
	Frequency       string                        `json:"frequency"`
	MaxImpressions  *int                          `json:"maxImpressions,omitempty"`
	BudgetLimit     *BudgetResponse               `json:"budgetLimit,omitempty"`
	ABTestSettings  *ABTestSettingsResponse       `json:"abTestSettings,omitempty"`
	SchedulingRules []SchedulingRuleResponse      `json:"schedulingRules"`
	Personalization PersonalizationConfigResponse `json:"personalization"`
}

// ABTestSettingsResponse represents A/B test settings in the response
type ABTestSettingsResponse struct {
	Enabled       bool              `json:"enabled"`
	Variants      []VariantResponse `json:"variants"`
	TrafficSplit  float64           `json:"trafficSplit"`
	SuccessMetric string            `json:"successMetric"`
	Duration      int               `json:"duration"`
}

// VariantResponse represents a variant in the response
type VariantResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Settings    map[string]interface{} `json:"settings"`
	Weight      float64                `json:"weight"`
}

// SchedulingRuleResponse represents a scheduling rule in the response
type SchedulingRuleResponse struct {
	ID          string                        `json:"id"`
	Name        string                        `json:"name"`
	Description string                        `json:"description"`
	Conditions  []SchedulingConditionResponse `json:"conditions"`
	Actions     []SchedulingActionResponse    `json:"actions"`
	IsActive    bool                          `json:"isActive"`
}

// SchedulingConditionResponse represents a scheduling condition in the response
type SchedulingConditionResponse struct {
	Type     string                 `json:"type"`
	Operator string                 `json:"operator"`
	Value    interface{}            `json:"value"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// SchedulingActionResponse represents a scheduling action in the response
type SchedulingActionResponse struct {
	Type       string                 `json:"type"`
	Parameters map[string]interface{} `json:"parameters"`
}

// PersonalizationConfigResponse represents personalization config in the response
type PersonalizationConfigResponse struct {
	Enabled     bool     `json:"enabled"`
	Rules       []string `json:"rules"`
	Fallback    string   `json:"fallback"`
	MaxVariants int      `json:"maxVariants"`
}

// CampaignMetricsResponse represents campaign metrics in the response
type CampaignMetricsResponse struct {
	Impressions       int64          `json:"impressions"`
	Clicks            int64          `json:"clicks"`
	Conversions       int64          `json:"conversions"`
	Revenue           BudgetResponse `json:"revenue"`
	Cost              BudgetResponse `json:"cost"`
	CTR               float64        `json:"ctr"`
	ConversionRate    float64        `json:"conversionRate"`
	CostPerClick      BudgetResponse `json:"costPerClick"`
	CostPerConversion BudgetResponse `json:"costPerConversion"`
	ROAS              float64        `json:"roas"`
	ROI               float64        `json:"roi"`
	LastUpdated       string         `json:"lastUpdated"`
}

// GetCampaignHandler handles get campaign queries
type GetCampaignHandler struct {
	campaignService *campaign.CampaignService
	validator       shared.Validator
}

// NewGetCampaignHandler creates a new GetCampaignHandler
func NewGetCampaignHandler(
	campaignService *campaign.CampaignService,
	validator shared.Validator,
) *GetCampaignHandler {
	return &GetCampaignHandler{
		campaignService: campaignService,
		validator:       validator,
	}
}

// Handle processes the get campaign query
func (h *GetCampaignHandler) Handle(ctx context.Context, query GetCampaignQuery) (*GetCampaignResult, error) {
	// Validate query input
	if err := h.validator.Validate(query); err != nil {
		return nil, shared.NewValidationError("invalid get campaign query", err)
	}

	// Parse campaign ID
	campaignID, err := campaign.NewCampaignIDFromString(query.CampaignID)
	if err != nil {
		return nil, shared.NewValidationError("invalid campaign ID", err)
	}

	// Get campaign using domain service
	campaign, err := h.campaignService.GetCampaign(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	result := h.mapCampaignToResponse(campaign)

	return result, nil
}

// mapCampaignToResponse maps a campaign domain entity to response format
func (h *GetCampaignHandler) mapCampaignToResponse(campaign *campaign.Campaign) *GetCampaignResult {
	// Map targeting rules
	targetingRules := make([]string, len(campaign.TargetingRules()))
	for i, rule := range campaign.TargetingRules() {
		targetingRules[i] = rule.String()
	}

	// Map budget
	var budget *BudgetResponse
	if campaign.Budget() != nil {
		budget = &BudgetResponse{
			Amount:   campaign.Budget().Amount,
			Currency: campaign.Budget().Currency,
		}
	}

	// Map end date
	var endDate *string
	if campaign.EndDate() != nil {
		endDateStr := campaign.EndDate().Format("2006-01-02T15:04:05Z")
		endDate = &endDateStr
	}

	// Map settings
	settings := h.mapSettingsToResponse(campaign.Settings())

	// Map metrics
	metrics := h.mapMetricsToResponse(campaign.Metrics())

	return &GetCampaignResult{
		CampaignID:     campaign.ID().String(),
		Name:           campaign.Name(),
		Description:    campaign.Description(),
		Status:         campaign.Status().String(),
		CampaignType:   campaign.CampaignType().String(),
		TargetingRules: targetingRules,
		StartDate:      campaign.StartDate().Format("2006-01-02T15:04:05Z"),
		EndDate:        endDate,
		Budget:         budget,
		CreatedBy:      campaign.CreatedBy().String(),
		CreatedAt:      campaign.CreatedAt().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      campaign.UpdatedAt().Format("2006-01-02T15:04:05Z"),
		Settings:       settings,
		Metrics:        metrics,
		Version:        campaign.Version(),
	}
}

// mapSettingsToResponse maps campaign settings to response format
func (h *GetCampaignHandler) mapSettingsToResponse(settings campaign.CampaignSettings) CampaignSettingsResponse {
	// Map channels
	channels := make([]string, len(settings.Channels))
	for i, channel := range settings.Channels {
		channels[i] = channel.String()
	}

	// Map budget limit
	var budgetLimit *BudgetResponse
	if settings.BudgetLimit != nil {
		budgetLimit = &BudgetResponse{
			Amount:   settings.BudgetLimit.Amount,
			Currency: settings.BudgetLimit.Currency,
		}
	}

	// Map A/B test settings
	var abTestSettings *ABTestSettingsResponse
	if settings.ABTestSettings != nil {
		variants := make([]VariantResponse, len(settings.ABTestSettings.Variants))
		for i, variant := range settings.ABTestSettings.Variants {
			variants[i] = VariantResponse{
				ID:          variant.ID,
				Name:        variant.Name,
				Description: variant.Description,
				Settings:    variant.Settings,
				Weight:      variant.Weight,
			}
		}

		abTestSettings = &ABTestSettingsResponse{
			Enabled:       settings.ABTestSettings.Enabled,
			Variants:      variants,
			TrafficSplit:  settings.ABTestSettings.TrafficSplit,
			SuccessMetric: settings.ABTestSettings.SuccessMetric,
			Duration:      settings.ABTestSettings.Duration,
		}
	}

	// Map scheduling rules
	schedulingRules := make([]SchedulingRuleResponse, len(settings.SchedulingRules))
	for i, rule := range settings.SchedulingRules {
		conditions := make([]SchedulingConditionResponse, len(rule.Conditions))
		for j, condition := range rule.Conditions {
			conditions[j] = SchedulingConditionResponse{
				Type:     condition.Type,
				Operator: condition.Operator,
				Value:    condition.Value,
				Metadata: condition.Metadata,
			}
		}

		actions := make([]SchedulingActionResponse, len(rule.Actions))
		for j, action := range rule.Actions {
			actions[j] = SchedulingActionResponse{
				Type:       action.Type,
				Parameters: action.Parameters,
			}
		}

		schedulingRules[i] = SchedulingRuleResponse{
			ID:          rule.ID,
			Name:        rule.Name,
			Description: rule.Description,
			Conditions:  conditions,
			Actions:     actions,
			IsActive:    rule.IsActive,
		}
	}

	return CampaignSettingsResponse{
		TargetAudience:  settings.TargetAudience,
		Channels:        channels,
		Frequency:       settings.Frequency.String(),
		MaxImpressions:  settings.MaxImpressions,
		BudgetLimit:     budgetLimit,
		ABTestSettings:  abTestSettings,
		SchedulingRules: schedulingRules,
		Personalization: PersonalizationConfigResponse{
			Enabled:     settings.Personalization.Enabled,
			Rules:       settings.Personalization.Rules,
			Fallback:    settings.Personalization.Fallback,
			MaxVariants: settings.Personalization.MaxVariants,
		},
	}
}

// mapMetricsToResponse maps campaign metrics to response format
func (h *GetCampaignHandler) mapMetricsToResponse(metrics campaign.CampaignMetrics) CampaignMetricsResponse {
	return CampaignMetricsResponse{
		Impressions: metrics.Impressions,
		Clicks:      metrics.Clicks,
		Conversions: metrics.Conversions,
		Revenue: BudgetResponse{
			Amount:   metrics.Revenue.Amount,
			Currency: metrics.Revenue.Currency,
		},
		Cost: BudgetResponse{
			Amount:   metrics.Cost.Amount,
			Currency: metrics.Cost.Currency,
		},
		CTR:            metrics.CTR,
		ConversionRate: metrics.ConversionRate,
		CostPerClick: BudgetResponse{
			Amount:   metrics.CostPerClick.Amount,
			Currency: metrics.CostPerClick.Currency,
		},
		CostPerConversion: BudgetResponse{
			Amount:   metrics.CostPerConversion.Amount,
			Currency: metrics.CostPerConversion.Currency,
		},
		ROAS:        metrics.ROAS,
		ROI:         metrics.ROI,
		LastUpdated: metrics.LastUpdated.Format("2006-01-02T15:04:05Z"),
	}
}
