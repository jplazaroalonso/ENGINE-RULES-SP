package commands

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/analytics"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
)

// CreateReportCommand represents the command to create a report
type CreateReportCommand struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	Type        string `json:"type" validate:"required,oneof=PERFORMANCE COMPLIANCE BUSINESS CUSTOM"`
	OwnerID     string `json:"ownerId" validate:"required,uuid"`
}

// CreateReportHandler handles the create report command
type CreateReportHandler struct {
	reportRepo analytics.ReportRepository
	eventBus   shared.EventBus
}

// NewCreateReportHandler creates a new create report handler
func NewCreateReportHandler(reportRepo analytics.ReportRepository, eventBus shared.EventBus) *CreateReportHandler {
	return &CreateReportHandler{
		reportRepo: reportRepo,
		eventBus:   eventBus,
	}
}

// Handle handles the create report command
func (h *CreateReportHandler) Handle(ctx context.Context, cmd CreateReportCommand) (*analytics.Report, error) {
	// Validate command
	if err := h.validateCommand(cmd); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Convert string IDs to domain types
	ownerID := shared.UserID(cmd.OwnerID)
	reportType := analytics.ReportType(cmd.Type)

	// Create report domain object
	report := analytics.NewReport(cmd.Name, cmd.Description, reportType, ownerID)

	// Save report
	if err := h.reportRepo.Save(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to save report: %w", err)
	}

	// Publish domain events
	if err := h.publishEvents(report); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to publish events: %v\n", err)
	}

	return report, nil
}

// validateCommand validates the create report command
func (h *CreateReportHandler) validateCommand(cmd CreateReportCommand) error {
	if cmd.Name == "" {
		return shared.NewDomainError("INVALID_NAME", "Report name is required", "")
	}

	if len(cmd.Name) > 255 {
		return shared.NewDomainError("INVALID_NAME", "Report name must be less than 255 characters", "")
	}

	if len(cmd.Description) > 1000 {
		return shared.NewDomainError("INVALID_DESCRIPTION", "Report description must be less than 1000 characters", "")
	}

	if cmd.Type == "" {
		return shared.NewDomainError("INVALID_TYPE", "Report type is required", "")
	}

	validTypes := []string{"PERFORMANCE", "COMPLIANCE", "BUSINESS", "CUSTOM"}
	if !shared.Contains(validTypes, cmd.Type) {
		return shared.NewDomainError("INVALID_TYPE", "Invalid report type", "")
	}

	if cmd.OwnerID == "" {
		return shared.NewDomainError("INVALID_OWNER", "Owner ID is required", "")
	}

	return nil
}

// publishEvents publishes domain events
func (h *CreateReportHandler) publishEvents(report *analytics.Report) error {
	for _, event := range report.Events {
		if err := h.eventBus.Publish(event); err != nil {
			return fmt.Errorf("failed to publish event %s: %w", event.Type, err)
		}
	}
	report.ClearEvents()
	return nil
}
