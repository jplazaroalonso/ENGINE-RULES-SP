package queries

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/analytics"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
)

// GetMetricsQuery represents the query to get metrics
type GetMetricsQuery struct {
	Category *string `json:"category,omitempty" validate:"omitempty,oneof=PERFORMANCE BUSINESS SYSTEM USER"`
	Type     *string `json:"type,omitempty" validate:"omitempty,oneof=COUNTER GAUGE HISTOGRAM SUMMARY"`
}

// GetMetricsHandler handles the get metrics query
type GetMetricsHandler struct {
	metricRepo analytics.MetricRepository
}

// NewGetMetricsHandler creates a new get metrics handler
func NewGetMetricsHandler(metricRepo analytics.MetricRepository) *GetMetricsHandler {
	return &GetMetricsHandler{
		metricRepo: metricRepo,
	}
}

// Handle handles the get metrics query
func (h *GetMetricsHandler) Handle(ctx context.Context, query GetMetricsQuery) ([]*analytics.Metric, error) {
	// Validate query
	if err := h.validateQuery(query); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	var metrics []*analytics.Metric
	var err error

	// Apply filters based on query parameters
	if query.Category != nil && query.Type != nil {
		// Filter by both category and type
		category := analytics.MetricCategory(*query.Category)
		metricType := analytics.MetricType(*query.Type)

		// Get metrics by category first, then filter by type
		categoryMetrics, err := h.metricRepo.FindByCategory(ctx, category)
		if err != nil {
			return nil, fmt.Errorf("failed to find metrics by category: %w", err)
		}

		// Filter by type
		for _, metric := range categoryMetrics {
			if metric.Type == metricType {
				metrics = append(metrics, metric)
			}
		}
	} else if query.Category != nil {
		// Filter by category only
		category := analytics.MetricCategory(*query.Category)
		metrics, err = h.metricRepo.FindByCategory(ctx, category)
		if err != nil {
			return nil, fmt.Errorf("failed to find metrics by category: %w", err)
		}
	} else if query.Type != nil {
		// Filter by type only
		metricType := analytics.MetricType(*query.Type)
		metrics, err = h.metricRepo.FindByType(ctx, metricType)
		if err != nil {
			return nil, fmt.Errorf("failed to find metrics by type: %w", err)
		}
	} else {
		// Get all metrics
		metrics, err = h.metricRepo.FindAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to find all metrics: %w", err)
		}
	}

	return metrics, nil
}

// validateQuery validates the get metrics query
func (h *GetMetricsHandler) validateQuery(query GetMetricsQuery) error {
	if query.Category != nil {
		validCategories := []string{"PERFORMANCE", "BUSINESS", "SYSTEM", "USER"}
		if !shared.Contains(validCategories, *query.Category) {
			return shared.NewDomainError("INVALID_CATEGORY", "Invalid metric category", "")
		}
	}

	if query.Type != nil {
		validTypes := []string{"COUNTER", "GAUGE", "HISTOGRAM", "SUMMARY"}
		if !shared.Contains(validTypes, *query.Type) {
			return shared.NewDomainError("INVALID_TYPE", "Invalid metric type", "")
		}
	}

	return nil
}
