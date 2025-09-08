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

// UserPreferenceHandler handles HTTP requests for user preference management
type UserPreferenceHandler struct {
	createUserPreferenceHandler *commands.CreateUserPreferenceHandler
	updateUserPreferenceHandler *commands.UpdateUserPreferenceHandler
	deleteUserPreferenceHandler *commands.DeleteUserPreferenceHandler
	getUserPreferenceHandler    *queries.GetUserPreferenceHandler
	listUserPreferencesHandler  *queries.ListUserPreferencesHandler
}

// NewUserPreferenceHandler creates a new UserPreferenceHandler
func NewUserPreferenceHandler(
	createUserPreferenceHandler *commands.CreateUserPreferenceHandler,
	updateUserPreferenceHandler *commands.UpdateUserPreferenceHandler,
	deleteUserPreferenceHandler *commands.DeleteUserPreferenceHandler,
	getUserPreferenceHandler *queries.GetUserPreferenceHandler,
	listUserPreferencesHandler *queries.ListUserPreferencesHandler,
) *UserPreferenceHandler {
	return &UserPreferenceHandler{
		createUserPreferenceHandler: createUserPreferenceHandler,
		updateUserPreferenceHandler: updateUserPreferenceHandler,
		deleteUserPreferenceHandler: deleteUserPreferenceHandler,
		getUserPreferenceHandler:    getUserPreferenceHandler,
		listUserPreferencesHandler:  listUserPreferencesHandler,
	}
}

// CreateUserPreference handles POST /user-preferences
func (h *UserPreferenceHandler) CreateUserPreference(c *gin.Context) {
	var req dto.CreateUserPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Convert DTO to command
	cmd := commands.CreateUserPreferenceCommand{
		UserID:         req.UserID,
		Category:       req.Category,
		Key:            req.Key,
		Value:          req.Value,
		OrganizationID: req.OrganizationID,
		Description:    req.Description,
		Tags:           req.Tags,
		Metadata:       req.Metadata,
		CreatedBy:      req.CreatedBy,
	}

	// Execute command
	userPreferenceID, err := h.createUserPreferenceHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		if shared.IsAlreadyExistsError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "User preference already exists", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user preference", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      userPreferenceID.String(),
		"message": "User preference created successfully",
	})
}

// UpdateUserPreference handles PUT /user-preferences/:id
func (h *UserPreferenceHandler) UpdateUserPreference(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User preference ID is required"})
		return
	}

	var req dto.UpdateUserPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Convert DTO to command
	cmd := commands.UpdateUserPreferenceCommand{
		ID:          id,
		Value:       req.Value,
		Description: req.Description,
		Tags:        req.Tags,
		Metadata:    req.Metadata,
		UpdatedBy:   req.UpdatedBy,
	}

	// Execute command
	if err := h.updateUserPreferenceHandler.Handle(c.Request.Context(), cmd); err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		if shared.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User preference not found", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user preference", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User preference updated successfully",
	})
}

// DeleteUserPreference handles DELETE /user-preferences/:id
func (h *UserPreferenceHandler) DeleteUserPreference(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User preference ID is required"})
		return
	}

	var req dto.DeleteUserPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Convert DTO to command
	cmd := commands.DeleteUserPreferenceCommand{
		ID:        id,
		DeletedBy: req.DeletedBy,
	}

	// Execute command
	if err := h.deleteUserPreferenceHandler.Handle(c.Request.Context(), cmd); err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		if shared.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User preference not found", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user preference", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User preference deleted successfully",
	})
}

// GetUserPreference handles GET /user-preferences/:id
func (h *UserPreferenceHandler) GetUserPreference(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User preference ID is required"})
		return
	}

	// Convert to query
	query := queries.GetUserPreferenceQuery{
		ID: id,
	}

	// Execute query
	userPreference, err := h.getUserPreferenceHandler.Handle(c.Request.Context(), query)
	if err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		if shared.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User preference not found", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user preference", "details": err.Error()})
		return
	}

	// Convert to response DTO
	response := dto.UserPreferenceResponse{
		ID:             userPreference.GetID().String(),
		UserID:         userPreference.GetUserID().String(),
		Category:       userPreference.GetCategory(),
		Key:            userPreference.GetKey(),
		Value:          userPreference.GetValue(),
		OrganizationID: userPreference.GetOrganizationID(),
		Description:    userPreference.GetDescription(),
		Tags:           userPreference.GetTags(),
		Metadata:       userPreference.GetMetadata(),
		CreatedBy:      userPreference.GetCreatedBy().String(),
		UpdatedBy:      userPreference.GetUpdatedBy(),
		CreatedAt:      userPreference.GetCreatedAt(),
		UpdatedAt:      userPreference.GetUpdatedAt(),
		Version:        userPreference.GetVersion(),
	}

	c.JSON(http.StatusOK, response)
}

// ListUserPreferences handles GET /user-preferences
func (h *UserPreferenceHandler) ListUserPreferences(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")
	userID := c.Query("userId")
	organizationID := c.Query("organizationId")
	category := c.Query("category")
	tags := c.QueryArray("tags")

	// Build filters
	filters := make(map[string]interface{})
	if userID != "" {
		filters["userId"] = userID
	}
	if organizationID != "" {
		filters["organizationId"] = organizationID
	}
	if category != "" {
		filters["category"] = category
	}
	if len(tags) > 0 {
		filters["tags"] = tags
	}

	// Convert to query
	query := queries.ListUserPreferencesQuery{
		Page:           page,
		Limit:          limit,
		SortBy:         sortBy,
		SortOrder:      sortOrder,
		UserID:         &userID,
		OrganizationID: &organizationID,
		Category:       &category,
		Tags:           tags,
		Filters:        filters,
	}

	// Execute query
	result, err := h.listUserPreferencesHandler.Handle(c.Request.Context(), query)
	if err != nil {
		if shared.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list user preferences", "details": err.Error()})
		return
	}

	// Convert to response DTO
	response := dto.ListUserPreferencesResponse{
		UserPreferences: make([]dto.UserPreferenceResponse, len(result.UserPreferences)),
		Total:           result.Total,
		Page:            result.Page,
		Limit:           result.Limit,
		TotalPages:      result.TotalPages,
	}

	for i, userPreference := range result.UserPreferences {
		response.UserPreferences[i] = dto.UserPreferenceResponse{
			ID:             userPreference.GetID().String(),
			UserID:         userPreference.GetUserID().String(),
			Category:       userPreference.GetCategory(),
			Key:            userPreference.GetKey(),
			Value:          userPreference.GetValue(),
			OrganizationID: userPreference.GetOrganizationID(),
			Description:    userPreference.GetDescription(),
			Tags:           userPreference.GetTags(),
			Metadata:       userPreference.GetMetadata(),
			CreatedBy:      userPreference.GetCreatedBy().String(),
			UpdatedBy:      userPreference.GetUpdatedBy(),
			CreatedAt:      userPreference.GetCreatedAt(),
			UpdatedAt:      userPreference.GetUpdatedAt(),
			Version:        userPreference.GetVersion(),
		}
	}

	c.JSON(http.StatusOK, response)
}
