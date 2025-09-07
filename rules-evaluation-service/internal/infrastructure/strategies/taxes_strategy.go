package strategies

import (
	"context"

	"rules-evaluation-service/internal/domain/evaluation"

	"go.opentelemetry.io/otel"
)

// TaxesStrategy is a strategy for evaluating taxes rules.
type TaxesStrategy struct{}

// Evaluate evaluates a taxes rule.
func (s *TaxesStrategy) Evaluate(dslContent string, evalContext evaluation.Context) (evaluation.Result, error) {
	ctx := context.Background()
	_, span := otel.Tracer("strategy").Start(ctx, "TaxesStrategy.Evaluate")
	defer span.End()
	// Placeholder implementation
	return evaluation.Result{"tax_rate": 0.21}, nil
}
