package application

import (
	"context"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-calculator-service/internal/domain/calculation"
)

// CalculateRulesCommand is the command for calculating rules.
type CalculateRulesCommand struct {
	RuleIDs []string               `json:"rule_ids"`
	Context map[string]interface{} `json:"context"`
}

// CalculateRulesResult is the result of calculating rules.
type CalculateRulesResult struct {
	CalculationID string             `json:"calculation_id"`
	Value         float64            `json:"value"`
	Breakdown     map[string]float64 `json:"breakdown"`
}

// RuleEvaluator is an interface for an external service that evaluates rules.
type RuleEvaluator interface {
	Evaluate(ctx context.Context, ruleID string, context map[string]interface{}) (float64, error)
}

// CalculateRulesHandler is the handler for the CalculateRulesCommand.
type CalculateRulesHandler struct {
	evaluator RuleEvaluator
}

// NewCalculateRulesHandler creates a new CalculateRulesHandler.
func NewCalculateRulesHandler(evaluator RuleEvaluator) *CalculateRulesHandler {
	return &CalculateRulesHandler{
		evaluator: evaluator,
	}
}

// Handle handles the CalculateRulesCommand.
func (h *CalculateRulesHandler) Handle(ctx context.Context, cmd CalculateRulesCommand) (*CalculateRulesResult, error) {
	calc, err := calculation.NewCalculation(cmd.RuleIDs, cmd.Context)
	if err != nil {
		return nil, err
	}

	totalValue := 0.0
	breakdown := make(map[string]float64)

	for _, ruleID := range cmd.RuleIDs {
		// This is a placeholder for calling the rule evaluation service.
		// In a real implementation, this would make an RPC or HTTP call.
		ruleValue, err := h.evaluator.Evaluate(ctx, ruleID, cmd.Context)
		if err != nil {
			// Decide on error handling: fail the whole calculation or skip the rule?
			// For now, we'll skip the failing rule.
			continue
		}
		totalValue += ruleValue
		breakdown[ruleID] = ruleValue
	}

	result := calculation.Result{
		Value:     totalValue,
		Breakdown: breakdown,
	}
	calc.Complete(result)

	return &CalculateRulesResult{
		CalculationID: calc.ID().String(),
		Value:         result.Value,
		Breakdown:     result.Breakdown,
	}, nil
}
