package queries

import (
	"context"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/campaign"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// ListCampaignsQuery represents the query to list campaigns with filtering and pagination
type ListCampaignsQuery struct {
	Status       *string    `json:"status,omitempty" validate:"omitempty,oneof=DRAFT ACTIVE PAUSED COMPLETED CANCELLED"`
	CampaignType *string    `json:"campaignType,omitempty" validate:"omitempty,oneof=PROMOTION LOYALTY COUPON SEGMENTATION RETARGETING"`
	CreatedBy    *string    `json:"createdBy,omitempty"`
	StartDate    *time.Time `json:"startDate,omitempty"`
	EndDate      *time.Time `json:"endDate,omitempty"`
	SearchQuery  string     `json:"searchQuery,omitempty" validate:"omitempty,max=100"`
	Page         int        `json:"page" validate:"min=1"`
	Limit        int        `json:"limit" validate:"min=1,max=100"`
	SortBy       string     `json:"sortBy,omitempty" validate:"omitempty,oneof=name status createdAt updatedAt"`
	SortOrder    string     `json:"sortOrder,omitempty" validate:"omitempty,oneof=asc desc"`
}

// ListCampaignsResult represents the result of listing campaigns
type ListCampaignsResult struct {
	Campaigns  []CampaignSummary `json:"campaigns"`
	Pagination PaginationInfo    `json:"pagination"`
}

// CampaignSummary represents a summary of a campaign for listing
type CampaignSummary struct {
	CampaignID   string                 `json:"campaignId"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Status       string                 `json:"status"`
	CampaignType string                 `json:"campaignType"`
	StartDate    string                 `json:"startDate"`
	EndDate      *string                `json:"endDate,omitempty"`
	Budget       *BudgetResponse        `json:"budget,omitempty"`
	CreatedBy    string                 `json:"createdBy"`
	CreatedAt    string                 `json:"createdAt"`
	UpdatedAt    string                 `json:"updatedAt"`
	Metrics      CampaignMetricsSummary `json:"metrics"`
	Version      int                    `json:"version"`
}

// CampaignMetricsSummary represents a summary of campaign metrics
type CampaignMetricsSummary struct {
	Impressions    int64   `json:"impressions"`
	Clicks         int64   `json:"clicks"`
	Conversions    int64   `json:"conversions"`
	CTR            float64 `json:"ctr"`
	ConversionRate float64 `json:"conversionRate"`
	ROI            float64 `json:"roi"`
}

// PaginationInfo represents pagination information
type PaginationInfo struct {
	Page        int   `json:"page"`
	Limit       int   `json:"limit"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"totalPages"`
	HasNext     bool  `json:"hasNext"`
	HasPrevious bool  `json:"hasPrevious"`
}

// ListCampaignsHandler handles list campaigns queries
type ListCampaignsHandler struct {
	campaignService *campaign.CampaignService
	validator       shared.Validator
}

// NewListCampaignsHandler creates a new ListCampaignsHandler
func NewListCampaignsHandler(
	campaignService *campaign.CampaignService,
	validator shared.Validator,
) *ListCampaignsHandler {
	return &ListCampaignsHandler{
		campaignService: campaignService,
		validator:       validator,
	}
}

// Handle processes the list campaigns query
func (h *ListCampaignsHandler) Handle(ctx context.Context, query ListCampaignsQuery) (*ListCampaignsResult, error) {
	// Set default values
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 20
	}
	if query.SortBy == "" {
		query.SortBy = "createdAt"
	}
	if query.SortOrder == "" {
		query.SortOrder = "desc"
	}

	// Validate query input
	if err := h.validator.Validate(query); err != nil {
		return nil, shared.NewValidationError("invalid list campaigns query", err)
	}

	// Build list criteria
	criteria := h.buildListCriteria(query)

	// Get campaigns using domain service
	campaigns, err := h.campaignService.ListCampaigns(ctx, criteria)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := h.campaignService.CountCampaigns(ctx, criteria)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	campaignSummaries := make([]CampaignSummary, len(campaigns))
	for i, campaign := range campaigns {
		campaignSummaries[i] = h.mapCampaignToSummary(campaign)
	}

	// Calculate pagination info
	totalPages := int((total + int64(query.Limit) - 1) / int64(query.Limit))
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	pagination := PaginationInfo{
		Page:        query.Page,
		Limit:       query.Limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     query.Page < totalPages,
		HasPrevious: query.Page > 1,
	}

	return &ListCampaignsResult{
		Campaigns:  campaignSummaries,
		Pagination: pagination,
	}, nil
}

// buildListCriteria builds list criteria from query
func (h *ListCampaignsHandler) buildListCriteria(query ListCampaignsQuery) campaign.ListCriteria {
	criteria := campaign.ListCriteria{
		PageSize:   query.Limit,
		PageOffset: (query.Page - 1) * query.Limit,
		SortBy:     query.SortBy,
		SortOrder:  query.SortOrder,
	}

	// Parse status if provided
	if query.Status != nil {
		if status, err := campaign.ParseCampaignStatus(*query.Status); err == nil {
			criteria.Status = &status
		}
	}

	// Parse campaign type if provided
	if query.CampaignType != nil {
		if campaignType, err := campaign.ParseCampaignType(*query.CampaignType); err == nil {
			criteria.CampaignType = &campaignType
		}
	}

	// Set created by if provided
	if query.CreatedBy != nil {
		createdBy, err := shared.NewUserIDFromString(*query.CreatedBy)
		if err != nil {
			// For now, skip invalid user IDs in criteria building
			// In a real implementation, you might want to return an error
		} else {
			criteria.CreatedBy = &createdBy
		}
	}

	// Set date range if provided
	if query.StartDate != nil {
		criteria.StartDate = query.StartDate
	}
	if query.EndDate != nil {
		criteria.EndDate = query.EndDate
	}

	// Set search query if provided
	if query.SearchQuery != "" {
		criteria.SearchQuery = query.SearchQuery
	}

	return criteria
}

// mapCampaignToSummary maps a campaign domain entity to summary format
func (h *ListCampaignsHandler) mapCampaignToSummary(campaign *campaign.Campaign) CampaignSummary {
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

	// Map metrics summary
	metrics := campaign.Metrics()
	metricsSummary := CampaignMetricsSummary{
		Impressions:    metrics.Impressions,
		Clicks:         metrics.Clicks,
		Conversions:    metrics.Conversions,
		CTR:            metrics.CTR,
		ConversionRate: metrics.ConversionRate,
		ROI:            metrics.ROI,
	}

	return CampaignSummary{
		CampaignID:   campaign.ID().String(),
		Name:         campaign.Name(),
		Description:  campaign.Description(),
		Status:       campaign.Status().String(),
		CampaignType: campaign.CampaignType().String(),
		StartDate:    campaign.StartDate().Format("2006-01-02T15:04:05Z"),
		EndDate:      endDate,
		Budget:       budget,
		CreatedBy:    campaign.CreatedBy().String(),
		CreatedAt:    campaign.CreatedAt().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    campaign.UpdatedAt().Format("2006-01-02T15:04:05Z"),
		Metrics:      metricsSummary,
		Version:      campaign.Version(),
	}
}
