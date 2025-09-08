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

// ConfigurationHandler handles HTTP requests for configuration management
type ConfigurationHandler struct {
	createConfigHandler *commands.CreateConfigurationHandler
	updateConfigHandler *commands.UpdateConfigurationHandler
	deleteConfigHandler *commands.DeleteConfigurationHandler
	getConfigHandler    *queries.GetConfigurationHandler
	listConfigsHandler  *queries.ListConfigurationsHandler
}

// NewConfigurationHandler creates a new ConfigurationHandler
func NewConfigurationHandler(
	createConfigHandler *commands.CreateConfigurationHandler,
	updateConfigHandler *commands.UpdateConfigurationHandler,
	deleteConfigHandler *commands.DeleteConfigurationHandler,
	getConfigHandler *queries.GetConfigurationHandler,
	listConfigsHandler *queries.ListConfigurationsHandler,
) *ConfigurationHandler {
	return &ConfigurationHandler{
		createConfigHandler: createConfigHandler,
		updateConfigHandler: updateConfigHandler,
		deleteConfigHandler: deleteConfigHandler,
		getConfigHandler:    getConfigHandler,
		listConfigsHandler:  listConfigsHandler,
	}
}

// CreateConfiguration handles POST /configurations
func (h *ConfigurationHandler) CreateConfiguration(c *gin.Context) {
	var req dto.CreateConfigurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Convert DTO to command
	cmd := commands.CreateConfigurationCommand{
		Key:            req.Key,
		Value:          req.Value,
		Environment:    req.Environment,
		OrganizationID: req.OrganizationID,
		Service:        req.Service,
		Category:       req.Category,
		Description:    req.Description,
		Tags:           req.Tags,
		Metadata:       req.Metadata,
		CreatedBy:      req.CreatedBy,
	}

	// Execute command
	configID, err := h.createConfigHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		if shared.IsAlreadyExistsError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "Configuration already exists", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create configuration", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      configID.String(),
		"message": "Configuration created successfully",
	})
}

// UpdateConfiguration handles PUT /configurations/:id
func (h *ConfigurationHandler) UpdateConfiguration(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Configuration ID is required"})
		return
	}

	var req dto.UpdateConfigurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Convert DTO to command
	cmd := commands.UpdateConfigurationCommand{
		ID:          id,
		Value:       req.Value,
		Description: req.Description,
		Tags:        req.Tags,
		Metadata:    req.Metadata,
		UpdatedBy:   req.UpdatedBy,
	}

	// Execute command
	if err := h.updateConfigHandler.Handle(c.Request.Context(), cmd); err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		if shared.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update configuration", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Configuration updated successfully",
	})
}

// DeleteConfiguration handles DELETE /configurations/:id
func (h *ConfigurationHandler) DeleteConfiguration(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Configuration ID is required"})
		return
	}

	var req dto.DeleteConfigurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Convert DTO to command
	cmd := commands.DeleteConfigurationCommand{
		ID:        id,
		DeletedBy: req.DeletedBy,
	}

	// Execute command
	if err := h.deleteConfigHandler.Handle(c.Request.Context(), cmd); err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		if shared.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete configuration", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Configuration deleted successfully",
	})
}

// GetConfiguration handles GET /configurations/:id
func (h *ConfigurationHandler) GetConfiguration(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Configuration ID is required"})
		return
	}

	// Convert to query
	query := queries.GetConfigurationQuery{
		ID: id,
	}

	// Execute query
	config, err := h.getConfigHandler.Handle(c.Request.Context(), query)
	if err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		if shared.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get configuration", "details": err.Error()})
		return
	}

	// Convert to response DTO
	response := dto.ConfigurationResponse{
		ID:             config.GetID().String(),
		Key:            config.GetKey(),
		Value:          config.GetValue(),
		Environment:    config.GetEnvironment().String(),
		OrganizationID: config.GetOrganizationID(),
		Service:        config.GetService(),
		Category:       config.GetCategory(),
		Description:    config.GetDescription(),
		Tags:           config.GetTags(),
		Metadata:       config.GetMetadata(),
		CreatedBy:      config.GetCreatedBy().String(),
		UpdatedBy:      config.GetUpdatedBy(),
		CreatedAt:      config.GetCreatedAt(),
		UpdatedAt:      config.GetUpdatedAt(),
		Version:        config.GetVersion(),
	}

	c.JSON(http.StatusOK, response)
}

// ListConfigurations handles GET /configurations
func (h *ConfigurationHandler) ListConfigurations(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")
	environment := c.Query("environment")
	organizationID := c.Query("organizationId")
	service := c.Query("service")
	category := c.Query("category")
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
	if len(tags) > 0 {
		filters["tags"] = tags
	}

	// Convert to query
	query := queries.ListConfigurationsQuery{
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
	result, err := h.listConfigsHandler.Handle(c.Request.Context(), query)
	if err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list configurations", "details": err.Error()})
		return
	}

	// Convert to response DTO
	response := dto.ListConfigurationsResponse{
		Configurations: make([]dto.ConfigurationResponse, len(result.Configurations)),
		Total:          result.Total,
		Page:           result.Page,
		Limit:          result.Limit,
		TotalPages:     result.TotalPages,
	}

	for i, config := range result.Configurations {
		response.Configurations[i] = dto.ConfigurationResponse{
			ID:             config.GetID().String(),
			Key:            config.GetKey(),
			Value:          config.GetValue(),
			Environment:    config.GetEnvironment().String(),
			OrganizationID: config.GetOrganizationID(),
			Service:        config.GetService(),
			Category:       config.GetCategory(),
			Description:    config.GetDescription(),
			Tags:           config.GetTags(),
			Metadata:       config.GetMetadata(),
			CreatedBy:      config.GetCreatedBy().String(),
			UpdatedBy:      config.GetUpdatedBy(),
			CreatedAt:      config.GetCreatedAt(),
			UpdatedAt:      config.GetUpdatedAt(),
			Version:        config.GetVersion(),
		}
	}

	c.JSON(http.StatusOK, response)
}
