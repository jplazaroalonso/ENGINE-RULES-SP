package dto

import (
	"time"
)

// CreateConfigurationRequest represents the request to create a configuration
type CreateConfigurationRequest struct {
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

// UpdateConfigurationRequest represents the request to update a configuration
type UpdateConfigurationRequest struct {
	Value       interface{}            `json:"value" validate:"required"`
	Description *string                `json:"description,omitempty" validate:"omitempty,max=1000"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	UpdatedBy   string                 `json:"updatedBy" validate:"required,uuid"`
}

// DeleteConfigurationRequest represents the request to delete a configuration
type DeleteConfigurationRequest struct {
	DeletedBy string `json:"deletedBy" validate:"required,uuid"`
}

// ConfigurationResponse represents the response for a configuration
type ConfigurationResponse struct {
	ID             string                 `json:"id"`
	Key            string                 `json:"key"`
	Value          interface{}            `json:"value"`
	Environment    string                 `json:"environment"`
	OrganizationID *string                `json:"organizationId,omitempty"`
	Service        *string                `json:"service,omitempty"`
	Category       string                 `json:"category"`
	Description    *string                `json:"description,omitempty"`
	Tags           []string               `json:"tags"`
	Metadata       map[string]interface{} `json:"metadata"`
	CreatedBy      string                 `json:"createdBy"`
	UpdatedBy      *string                `json:"updatedBy,omitempty"`
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt"`
	Version        int                    `json:"version"`
}

// ListConfigurationsResponse represents the response for listing configurations
type ListConfigurationsResponse struct {
	Configurations []ConfigurationResponse `json:"configurations"`
	Total          int                     `json:"total"`
	Page           int                     `json:"page"`
	Limit          int                     `json:"limit"`
	TotalPages     int                     `json:"totalPages"`
}
