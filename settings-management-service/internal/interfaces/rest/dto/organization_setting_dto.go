package dto

import (
	"time"
)

// CreateOrganizationSettingRequest represents the request to create an organization setting
type CreateOrganizationSettingRequest struct {
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

// UpdateOrganizationSettingRequest represents the request to update an organization setting
type UpdateOrganizationSettingRequest struct {
	Value       interface{}            `json:"value" validate:"required"`
	Description *string                `json:"description,omitempty" validate:"omitempty,max=1000"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	UpdatedBy   string                 `json:"updatedBy" validate:"required,uuid"`
}

// DeleteOrganizationSettingRequest represents the request to delete an organization setting
type DeleteOrganizationSettingRequest struct {
	DeletedBy string `json:"deletedBy" validate:"required,uuid"`
}

// OrganizationSettingResponse represents the response for an organization setting
type OrganizationSettingResponse struct {
	ID             string                 `json:"id"`
	OrganizationID string                 `json:"organizationId"`
	Category       string                 `json:"category"`
	Key            string                 `json:"key"`
	Value          interface{}            `json:"value"`
	ParentID       *string                `json:"parentId,omitempty"`
	Description    *string                `json:"description,omitempty"`
	Tags           []string               `json:"tags"`
	Metadata       map[string]interface{} `json:"metadata"`
	CreatedBy      string                 `json:"createdBy"`
	UpdatedBy      *string                `json:"updatedBy,omitempty"`
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt"`
	Version        int                    `json:"version"`
}

// ListOrganizationSettingsResponse represents the response for listing organization settings
type ListOrganizationSettingsResponse struct {
	OrganizationSettings []OrganizationSettingResponse `json:"organizationSettings"`
	Total                int                           `json:"total"`
	Page                 int                           `json:"page"`
	Limit                int                           `json:"limit"`
	TotalPages           int                           `json:"totalPages"`
}
