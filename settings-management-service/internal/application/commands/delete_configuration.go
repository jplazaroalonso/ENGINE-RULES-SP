package commands

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// DeleteConfigurationCommand represents a command to delete a configuration
type DeleteConfigurationCommand struct {
	ID        string `json:"id" validate:"required,uuid"`
	DeletedBy string `json:"deletedBy" validate:"required,uuid"`
}

// DeleteConfigurationHandler handles the DeleteConfigurationCommand
type DeleteConfigurationHandler struct {
	configRepo settings.ConfigurationRepository
	eventBus   shared.EventBus
	validator  shared.StructValidator
}

// NewDeleteConfigurationHandler creates a new DeleteConfigurationHandler
func NewDeleteConfigurationHandler(
	configRepo settings.ConfigurationRepository,
	eventBus shared.EventBus,
	validator shared.StructValidator,
) *DeleteConfigurationHandler {
	return &DeleteConfigurationHandler{
		configRepo: configRepo,
		eventBus:   eventBus,
		validator:  validator,
	}
}

// Handle executes the DeleteConfigurationCommand
func (h *DeleteConfigurationHandler) Handle(ctx context.Context, cmd DeleteConfigurationCommand) error {
	// Validate the command
	if err := h.validator.Struct(cmd); err != nil {
		return fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Parse configuration ID
	configID, err := settings.NewConfigurationIDFromString(cmd.ID)
	if err != nil {
		return err
	}

	// Parse deleted by user ID
	deletedBy, err := settings.NewUserIDFromString(cmd.DeletedBy)
	if err != nil {
		return err
	}

	// Find existing configuration
	existingConfig, err := h.configRepo.FindByID(ctx, configID)
	if err != nil {
		return fmt.Errorf("%w: configuration not found: %v", shared.ErrNotFound, err)
	}

	// Delete the configuration
	if err := h.configRepo.Delete(ctx, configID); err != nil {
		return fmt.Errorf("%w: failed to delete configuration: %v", shared.ErrInternalService, err)
	}

	// Publish domain events
	for _, event := range existingConfig.GetEvents() {
		if err := h.eventBus.Publish(ctx, event); err != nil {
			fmt.Printf("Warning: Failed to publish event %s for configuration %s: %v\n", event.EventType(), existingConfig.GetID().String(), err)
		}
	}
	existingConfig.ClearEvents()

	return nil
}
