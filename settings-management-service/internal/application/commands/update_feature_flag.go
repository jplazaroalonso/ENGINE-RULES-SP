package commands

import (
	"context"
	"fmt"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// UpdateFeatureFlagCommand represents a command to update an existing feature flag
type UpdateFeatureFlagCommand struct {
	ID             string                   `json:"id" validate:"required,uuid"`
	IsEnabled      *bool                    `json:"isEnabled,omitempty"`
	Description    *string                  `json:"description,omitempty" validate:"omitempty,max=1000"`
	Variants       []settings.Variant       `json:"variants,omitempty"`
	TargetingRules []settings.TargetingRule `json:"targetingRules,omitempty"`
	Tags           []string                 `json:"tags,omitempty"`
	Metadata       map[string]interface{}   `json:"metadata,omitempty"`
	UpdatedBy      string                   `json:"updatedBy" validate:"required,uuid"`
}

// UpdateFeatureFlagHandler handles the UpdateFeatureFlagCommand
type UpdateFeatureFlagHandler struct {
	featureFlagRepo settings.FeatureFlagRepository
	eventBus        shared.EventBus
	validator       shared.StructValidator
}

// NewUpdateFeatureFlagHandler creates a new UpdateFeatureFlagHandler
func NewUpdateFeatureFlagHandler(
	featureFlagRepo settings.FeatureFlagRepository,
	eventBus shared.EventBus,
	validator shared.StructValidator,
) *UpdateFeatureFlagHandler {
	return &UpdateFeatureFlagHandler{
		featureFlagRepo: featureFlagRepo,
		eventBus:        eventBus,
		validator:       validator,
	}
}

// Handle executes the UpdateFeatureFlagCommand
func (h *UpdateFeatureFlagHandler) Handle(ctx context.Context, cmd UpdateFeatureFlagCommand) error {
	// Validate the command
	if err := h.validator.Struct(cmd); err != nil {
		return fmt.Errorf("%w: %s", shared.ErrValidation, err.Error())
	}

	// Parse feature flag ID
	featureFlagID, err := settings.NewFeatureFlagIDFromString(cmd.ID)
	if err != nil {
		return err
	}

	// Parse updated by user ID
	updatedBy, err := settings.NewUserIDFromString(cmd.UpdatedBy)
	if err != nil {
		return err
	}

	// Find existing feature flag
	existingFeatureFlag, err := h.featureFlagRepo.FindByID(ctx, featureFlagID)
	if err != nil {
		return fmt.Errorf("%w: feature flag not found: %v", shared.ErrNotFound, err)
	}

	// Track changes for event
	changes := make(map[string]interface{})

	// Update the feature flag
	if cmd.IsEnabled != nil {
		oldEnabled := existingFeatureFlag.GetIsEnabled()
		if err := existingFeatureFlag.SetEnabled(*cmd.IsEnabled, updatedBy); err != nil {
			return fmt.Errorf("%w: failed to update feature flag enabled state: %v", shared.ErrInvalidInput, err)
		}
		changes["isEnabled"] = map[string]interface{}{
			"old": oldEnabled,
			"new": *cmd.IsEnabled,
		}
	}

	if cmd.Description != nil {
		oldDescription := existingFeatureFlag.GetDescription()
		if err := existingFeatureFlag.SetDescription(*cmd.Description, updatedBy); err != nil {
			return fmt.Errorf("%w: failed to update feature flag description: %v", shared.ErrInvalidInput, err)
		}
		changes["description"] = map[string]interface{}{
			"old": oldDescription,
			"new": *cmd.Description,
		}
	}

	if cmd.Variants != nil {
		oldVariants := existingFeatureFlag.GetVariants()
		if err := existingFeatureFlag.SetVariants(cmd.Variants, updatedBy); err != nil {
			return fmt.Errorf("%w: failed to update feature flag variants: %v", shared.ErrInvalidInput, err)
		}
		changes["variants"] = map[string]interface{}{
			"old": oldVariants,
			"new": cmd.Variants,
		}
	}

	if cmd.TargetingRules != nil {
		oldTargetingRules := existingFeatureFlag.GetTargetingRules()
		if err := existingFeatureFlag.SetTargetingRules(cmd.TargetingRules, updatedBy); err != nil {
			return fmt.Errorf("%w: failed to update feature flag targeting rules: %v", shared.ErrInvalidInput, err)
		}
		changes["targetingRules"] = map[string]interface{}{
			"old": oldTargetingRules,
			"new": cmd.TargetingRules,
		}
	}

	if cmd.Tags != nil {
		oldTags := existingFeatureFlag.GetTags()
		if err := existingFeatureFlag.SetTags(cmd.Tags, updatedBy); err != nil {
			return fmt.Errorf("%w: failed to update feature flag tags: %v", shared.ErrInvalidInput, err)
		}
		changes["tags"] = map[string]interface{}{
			"old": oldTags,
			"new": cmd.Tags,
		}
	}

	if cmd.Metadata != nil {
		oldMetadata := existingFeatureFlag.GetMetadata()
		if err := existingFeatureFlag.SetMetadata(cmd.Metadata, updatedBy); err != nil {
			return fmt.Errorf("%w: failed to update feature flag metadata: %v", shared.ErrInvalidInput, err)
		}
		changes["metadata"] = map[string]interface{}{
			"old": oldMetadata,
			"new": cmd.Metadata,
		}
	}

	// Save the updated feature flag
	if err := h.featureFlagRepo.Update(ctx, existingFeatureFlag); err != nil {
		return fmt.Errorf("%w: failed to save updated feature flag: %v", shared.ErrInternalService, err)
	}

	// Publish domain events
	for _, event := range existingFeatureFlag.GetEvents() {
		if err := h.eventBus.Publish(ctx, event); err != nil {
			fmt.Printf("Warning: Failed to publish event %s for feature flag %s: %v\n", event.EventType(), existingFeatureFlag.GetID().String(), err)
		}
	}
	existingFeatureFlag.ClearEvents()

	return nil
}
