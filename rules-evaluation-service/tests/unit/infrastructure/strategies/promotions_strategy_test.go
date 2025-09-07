package strategies_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"rules-evaluation-service/internal/domain/evaluation"
	"rules-evaluation-service/internal/infrastructure/strategies"
)

func TestPromotionsStrategy(t *testing.T) {
	strategy := strategies.NewPromotionsStrategy()

	t.Run("should return eligible and discount when amount is over threshold", func(t *testing.T) {
		dsl := "IF order.amount > 100 THEN discount.percentage = 10"
		context := evaluation.Context{"order_amount": 150.0}

		result, err := strategy.Evaluate(dsl, context)
		require.NoError(t, err)

		assert.Equal(t, evaluation.Result{"eligible": true, "discount_percentage": 10.0}, result)
	})

	t.Run("should return not eligible when amount is under threshold", func(t *testing.T) {
		dsl := "IF order.amount > 100 THEN discount.percentage = 10"
		context := evaluation.Context{"order_amount": 50.0}

		result, err := strategy.Evaluate(dsl, context)
		require.NoError(t, err)

		assert.Equal(t, evaluation.Result{"eligible": false}, result)
	})

	t.Run("should return an error for invalid DSL", func(t *testing.T) {
		dsl := "invalid dsl"
		context := evaluation.Context{"order_amount": 150.0}

		_, err := strategy.Evaluate(dsl, context)
		require.Error(t, err)
	})

	t.Run("should return not eligible for missing context", func(t *testing.T) {
		dsl := "IF order.amount > 100 THEN discount.percentage = 10"
		context := evaluation.Context{} // Missing order_amount

		result, err := strategy.Evaluate(dsl, context)
		require.NoError(t, err)

		assert.Equal(t, evaluation.Result{"eligible": false, "reason": "Missing order_amount in context"}, result)
	})
}
