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

// OrganizationSettingHandler handles HTTP requests for organization setting management
type OrganizationSettingHandler struct {
	createOrganizationSettingHandler *commands.CreateOrganizationSettingHandler
	updateOrganizationSettingHandler *commands.UpdateOrganizationSettingHandler
	deleteOrganizationSettingHandler *commands.DeleteOrganizationSettingHandler
	getOrganizationSettingHandler    *queries.GetOrganizationSettingHandler
	listOrganizationSettingsHandler  *queries.ListOrganizationSettingsHandler
}

// NewOrganizationSettingHandler creates a new OrganizationSettingHandler
func NewOrganizationSettingHandler(
	createOrganizationSettingHandler *commands.CreateOrganizationSettingHandler,
	updateOrganizationSettingHandler *commands.UpdateOrganizationSettingHandler,
	deleteOrganizationSettingHandler *commands.DeleteOrganizationSettingHandler,
	getOrganizationSettingHandler *queries.GetOrganizationSettingHandler,
	listOrganizationSettingsHandler *queries.ListOrganizationSettingsHandler,
) *OrganizationSettingHandler {
	return &OrganizationSettingHandler{
		createOrganizationSettingHandler: createOrganizationSettingHandler,
		updateOrganizationSettingHandler: updateOrganizationSettingHandler,
		deleteOrganizationSettingHandler: deleteOrganizationSettingHandler,
		getOrganizationSettingHandler:    getOrganizationSettingHandler,
		listOrganizationSettingsHandler:  listOrganizationSettingsHandler,
	}
}

// CreateOrganizationSetting handles POST /organization-settings
func (h *OrganizationSettingHandler) CreateOrganizationSetting(c *gin.Context) {
	var req dto.CreateOrganizationSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Convert DTO to command
	cmd := commands.CreateOrganizationSettingCommand{
		OrganizationID: req.OrganizationID,
		Category:       req.Category,
		Key:            req.Key,
		Value:          req.Value,
		ParentID:       req.ParentID,
		Description:    req.Description,
		Tags:           req.Tags,
		Metadata:       req.Metadata,
		CreatedBy:      req.CreatedBy,
	}

	// Execute command
	organizationSettingID, err := h.createOrganizationSettingHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		if shared.IsAlreadyExistsError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "Organization setting already exists", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization setting", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      organizationSettingID.String(),
		"message": "Organization setting created successfully",
	})
}

// UpdateOrganizationSetting handles PUT /organization-settings/:id
func (h *OrganizationSettingHandler) UpdateOrganizationSetting(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization setting ID is required"})
		return
	}

	var req dto.UpdateOrganizationSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Convert DTO to command
	cmd := commands.UpdateOrganizationSettingCommand{
		ID:          id,
		Value:       req.Value,
		Description: req.Description,
		Tags:        req.Tags,
		Metadata:    req.Metadata,
		UpdatedBy:   req.UpdatedBy,
	}

	// Execute command
	if err := h.updateOrganizationSettingHandler.Handle(c.Request.Context(), cmd); err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		if shared.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization setting not found", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization setting", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organization setting updated successfully",
	})
}

// DeleteOrganizationSetting handles DELETE /organization-settings/:id
func (h *OrganizationSettingHandler) DeleteOrganizationSetting(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization setting ID is required"})
		return
	}

	var req dto.DeleteOrganizationSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Convert DTO to command
	cmd := commands.DeleteOrganizationSettingCommand{
		ID:        id,
		DeletedBy: req.DeletedBy,
	}

	// Execute command
	if err := h.deleteOrganizationSettingHandler.Handle(c.Request.Context(), cmd); err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		if shared.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization setting not found", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization setting", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organization setting deleted successfully",
	})
}

// GetOrganizationSetting handles GET /organization-settings/:id
func (h *OrganizationSettingHandler) GetOrganizationSetting(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization setting ID is required"})
		return
	}

	// Convert to query
	query := queries.GetOrganizationSettingQuery{
		ID: id,
	}

	// Execute query
	organizationSetting, err := h.getOrganizationSettingHandler.Handle(c.Request.Context(), query)
	if err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		if shared.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization setting not found", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization setting", "details": err.Error()})
		return
	}

	// Convert to response DTO
	response := dto.OrganizationSettingResponse{
		ID:             organizationSetting.GetID().String(),
		OrganizationID: organizationSetting.GetOrganizationID().String(),
		Category:       organizationSetting.GetCategory(),
		Key:            organizationSetting.GetKey(),
		Value:          organizationSetting.GetValue(),
		ParentID:       organizationSetting.GetParentID(),
		Description:    organizationSetting.GetDescription(),
		Tags:           organizationSetting.GetTags(),
		Metadata:       organizationSetting.GetMetadata(),
		CreatedBy:      organizationSetting.GetCreatedBy().String(),
		UpdatedBy:      organizationSetting.GetUpdatedBy(),
		CreatedAt:      organizationSetting.GetCreatedAt(),
		UpdatedAt:      organizationSetting.GetUpdatedAt(),
		Version:        organizationSetting.GetVersion(),
	}

	c.JSON(http.StatusOK, response)
}

// ListOrganizationSettings handles GET /organization-settings
func (h *OrganizationSettingHandler) ListOrganizationSettings(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")
	organizationID := c.Query("organizationId")
	category := c.Query("category")
	parentID := c.Query("parentId")
	tags := c.QueryArray("tags")

	// Build filters
	filters := make(map[string]interface{})
	if organizationID != "" {
		filters["organizationId"] = organizationID
	}
	if category != "" {
		filters["category"] = category
	}
	if parentID != "" {
		filters["parentId"] = parentID
	}
	if len(tags) > 0 {
		filters["tags"] = tags
	}

	// Convert to query
	query := queries.ListOrganizationSettingsQuery{
		Page:           page,
		Limit:          limit,
		SortBy:         sortBy,
		SortOrder:      sortOrder,
		OrganizationID: &organizationID,
		Category:       &category,
		ParentID:       &parentID,
		Tags:           tags,
		Filters:        filters,
	}

	// Execute query
	result, err := h.listOrganizationSettingsHandler.Handle(c.Request.Context(), query)
	if err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list organization settings", "details": err.Error()})
		return
	}

	// Convert to response DTO
	response := dto.ListOrganizationSettingsResponse{
		OrganizationSettings: make([]dto.OrganizationSettingResponse, len(result.OrganizationSettings)),
		Total:                result.Total,
		Page:                 result.Page,
		Limit:                result.Limit,
		TotalPages:           result.TotalPages,
	}

	for i, organizationSetting := range result.OrganizationSettings {
		response.OrganizationSettings[i] = dto.OrganizationSettingResponse{
			ID:             organizationSetting.GetID().String(),
			OrganizationID: organizationSetting.GetOrganizationID().String(),
			Category:       organizationSetting.GetCategory(),
			Key:            organizationSetting.GetKey(),
			Value:          organizationSetting.GetValue(),
			ParentID:       organizationSetting.GetParentID(),
			Description:    organizationSetting.GetDescription(),
			Tags:           organizationSetting.GetTags(),
			Metadata:       organizationSetting.GetMetadata(),
			CreatedBy:      organizationSetting.GetCreatedBy().String(),
			UpdatedBy:      organizationSetting.GetUpdatedBy(),
			CreatedAt:      organizationSetting.GetCreatedAt(),
			UpdatedAt:      organizationSetting.GetUpdatedAt(),
			Version:        organizationSetting.GetVersion(),
		}
	}

	c.JSON(http.StatusOK, response)
}
