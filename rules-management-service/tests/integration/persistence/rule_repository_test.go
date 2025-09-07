package persistence_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/domain/rule"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/infrastructure/persistence/postgres"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&postgres.RuleDBModel{})
	require.NoError(t, err)

	return db
}

func TestRuleRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db := setupTestDB(t)
	repo := postgres.NewRuleRepository(db)
	ctx := context.Background()

	t.Run("should save and find a rule by ID", func(t *testing.T) {
		newRule, err := rule.NewRule("Test Rule", "Desc", "IF true THEN false", "user", rule.PriorityMedium, "cat", []string{"tag1"})
		require.NoError(t, err)

		err = repo.Save(ctx, newRule)
		require.NoError(t, err)

		foundRule, err := repo.FindByID(ctx, newRule.ID())
		require.NoError(t, err)
		assert.Equal(t, newRule.ID(), foundRule.ID())
		assert.Equal(t, newRule.Name(), foundRule.Name())
	})

	t.Run("should return not found error for non-existent rule", func(t *testing.T) {
		nonExistentID, _ := rule.RuleIDFromStr("123e4567-e89b-12d3-a456-426614174000")
		_, err := repo.FindByID(ctx, nonExistentID)
		require.Error(t, err)
	})

	t.Run("should correctly check if a rule exists by name", func(t *testing.T) {
		exists, err := repo.ExistsByName(ctx, "Test Rule")
		require.NoError(t, err)
		assert.True(t, exists)

		exists, err = repo.ExistsByName(ctx, "NonExistent Rule")
		require.NoError(t, err)
		assert.False(t, exists)
	})
}
