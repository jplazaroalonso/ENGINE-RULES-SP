package commands

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/analytics"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
)

// CreateDashboardCommand represents the command to create a dashboard
type CreateDashboardCommand struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	OwnerID     string `json:"ownerId" validate:"required,uuid"`
}

// CreateDashboardHandler handles the create dashboard command
type CreateDashboardHandler struct {
	dashboardRepo analytics.DashboardRepository
	eventBus      shared.EventBus
}

// NewCreateDashboardHandler creates a new create dashboard handler
func NewCreateDashboardHandler(dashboardRepo analytics.DashboardRepository, eventBus shared.EventBus) *CreateDashboardHandler {
	return &CreateDashboardHandler{
		dashboardRepo: dashboardRepo,
		eventBus:      eventBus,
	}
}

// Handle handles the create dashboard command
func (h *CreateDashboardHandler) Handle(ctx context.Context, cmd CreateDashboardCommand) (*analytics.Dashboard, error) {
	// Validate command
	if err := h.validateCommand(cmd); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Convert string IDs to domain types
	ownerID := shared.UserID(cmd.OwnerID)

	// Create dashboard domain object
	dashboard := analytics.NewDashboard(cmd.Name, cmd.Description, ownerID)

	// Save dashboard
	if err := h.dashboardRepo.Save(ctx, dashboard); err != nil {
		return nil, fmt.Errorf("failed to save dashboard: %w", err)
	}

	// Publish domain events
	if err := h.publishEvents(dashboard); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to publish events: %v\n", err)
	}

	return dashboard, nil
}

// validateCommand validates the create dashboard command
func (h *CreateDashboardHandler) validateCommand(cmd CreateDashboardCommand) error {
	if cmd.Name == "" {
		return shared.NewDomainError("INVALID_NAME", "Dashboard name is required", "")
	}

	if len(cmd.Name) > 255 {
		return shared.NewDomainError("INVALID_NAME", "Dashboard name must be less than 255 characters", "")
	}

	if len(cmd.Description) > 1000 {
		return shared.NewDomainError("INVALID_DESCRIPTION", "Dashboard description must be less than 1000 characters", "")
	}

	if cmd.OwnerID == "" {
		return shared.NewDomainError("INVALID_OWNER", "Owner ID is required", "")
	}

	return nil
}

// publishEvents publishes domain events
func (h *CreateDashboardHandler) publishEvents(dashboard *analytics.Dashboard) error {
	for _, event := range dashboard.Events {
		if err := h.eventBus.Publish(event); err != nil {
			return fmt.Errorf("failed to publish event %s: %w", event.Type, err)
		}
	}
	dashboard.ClearEvents()
	return nil
}
