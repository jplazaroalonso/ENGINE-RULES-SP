package strategies

import (
	"fmt"
	"strings"

	"rules-evaluation-service/internal/domain/evaluation"
)

// TaxesStrategy implements the evaluation logic for tax rules.
type TaxesStrategy struct{}

func NewTaxesStrategy() *TaxesStrategy {
	return &TaxesStrategy{}
}

// Evaluate for taxes will check if the customer is in a specific region.
// This is a simplified example.
func (s *TaxesStrategy) Evaluate(dslContent string, context evaluation.Context) (evaluation.Result, error) {
	// Example DSL: "IF customer.region == 'CA' THEN tax.percentage = 9.5"
	if !strings.Contains(dslContent, "customer.region ==") {
		return nil, fmt.Errorf("invalid tax DSL: missing 'customer.region =='")
	}

	customerRegion, ok := context["customer_region"].(string)
	if !ok {
		return evaluation.Result{"taxable": false, "reason": "Missing customer_region in context"}, nil
	}

	// Simplified parsing logic
	parts := strings.Split(dslContent, "==")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid tax DSL")
	}
	regionStr := strings.Split(strings.TrimSpace(parts[1]), " ")[0]
	requiredRegion := strings.Trim(regionStr, "'")

	if customerRegion == requiredRegion {
		// Another simplified extraction
		taxParts := strings.Split(dslContent, "= ")
		var tax float64
		fmt.Sscanf(strings.TrimSpace(taxParts[1]), "%f", &tax)

		return evaluation.Result{"taxable": true, "tax_percentage": tax}, nil
	}

	return evaluation.Result{"taxable": false}, nil
}
