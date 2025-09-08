package dto

import (
	"time"
)

// CreateUserPreferenceRequest represents the request to create a user preference
type CreateUserPreferenceRequest struct {
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

// UpdateUserPreferenceRequest represents the request to update a user preference
type UpdateUserPreferenceRequest struct {
	Value       interface{}            `json:"value" validate:"required"`
	Description *string                `json:"description,omitempty" validate:"omitempty,max=1000"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	UpdatedBy   string                 `json:"updatedBy" validate:"required,uuid"`
}

// DeleteUserPreferenceRequest represents the request to delete a user preference
type DeleteUserPreferenceRequest struct {
	DeletedBy string `json:"deletedBy" validate:"required,uuid"`
}

// UserPreferenceResponse represents the response for a user preference
type UserPreferenceResponse struct {
	ID             string                 `json:"id"`
	UserID         string                 `json:"userId"`
	Category       string                 `json:"category"`
	Key            string                 `json:"key"`
	Value          interface{}            `json:"value"`
	OrganizationID *string                `json:"organizationId,omitempty"`
	Description    *string                `json:"description,omitempty"`
	Tags           []string               `json:"tags"`
	Metadata       map[string]interface{} `json:"metadata"`
	CreatedBy      string                 `json:"createdBy"`
	UpdatedBy      *string                `json:"updatedBy,omitempty"`
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt"`
	Version        int                    `json:"version"`
}

// ListUserPreferencesResponse represents the response for listing user preferences
type ListUserPreferencesResponse struct {
	UserPreferences []UserPreferenceResponse `json:"userPreferences"`
	Total           int                      `json:"total"`
	Page            int                      `json:"page"`
	Limit           int                      `json:"limit"`
	TotalPages      int                      `json:"totalPages"`
}
