package queries

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// GetConfigurationQuery represents a query to get a configuration by ID
type GetConfigurationQuery struct {
	ID string `json:"id" validate:"required,uuid"`
}

// GetConfigurationHandler handles the GetConfigurationQuery
type GetConfigurationHandler struct {
	configRepo settings.ConfigurationRepository
	validator  shared.StructValidator
}

// NewGetConfigurationHandler creates a new GetConfigurationHandler
func NewGetConfigurationHandler(
	configRepo settings.ConfigurationRepository,
	validator shared.StructValidator,
) *GetConfigurationHandler {
	return &GetConfigurationHandler{
		configRepo: configRepo,
		validator:  validator,
	}
}

// Handle executes the GetConfigurationQuery
func (h *GetConfigurationHandler) Handle(ctx context.Context, query GetConfigurationQuery) (*settings.Configuration, error) {
	// Validate the query
	if err := h.validator.Struct(query); err != nil {
		return nil, fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Parse configuration ID
	configID, err := settings.NewConfigurationIDFromString(query.ID)
	if err != nil {
		return nil, err
	}

	// Find the configuration
	configuration, err := h.configRepo.FindByID(ctx, configID)
	if err != nil {
		return nil, fmt.Errorf("%w: configuration not found: %v", shared.ErrNotFound, err)
	}

	return configuration, nil
}
