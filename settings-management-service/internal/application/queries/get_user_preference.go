package queries

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// GetUserPreferenceQuery represents a query to get a user preference by ID
type GetUserPreferenceQuery struct {
	ID string `json:"id" validate:"required,uuid"`
}

// GetUserPreferenceHandler handles the GetUserPreferenceQuery
type GetUserPreferenceHandler struct {
	userPreferenceRepo settings.UserPreferenceRepository
	validator          shared.StructValidator
}

// NewGetUserPreferenceHandler creates a new GetUserPreferenceHandler
func NewGetUserPreferenceHandler(
	userPreferenceRepo settings.UserPreferenceRepository,
	validator shared.StructValidator,
) *GetUserPreferenceHandler {
	return &GetUserPreferenceHandler{
		userPreferenceRepo: userPreferenceRepo,
		validator:          validator,
	}
}

// Handle executes the GetUserPreferenceQuery
func (h *GetUserPreferenceHandler) Handle(ctx context.Context, query GetUserPreferenceQuery) (*settings.UserPreference, error) {
	// Validate the query
	if err := h.validator.Struct(query); err != nil {
		return nil, fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Parse user preference ID
	userPreferenceID, err := settings.NewUserPreferenceIDFromString(query.ID)
	if err != nil {
		return nil, err
	}

	// Find the user preference
	userPreference, err := h.userPreferenceRepo.FindByID(ctx, userPreferenceID)
	if err != nil {
		return nil, fmt.Errorf("%w: user preference not found: %v", shared.ErrNotFound, err)
	}

	return userPreference, nil
}
