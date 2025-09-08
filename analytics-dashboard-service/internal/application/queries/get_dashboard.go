package queries

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/analytics"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
)

// GetDashboardQuery represents the query to get a dashboard
type GetDashboardQuery struct {
	ID string `json:"id" validate:"required,uuid"`
}

// GetDashboardHandler handles the get dashboard query
type GetDashboardHandler struct {
	dashboardRepo analytics.DashboardRepository
}

// NewGetDashboardHandler creates a new get dashboard handler
func NewGetDashboardHandler(dashboardRepo analytics.DashboardRepository) *GetDashboardHandler {
	return &GetDashboardHandler{
		dashboardRepo: dashboardRepo,
	}
}

// Handle handles the get dashboard query
func (h *GetDashboardHandler) Handle(ctx context.Context, query GetDashboardQuery) (*analytics.Dashboard, error) {
	// Validate query
	if err := h.validateQuery(query); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Convert string ID to domain type
	dashboardID := shared.DashboardID(query.ID)

	// Find dashboard
	dashboard, err := h.dashboardRepo.FindByID(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to find dashboard: %w", err)
	}

	return dashboard, nil
}

// validateQuery validates the get dashboard query
func (h *GetDashboardHandler) validateQuery(query GetDashboardQuery) error {
	if query.ID == "" {
		return shared.NewDomainError("INVALID_ID", "Dashboard ID is required", "")
	}

	return nil
}
