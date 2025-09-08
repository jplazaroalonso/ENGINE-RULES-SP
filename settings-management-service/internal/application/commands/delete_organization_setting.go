package commands

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// DeleteOrganizationSettingCommand represents a command to delete an organization setting
type DeleteOrganizationSettingCommand struct {
	ID        string `json:"id" validate:"required,uuid"`
	DeletedBy string `json:"deletedBy" validate:"required,uuid"`
}

// DeleteOrganizationSettingHandler handles the DeleteOrganizationSettingCommand
type DeleteOrganizationSettingHandler struct {
	organizationSettingRepo settings.OrganizationSettingRepository
	eventBus                shared.EventBus
	validator               shared.StructValidator
}

// NewDeleteOrganizationSettingHandler creates a new DeleteOrganizationSettingHandler
func NewDeleteOrganizationSettingHandler(
	organizationSettingRepo settings.OrganizationSettingRepository,
	eventBus shared.EventBus,
	validator shared.StructValidator,
) *DeleteOrganizationSettingHandler {
	return &DeleteOrganizationSettingHandler{
		organizationSettingRepo: organizationSettingRepo,
		eventBus:                eventBus,
		validator:               validator,
	}
}

// Handle executes the DeleteOrganizationSettingCommand
func (h *DeleteOrganizationSettingHandler) Handle(ctx context.Context, cmd DeleteOrganizationSettingCommand) error {
	// Validate the command
	if err := h.validator.Struct(cmd); err != nil {
		return fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Parse organization setting ID
	organizationSettingID, err := settings.NewOrganizationSettingIDFromString(cmd.ID)
	if err != nil {
		return err
	}

	// Parse deleted by user ID
	deletedBy, err := settings.NewUserIDFromString(cmd.DeletedBy)
	if err != nil {
		return err
	}

	// Find existing organization setting
	existingOrganizationSetting, err := h.organizationSettingRepo.FindByID(ctx, organizationSettingID)
	if err != nil {
		return fmt.Errorf("%w: organization setting not found: %v", shared.ErrNotFound, err)
	}

	// Delete the organization setting
	if err := h.organizationSettingRepo.Delete(ctx, organizationSettingID); err != nil {
		return fmt.Errorf("%w: failed to delete organization setting: %v", shared.ErrInternalService, err)
	}

	// Publish domain events
	for _, event := range existingOrganizationSetting.GetEvents() {
		if err := h.eventBus.Publish(ctx, event); err != nil {
			fmt.Printf("Warning: Failed to publish event %s for organization setting %s: %v\n", event.EventType(), existingOrganizationSetting.GetID().String(), err)
		}
	}
	existingOrganizationSetting.ClearEvents()

	return nil
}
