package commands

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// CreateConfigurationCommand represents a command to create a new configuration
type CreateConfigurationCommand struct {
	Key            string                 `json:"key" validate:"required,min=1,max=255"`
	Value          interface{}            `json:"value" validate:"required"`
	Environment    string                 `json:"environment" validate:"required,oneof=DEVELOPMENT STAGING PRODUCTION"`
	OrganizationID *string                `json:"organizationId,omitempty" validate:"omitempty,uuid"`
	Service        *string                `json:"service,omitempty" validate:"omitempty,min=1,max=100"`
	Category       string                 `json:"category" validate:"required,min=1,max=100"`
	Description    *string                `json:"description,omitempty" validate:"omitempty,max=1000"`
	Tags           []string               `json:"tags,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedBy      string                 `json:"createdBy" validate:"required,uuid"`
}

// CreateConfigurationHandler handles the CreateConfigurationCommand
type CreateConfigurationHandler struct {
	configRepo settings.ConfigurationRepository
	eventBus   shared.EventBus
	validator  shared.StructValidator
}

// NewCreateConfigurationHandler creates a new CreateConfigurationHandler
func NewCreateConfigurationHandler(
	configRepo settings.ConfigurationRepository,
	eventBus shared.EventBus,
	validator shared.StructValidator,
) *CreateConfigurationHandler {
	return &CreateConfigurationHandler{
		configRepo: configRepo,
		eventBus:   eventBus,
		validator:  validator,
	}
}

// Handle executes the CreateConfigurationCommand
func (h *CreateConfigurationHandler) Handle(ctx context.Context, cmd CreateConfigurationCommand) (settings.ConfigurationID, error) {
	// Validate the command
	if err := h.validator.Struct(cmd); err != nil {
		return settings.ConfigurationID{}, fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Parse environment
	environment, err := settings.ParseEnvironment(cmd.Environment)
	if err != nil {
		return settings.ConfigurationID{}, err
	}

	// Parse organization ID if provided
	var organizationID *settings.OrganizationID
	if cmd.OrganizationID != nil {
		orgID, err := settings.NewOrganizationIDFromString(*cmd.OrganizationID)
		if err != nil {
			return settings.ConfigurationID{}, err
		}
		organizationID = &orgID
	}

	// Parse service name if provided
	var service *settings.ServiceName
	if cmd.Service != nil {
		svc, err := settings.NewServiceName(*cmd.Service)
		if err != nil {
			return settings.ConfigurationID{}, err
		}
		service = &svc
	}

	// Parse created by user ID
	createdBy, err := settings.NewUserIDFromString(cmd.CreatedBy)
	if err != nil {
		return settings.ConfigurationID{}, err
	}

	// Check if configuration already exists
	exists, err := h.configRepo.ExistsByKey(ctx, cmd.Key, environment, organizationID, service)
	if err != nil {
		return settings.ConfigurationID{}, fmt.Errorf("%w: failed to check configuration existence: %v", shared.ErrInternalService, err)
	}
	if exists {
		return settings.ConfigurationID{}, fmt.Errorf("%w: configuration with key '%s' already exists for environment '%s'", shared.ErrAlreadyExists, cmd.Key, environment.String())
	}

	// Create new configuration
	configuration, err := settings.NewConfiguration(
		cmd.Key,
		cmd.Value,
		environment,
		organizationID,
		service,
		cmd.Category,
		cmd.Description,
		cmd.Tags,
		cmd.Metadata,
		createdBy,
	)
	if err != nil {
		return settings.ConfigurationID{}, fmt.Errorf("%w: failed to create configuration: %v", shared.ErrInvalidInput, err)
	}

	// Save the configuration
	if err := h.configRepo.Save(ctx, configuration); err != nil {
		return settings.ConfigurationID{}, fmt.Errorf("%w: failed to save configuration: %v", shared.ErrInternalService, err)
	}

	// Publish domain events
	for _, event := range configuration.GetEvents() {
		if err := h.eventBus.Publish(ctx, event); err != nil {
			// Log the error but don't block the main flow
			fmt.Printf("Warning: Failed to publish event %s for configuration %s: %v\n", event.EventType(), configuration.GetID().String(), err)
		}
	}
	configuration.ClearEvents()

	return configuration.GetID(), nil
}
