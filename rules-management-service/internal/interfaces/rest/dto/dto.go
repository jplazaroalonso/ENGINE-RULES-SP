package dto

import "time"

// CreateRuleRequest defines the request body for creating a rule.
type CreateRuleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	DSLContent  string   `json:"dsl_content" binding:"required"`
	Priority    string   `json:"priority" binding:"required,oneof=LOW MEDIUM HIGH CRITICAL"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
}

// ValidateRuleRequest defines the request body for validating a rule's DSL.
type ValidateRuleRequest struct {
	DSLContent string `json:"dsl_content" binding:"required"`
}

// RuleResponse defines the structure for a rule in an API response.
type RuleResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DSLContent  string    `json:"dsl_content"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   string    `json:"created_by"`
}

// ErrorResponse defines the structure for a generic error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
