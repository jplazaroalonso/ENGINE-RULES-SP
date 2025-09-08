package commands

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// UpdateConfigurationCommand represents a command to update an existing configuration
type UpdateConfigurationCommand struct {
	ID          string                 `json:"id" validate:"required,uuid"`
	Value       interface{}            `json:"value" validate:"required"`
	Description *string                `json:"description,omitempty" validate:"omitempty,max=1000"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	UpdatedBy   string                 `json:"updatedBy" validate:"required,uuid"`
}

// UpdateConfigurationHandler handles the UpdateConfigurationCommand
type UpdateConfigurationHandler struct {
	configRepo settings.ConfigurationRepository
	eventBus   shared.EventBus
	validator  shared.StructValidator
}

// NewUpdateConfigurationHandler creates a new UpdateConfigurationHandler
func NewUpdateConfigurationHandler(
	configRepo settings.ConfigurationRepository,
	eventBus shared.EventBus,
	validator shared.StructValidator,
) *UpdateConfigurationHandler {
	return &UpdateConfigurationHandler{
		configRepo: configRepo,
		eventBus:   eventBus,
		validator:  validator,
	}
}

// Handle executes the UpdateConfigurationCommand
func (h *UpdateConfigurationHandler) Handle(ctx context.Context, cmd UpdateConfigurationCommand) error {
	// Validate the command
	if err := h.validator.Struct(cmd); err != nil {
		return fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Parse configuration ID
	configID, err := settings.NewConfigurationIDFromString(cmd.ID)
	if err != nil {
		return err
	}

	// Parse updated by user ID
	updatedBy, err := settings.NewUserIDFromString(cmd.UpdatedBy)
	if err != nil {
		return err
	}

	// Find existing configuration
	existingConfig, err := h.configRepo.FindByID(ctx, configID)
	if err != nil {
		return fmt.Errorf("%w: configuration not found: %v", shared.ErrNotFound, err)
	}

	// Store old value for event
	oldValue := existingConfig.GetValue()

	// Update the configuration
	if err := existingConfig.UpdateConfiguration(
		cmd.Value,
		cmd.Description,
		cmd.Tags,
		cmd.Metadata,
		updatedBy,
	); err != nil {
		return fmt.Errorf("%w: failed to update configuration: %v", shared.ErrInvalidInput, err)
	}

	// Save the updated configuration
	if err := h.configRepo.Update(ctx, existingConfig); err != nil {
		return fmt.Errorf("%w: failed to save updated configuration: %v", shared.ErrInternalService, err)
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
