package queries

import (
	"context"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/domain/rule"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
func (h *GetRuleHandler) Handle(ctx context.Context, query GetRuleQuery) (*rule.Rule, error) {
	tr := otel.Tracer("application")
	ctx, span := tr.Start(ctx, "GetRuleHandler.Handle")
	defer span.End()

	span.SetAttributes(attribute.String("rule.id", query.RuleID))

	ruleID, err := rule.RuleIDFromStr(query.RuleID)
	if err != nil {
		return nil, err // Can be NotFoundError or InfrastructureError
	}

	foundRule, err := h.ruleRepo.FindByID(ctx, ruleID)
	if err != nil {
		return nil, err // Can be NotFoundError or InfrastructureError
	}

	return foundRule, nil
}
