package strategies

import (
	"fmt"
	"strings"

	"rules-evaluation-service/internal/domain/evaluation"
)

// PromotionsStrategy implements the evaluation logic for promotional rules.
type PromotionsStrategy struct{}

func NewPromotionsStrategy() *PromotionsStrategy {
	return &PromotionsStrategy{}
}

// Evaluate for promotions will check if the transaction amount is over a certain threshold.
// This is a simplified example. A real implementation would parse the DSL properly.
func (s *PromotionsStrategy) Evaluate(dslContent string, context evaluation.Context) (evaluation.Result, error) {
	// Example DSL: "IF order.amount > 100 THEN discount.percentage = 10"
	if !strings.Contains(dslContent, "order.amount >") {
		return nil, fmt.Errorf("invalid promotions DSL: missing 'order.amount >'")
	}

	orderAmount, ok := context["order_amount"].(float64)
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
