package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/application/commands"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/application/queries"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/interfaces/rest/dto"
)

// CampaignHandler handles HTTP requests for campaign operations
type CampaignHandler struct {
	createHandler   *commands.CreateCampaignHandler
	updateHandler   *commands.UpdateCampaignHandler
	activateHandler *commands.ActivateCampaignHandler
	pauseHandler    *commands.PauseCampaignHandler
	deleteHandler   *commands.DeleteCampaignHandler
	getHandler      *queries.GetCampaignHandler
	listHandler     *queries.ListCampaignsHandler
	metricsHandler  *queries.GetCampaignMetricsHandler
}

// NewCampaignHandler creates a new CampaignHandler
func NewCampaignHandler(
	createHandler *commands.CreateCampaignHandler,
	updateHandler *commands.UpdateCampaignHandler,
	activateHandler *commands.ActivateCampaignHandler,
	pauseHandler *commands.PauseCampaignHandler,
	deleteHandler *commands.DeleteCampaignHandler,
	getHandler *queries.GetCampaignHandler,
	listHandler *queries.ListCampaignsHandler,
	metricsHandler *queries.GetCampaignMetricsHandler,
) *CampaignHandler {
	return &CampaignHandler{
		createHandler:   createHandler,
		updateHandler:   updateHandler,
		activateHandler: activateHandler,
		pauseHandler:    pauseHandler,
		deleteHandler:   deleteHandler,
		getHandler:      getHandler,
		listHandler:     listHandler,
		metricsHandler:  metricsHandler,
	}
}

// CreateCampaign handles POST /campaigns
func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var req dto.CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Convert budget from DTO to command format
	var budget *commands.BudgetRequest
	if req.Budget != nil {
		budget = &commands.BudgetRequest{
			Amount:   req.Budget.Amount,
			Currency: req.Budget.Currency,
		}
	}

	// Convert DTO to command
	cmd := commands.CreateCampaignCommand{
		Name:           req.Name,
		Description:    req.Description,
		CampaignType:   req.Type,
		TargetingRules: req.TargetingRules,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Budget:         budget,
		CreatedBy:      req.CreatedBy,
		Settings:       h.convertToCampaignSettingsRequest(req.Settings),
	}

	result, err := h.createHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    result,
	})
}

// UpdateCampaign handles PUT /campaigns/:id
func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campaign ID is required"})
		return
	}

	var req dto.UpdateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Convert budget from DTO to command format
	var budget *commands.BudgetRequest
	if req.Budget != nil {
		budget = &commands.BudgetRequest{
			Amount:   req.Budget.Amount,
			Currency: req.Budget.Currency,
		}
	}

	// Convert DTO to command
	cmd := commands.UpdateCampaignCommand{
		CampaignID:     id,
		Name:           &req.Name,
		Description:    &req.Description,
		TargetingRules: req.TargetingRules,
		StartDate:      &req.StartDate,
		EndDate:        req.EndDate,
		Budget:         budget,
		Settings:       h.convertToCampaignSettingsRequestPtr(&req.Settings),
		UpdatedBy:      req.UpdatedBy,
	}

	result, err := h.updateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ActivateCampaign handles POST /campaigns/:id/activate
func (h *CampaignHandler) ActivateCampaign(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campaign ID is required"})
		return
	}

	var req dto.ActivateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	cmd := commands.ActivateCampaignCommand{
		CampaignID:  id,
		ActivatedBy: req.ActivatedBy,
	}

	result, err := h.activateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// PauseCampaign handles POST /campaigns/:id/pause
func (h *CampaignHandler) PauseCampaign(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campaign ID is required"})
		return
	}

	var req dto.PauseCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	cmd := commands.PauseCampaignCommand{
		CampaignID: id,
		PausedBy:   req.PausedBy,
	}

	result, err := h.pauseHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// DeleteCampaign handles DELETE /campaigns/:id
func (h *CampaignHandler) DeleteCampaign(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campaign ID is required"})
		return
	}

	var req dto.DeleteCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	cmd := commands.DeleteCampaignCommand{
		CampaignID: id,
		DeletedBy:  req.DeletedBy,
		Reason:     req.Reason,
	}

	result, err := h.deleteHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetCampaign handles GET /campaigns/:id
func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campaign ID is required"})
		return
	}

	query := queries.GetCampaignQuery{CampaignID: id}
	result, err := h.getHandler.Handle(c.Request.Context(), query)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ListCampaigns handles GET /campaigns
func (h *CampaignHandler) ListCampaigns(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	sortBy := c.DefaultQuery("sortBy", "createdAt")
	sortOrder := c.DefaultQuery("sortOrder", "desc")
	status := c.Query("status")
	campaignType := c.Query("type")
	search := c.Query("search")
	createdBy := c.Query("createdBy")

	query := queries.ListCampaignsQuery{
		Page:         page,
		Limit:        limit,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
		Status:       status,
		CampaignType: campaignType,
		SearchQuery:  search,
		CreatedBy:    createdBy,
	}

	result, err := h.listHandler.Handle(c.Request.Context(), query)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetCampaignMetrics handles GET /campaigns/:id/metrics
func (h *CampaignHandler) GetCampaignMetrics(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campaign ID is required"})
		return
	}

	query := queries.GetCampaignMetricsQuery{CampaignID: id}
	result, err := h.metricsHandler.Handle(c.Request.Context(), query)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// handleError handles application errors and converts them to appropriate HTTP responses
func (h *CampaignHandler) handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *shared.ValidationError:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": e.Error()})
	case *shared.NotFoundError:
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "details": e.Error()})
	case *shared.BusinessError:
		c.JSON(http.StatusConflict, gin.H{"error": "Business error", "details": e.Error()})
	case *shared.InfrastructureError:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Infrastructure error", "details": e.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": err.Error()})
	}
}

// convertToCampaignSettingsRequest converts DTO settings to command settings
func (h *CampaignHandler) convertToCampaignSettingsRequest(dtoSettings dto.CreateCampaignSettings) commands.CampaignSettingsRequest {
	return commands.CampaignSettingsRequest{
		TargetAudience:  dtoSettings.TargetAudience,
		Channels:        dtoSettings.Channels,
		Frequency:       dtoSettings.Frequency,
		MaxImpressions:  dtoSettings.MaxImpressions,
		BudgetLimit:     h.convertBudgetToRequest(dtoSettings.BudgetLimit),
		ABTestSettings:  h.convertABTestSettings(dtoSettings.ABTestSettings),
		SchedulingRules: h.convertSchedulingRules(dtoSettings.SchedulingRules),
		Personalization: commands.PersonalizationConfigRequest{
			Enabled:     dtoSettings.Personalization.Enabled,
			Rules:       dtoSettings.Personalization.Rules,
			Fallback:    dtoSettings.Personalization.Fallback,
			MaxVariants: dtoSettings.Personalization.MaxVariants,
		},
	}
}

// convertToCampaignSettingsRequestPtr converts DTO settings pointer to command settings pointer
func (h *CampaignHandler) convertToCampaignSettingsRequestPtr(dtoSettings *dto.UpdateCampaignSettings) *commands.CampaignSettingsRequest {
	if dtoSettings == nil {
		return nil
	}

	settings := commands.CampaignSettingsRequest{
		TargetAudience:  dtoSettings.TargetAudience,
		Channels:        dtoSettings.Channels,
		Frequency:       dtoSettings.Frequency,
		MaxImpressions:  dtoSettings.MaxImpressions,
		BudgetLimit:     h.convertBudgetToRequest(dtoSettings.BudgetLimit),
		ABTestSettings:  h.convertABTestSettings(dtoSettings.ABTestSettings),
		SchedulingRules: h.convertSchedulingRules(dtoSettings.SchedulingRules),
		Personalization: commands.PersonalizationConfigRequest{
			Enabled:     dtoSettings.Personalization.Enabled,
			Rules:       dtoSettings.Personalization.Rules,
			Fallback:    dtoSettings.Personalization.Fallback,
			MaxVariants: dtoSettings.Personalization.MaxVariants,
		},
	}
	return &settings
}

// convertBudgetToRequest converts DTO budget to command budget request
func (h *CampaignHandler) convertBudgetToRequest(budget *shared.Money) *commands.BudgetRequest {
	if budget == nil {
		return nil
	}
	return &commands.BudgetRequest{
		Amount:   budget.Amount,
		Currency: budget.Currency,
	}
}

// convertABTestSettings converts DTO AB test settings to command AB test settings
func (h *CampaignHandler) convertABTestSettings(settings *dto.ABTestSettingsRequest) *commands.ABTestSettingsRequest {
	if settings == nil {
		return nil
	}

	variants := make([]commands.VariantRequest, len(settings.Variants))
	for i, variant := range settings.Variants {
		variants[i] = commands.VariantRequest{
			ID:          variant.ID,
			Name:        variant.Name,
			Description: variant.Description,
			Settings:    variant.Settings,
			Weight:      variant.Weight,
		}
	}

	return &commands.ABTestSettingsRequest{
		Enabled:       settings.Enabled,
		Variants:      variants,
		TrafficSplit:  settings.TrafficSplit,
		SuccessMetric: settings.SuccessMetric,
		Duration:      settings.Duration,
	}
}

// convertSchedulingRules converts DTO scheduling rules to command scheduling rules
func (h *CampaignHandler) convertSchedulingRules(rules []dto.SchedulingRuleRequest) []commands.SchedulingRuleRequest {
	if rules == nil {
		return nil
	}

	converted := make([]commands.SchedulingRuleRequest, len(rules))
	for i, rule := range rules {
		conditions := make([]commands.SchedulingConditionRequest, len(rule.Conditions))
		for j, condition := range rule.Conditions {
			conditions[j] = commands.SchedulingConditionRequest{
				Type:     condition.Type,
				Operator: condition.Operator,
				Value:    condition.Value,
				Metadata: condition.Metadata,
			}
		}

		actions := make([]commands.SchedulingActionRequest, len(rule.Actions))
		for j, action := range rule.Actions {
			actions[j] = commands.SchedulingActionRequest{
				Type:       action.Type,
				Parameters: action.Parameters,
			}
		}

		converted[i] = commands.SchedulingRuleRequest{
			ID:          rule.ID,
			Name:        rule.Name,
			Description: rule.Description,
			Conditions:  conditions,
			Actions:     actions,
			IsActive:    rule.IsActive,
		}
	}

	return converted
}
