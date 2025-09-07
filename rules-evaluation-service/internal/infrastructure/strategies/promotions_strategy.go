package strategies

import (
	"context"
	"fmt"
	"strings"

	"rules-evaluation-service/internal/domain/evaluation"

	"go.opentelemetry.io/otel"
)

// PromotionsStrategy is a strategy for evaluating promotions rules.
type PromotionsStrategy struct{}

func NewPromotionsStrategy() *PromotionsStrategy {
	return &PromotionsStrategy{}
}

// Evaluate evaluates a promotions rule.
func (s *PromotionsStrategy) Evaluate(dslContent string, evalContext evaluation.Context) (evaluation.Result, error) {
	ctx := context.Background()
	_, span := otel.Tracer("strategy").Start(ctx, "PromotionsStrategy.Evaluate")
	defer span.End()
	// Example DSL: "IF order.amount > 100 THEN discount.percentage = 10"
	if !strings.Contains(dslContent, "order.amount >") {
		return nil, fmt.Errorf("invalid promotions DSL: missing 'order.amount >'")
	}

	orderAmount, ok := evalContext["order_amount"].(float64)
	if !ok {
		return evaluation.Result{"eligible": false, "reason": "Missing order_amount in context"}, nil
	}

	// Simplified parsing logic
	parts := strings.Split(dslContent, ">")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid promotions DSL")
	}
	thresholdStr := strings.Split(strings.TrimSpace(parts[1]), " ")[0]
	var threshold float64
	fmt.Sscanf(thresholdStr, "%f", &threshold)

	if orderAmount > threshold {
		// Another simplified extraction
		discountParts := strings.Split(dslContent, "= ")
		var discount float64
		fmt.Sscanf(strings.TrimSpace(discountParts[1]), "%f", &discount)
		return evaluation.Result{"eligible": true, "discount_percentage": discount}, nil
	}

	return evaluation.Result{"eligible": false}, nil
}
