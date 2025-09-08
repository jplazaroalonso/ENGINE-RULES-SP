package commands

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// DeleteFeatureFlagCommand represents a command to delete a feature flag
type DeleteFeatureFlagCommand struct {
	ID        string `json:"id" validate:"required,uuid"`
	DeletedBy string `json:"deletedBy" validate:"required,uuid"`
}

// DeleteFeatureFlagHandler handles the DeleteFeatureFlagCommand
type DeleteFeatureFlagHandler struct {
	featureFlagRepo settings.FeatureFlagRepository
	eventBus        shared.EventBus
	validator       shared.StructValidator
}

// NewDeleteFeatureFlagHandler creates a new DeleteFeatureFlagHandler
func NewDeleteFeatureFlagHandler(
	featureFlagRepo settings.FeatureFlagRepository,
	eventBus shared.EventBus,
	validator shared.StructValidator,
) *DeleteFeatureFlagHandler {
	return &DeleteFeatureFlagHandler{
		featureFlagRepo: featureFlagRepo,
		eventBus:        eventBus,
		validator:       validator,
	}
}

// Handle executes the DeleteFeatureFlagCommand
func (h *DeleteFeatureFlagHandler) Handle(ctx context.Context, cmd DeleteFeatureFlagCommand) error {
	// Validate the command
	if err := h.validator.Struct(cmd); err != nil {
		return fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Parse feature flag ID
	featureFlagID, err := settings.NewFeatureFlagIDFromString(cmd.ID)
	if err != nil {
		return err
	}

	// Parse deleted by user ID
	deletedBy, err := settings.NewUserIDFromString(cmd.DeletedBy)
	if err != nil {
		return err
	}

	// Find existing feature flag
	existingFeatureFlag, err := h.featureFlagRepo.FindByID(ctx, featureFlagID)
	if err != nil {
		return fmt.Errorf("%w: feature flag not found: %v", shared.ErrNotFound, err)
	}

	// Delete the feature flag
	if err := h.featureFlagRepo.Delete(ctx, featureFlagID); err != nil {
		return fmt.Errorf("%w: failed to delete feature flag: %v", shared.ErrInternalService, err)
	}

	// Publish domain events
	for _, event := range existingFeatureFlag.GetEvents() {
		if err := h.eventBus.Publish(ctx, event); err != nil {
			fmt.Printf("Warning: Failed to publish event %s for feature flag %s: %v\n", event.EventType(), existingFeatureFlag.GetID().String(), err)
		}
	}
	existingFeatureFlag.ClearEvents()

	return nil
}
