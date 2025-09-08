package queries

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/analytics"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
)

// ListDashboardsQuery represents the query to list dashboards
type ListDashboardsQuery struct {
	OwnerID string `json:"ownerId" validate:"required,uuid"`
	Public  *bool  `json:"public,omitempty"`
}

// ListDashboardsHandler handles the list dashboards query
type ListDashboardsHandler struct {
	dashboardRepo analytics.DashboardRepository
}

// NewListDashboardsHandler creates a new list dashboards handler
func NewListDashboardsHandler(dashboardRepo analytics.DashboardRepository) *ListDashboardsHandler {
	return &ListDashboardsHandler{
		dashboardRepo: dashboardRepo,
	}
}

// Handle handles the list dashboards query
func (h *ListDashboardsHandler) Handle(ctx context.Context, query ListDashboardsQuery) ([]*analytics.Dashboard, error) {
	// Validate query
	if err := h.validateQuery(query); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Convert string ID to domain type
	ownerID := shared.UserID(query.OwnerID)

	var dashboards []*analytics.Dashboard
	var err error

	// If public filter is specified, handle accordingly
	if query.Public != nil {
		if *query.Public {
			// Get public dashboards
			dashboards, err = h.dashboardRepo.FindPublic(ctx)
		} else {
			// Get private dashboards for owner
			dashboards, err = h.dashboardRepo.FindByOwnerID(ctx, ownerID)
		}
	} else {
		// Get all dashboards for owner
		dashboards, err = h.dashboardRepo.FindByOwnerID(ctx, ownerID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find dashboards: %w", err)
	}

	return dashboards, nil
}

// validateQuery validates the list dashboards query
func (h *ListDashboardsHandler) validateQuery(query ListDashboardsQuery) error {
	if query.OwnerID == "" {
		return shared.NewDomainError("INVALID_OWNER", "Owner ID is required", "")
	}

	return nil
}
