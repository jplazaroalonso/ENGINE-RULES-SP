package queries

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/domain/rule"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/domain/shared"
)

// ListRulesQuery represents a query to list rules with pagination and filtering
type ListRulesQuery struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	SortBy   string `json:"sort_by"`
	SortOrder string `json:"sort_order"` // "asc" or "desc"
	Status   string `json:"status"`
	Category string `json:"category"`
	Search   string `json:"search"`
}

// ListRulesResult represents the result of listing rules
type ListRulesResult struct {
	Rules      []rule.Rule `json:"rules"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ListRulesHandler handles the list rules query
type ListRulesHandler struct {
	ruleRepo rule.Repository
}

// NewListRulesHandler creates a new ListRulesHandler
func NewListRulesHandler(ruleRepo rule.Repository) *ListRulesHandler {
	return &ListRulesHandler{
		ruleRepo: ruleRepo,
	}
}

// Handle executes the list rules query
func (h *ListRulesHandler) Handle(ctx context.Context, query ListRulesQuery) (*ListRulesResult, error) {
	// Set default values
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100 // Max limit to prevent abuse
	}
	if query.SortBy == "" {
		query.SortBy = "created_at"
	}
	if query.SortOrder == "" {
		query.SortOrder = "desc"
	}

	// Validate sort order
	if query.SortOrder != "asc" && query.SortOrder != "desc" {
		return nil, shared.NewValidationError("invalid sort order, must be 'asc' or 'desc'", nil)
	}

	// Validate sort field
	allowedSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"name":       true,
		"status":     true,
		"priority":   true,
		"category":   true,
	}
	if !allowedSortFields[query.SortBy] {
		return nil, shared.NewValidationError(fmt.Sprintf("invalid sort field: %s", query.SortBy), nil)
	}

	// Get total count for pagination
	total, err := h.ruleRepo.Count(ctx, rule.ListFilters{
		Status:   query.Status,
		Category: query.Category,
		Search:   query.Search,
	})
	if err != nil {
		return nil, shared.NewInfrastructureError("failed to count rules", err)
	}

	// Calculate total pages
	totalPages := (total + query.Limit - 1) / query.Limit

	// Get rules with pagination
	rules, err := h.ruleRepo.List(ctx, rule.ListOptions{
		Page:      query.Page,
		Limit:     query.Limit,
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
		Filters: rule.ListFilters{
			Status:   query.Status,
			Category: query.Category,
			Search:   query.Search,
		},
	})
	if err != nil {
		return nil, shared.NewInfrastructureError("failed to list rules", err)
	}

	return &ListRulesResult{
		Rules: rules,
		Pagination: Pagination{
			Page:       query.Page,
			Limit:      query.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
