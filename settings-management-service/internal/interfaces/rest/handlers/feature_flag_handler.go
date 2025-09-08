package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/application/commands"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/application/queries"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/interfaces/rest/dto"
)

// FeatureFlagHandler handles HTTP requests for feature flag management
type FeatureFlagHandler struct {
	createFeatureFlagHandler *commands.CreateFeatureFlagHandler
	updateFeatureFlagHandler *commands.UpdateFeatureFlagHandler
	deleteFeatureFlagHandler *commands.DeleteFeatureFlagHandler
	getFeatureFlagHandler    *queries.GetFeatureFlagHandler
	listFeatureFlagsHandler  *queries.ListFeatureFlagsHandler
}

// NewFeatureFlagHandler creates a new FeatureFlagHandler
func NewFeatureFlagHandler(
	createFeatureFlagHandler *commands.CreateFeatureFlagHandler,
	updateFeatureFlagHandler *commands.UpdateFeatureFlagHandler,
	deleteFeatureFlagHandler *commands.DeleteFeatureFlagHandler,
	getFeatureFlagHandler *queries.GetFeatureFlagHandler,
	listFeatureFlagsHandler *queries.ListFeatureFlagsHandler,
) *FeatureFlagHandler {
	return &FeatureFlagHandler{
		createFeatureFlagHandler: createFeatureFlagHandler,
		updateFeatureFlagHandler: updateFeatureFlagHandler,
		deleteFeatureFlagHandler: deleteFeatureFlagHandler,
		getFeatureFlagHandler:    getFeatureFlagHandler,
		listFeatureFlagsHandler:  listFeatureFlagsHandler,
	}
}

// CreateFeatureFlag handles POST /feature-flags
func (h *FeatureFlagHandler) CreateFeatureFlag(c *gin.Context) {
	var req dto.CreateFeatureFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Convert DTO to command
	cmd := commands.CreateFeatureFlagCommand{
		Key:            req.Key,
		IsEnabled:      req.IsEnabled,
		Environment:    req.Environment,
		OrganizationID: req.OrganizationID,
		Service:        req.Service,
		Category:       req.Category,
		Description:    req.Description,
		Variants:       req.Variants,
		TargetingRules: req.TargetingRules,
		Tags:           req.Tags,
		Metadata:       req.Metadata,
		CreatedBy:      req.CreatedBy,
	}

	// Execute command
	featureFlagID, err := h.createFeatureFlagHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		if shared.IsAlreadyExistsError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "Feature flag already exists", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create feature flag", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      featureFlagID.String(),
		"message": "Feature flag created successfully",
	})
}

// UpdateFeatureFlag handles PUT /feature-flags/:id
func (h *FeatureFlagHandler) UpdateFeatureFlag(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Feature flag ID is required"})
		return
	}

	var req dto.UpdateFeatureFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Convert DTO to command
	cmd := commands.UpdateFeatureFlagCommand{
		ID:             id,
		IsEnabled:      req.IsEnabled,
		Description:    req.Description,
		Variants:       req.Variants,
		TargetingRules: req.TargetingRules,
		Tags:           req.Tags,
		Metadata:       req.Metadata,
		UpdatedBy:      req.UpdatedBy,
	}

	// Execute command
	if err := h.updateFeatureFlagHandler.Handle(c.Request.Context(), cmd); err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		if shared.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Feature flag not found", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update feature flag", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Feature flag updated successfully",
	})
}

// DeleteFeatureFlag handles DELETE /feature-flags/:id
func (h *FeatureFlagHandler) DeleteFeatureFlag(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Feature flag ID is required"})
		return
	}

	var req dto.DeleteFeatureFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Convert DTO to command
	cmd := commands.DeleteFeatureFlagCommand{
		ID:        id,
		DeletedBy: req.DeletedBy,
	}

	// Execute command
	if err := h.deleteFeatureFlagHandler.Handle(c.Request.Context(), cmd); err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		if shared.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Feature flag not found", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete feature flag", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Feature flag deleted successfully",
	})
}

// GetFeatureFlag handles GET /feature-flags/:id
func (h *FeatureFlagHandler) GetFeatureFlag(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Feature flag ID is required"})
		return
	}

	// Convert to query
	query := queries.GetFeatureFlagQuery{
		ID: id,
	}

	// Execute query
	featureFlag, err := h.getFeatureFlagHandler.Handle(c.Request.Context(), query)
	if err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		if shared.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Feature flag not found", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feature flag", "details": err.Error()})
		return
	}

	// Convert to response DTO
	response := dto.FeatureFlagResponse{
		ID:             featureFlag.GetID().String(),
		Key:            featureFlag.GetKey(),
		IsEnabled:      featureFlag.GetIsEnabled(),
		Environment:    featureFlag.GetEnvironment().String(),
		OrganizationID: featureFlag.GetOrganizationID(),
		Service:        featureFlag.GetService(),
		Category:       featureFlag.GetCategory(),
		Description:    featureFlag.GetDescription(),
		Variants:       featureFlag.GetVariants(),
		TargetingRules: featureFlag.GetTargetingRules(),
		Tags:           featureFlag.GetTags(),
		Metadata:       featureFlag.GetMetadata(),
		CreatedBy:      featureFlag.GetCreatedBy().String(),
		UpdatedBy:      featureFlag.GetUpdatedBy(),
		CreatedAt:      featureFlag.GetCreatedAt(),
		UpdatedAt:      featureFlag.GetUpdatedAt(),
		Version:        featureFlag.GetVersion(),
	}

	c.JSON(http.StatusOK, response)
}

// ListFeatureFlags handles GET /feature-flags
func (h *FeatureFlagHandler) ListFeatureFlags(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")
	environment := c.Query("environment")
	organizationID := c.Query("organizationId")
	service := c.Query("service")
	category := c.Query("category")
	isEnabledStr := c.Query("isEnabled")
	tags := c.QueryArray("tags")

	// Build filters
	filters := make(map[string]interface{})
	if environment != "" {
		filters["environment"] = environment
	}
	if organizationID != "" {
		filters["organizationId"] = organizationID
	}
	if service != "" {
		filters["service"] = service
	}
	if category != "" {
		filters["category"] = category
	}
	if isEnabledStr != "" {
		if isEnabled, err := strconv.ParseBool(isEnabledStr); err == nil {
			filters["isEnabled"] = isEnabled
		}
	}
	if len(tags) > 0 {
		filters["tags"] = tags
	}

	// Convert to query
	query := queries.ListFeatureFlagsQuery{
		Page:           page,
		Limit:          limit,
		SortBy:         sortBy,
		SortOrder:      sortOrder,
		Environment:    &environment,
		OrganizationID: &organizationID,
		Service:        &service,
		Category:       &category,
		Tags:           tags,
		Filters:        filters,
	}

	// Execute query
	result, err := h.listFeatureFlagsHandler.Handle(c.Request.Context(), query)
	if err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list feature flags", "details": err.Error()})
		return
	}

	// Convert to response DTO
	response := dto.ListFeatureFlagsResponse{
		FeatureFlags: make([]dto.FeatureFlagResponse, len(result.FeatureFlags)),
		Total:        result.Total,
		Page:         result.Page,
		Limit:        result.Limit,
		TotalPages:   result.TotalPages,
	}

	for i, featureFlag := range result.FeatureFlags {
		response.FeatureFlags[i] = dto.FeatureFlagResponse{
			ID:             featureFlag.GetID().String(),
			Key:            featureFlag.GetKey(),
			IsEnabled:      featureFlag.GetIsEnabled(),
			Environment:    featureFlag.GetEnvironment().String(),
			OrganizationID: featureFlag.GetOrganizationID(),
			Service:        featureFlag.GetService(),
			Category:       featureFlag.GetCategory(),
			Description:    featureFlag.GetDescription(),
			Variants:       featureFlag.GetVariants(),
			TargetingRules: featureFlag.GetTargetingRules(),
			Tags:           featureFlag.GetTags(),
			Metadata:       featureFlag.GetMetadata(),
			CreatedBy:      featureFlag.GetCreatedBy().String(),
			UpdatedBy:      featureFlag.GetUpdatedBy(),
			CreatedAt:      featureFlag.GetCreatedAt(),
			UpdatedAt:      featureFlag.GetUpdatedAt(),
			Version:        featureFlag.GetVersion(),
		}
	}

	c.JSON(http.StatusOK, response)
}
