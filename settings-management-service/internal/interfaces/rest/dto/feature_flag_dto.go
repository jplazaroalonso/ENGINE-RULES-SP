package dto

import (
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
)

// CreateFeatureFlagRequest represents the request to create a feature flag
type CreateFeatureFlagRequest struct {
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

// UpdateFeatureFlagRequest represents the request to update a feature flag
type UpdateFeatureFlagRequest struct {
	IsEnabled      *bool                    `json:"isEnabled,omitempty"`
	Description    *string                  `json:"description,omitempty" validate:"omitempty,max=1000"`
	Variants       []settings.Variant       `json:"variants,omitempty"`
	TargetingRules []settings.TargetingRule `json:"targetingRules,omitempty"`
	Tags           []string                 `json:"tags,omitempty"`
	Metadata       map[string]interface{}   `json:"metadata,omitempty"`
	UpdatedBy      string                   `json:"updatedBy" validate:"required,uuid"`
}

// DeleteFeatureFlagRequest represents the request to delete a feature flag
type DeleteFeatureFlagRequest struct {
	DeletedBy string `json:"deletedBy" validate:"required,uuid"`
}

// FeatureFlagResponse represents the response for a feature flag
type FeatureFlagResponse struct {
	ID             string                   `json:"id"`
	Key            string                   `json:"key"`
	IsEnabled      bool                     `json:"isEnabled"`
	Environment    string                   `json:"environment"`
	OrganizationID *string                  `json:"organizationId,omitempty"`
	Service        *string                  `json:"service,omitempty"`
	Category       string                   `json:"category"`
	Description    *string                  `json:"description,omitempty"`
	Variants       []settings.Variant       `json:"variants"`
	TargetingRules []settings.TargetingRule `json:"targetingRules"`
	Tags           []string                 `json:"tags"`
	Metadata       map[string]interface{}   `json:"metadata"`
	CreatedBy      string                   `json:"createdBy"`
	UpdatedBy      *string                  `json:"updatedBy,omitempty"`
	CreatedAt      time.Time                `json:"createdAt"`
	UpdatedAt      time.Time                `json:"updatedAt"`
	Version        int                      `json:"version"`
}

// ListFeatureFlagsResponse represents the response for listing feature flags
type ListFeatureFlagsResponse struct {
	FeatureFlags []FeatureFlagResponse `json:"featureFlags"`
	Total        int                   `json:"total"`
	Page         int                   `json:"page"`
	Limit        int                   `json:"limit"`
	TotalPages   int                   `json:"totalPages"`
}
