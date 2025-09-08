package queries

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// ListUserPreferencesQuery represents a query to list user preferences with pagination and filtering
type ListUserPreferencesQuery struct {
	Page           int                    `json:"page" validate:"min=1"`
	Limit          int                    `json:"limit" validate:"min=1,max=1000"`
	SortBy         string                 `json:"sortBy"`
	SortOrder      string                 `json:"sortOrder" validate:"omitempty,oneof=asc desc"`
	UserID         *string                `json:"userId,omitempty" validate:"omitempty,uuid"`
	OrganizationID *string                `json:"organizationId,omitempty" validate:"omitempty,uuid"`
	Category       *string                `json:"category,omitempty" validate:"omitempty,min=1,max=100"`
	Tags           []string               `json:"tags,omitempty"`
	Filters        map[string]interface{} `json:"filters,omitempty"`
}

// ListUserPreferencesResult represents the result of listing user preferences
type ListUserPreferencesResult struct {
	UserPreferences []*settings.UserPreference `json:"userPreferences"`
	Total           int                        `json:"total"`
	Page            int                        `json:"page"`
	Limit           int                        `json:"limit"`
	TotalPages      int                        `json:"totalPages"`
}

// ListUserPreferencesHandler handles the ListUserPreferencesQuery
type ListUserPreferencesHandler struct {
	userPreferenceRepo settings.UserPreferenceRepository
	validator          shared.StructValidator
}

// NewListUserPreferencesHandler creates a new ListUserPreferencesHandler
func NewListUserPreferencesHandler(
	userPreferenceRepo settings.UserPreferenceRepository,
	validator shared.StructValidator,
) *ListUserPreferencesHandler {
	return &ListUserPreferencesHandler{
		userPreferenceRepo: userPreferenceRepo,
		validator:          validator,
	}
}

// Handle executes the ListUserPreferencesQuery
func (h *ListUserPreferencesHandler) Handle(ctx context.Context, query ListUserPreferencesQuery) (*ListUserPreferencesResult, error) {
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
	if query.UserID != nil {
		filters["user_id"] = *query.UserID
	}
	if query.OrganizationID != nil {
		filters["organization_id"] = *query.OrganizationID
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

	// Get user preferences
	userPreferences, err := h.userPreferenceRepo.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to list user preferences: %v", shared.ErrInternalService, err)
	}

	// Get total count
	total, err := h.userPreferenceRepo.Count(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to count user preferences: %v", shared.ErrInternalService, err)
	}

	// Calculate total pages
	totalPages := (total + query.Limit - 1) / query.Limit

	return &ListUserPreferencesResult{
		UserPreferences: userPreferences,
		Total:           total,
		Page:            query.Page,
		Limit:           query.Limit,
		TotalPages:      totalPages,
	}, nil
}
