package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/application/commands"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/application/queries"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/interfaces/rest/dto"
)

// DashboardHandler handles dashboard-related HTTP requests
type DashboardHandler struct {
	createDashboardHandler *commands.CreateDashboardHandler
	getDashboardHandler    *queries.GetDashboardHandler
	listDashboardsHandler  *queries.ListDashboardsHandler
	validator              *validator.Validate
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(
	createDashboardHandler *commands.CreateDashboardHandler,
	getDashboardHandler *queries.GetDashboardHandler,
	listDashboardsHandler *queries.ListDashboardsHandler,
) *DashboardHandler {
	return &DashboardHandler{
		createDashboardHandler: createDashboardHandler,
		getDashboardHandler:    getDashboardHandler,
		listDashboardsHandler:  listDashboardsHandler,
		validator:              validator.New(),
	}
}

// CreateDashboard handles POST /api/v1/dashboards
func (h *DashboardHandler) CreateDashboard(c *gin.Context) {
	var req dto.CreateDashboardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := commands.CreateDashboardCommand{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     req.OwnerID,
	}

	dashboard, err := h.createDashboardHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		if domainErr, ok := err.(*shared.DomainError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": domainErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	response := dto.DashboardResponse{
		ID:              dashboard.ID.String(),
		Name:            dashboard.Name,
		Description:     dashboard.Description,
		Layout:          dashboard.Layout,
		Widgets:         dashboard.Widgets,
		Filters:         dashboard.Filters,
		RefreshInterval: dashboard.RefreshInterval,
		IsPublic:        dashboard.IsPublic,
		OwnerID:         dashboard.OwnerID.String(),
		CreatedAt:       dashboard.CreatedAt,
		UpdatedAt:       dashboard.UpdatedAt,
		Version:         dashboard.Version,
	}

	c.JSON(http.StatusCreated, response)
}

// GetDashboard handles GET /api/v1/dashboards/:id
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dashboard ID is required"})
		return
	}

	query := queries.GetDashboardQuery{ID: id}
	dashboard, err := h.getDashboardHandler.Handle(c.Request.Context(), query)
	if err != nil {
		if err == shared.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Dashboard not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	response := dto.DashboardResponse{
		ID:              dashboard.ID.String(),
		Name:            dashboard.Name,
		Description:     dashboard.Description,
		Layout:          dashboard.Layout,
		Widgets:         dashboard.Widgets,
		Filters:         dashboard.Filters,
		RefreshInterval: dashboard.RefreshInterval,
		IsPublic:        dashboard.IsPublic,
		OwnerID:         dashboard.OwnerID.String(),
		CreatedAt:       dashboard.CreatedAt,
		UpdatedAt:       dashboard.UpdatedAt,
		Version:         dashboard.Version,
	}

	c.JSON(http.StatusOK, response)
}

// ListDashboards handles GET /api/v1/dashboards
func (h *DashboardHandler) ListDashboards(c *gin.Context) {
	ownerID := c.Query("ownerId")
	if ownerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Owner ID is required"})
		return
	}

	publicStr := c.Query("public")
	var public *bool
	if publicStr != "" {
		if publicVal, err := strconv.ParseBool(publicStr); err == nil {
			public = &publicVal
		}
	}

	query := queries.ListDashboardsQuery{
		OwnerID: ownerID,
		Public:  public,
	}

	dashboards, err := h.listDashboardsHandler.Handle(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	responses := make([]dto.DashboardResponse, len(dashboards))
	for i, dashboard := range dashboards {
		responses[i] = dto.DashboardResponse{
			ID:              dashboard.ID.String(),
			Name:            dashboard.Name,
			Description:     dashboard.Description,
			Layout:          dashboard.Layout,
			Widgets:         dashboard.Widgets,
			Filters:         dashboard.Filters,
			RefreshInterval: dashboard.RefreshInterval,
			IsPublic:        dashboard.IsPublic,
			OwnerID:         dashboard.OwnerID.String(),
			CreatedAt:       dashboard.CreatedAt,
			UpdatedAt:       dashboard.UpdatedAt,
			Version:         dashboard.Version,
		}
	}

	c.JSON(http.StatusOK, gin.H{"dashboards": responses})
}

// UpdateDashboard handles PUT /api/v1/dashboards/:id
func (h *DashboardHandler) UpdateDashboard(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// DeleteDashboard handles DELETE /api/v1/dashboards/:id
func (h *DashboardHandler) DeleteDashboard(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// GetRealTimeAnalytics handles GET /api/v1/analytics/real-time
func (h *DashboardHandler) GetRealTimeAnalytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"timestamp": "2024-01-01T00:00:00Z",
		"metrics": map[string]interface{}{
			"activeUsers":    1250,
			"requestsPerSec": 450,
			"responseTime":   "120ms",
			"errorRate":      "0.1%",
		},
	})
}

// GetPerformanceMetrics handles GET /api/v1/analytics/performance
func (h *DashboardHandler) GetPerformanceMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"timestamp": "2024-01-01T00:00:00Z",
		"metrics": map[string]interface{}{
			"cpuUsage":     "45%",
			"memoryUsage":  "2.1GB",
			"diskUsage":    "78%",
			"networkIO":    "125MB/s",
			"responseTime": "120ms",
			"throughput":   "450 req/s",
		},
	})
}

// GetBusinessMetrics handles GET /api/v1/analytics/business
func (h *DashboardHandler) GetBusinessMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"timestamp": "2024-01-01T00:00:00Z",
		"metrics": map[string]interface{}{
			"totalRevenue":    "$125,000",
			"activeCustomers": 1250,
			"conversionRate":  "3.2%",
			"avgOrderValue":   "$89.50",
			"churnRate":       "2.1%",
		},
	})
}

// GetComplianceMetrics handles GET /api/v1/analytics/compliance
func (h *DashboardHandler) GetComplianceMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"timestamp": "2024-01-01T00:00:00Z",
		"metrics": map[string]interface{}{
			"auditScore":     "98%",
			"complianceRate": "99.5%",
			"violations":     2,
			"lastAudit":      "2024-01-01",
			"nextAudit":      "2024-04-01",
		},
	})
}
