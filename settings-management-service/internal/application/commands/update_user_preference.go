package commands

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// UpdateUserPreferenceCommand represents a command to update an existing user preference
type UpdateUserPreferenceCommand struct {
	ID          string                 `json:"id" validate:"required,uuid"`
	Value       interface{}            `json:"value" validate:"required"`
	Description *string                `json:"description,omitempty" validate:"omitempty,max=1000"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	UpdatedBy   string                 `json:"updatedBy" validate:"required,uuid"`
}

// UpdateUserPreferenceHandler handles the UpdateUserPreferenceCommand
type UpdateUserPreferenceHandler struct {
	userPreferenceRepo settings.UserPreferenceRepository
	eventBus           shared.EventBus
	validator          shared.StructValidator
}

// NewUpdateUserPreferenceHandler creates a new UpdateUserPreferenceHandler
func NewUpdateUserPreferenceHandler(
	userPreferenceRepo settings.UserPreferenceRepository,
	eventBus shared.EventBus,
	validator shared.StructValidator,
) *UpdateUserPreferenceHandler {
	return &UpdateUserPreferenceHandler{
		userPreferenceRepo: userPreferenceRepo,
		eventBus:           eventBus,
		validator:          validator,
	}
}

// Handle executes the UpdateUserPreferenceCommand
func (h *UpdateUserPreferenceHandler) Handle(ctx context.Context, cmd UpdateUserPreferenceCommand) error {
	// Validate the command
	if err := h.validator.Struct(cmd); err != nil {
		return fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Parse user preference ID
	userPreferenceID, err := settings.NewUserPreferenceIDFromString(cmd.ID)
	if err != nil {
		return err
	}

	// Parse updated by user ID
	updatedBy, err := settings.NewUserIDFromString(cmd.UpdatedBy)
	if err != nil {
		return err
	}

	// Find existing user preference
	existingUserPreference, err := h.userPreferenceRepo.FindByID(ctx, userPreferenceID)
	if err != nil {
		return fmt.Errorf("%w: user preference not found: %v", shared.ErrNotFound, err)
	}

	// Store old value for event
	oldValue := existingUserPreference.GetValue()

	// Update the user preference
	if err := existingUserPreference.UpdateUserPreference(
		cmd.Value,
		cmd.Description,
		cmd.Tags,
		cmd.Metadata,
		updatedBy,
	); err != nil {
		return fmt.Errorf("%w: failed to update user preference: %v", shared.ErrInvalidInput, err)
	}

	// Save the updated user preference
	if err := h.userPreferenceRepo.Update(ctx, existingUserPreference); err != nil {
		return fmt.Errorf("%w: failed to save updated user preference: %v", shared.ErrInternalService, err)
	}

	// Publish domain events
	for _, event := range existingUserPreference.GetEvents() {
		if err := h.eventBus.Publish(ctx, event); err != nil {
			fmt.Printf("Warning: Failed to publish event %s for user preference %s: %v\n", event.EventType(), existingUserPreference.GetID().String(), err)
		}
	}
	existingUserPreference.ClearEvents()

	return nil
}
