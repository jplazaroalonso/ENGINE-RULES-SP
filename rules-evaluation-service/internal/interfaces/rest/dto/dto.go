package dto

import "rules-evaluation-service/internal/domain/evaluation"

// EvaluationRequest defines the request body for an evaluation.
type EvaluationRequest struct {
	RuleCategory string             `json:"rule_category" binding:"required"`
	DSLContent   string             `json:"dsl_content" binding:"required"`
	Context      evaluation.Context `json:"context" binding:"required"`
}

// EvaluationResponse defines the API response for an evaluation.
type EvaluationResponse struct {
	Result evaluation.Result `json:"result"`
}

// ErrorResponse defines the structure for a generic error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
