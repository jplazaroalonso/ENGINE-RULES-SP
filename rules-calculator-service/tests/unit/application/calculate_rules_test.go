package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-calculator-service/internal/application"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRuleEvaluator is a mock implementation of the RuleEvaluator interface.
type MockRuleEvaluator struct {
	mock.Mock
}

func (m *MockRuleEvaluator) Evaluate(ctx context.Context, ruleID string, context map[string]interface{}) (float64, error) {
	args := m.Called(ctx, ruleID, context)
	return args.Get(0).(float64), args.Error(1)
}

func TestCalculateRulesHandler_Handle_Success(t *testing.T) {
	// Arrange
	mockEvaluator := new(MockRuleEvaluator)
	handler := application.NewCalculateRulesHandler(mockEvaluator)

	cmd := application.CalculateRulesCommand{
		RuleIDs: []string{"rule1", "rule2"},
		Context: map[string]interface{}{"customer_tier": "gold"},
	}

	// Setup mock expectations
	mockEvaluator.On("Evaluate", mock.Anything, "rule1", cmd.Context).Return(100.0, nil)
	mockEvaluator.On("Evaluate", mock.Anything, "rule2", cmd.Context).Return(50.5, nil)

	// Act
	result, err := handler.Handle(context.Background(), cmd)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.CalculationID)
	assert.Equal(t, 150.5, result.Value)
	assert.Equal(t, 100.0, result.Breakdown["rule1"])
	assert.Equal(t, 50.5, result.Breakdown["rule2"])

	mockEvaluator.AssertExpectations(t)
}

func TestCalculateRulesHandler_Handle_WithFailingRule(t *testing.T) {
	// Arrange
	mockEvaluator := new(MockRuleEvaluator)
	handler := application.NewCalculateRulesHandler(mockEvaluator)

	cmd := application.CalculateRulesCommand{
		RuleIDs: []string{"rule1", "failing_rule", "rule2"},
		Context: map[string]interface{}{"customer_tier": "gold"},
	}

	// Setup mock expectations
	mockEvaluator.On("Evaluate", mock.Anything, "rule1", cmd.Context).Return(100.0, nil)
	mockEvaluator.On("Evaluate", mock.Anything, "failing_rule", cmd.Context).Return(0.0, errors.New("evaluation failed"))
	mockEvaluator.On("Evaluate", mock.Anything, "rule2", cmd.Context).Return(50.5, nil)

	// Act
	result, err := handler.Handle(context.Background(), cmd)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 150.5, result.Value) // failing_rule is skipped
	assert.Len(t, result.Breakdown, 2)
	assert.Equal(t, 100.0, result.Breakdown["rule1"])
	assert.Equal(t, 50.5, result.Breakdown["rule2"])
	_, exists := result.Breakdown["failing_rule"]
	assert.False(t, exists)

	mockEvaluator.AssertExpectations(t)
}

func TestCalculateRulesHandler_Handle_NoRules(t *testing.T) {
	// Arrange
	mockEvaluator := new(MockRuleEvaluator)
	handler := application.NewCalculateRulesHandler(mockEvaluator)

	cmd := application.CalculateRulesCommand{
		RuleIDs: []string{},
		Context: map[string]interface{}{"customer_tier": "gold"},
	}

	// Act
	result, err := handler.Handle(context.Background(), cmd)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0.0, result.Value)
	assert.Empty(t, result.Breakdown)

	mockEvaluator.AssertNotCalled(t, "Evaluate")
}
