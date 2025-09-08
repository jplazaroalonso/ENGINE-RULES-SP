package commands

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// DeleteUserPreferenceCommand represents a command to delete a user preference
type DeleteUserPreferenceCommand struct {
	ID        string `json:"id" validate:"required,uuid"`
	DeletedBy string `json:"deletedBy" validate:"required,uuid"`
}

// DeleteUserPreferenceHandler handles the DeleteUserPreferenceCommand
type DeleteUserPreferenceHandler struct {
	userPreferenceRepo settings.UserPreferenceRepository
	eventBus           shared.EventBus
	validator          shared.StructValidator
}

// NewDeleteUserPreferenceHandler creates a new DeleteUserPreferenceHandler
func NewDeleteUserPreferenceHandler(
	userPreferenceRepo settings.UserPreferenceRepository,
	eventBus shared.EventBus,
	validator shared.StructValidator,
) *DeleteUserPreferenceHandler {
	return &DeleteUserPreferenceHandler{
		userPreferenceRepo: userPreferenceRepo,
		eventBus:           eventBus,
		validator:          validator,
	}
}

// Handle executes the DeleteUserPreferenceCommand
func (h *DeleteUserPreferenceHandler) Handle(ctx context.Context, cmd DeleteUserPreferenceCommand) error {
	// Validate the command
	if err := h.validator.Struct(cmd); err != nil {
		return fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Parse user preference ID
	userPreferenceID, err := settings.NewUserPreferenceIDFromString(cmd.ID)
	if err != nil {
		return err
	}

	// Parse deleted by user ID
	deletedBy, err := settings.NewUserIDFromString(cmd.DeletedBy)
	if err != nil {
		return err
	}

	// Find existing user preference
	existingUserPreference, err := h.userPreferenceRepo.FindByID(ctx, userPreferenceID)
	if err != nil {
		return fmt.Errorf("%w: user preference not found: %v", shared.ErrNotFound, err)
	}

	// Delete the user preference
	if err := h.userPreferenceRepo.Delete(ctx, userPreferenceID); err != nil {
		return fmt.Errorf("%w: failed to delete user preference: %v", shared.ErrInternalService, err)
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
