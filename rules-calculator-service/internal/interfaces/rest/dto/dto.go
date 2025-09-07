package dto

// CalculationRequest is the DTO for a calculation request.
type CalculationRequest struct {
	RuleIDs []string               `json:"rule_ids" binding:"required"`
	Context map[string]interface{} `json:"context" binding:"required"`
}

// CalculationResponse is the DTO for a calculation response.
type CalculationResponse struct {
	CalculationID string             `json:"calculation_id"`
	Value         float64            `json:"value"`
	Breakdown     map[string]float64 `json:"breakdown"`
}

// ErrorResponse is the DTO for an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}
