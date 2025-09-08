package dto

import (
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/campaign"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// CreateCampaignRequest represents the request to create a new campaign
type CreateCampaignRequest struct {
	Name           string                 `json:"name" binding:"required,min=3,max=255"`
	Description    string                 `json:"description" binding:"max=1000"`
	Type           string                 `json:"type" binding:"required,oneof=PROMOTION LOYALTY COUPON SEGMENTATION RETARGETING"`
	TargetingRules []string               `json:"targetingRules" binding:"required,min=1"`
	StartDate      time.Time              `json:"startDate" binding:"required"`
	EndDate        *time.Time             `json:"endDate,omitempty"`
	Budget         *shared.Money          `json:"budget,omitempty"`
	CreatedBy      string                 `json:"createdBy" binding:"required"`
	Settings       CreateCampaignSettings `json:"settings" binding:"required"`
}

// UpdateCampaignRequest represents the request to update an existing campaign
type UpdateCampaignRequest struct {
	Name           string                 `json:"name" binding:"required,min=3,max=255"`
	Description    string                 `json:"description" binding:"max=1000"`
	Type           string                 `json:"type" binding:"required,oneof=PROMOTION LOYALTY COUPON SEGMENTATION RETARGETING"`
	TargetingRules []string               `json:"targetingRules" binding:"required,min=1"`
	StartDate      time.Time              `json:"startDate" binding:"required"`
	EndDate        *time.Time             `json:"endDate,omitempty"`
	Budget         *shared.Money          `json:"budget,omitempty"`
	UpdatedBy      string                 `json:"updatedBy" binding:"required"`
	Settings       UpdateCampaignSettings `json:"settings" binding:"required"`
}

// ActivateCampaignRequest represents the request to activate a campaign
type ActivateCampaignRequest struct {
	ActivatedBy string `json:"activatedBy" binding:"required"`
}

// PauseCampaignRequest represents the request to pause a campaign
type PauseCampaignRequest struct {
	PausedBy string `json:"pausedBy" binding:"required"`
}

// DeleteCampaignRequest represents the request to delete a campaign
type DeleteCampaignRequest struct {
	DeletedBy string `json:"deletedBy" binding:"required"`
	Reason    string `json:"reason,omitempty" binding:"max=500"`
}

// CreateCampaignSettings represents campaign settings for creation
type CreateCampaignSettings struct {
	TargetAudience  []string                       `json:"targetAudience"`
	Channels        []string                       `json:"channels" binding:"required,min=1"`
	Frequency       string                         `json:"frequency" binding:"required,oneof=ONCE DAILY WEEKLY MONTHLY"`
	MaxImpressions  *int                           `json:"maxImpressions,omitempty"`
	BudgetLimit     *shared.Money                  `json:"budgetLimit,omitempty"`
	ABTestSettings  *campaign.ABTestSettings       `json:"abTestSettings,omitempty"`
	SchedulingRules []campaign.SchedulingRule      `json:"schedulingRules"`
	Personalization campaign.PersonalizationConfig `json:"personalization"`
}

// UpdateCampaignSettings represents campaign settings for updates
type UpdateCampaignSettings struct {
	TargetAudience  []string                       `json:"targetAudience"`
	Channels        []string                       `json:"channels" binding:"required,min=1"`
	Frequency       string                         `json:"frequency" binding:"required,oneof=ONCE DAILY WEEKLY MONTHLY"`
	MaxImpressions  *int                           `json:"maxImpressions,omitempty"`
	BudgetLimit     *shared.Money                  `json:"budgetLimit,omitempty"`
	ABTestSettings  *campaign.ABTestSettings       `json:"abTestSettings,omitempty"`
	SchedulingRules []campaign.SchedulingRule      `json:"schedulingRules"`
	Personalization campaign.PersonalizationConfig `json:"personalization"`
}

// CampaignResponse represents a campaign in API responses
type CampaignResponse struct {
	ID             string                    `json:"id"`
	Name           string                    `json:"name"`
	Description    string                    `json:"description"`
	Status         string                    `json:"status"`
	Type           string                    `json:"type"`
	TargetingRules []string                  `json:"targetingRules"`
	StartDate      time.Time                 `json:"startDate"`
	EndDate        *time.Time                `json:"endDate,omitempty"`
	Budget         *shared.Money             `json:"budget,omitempty"`
	CreatedBy      string                    `json:"createdBy"`
	CreatedAt      time.Time                 `json:"createdAt"`
	UpdatedAt      time.Time                 `json:"updatedAt"`
	Settings       campaign.CampaignSettings `json:"settings"`
	Metrics        campaign.CampaignMetrics  `json:"metrics"`
	Version        int                       `json:"version"`
}

// CampaignMetricsResponse represents campaign metrics in API responses
type CampaignMetricsResponse struct {
	Impressions       int64        `json:"impressions"`
	Clicks            int64        `json:"clicks"`
	Conversions       int64        `json:"conversions"`
	Revenue           shared.Money `json:"revenue"`
	Cost              shared.Money `json:"cost"`
	CTR               float64      `json:"ctr"`
	ConversionRate    float64      `json:"conversionRate"`
	CostPerClick      shared.Money `json:"costPerClick"`
	CostPerConversion shared.Money `json:"costPerConversion"`
	ROAS              float64      `json:"roas"`
	ROI               float64      `json:"roi"`
	LastUpdated       time.Time    `json:"lastUpdated"`
}

// ListCampaignsResponse represents the response for listing campaigns
type ListCampaignsResponse struct {
	Campaigns  []*CampaignResponse `json:"campaigns"`
	Total      int                 `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"totalPages"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a successful response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
