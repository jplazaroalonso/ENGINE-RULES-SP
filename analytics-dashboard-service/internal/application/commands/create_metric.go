package commands

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/analytics"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
)

// CreateMetricCommand represents the command to create a metric
type CreateMetricCommand struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	Type        string `json:"type" validate:"required,oneof=COUNTER GAUGE HISTOGRAM SUMMARY"`
	Category    string `json:"category" validate:"required,oneof=PERFORMANCE BUSINESS SYSTEM USER"`
	Unit        string `json:"unit" validate:"max=50"`
}

// CreateMetricHandler handles the create metric command
type CreateMetricHandler struct {
	metricRepo analytics.MetricRepository
	eventBus   shared.EventBus
}

// NewCreateMetricHandler creates a new create metric handler
func NewCreateMetricHandler(metricRepo analytics.MetricRepository, eventBus shared.EventBus) *CreateMetricHandler {
	return &CreateMetricHandler{
		metricRepo: metricRepo,
		eventBus:   eventBus,
	}
}

// Handle handles the create metric command
func (h *CreateMetricHandler) Handle(ctx context.Context, cmd CreateMetricCommand) (*analytics.Metric, error) {
	// Validate command
	if err := h.validateCommand(cmd); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Convert string types to domain types
	metricType := analytics.MetricType(cmd.Type)
	category := analytics.MetricCategory(cmd.Category)

	// Create metric domain object
	metric := analytics.NewMetric(cmd.Name, cmd.Description, metricType, category)

	if cmd.Unit != "" {
		metric.SetUnit(cmd.Unit)
	}

	// Save metric
	if err := h.metricRepo.Save(ctx, metric); err != nil {
		return nil, fmt.Errorf("failed to save metric: %w", err)
	}

	// Publish domain events
	if err := h.publishEvents(metric); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to publish events: %v\n", err)
	}

	return metric, nil
}

// validateCommand validates the create metric command
func (h *CreateMetricHandler) validateCommand(cmd CreateMetricCommand) error {
	if cmd.Name == "" {
		return shared.NewDomainError("INVALID_NAME", "Metric name is required", "")
	}

	if len(cmd.Name) > 255 {
		return shared.NewDomainError("INVALID_NAME", "Metric name must be less than 255 characters", "")
	}

	if len(cmd.Description) > 1000 {
		return shared.NewDomainError("INVALID_DESCRIPTION", "Metric description must be less than 1000 characters", "")
	}

	if cmd.Type == "" {
		return shared.NewDomainError("INVALID_TYPE", "Metric type is required", "")
	}

	validTypes := []string{"COUNTER", "GAUGE", "HISTOGRAM", "SUMMARY"}
	if !shared.Contains(validTypes, cmd.Type) {
		return shared.NewDomainError("INVALID_TYPE", "Invalid metric type", "")
	}

	if cmd.Category == "" {
		return shared.NewDomainError("INVALID_CATEGORY", "Metric category is required", "")
	}

	validCategories := []string{"PERFORMANCE", "BUSINESS", "SYSTEM", "USER"}
	if !shared.Contains(validCategories, cmd.Category) {
		return shared.NewDomainError("INVALID_CATEGORY", "Invalid metric category", "")
	}

	if len(cmd.Unit) > 50 {
		return shared.NewDomainError("INVALID_UNIT", "Unit must be less than 50 characters", "")
	}

	return nil
}

// publishEvents publishes domain events
func (h *CreateMetricHandler) publishEvents(metric *analytics.Metric) error {
	for _, event := range metric.Events {
		if err := h.eventBus.Publish(event); err != nil {
			return fmt.Errorf("failed to publish event %s: %w", event.Type, err)
		}
	}
	metric.ClearEvents()
	return nil
}
