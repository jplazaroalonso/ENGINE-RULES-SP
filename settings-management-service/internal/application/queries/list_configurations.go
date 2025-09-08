package queries

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// ListConfigurationsQuery represents a query to list configurations with pagination and filtering
type ListConfigurationsQuery struct {
	Page           int                    `json:"page" validate:"min=1"`
	Limit          int                    `json:"limit" validate:"min=1,max=1000"`
	SortBy         string                 `json:"sortBy"`
	SortOrder      string                 `json:"sortOrder" validate:"omitempty,oneof=asc desc"`
	Environment    *string                `json:"environment,omitempty" validate:"omitempty,oneof=DEVELOPMENT STAGING PRODUCTION"`
	OrganizationID *string                `json:"organizationId,omitempty" validate:"omitempty,uuid"`
	Service        *string                `json:"service,omitempty" validate:"omitempty,min=1,max=100"`
	Category       *string                `json:"category,omitempty" validate:"omitempty,min=1,max=100"`
	Tags           []string               `json:"tags,omitempty"`
	Filters        map[string]interface{} `json:"filters,omitempty"`
}

// ListConfigurationsResult represents the result of listing configurations
type ListConfigurationsResult struct {
	Configurations []*settings.Configuration `json:"configurations"`
	Total          int                       `json:"total"`
	Page           int                       `json:"page"`
	Limit          int                       `json:"limit"`
	TotalPages     int                       `json:"totalPages"`
}

// ListConfigurationsHandler handles the ListConfigurationsQuery
type ListConfigurationsHandler struct {
	configRepo settings.ConfigurationRepository
	validator  shared.StructValidator
}

// NewListConfigurationsHandler creates a new ListConfigurationsHandler
func NewListConfigurationsHandler(
	configRepo settings.ConfigurationRepository,
	validator shared.StructValidator,
) *ListConfigurationsHandler {
	return &ListConfigurationsHandler{
		configRepo: configRepo,
		validator:  validator,
	}
}

// Handle executes the ListConfigurationsQuery
func (h *ListConfigurationsHandler) Handle(ctx context.Context, query ListConfigurationsQuery) (*ListConfigurationsResult, error) {
	// Validate the query
	if err := h.validator.Struct(query); err != nil {
		return nil, fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Set defaults
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 10
	}
	if query.SortBy == "" {
		query.SortBy = "created_at"
	}
	if query.SortOrder == "" {
		query.SortOrder = "desc"
	}

	// Build filters
	filters := make(settings.ListFilters)
	if query.Environment != nil {
		filters["environment"] = *query.Environment
	}
	if query.OrganizationID != nil {
		filters["organization_id"] = *query.OrganizationID
	}
	if query.Service != nil {
		filters["service"] = *query.Service
	}
	if query.Category != nil {
		filters["category"] = *query.Category
	}
	if len(query.Tags) > 0 {
		filters["tags"] = query.Tags
	}

	// Add custom filters
	for key, value := range query.Filters {
		filters[key] = value
	}

	// Create list options
	options := settings.NewListOptions(query.Page, query.Limit, query.SortBy, query.SortOrder, filters)

	// Get configurations
	configurations, err := h.configRepo.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to list configurations: %v", shared.ErrInternalService, err)
	}

	// Get total count
	total, err := h.configRepo.Count(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to count configurations: %v", shared.ErrInternalService, err)
	}

	// Calculate total pages
	totalPages := (total + query.Limit - 1) / query.Limit

	return &ListConfigurationsResult{
		Configurations: configurations,
		Total:          total,
		Page:           query.Page,
		Limit:          query.Limit,
		TotalPages:     totalPages,
	}, nil
}
