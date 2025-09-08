package commands

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// CreateUserPreferenceCommand represents a command to create a new user preference
type CreateUserPreferenceCommand struct {
	UserID         string                 `json:"userId" validate:"required,uuid"`
	Category       string                 `json:"category" validate:"required,min=1,max=100"`
	Key            string                 `json:"key" validate:"required,min=1,max=255"`
	Value          interface{}            `json:"value" validate:"required"`
	OrganizationID *string                `json:"organizationId,omitempty" validate:"omitempty,uuid"`
	Description    *string                `json:"description,omitempty" validate:"omitempty,max=1000"`
	Tags           []string               `json:"tags,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedBy      string                 `json:"createdBy" validate:"required,uuid"`
}

// CreateUserPreferenceHandler handles the CreateUserPreferenceCommand
type CreateUserPreferenceHandler struct {
	userPreferenceRepo settings.UserPreferenceRepository
	eventBus           shared.EventBus
	validator          shared.StructValidator
}

// NewCreateUserPreferenceHandler creates a new CreateUserPreferenceHandler
func NewCreateUserPreferenceHandler(
	userPreferenceRepo settings.UserPreferenceRepository,
	eventBus shared.EventBus,
	validator shared.StructValidator,
) *CreateUserPreferenceHandler {
	return &CreateUserPreferenceHandler{
		userPreferenceRepo: userPreferenceRepo,
		eventBus:           eventBus,
		validator:          validator,
	}
}

// Handle executes the CreateUserPreferenceCommand
func (h *CreateUserPreferenceHandler) Handle(ctx context.Context, cmd CreateUserPreferenceCommand) (settings.UserPreferenceID, error) {
	// Validate the command
	if err := h.validator.Struct(cmd); err != nil {
		return settings.UserPreferenceID{}, fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Parse user ID
	userID, err := settings.NewUserIDFromString(cmd.UserID)
	if err != nil {
		return settings.UserPreferenceID{}, err
	}

	// Parse organization ID if provided
	var organizationID *settings.OrganizationID
	if cmd.OrganizationID != nil {
		orgID, err := settings.NewOrganizationIDFromString(*cmd.OrganizationID)
		if err != nil {
			return settings.UserPreferenceID{}, err
		}
		organizationID = &orgID
	}

	// Parse created by user ID
	createdBy, err := settings.NewUserIDFromString(cmd.CreatedBy)
	if err != nil {
		return settings.UserPreferenceID{}, err
	}

	// Check if user preference already exists
	exists, err := h.userPreferenceRepo.ExistsByKey(ctx, userID, cmd.Category, cmd.Key, organizationID)
	if err != nil {
		return settings.UserPreferenceID{}, fmt.Errorf("%w: failed to check user preference existence: %v", shared.ErrInternalService, err)
	}
	if exists {
		return settings.UserPreferenceID{}, fmt.Errorf("%w: user preference with key '%s' already exists for user '%s' in category '%s'", shared.ErrAlreadyExists, cmd.Key, userID.String(), cmd.Category)
	}

	// Create new user preference
	userPreference, err := settings.NewUserPreference(
		userID,
		cmd.Category,
		cmd.Key,
		cmd.Value,
		organizationID,
		cmd.Description,
		cmd.Tags,
		cmd.Metadata,
		createdBy,
	)
	if err != nil {
		return settings.UserPreferenceID{}, fmt.Errorf("%w: failed to create user preference: %v", shared.ErrInvalidInput, err)
	}

	// Save the user preference
	if err := h.userPreferenceRepo.Save(ctx, userPreference); err != nil {
		return settings.UserPreferenceID{}, fmt.Errorf("%w: failed to save user preference: %v", shared.ErrInternalService, err)
	}

	// Publish domain events
	for _, event := range userPreference.GetEvents() {
		if err := h.eventBus.Publish(ctx, event); err != nil {
			// Log the error but don't block the main flow
			fmt.Printf("Warning: Failed to publish event %s for user preference %s: %v\n", event.EventType(), userPreference.GetID().String(), err)
		}
	}
	userPreference.ClearEvents()

	return userPreference.GetID(), nil
}
