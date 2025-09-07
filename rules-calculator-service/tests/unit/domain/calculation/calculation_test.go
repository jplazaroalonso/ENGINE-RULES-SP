package calculation_test

import (
	"testing"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-calculator-service/internal/domain/calculation"
	"github.com/stretchr/testify/assert"
)

func TestNewCalculation(t *testing.T) {
	ruleIDs := []string{"rule1", "rule2"}
	context := map[string]interface{}{"key": "value"}

	calc, err := calculation.NewCalculation(ruleIDs, context)

	assert.NoError(t, err)
	assert.NotNil(t, calc)
	assert.NotEqual(t, calculation.CalculationID{}, calc.ID())
	assert.Equal(t, calculation.StatusPending, calc.Status())
	assert.Nil(t, calc.Result())
	assert.Nil(t, calc.CompletedAt())
}

func TestCalculation_Complete(t *testing.T) {
	// Arrange
	ruleIDs := []string{"rule1"}
	context := map[string]interface{}{"key": "value"}
	calc, _ := calculation.NewCalculation(ruleIDs, context)
	result := calculation.Result{
		Value:     123.45,
		Breakdown: map[string]float64{"rule1": 123.45},
	}

	// Act
	calc.Complete(result)

	// Assert
	assert.Equal(t, calculation.StatusCompleted, calc.Status())
	assert.NotNil(t, calc.Result())
	assert.Equal(t, 123.45, calc.Result().Value)
	assert.NotNil(t, calc.CompletedAt())
}

func TestCalculation_Fail(t *testing.T) {
	// Arrange
	ruleIDs := []string{"rule1"}
	context := map[string]interface{}{"key": "value"}
	calc, _ := calculation.NewCalculation(ruleIDs, context)

	// Act
	calc.Fail()

	// Assert
	assert.Equal(t, calculation.StatusFailed, calc.Status())
	assert.Nil(t, calc.Result()) // No result on failure
	assert.NotNil(t, calc.CompletedAt())
}
