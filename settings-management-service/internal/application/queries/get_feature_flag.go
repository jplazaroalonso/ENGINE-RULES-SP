package queries

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// GetFeatureFlagQuery represents a query to get a feature flag by ID
type GetFeatureFlagQuery struct {
	ID string `json:"id" validate:"required,uuid"`
}

// GetFeatureFlagHandler handles the GetFeatureFlagQuery
type GetFeatureFlagHandler struct {
	featureFlagRepo settings.FeatureFlagRepository
	validator       shared.StructValidator
}

// NewGetFeatureFlagHandler creates a new GetFeatureFlagHandler
func NewGetFeatureFlagHandler(
	featureFlagRepo settings.FeatureFlagRepository,
	validator shared.StructValidator,
) *GetFeatureFlagHandler {
	return &GetFeatureFlagHandler{
		featureFlagRepo: featureFlagRepo,
		validator:       validator,
	}
}

// Handle executes the GetFeatureFlagQuery
func (h *GetFeatureFlagHandler) Handle(ctx context.Context, query GetFeatureFlagQuery) (*settings.FeatureFlag, error) {
	// Validate the query
	if err := h.validator.Struct(query); err != nil {
		return nil, fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Parse feature flag ID
	featureFlagID, err := settings.NewFeatureFlagIDFromString(query.ID)
	if err != nil {
		return nil, err
	}

	// Find the feature flag
	featureFlag, err := h.featureFlagRepo.FindByID(ctx, featureFlagID)
	if err != nil {
		return nil, fmt.Errorf("%w: feature flag not found: %v", shared.ErrNotFound, err)
	}

	return featureFlag, nil
}
