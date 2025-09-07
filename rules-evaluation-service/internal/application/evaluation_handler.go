package application

import (
	"context"

	"rules-evaluation-service/internal/domain/evaluation"
)

// EvaluateRuleCommand represents the command to evaluate a rule.
type EvaluateRuleCommand struct {
	RuleCategory string
	DSLContent   string
	Context      evaluation.Context
}

// EvaluateRuleResult represents the result of a rule evaluation.
type EvaluateRuleResult struct {
	Result evaluation.Result
}

// EvaluateRuleHandler handles the evaluation of a rule.
type EvaluateRuleHandler struct {
	evaluationService *evaluation.Service
}

// NewEvaluateRuleHandler creates a new handler.
func NewEvaluateRuleHandler(evaluationService *evaluation.Service) *EvaluateRuleHandler {
	return &EvaluateRuleHandler{evaluationService: evaluationService}
}

// Handle executes the command.
func (h *EvaluateRuleHandler) Handle(ctx context.Context, cmd EvaluateRuleCommand) (*EvaluateRuleResult, error) {
	// 1. Get the strategy for the rule's category.
	strategy, err := h.evaluationService.GetStrategyForCategory(cmd.RuleCategory)
	if err != nil {
		return nil, err
	}

	// 2. Execute the strategy.
	result, err := strategy.Evaluate(cmd.DSLContent, cmd.Context)
	if err != nil {
		return nil, err
	}

	return &EvaluateRuleResult{Result: result}, nil
}
