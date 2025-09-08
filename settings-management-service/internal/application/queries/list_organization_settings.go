package queries

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// ListOrganizationSettingsQuery represents a query to list organization settings with pagination and filtering
type ListOrganizationSettingsQuery struct {
	Page           int                    `json:"page" validate:"min=1"`
	Limit          int                    `json:"limit" validate:"min=1,max=1000"`
	SortBy         string                 `json:"sortBy"`
	SortOrder      string                 `json:"sortOrder" validate:"omitempty,oneof=asc desc"`
	OrganizationID *string                `json:"organizationId,omitempty" validate:"omitempty,uuid"`
	Category       *string                `json:"category,omitempty" validate:"omitempty,min=1,max=100"`
	ParentID       *string                `json:"parentId,omitempty" validate:"omitempty,uuid"`
	Tags           []string               `json:"tags,omitempty"`
	Filters        map[string]interface{} `json:"filters,omitempty"`
}

// ListOrganizationSettingsResult represents the result of listing organization settings
type ListOrganizationSettingsResult struct {
	OrganizationSettings []*settings.OrganizationSetting `json:"organizationSettings"`
	Total                int                             `json:"total"`
	Page                 int                             `json:"page"`
	Limit                int                             `json:"limit"`
	TotalPages           int                             `json:"totalPages"`
}

// ListOrganizationSettingsHandler handles the ListOrganizationSettingsQuery
type ListOrganizationSettingsHandler struct {
	organizationSettingRepo settings.OrganizationSettingRepository
	validator               shared.StructValidator
}

// NewListOrganizationSettingsHandler creates a new ListOrganizationSettingsHandler
func NewListOrganizationSettingsHandler(
	organizationSettingRepo settings.OrganizationSettingRepository,
	validator shared.StructValidator,
) *ListOrganizationSettingsHandler {
	return &ListOrganizationSettingsHandler{
		organizationSettingRepo: organizationSettingRepo,
		validator:               validator,
	}
}

// Handle executes the ListOrganizationSettingsQuery
func (h *ListOrganizationSettingsHandler) Handle(ctx context.Context, query ListOrganizationSettingsQuery) (*ListOrganizationSettingsResult, error) {
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
	if query.OrganizationID != nil {
		filters["organization_id"] = *query.OrganizationID
	}
	if query.Category != nil {
		filters["category"] = *query.Category
	}
	if query.ParentID != nil {
		filters["parent_id"] = *query.ParentID
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

	// Get organization settings
	organizationSettings, err := h.organizationSettingRepo.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to list organization settings: %v", shared.ErrInternalService, err)
	}

	// Get total count
	total, err := h.organizationSettingRepo.Count(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to count organization settings: %v", shared.ErrInternalService, err)
	}

	// Calculate total pages
	totalPages := (total + query.Limit - 1) / query.Limit

	return &ListOrganizationSettingsResult{
		OrganizationSettings: organizationSettings,
		Total:                total,
		Page:                 query.Page,
		Limit:                query.Limit,
		TotalPages:           totalPages,
	}, nil
}
