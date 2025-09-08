package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/application/commands"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/application/queries"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/interfaces/rest/dto"
)

// MetricHandler handles metric-related HTTP requests
type MetricHandler struct {
	createMetricHandler *commands.CreateMetricHandler
	getMetricsHandler   *queries.GetMetricsHandler
	validator           *validator.Validate
}

// NewMetricHandler creates a new metric handler
func NewMetricHandler(
	createMetricHandler *commands.CreateMetricHandler,
	getMetricsHandler *queries.GetMetricsHandler,
) *MetricHandler {
	return &MetricHandler{
		createMetricHandler: createMetricHandler,
		getMetricsHandler:   getMetricsHandler,
		validator:           validator.New(),
	}
}

// CreateMetric handles POST /api/v1/metrics
func (h *MetricHandler) CreateMetric(c *gin.Context) {
	var req dto.CreateMetricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := commands.CreateMetricCommand{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Category:    req.Category,
		Unit:        req.Unit,
	}

	metric, err := h.createMetricHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		if domainErr, ok := err.(*shared.DomainError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": domainErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	response := dto.MetricResponse{
		ID:           metric.ID.String(),
		Name:         metric.Name,
		Description:  metric.Description,
		Type:         string(metric.Type),
		Category:     string(metric.Category),
		Unit:         metric.Unit,
		Aggregation:  string(metric.Aggregation),
		DataSource:   metric.DataSource,
		Dimensions:   metric.Dimensions,
		Filters:      metric.Filters,
		Calculation:  metric.Calculation,
		IsCalculated: metric.IsCalculated,
		CreatedAt:    metric.CreatedAt,
		UpdatedAt:    metric.UpdatedAt,
		Version:      metric.Version,
	}

	c.JSON(http.StatusCreated, response)
}

// ListMetrics handles GET /api/v1/metrics
func (h *MetricHandler) ListMetrics(c *gin.Context) {
	category := c.Query("category")
	typeFilter := c.Query("type")

	query := queries.GetMetricsQuery{}
	if category != "" {
		query.Category = &category
	}
	if typeFilter != "" {
		query.Type = &typeFilter
	}

	metrics, err := h.getMetricsHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	responses := make([]dto.MetricResponse, len(metrics))
	for i, metric := range metrics {
		responses[i] = dto.MetricResponse{
			ID:           metric.ID.String(),
			Name:         metric.Name,
			Description:  metric.Description,
			Type:         string(metric.Type),
			Category:     string(metric.Category),
			Unit:         metric.Unit,
			Aggregation:  string(metric.Aggregation),
			DataSource:   metric.DataSource,
			Dimensions:   metric.Dimensions,
			Filters:      metric.Filters,
			Calculation:  metric.Calculation,
			IsCalculated: metric.IsCalculated,
			CreatedAt:    metric.CreatedAt,
			UpdatedAt:    metric.UpdatedAt,
			Version:      metric.Version,
		}
	}

	c.JSON(http.StatusOK, gin.H{"metrics": responses})
}

// GetMetric handles GET /api/v1/metrics/:id
func (h *MetricHandler) GetMetric(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// UpdateMetric handles PUT /api/v1/metrics/:id
func (h *MetricHandler) UpdateMetric(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// DeleteMetric handles DELETE /api/v1/metrics/:id
func (h *MetricHandler) DeleteMetric(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// GetMetricData handles GET /api/v1/metrics/:id/data
func (h *MetricHandler) GetMetricData(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}
