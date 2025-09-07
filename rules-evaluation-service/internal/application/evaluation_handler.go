package application

import (
	"context"
	"strconv"
	"time"

	"rules-evaluation-service/internal/domain/evaluation"
	"rules-evaluation-service/internal/infrastructure/telemetry"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	tr := otel.Tracer("application")
	ctx, span := tr.Start(ctx, "EvaluateRuleHandler.Handle")
	defer span.End()

	span.SetAttributes(attribute.String("rule.category", cmd.RuleCategory))

	startTime := time.Now()
	strategy, err := h.evaluationService.GetStrategyForCategory(cmd.RuleCategory)
	if err != nil {
		telemetry.EvaluationsTotal.WithLabelValues(cmd.RuleCategory, "false").Inc()
		telemetry.EvaluationDuration.WithLabelValues(cmd.RuleCategory).Observe(time.Since(startTime).Seconds())
		return nil, err
	}

	result, err := strategy.Evaluate(cmd.DSLContent, cmd.Context)
	success := err == nil
	telemetry.EvaluationsTotal.WithLabelValues(cmd.RuleCategory, strconv.FormatBool(success)).Inc()
	telemetry.EvaluationDuration.WithLabelValues(cmd.RuleCategory).Observe(time.Since(startTime).Seconds())

	if err != nil {
		return nil, err
	}

	return &EvaluateRuleResult{
		Result: result,
	}, nil
}
