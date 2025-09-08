package commands

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// CreateOrganizationSettingCommand represents a command to create a new organization setting
type CreateOrganizationSettingCommand struct {
	OrganizationID string                 `json:"organizationId" validate:"required,uuid"`
	Category       string                 `json:"category" validate:"required,min=1,max=100"`
	Key            string                 `json:"key" validate:"required,min=1,max=255"`
	Value          interface{}            `json:"value" validate:"required"`
	ParentID       *string                `json:"parentId,omitempty" validate:"omitempty,uuid"`
	Description    *string                `json:"description,omitempty" validate:"omitempty,max=1000"`
	Tags           []string               `json:"tags,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedBy      string                 `json:"createdBy" validate:"required,uuid"`
}

// CreateOrganizationSettingHandler handles the CreateOrganizationSettingCommand
type CreateOrganizationSettingHandler struct {
	organizationSettingRepo settings.OrganizationSettingRepository
	eventBus                shared.EventBus
	validator               shared.StructValidator
}

// NewCreateOrganizationSettingHandler creates a new CreateOrganizationSettingHandler
func NewCreateOrganizationSettingHandler(
	organizationSettingRepo settings.OrganizationSettingRepository,
	eventBus shared.EventBus,
	validator shared.StructValidator,
) *CreateOrganizationSettingHandler {
	return &CreateOrganizationSettingHandler{
		organizationSettingRepo: organizationSettingRepo,
		eventBus:                eventBus,
		validator:               validator,
	}
}

// Handle executes the CreateOrganizationSettingCommand
func (h *CreateOrganizationSettingHandler) Handle(ctx context.Context, cmd CreateOrganizationSettingCommand) (settings.OrganizationSettingID, error) {
	// Validate the command
	if err := h.validator.Struct(cmd); err != nil {
		return settings.OrganizationSettingID{}, fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Parse organization ID
	organizationID, err := settings.NewOrganizationIDFromString(cmd.OrganizationID)
	if err != nil {
		return settings.OrganizationSettingID{}, err
	}

	// Parse parent ID if provided
	var parentID *settings.OrganizationSettingID
	if cmd.ParentID != nil {
		parent, err := settings.NewOrganizationSettingIDFromString(*cmd.ParentID)
		if err != nil {
			return settings.OrganizationSettingID{}, err
		}
		parentID = &parent
	}

	// Parse created by user ID
	createdBy, err := settings.NewUserIDFromString(cmd.CreatedBy)
	if err != nil {
		return settings.OrganizationSettingID{}, err
	}

	// Check if organization setting already exists
	exists, err := h.organizationSettingRepo.ExistsByKey(ctx, organizationID, cmd.Category, cmd.Key)
	if err != nil {
		return settings.OrganizationSettingID{}, fmt.Errorf("%w: failed to check organization setting existence: %v", shared.ErrInternalService, err)
	}
	if exists {
		return settings.OrganizationSettingID{}, fmt.Errorf("%w: organization setting with key '%s' already exists for organization '%s' in category '%s'", shared.ErrAlreadyExists, cmd.Key, organizationID.String(), cmd.Category)
	}

	// Create new organization setting
	organizationSetting, err := settings.NewOrganizationSetting(
		organizationID,
		cmd.Category,
		cmd.Key,
		cmd.Value,
		parentID,
		cmd.Description,
		cmd.Tags,
		cmd.Metadata,
		createdBy,
	)
	if err != nil {
		return settings.OrganizationSettingID{}, fmt.Errorf("%w: failed to create organization setting: %v", shared.ErrInvalidInput, err)
	}

	// Save the organization setting
	if err := h.organizationSettingRepo.Save(ctx, organizationSetting); err != nil {
		return settings.OrganizationSettingID{}, fmt.Errorf("%w: failed to save organization setting: %v", shared.ErrInternalService, err)
	}

	// Publish domain events
	for _, event := range organizationSetting.GetEvents() {
		if err := h.eventBus.Publish(ctx, event); err != nil {
			// Log the error but don't block the main flow
			fmt.Printf("Warning: Failed to publish event %s for organization setting %s: %v\n", event.EventType(), organizationSetting.GetID().String(), err)
		}
	}
	organizationSetting.ClearEvents()

	return organizationSetting.GetID(), nil
}
