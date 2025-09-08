package commands

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// UpdateOrganizationSettingCommand represents a command to update an existing organization setting
type UpdateOrganizationSettingCommand struct {
	ID          string                 `json:"id" validate:"required,uuid"`
	Value       interface{}            `json:"value" validate:"required"`
	Description *string                `json:"description,omitempty" validate:"omitempty,max=1000"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	UpdatedBy   string                 `json:"updatedBy" validate:"required,uuid"`
}

// UpdateOrganizationSettingHandler handles the UpdateOrganizationSettingCommand
type UpdateOrganizationSettingHandler struct {
	organizationSettingRepo settings.OrganizationSettingRepository
	eventBus                shared.EventBus
	validator               shared.StructValidator
}

// NewUpdateOrganizationSettingHandler creates a new UpdateOrganizationSettingHandler
func NewUpdateOrganizationSettingHandler(
	organizationSettingRepo settings.OrganizationSettingRepository,
	eventBus shared.EventBus,
	validator shared.StructValidator,
) *UpdateOrganizationSettingHandler {
	return &UpdateOrganizationSettingHandler{
		organizationSettingRepo: organizationSettingRepo,
		eventBus:                eventBus,
		validator:               validator,
	}
}

// Handle executes the UpdateOrganizationSettingCommand
func (h *UpdateOrganizationSettingHandler) Handle(ctx context.Context, cmd UpdateOrganizationSettingCommand) error {
	// Validate the command
	if err := h.validator.Struct(cmd); err != nil {
		return fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Parse organization setting ID
	organizationSettingID, err := settings.NewOrganizationSettingIDFromString(cmd.ID)
	if err != nil {
		return err
	}

	// Parse updated by user ID
	updatedBy, err := settings.NewUserIDFromString(cmd.UpdatedBy)
	if err != nil {
		return err
	}

	// Find existing organization setting
	existingOrganizationSetting, err := h.organizationSettingRepo.FindByID(ctx, organizationSettingID)
	if err != nil {
		return fmt.Errorf("%w: organization setting not found: %v", shared.ErrNotFound, err)
	}

	// Store old value for event
	oldValue := existingOrganizationSetting.GetValue()

	// Update the organization setting
	if err := existingOrganizationSetting.UpdateOrganizationSetting(
		cmd.Value,
		cmd.Description,
		cmd.Tags,
		cmd.Metadata,
		updatedBy,
	); err != nil {
		return fmt.Errorf("%w: failed to update organization setting: %v", shared.ErrInvalidInput, err)
	}

	// Save the updated organization setting
	if err := h.organizationSettingRepo.Update(ctx, existingOrganizationSetting); err != nil {
		return fmt.Errorf("%w: failed to save updated organization setting: %v", shared.ErrInternalService, err)
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
