package queries

import (
	"context"

	"rules-management-service/internal/domain/rule"
	"rules-management-service/internal/domain/shared"
)

// GetRuleQuery represents the query to get a rule by ID
type GetRuleQuery struct {
	RuleID string
}

// GetRuleResult represents the result of getting a rule
type GetRuleResult struct {
	RuleID      string `json:"rule_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DSLContent  string `json:"dsl_content"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	Version     int    `json:"version"`
}

// GetRuleHandler handles get rule queries
type GetRuleHandler struct {
	ruleRepo rule.Repository
}

func NewGetRuleHandler(ruleRepo rule.Repository) *GetRuleHandler {
	return &GetRuleHandler{ruleRepo: ruleRepo}
}

// Handle processes the get rule query
func (h *GetRuleHandler) Handle(ctx context.Context, query GetRuleQuery) (*GetRuleResult, error) {
	ruleID, err := rule.RuleIDFromStr(query.RuleID)
	if err != nil {
		return nil, shared.NewValidationError("invalid rule id", err)
	}

	foundRule, err := h.ruleRepo.FindByID(ctx, ruleID)
	if err != nil {
		return nil, err // Can be NotFoundError or InfrastructureError
	}

	return &GetRuleResult{
		RuleID:      foundRule.ID().String(),
		Name:        foundRule.Name(),
		Description: foundRule.Description(),
		DSLContent:  foundRule.DSLContent(),
		Status:      string(foundRule.Status()),
		Priority:    string(foundRule.Priority()),
		Version:     foundRule.Version(),
	}, nil
}
