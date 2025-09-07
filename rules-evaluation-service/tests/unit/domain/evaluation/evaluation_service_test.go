package evaluation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"rules-evaluation-service/internal/domain/evaluation"
)

// MockStrategy is a mock implementation of the EvaluationStrategy for testing.
type MockStrategy struct {
	Name string
}

func (s *MockStrategy) Evaluate(dslContent string, context evaluation.Context) (evaluation.Result, error) {
	return evaluation.Result{"strategy": s.Name}, nil
}

func TestEvaluationService(t *testing.T) {
	promoStrategy := &MockStrategy{Name: "PROMOTIONS"}
	taxStrategy := &MockStrategy{Name: "TAXES"}

	strategyMap := map[string]evaluation.EvaluationStrategy{
		"PROMOTIONS": promoStrategy,
		"TAXES":      taxStrategy,
	}

	service := evaluation.NewService(strategyMap)

	t.Run("should return the correct strategy for a category", func(t *testing.T) {
		strategy, err := service.GetStrategyForCategory("PROMOTIONS")
		require.NoError(t, err)
		assert.Equal(t, promoStrategy, strategy)

		strategy, err = service.GetStrategyForCategory("TAXES")
		require.NoError(t, err)
		assert.Equal(t, taxStrategy, strategy)
	})

	t.Run("should return an error for an unknown category", func(t *testing.T) {
		_, err := service.GetStrategyForCategory("UNKNOWN")
		require.Error(t, err)
		assert.EqualError(t, err, "no evaluation strategy found for category: UNKNOWN")
	})
}
