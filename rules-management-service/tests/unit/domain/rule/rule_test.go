package rule_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/domain/rule"
)

func TestNewRule(t *testing.T) {
	t.Run("should create a new rule successfully", func(t *testing.T) {
		name := "Test Rule"
		description := "A rule for testing"
		dslContent := "IF condition THEN action"
		createdBy := "test-user"
		priority := rule.PriorityHigh
		category := "TEST"
		tags := []string{"testing"}

		r, err := rule.NewRule(name, description, dslContent, createdBy, priority, category, tags)

		require.NoError(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, name, r.Name())
		assert.Equal(t, description, r.Description())
		assert.Equal(t, dslContent, r.DSLContent())
		assert.Equal(t, createdBy, r.CreatedBy())
		assert.Equal(t, priority, r.Priority())
		assert.Equal(t, category, r.Category())
		assert.Equal(t, tags, r.Tags())
		assert.Equal(t, rule.StatusDraft, r.Status())
		assert.Equal(t, 1, r.Version())
		assert.NotZero(t, r.ID())
		assert.False(t, r.CreatedAt().IsZero())
	})

	t.Run("should return an error for an empty name", func(t *testing.T) {
		_, err := rule.NewRule("", "desc", "dsl", "user", rule.PriorityLow, "cat", nil)
		require.Error(t, err)
	})

	t.Run("should return an error for empty DSL content", func(t *testing.T) {
		_, err := rule.NewRule("name", "desc", "", "user", rule.PriorityLow, "cat", nil)
		require.Error(t, err)
	})
}
