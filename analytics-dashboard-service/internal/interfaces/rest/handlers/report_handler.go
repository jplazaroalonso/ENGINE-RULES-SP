package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/application/commands"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/interfaces/rest/dto"
)

// ReportHandler handles report-related HTTP requests
type ReportHandler struct {
	createReportHandler *commands.CreateReportHandler
	validator           *validator.Validate
}

// NewReportHandler creates a new report handler
func NewReportHandler(createReportHandler *commands.CreateReportHandler) *ReportHandler {
	return &ReportHandler{
		createReportHandler: createReportHandler,
		validator:           validator.New(),
	}
}

// CreateReport handles POST /api/v1/reports
func (h *ReportHandler) CreateReport(c *gin.Context) {
	var req dto.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := commands.CreateReportCommand{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		OwnerID:     req.OwnerID,
	}

	report, err := h.createReportHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		if domainErr, ok := err.(*shared.DomainError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": domainErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	response := dto.ReportResponse{
		ID:            report.ID.String(),
		Name:          report.Name,
		Description:   report.Description,
		Type:          string(report.Type),
		Template:      report.Template,
		Parameters:    report.Parameters,
		Schedule:      report.Schedule,
		OutputFormat:  string(report.OutputFormat),
		Recipients:    report.Recipients,
		Status:        string(report.Status),
		LastGenerated: report.LastGenerated,
		NextRun:       report.NextRun,
		OwnerID:       report.OwnerID.String(),
		CreatedAt:     report.CreatedAt,
		UpdatedAt:     report.UpdatedAt,
		Version:       report.Version,
	}

	c.JSON(http.StatusCreated, response)
}

// ListReports handles GET /api/v1/reports
func (h *ReportHandler) ListReports(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// GetReport handles GET /api/v1/reports/:id
func (h *ReportHandler) GetReport(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// UpdateReport handles PUT /api/v1/reports/:id
func (h *ReportHandler) UpdateReport(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// DeleteReport handles DELETE /api/v1/reports/:id
func (h *ReportHandler) DeleteReport(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// GenerateReport handles POST /api/v1/reports/:id/generate
func (h *ReportHandler) GenerateReport(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}
