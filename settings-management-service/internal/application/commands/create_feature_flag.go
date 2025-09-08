package commands

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// CreateFeatureFlagCommand represents a command to create a new feature flag
type CreateFeatureFlagCommand struct {
	Key            string                   `json:"key" validate:"required,min=1,max=255"`
	IsEnabled      bool                     `json:"isEnabled"`
	Environment    string                   `json:"environment" validate:"required,oneof=DEVELOPMENT STAGING PRODUCTION"`
	OrganizationID *string                  `json:"organizationId,omitempty" validate:"omitempty,uuid"`
	Service        *string                  `json:"service,omitempty" validate:"omitempty,min=1,max=100"`
	Category       string                   `json:"category" validate:"required,min=1,max=100"`
	Description    *string                  `json:"description,omitempty" validate:"omitempty,max=1000"`
	Variants       []settings.Variant       `json:"variants,omitempty"`
	TargetingRules []settings.TargetingRule `json:"targetingRules,omitempty"`
	Tags           []string                 `json:"tags,omitempty"`
	Metadata       map[string]interface{}   `json:"metadata,omitempty"`
	CreatedBy      string                   `json:"createdBy" validate:"required,uuid"`
}

// CreateFeatureFlagHandler handles the CreateFeatureFlagCommand
type CreateFeatureFlagHandler struct {
	featureFlagRepo settings.FeatureFlagRepository
	eventBus        shared.EventBus
	validator       shared.StructValidator
}

// NewCreateFeatureFlagHandler creates a new CreateFeatureFlagHandler
func NewCreateFeatureFlagHandler(
	featureFlagRepo settings.FeatureFlagRepository,
	eventBus shared.EventBus,
	validator shared.StructValidator,
) *CreateFeatureFlagHandler {
	return &CreateFeatureFlagHandler{
		featureFlagRepo: featureFlagRepo,
		eventBus:        eventBus,
		validator:       validator,
	}
}

// Handle executes the CreateFeatureFlagCommand
func (h *CreateFeatureFlagHandler) Handle(ctx context.Context, cmd CreateFeatureFlagCommand) (settings.FeatureFlagID, error) {
	// Validate the command
	if err := h.validator.Struct(cmd); err != nil {
		return settings.FeatureFlagID{}, fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Parse environment
	environment, err := settings.ParseEnvironment(cmd.Environment)
	if err != nil {
		return settings.FeatureFlagID{}, err
	}

	// Parse organization ID if provided
	var organizationID *settings.OrganizationID
	if cmd.OrganizationID != nil {
		orgID, err := settings.NewOrganizationIDFromString(*cmd.OrganizationID)
		if err != nil {
			return settings.FeatureFlagID{}, err
		}
		organizationID = &orgID
	}

	// Parse service name if provided
	var service *settings.ServiceName
	if cmd.Service != nil {
		svc, err := settings.NewServiceName(*cmd.Service)
		if err != nil {
			return settings.FeatureFlagID{}, err
		}
		service = &svc
	}

	// Parse created by user ID
	createdBy, err := settings.NewUserIDFromString(cmd.CreatedBy)
	if err != nil {
		return settings.FeatureFlagID{}, err
	}

	// Check if feature flag already exists
	exists, err := h.featureFlagRepo.ExistsByKey(ctx, cmd.Key, environment, organizationID, service)
	if err != nil {
		return settings.FeatureFlagID{}, fmt.Errorf("%w: failed to check feature flag existence: %v", shared.ErrInternalService, err)
	}
	if exists {
		return settings.FeatureFlagID{}, fmt.Errorf("%w: feature flag with key '%s' already exists for environment '%s'", shared.ErrAlreadyExists, cmd.Key, environment.String())
	}

	// Create new feature flag
	featureFlag, err := settings.NewFeatureFlag(
		cmd.Key,
		cmd.IsEnabled,
		environment,
		organizationID,
		service,
		cmd.Category,
		cmd.Description,
		cmd.Variants,
		cmd.TargetingRules,
		cmd.Tags,
		cmd.Metadata,
		createdBy,
	)
	if err != nil {
		return settings.FeatureFlagID{}, fmt.Errorf("%w: failed to create feature flag: %v", shared.ErrInvalidInput, err)
	}

	// Save the feature flag
	if err := h.featureFlagRepo.Save(ctx, featureFlag); err != nil {
		return settings.FeatureFlagID{}, fmt.Errorf("%w: failed to save feature flag: %v", shared.ErrInternalService, err)
	}

	// Publish domain events
	for _, event := range featureFlag.GetEvents() {
		if err := h.eventBus.Publish(ctx, event); err != nil {
			// Log the error but don't block the main flow
			fmt.Printf("Warning: Failed to publish event %s for feature flag %s: %v\n", event.EventType(), featureFlag.GetID().String(), err)
		}
	}
	featureFlag.ClearEvents()

	return featureFlag.GetID(), nil
}
