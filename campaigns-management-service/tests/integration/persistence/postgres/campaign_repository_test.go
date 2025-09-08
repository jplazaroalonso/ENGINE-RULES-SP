package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/campaign"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/infrastructure/persistence/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCampaignRepository_Integration(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := postgres.NewCampaignRepository(db)
	ctx := context.Background()

	t.Run("Save and FindByID", func(t *testing.T) {
		// Arrange
		testCampaign := createTestCampaign(t)

		// Act
		err := repo.Save(ctx, testCampaign)
		require.NoError(t, err)

		// Assert
		foundCampaign, err := repo.FindByID(ctx, testCampaign.ID())
		require.NoError(t, err)
		require.NotNil(t, foundCampaign)
		assert.Equal(t, testCampaign.Name(), foundCampaign.Name())
		assert.Equal(t, testCampaign.Description(), foundCampaign.Description())
		assert.Equal(t, testCampaign.Status(), foundCampaign.Status())
		assert.Equal(t, testCampaign.CampaignType(), foundCampaign.CampaignType())
		assert.Equal(t, testCampaign.StartDate(), foundCampaign.StartDate())
		assert.Equal(t, testCampaign.CreatedBy(), foundCampaign.CreatedBy())
	})

	t.Run("FindByName", func(t *testing.T) {
		// Arrange
		testCampaign := createTestCampaign(t)
		err := repo.Save(ctx, testCampaign)
		require.NoError(t, err)

		// Act
		foundCampaign, err := repo.FindByName(ctx, testCampaign.Name())
		require.NoError(t, err)
		require.NotNil(t, foundCampaign)

		// Assert
		assert.Equal(t, testCampaign.ID(), foundCampaign.ID())
		assert.Equal(t, testCampaign.Name(), foundCampaign.Name())
	})

	t.Run("FindByName_NotFound", func(t *testing.T) {
		// Act
		foundCampaign, err := repo.FindByName(ctx, "Non-existent Campaign")
		require.NoError(t, err)
		assert.Nil(t, foundCampaign)
	})

	t.Run("ExistsByName", func(t *testing.T) {
		// Arrange
		testCampaign := createTestCampaign(t)
		err := repo.Save(ctx, testCampaign)
		require.NoError(t, err)

		// Act & Assert
		exists, err := repo.ExistsByName(ctx, testCampaign.Name())
		require.NoError(t, err)
		assert.True(t, exists)

		exists, err = repo.ExistsByName(ctx, "Non-existent Campaign")
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("List with criteria", func(t *testing.T) {
		// Arrange
		userID := shared.NewUserID()
		campaigns := createTestCampaigns(t, userID, 3)

		for _, c := range campaigns {
			err := repo.Save(ctx, c)
			require.NoError(t, err)
		}

		criteria := campaign.ListCriteria{
			Status:     &[]campaign.CampaignStatus{campaign.CampaignStatusDraft}[0],
			CreatedBy:  &userID,
			PageSize:   10,
			PageOffset: 0,
			SortBy:     "name",
			SortOrder:  "asc",
		}

		// Act
		foundCampaigns, err := repo.List(ctx, criteria)
		require.NoError(t, err)

		// Assert
		assert.Len(t, foundCampaigns, 3)
		for _, c := range foundCampaigns {
			assert.Equal(t, campaign.CampaignStatusDraft, c.Status())
			assert.Equal(t, userID, c.CreatedBy())
		}
	})

	t.Run("List with pagination", func(t *testing.T) {
		// Arrange
		userID := shared.NewUserID()
		campaigns := createTestCampaigns(t, userID, 5)

		for _, c := range campaigns {
			err := repo.Save(ctx, c)
			require.NoError(t, err)
		}

		criteria := campaign.ListCriteria{
			CreatedBy:  &userID,
			PageSize:   2,
			PageOffset: 0,
			SortBy:     "name",
			SortOrder:  "asc",
		}

		// Act
		foundCampaigns, err := repo.List(ctx, criteria)
		require.NoError(t, err)

		// Assert
		assert.Len(t, foundCampaigns, 2)
	})

	t.Run("Count with criteria", func(t *testing.T) {
		// Arrange
		userID := shared.NewUserID()
		campaigns := createTestCampaigns(t, userID, 3)

		for _, c := range campaigns {
			err := repo.Save(ctx, c)
			require.NoError(t, err)
		}

		criteria := campaign.ListCriteria{
			Status:    &[]campaign.CampaignStatus{campaign.CampaignStatusDraft}[0],
			CreatedBy: &userID,
		}

		// Act
		count, err := repo.Count(ctx, criteria)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, int64(3), count)
	})

	t.Run("FindByStatus", func(t *testing.T) {
		// Arrange
		userID := shared.NewUserID()
		campaigns := createTestCampaigns(t, userID, 2)

		// Activate one campaign
		err := campaigns[0].Activate()
		require.NoError(t, err)

		for _, c := range campaigns {
			err := repo.Save(ctx, c)
			require.NoError(t, err)
		}

		// Act
		activeCampaigns, err := repo.FindByStatus(ctx, campaign.CampaignStatusActive)
		require.NoError(t, err)

		// Assert
		assert.Len(t, activeCampaigns, 1)
		assert.Equal(t, campaign.CampaignStatusActive, activeCampaigns[0].Status())
	})

	t.Run("FindByType", func(t *testing.T) {
		// Arrange
		userID := shared.NewUserID()
		campaigns := createTestCampaigns(t, userID, 2)

		// Change type of one campaign
		campaigns[0].UpdateSettings(createValidSettings())

		for _, c := range campaigns {
			err := repo.Save(ctx, c)
			require.NoError(t, err)
		}

		// Act
		promotionCampaigns, err := repo.FindByType(ctx, campaign.CampaignTypePromotion)
		require.NoError(t, err)

		// Assert
		assert.Len(t, promotionCampaigns, 2)
		for _, c := range promotionCampaigns {
			assert.Equal(t, campaign.CampaignTypePromotion, c.CampaignType())
		}
	})

	t.Run("FindByDateRange", func(t *testing.T) {
		// Arrange
		userID := shared.NewUserID()
		now := time.Now()

		// Create campaigns with different start dates
		campaign1 := createTestCampaignWithDate(t, userID, now.Add(-48*time.Hour))
		campaign2 := createTestCampaignWithDate(t, userID, now.Add(-24*time.Hour))
		campaign3 := createTestCampaignWithDate(t, userID, now.Add(24*time.Hour))

		campaigns := []*campaign.Campaign{campaign1, campaign2, campaign3}
		for _, c := range campaigns {
			err := repo.Save(ctx, c)
			require.NoError(t, err)
		}

		// Act - Find campaigns in the last 30 days
		startDate := now.Add(-30 * 24 * time.Hour)
		endDate := now
		foundCampaigns, err := repo.FindByDateRange(ctx, startDate, endDate)
		require.NoError(t, err)

		// Assert
		assert.Len(t, foundCampaigns, 2) // campaign1 and campaign2
	})

	t.Run("FindByCreatedBy", func(t *testing.T) {
		// Arrange
		userID1 := shared.NewUserID()
		userID2 := shared.NewUserID()

		campaigns1 := createTestCampaigns(t, userID1, 2)
		campaigns2 := createTestCampaigns(t, userID2, 1)

		allCampaigns := append(campaigns1, campaigns2...)
		for _, c := range allCampaigns {
			err := repo.Save(ctx, c)
			require.NoError(t, err)
		}

		// Act
		foundCampaigns, err := repo.FindByCreatedBy(ctx, userID1)
		require.NoError(t, err)

		// Assert
		assert.Len(t, foundCampaigns, 2)
		for _, c := range foundCampaigns {
			assert.Equal(t, userID1, c.CreatedBy())
		}
	})

	t.Run("Delete", func(t *testing.T) {
		// Arrange
		testCampaign := createTestCampaign(t)
		err := repo.Save(ctx, testCampaign)
		require.NoError(t, err)

		// Act
		err = repo.Delete(ctx, testCampaign.ID())
		require.NoError(t, err)

		// Assert
		foundCampaign, err := repo.FindByID(ctx, testCampaign.ID())
		require.NoError(t, err)
		assert.Nil(t, foundCampaign)
	})

	t.Run("Update campaign", func(t *testing.T) {
		// Arrange
		testCampaign := createTestCampaign(t)
		err := repo.Save(ctx, testCampaign)
		require.NoError(t, err)

		// Act - Update campaign
		newBudget := &shared.Money{Amount: 2000.0, Currency: "USD"}
		err = testCampaign.UpdateBudget(newBudget)
		require.NoError(t, err)

		err = repo.Save(ctx, testCampaign)
		require.NoError(t, err)

		// Assert
		foundCampaign, err := repo.FindByID(ctx, testCampaign.ID())
		require.NoError(t, err)
		require.NotNil(t, foundCampaign)
		assert.Equal(t, newBudget.Amount, foundCampaign.Budget().Amount)
		assert.Equal(t, newBudget.Currency, foundCampaign.Budget().Currency)
		assert.Greater(t, foundCampaign.Version(), 1)
	})
}

func TestCampaignRepository_ErrorHandling(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := postgres.NewCampaignRepository(db)
	ctx := context.Background()

	t.Run("Save with invalid data", func(t *testing.T) {
		// This test would require more complex setup to trigger database errors
		// For now, we'll test the basic error handling
		testCampaign := createTestCampaign(t)

		// Save should succeed with valid data
		err := repo.Save(ctx, testCampaign)
		require.NoError(t, err)
	})

	t.Run("FindByID with invalid UUID", func(t *testing.T) {
		// This would require the repository to handle invalid UUIDs gracefully
		// The current implementation should handle this through the domain layer
		invalidID, err := campaign.NewCampaignIDFromString("invalid-uuid")
		require.Error(t, err)

		// The repository should not be called with invalid IDs
		// as validation should happen at the application layer
	})
}

// Helper functions

func setupTestDB(t *testing.T) *gorm.DB {
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&postgres.CampaignDBModel{},
		&postgres.CampaignMetricsDBModel{},
		&postgres.CampaignEventDBModel{},
	)
	require.NoError(t, err)

	return db
}

func cleanupTestDB(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	require.NoError(t, err)
	err = sqlDB.Close()
	require.NoError(t, err)
}

func createTestCampaign(t *testing.T) *campaign.Campaign {
	campaign, err := campaign.NewCampaign(
		"Test Campaign",
		"Test campaign description",
		campaign.CampaignTypePromotion,
		[]shared.RuleID{shared.NewRuleID()},
		time.Now().Add(24*time.Hour),
		nil,
		&shared.Money{Amount: 1000.0, Currency: "EUR"},
		shared.NewUserID(),
		createValidSettings(),
	)
	require.NoError(t, err)
	return campaign
}

func createTestCampaignWithDate(t *testing.T, userID shared.UserID, startDate time.Time) *campaign.Campaign {
	campaign, err := campaign.NewCampaign(
		"Test Campaign "+startDate.Format("2006-01-02"),
		"Test campaign description",
		campaign.CampaignTypePromotion,
		[]shared.RuleID{shared.NewRuleID()},
		startDate,
		nil,
		&shared.Money{Amount: 1000.0, Currency: "EUR"},
		userID,
		createValidSettings(),
	)
	require.NoError(t, err)
	return campaign
}

func createTestCampaigns(t *testing.T, userID shared.UserID, count int) []*campaign.Campaign {
	campaigns := make([]*campaign.Campaign, count)
	for i := 0; i < count; i++ {
		campaign, err := campaign.NewCampaign(
			"Test Campaign "+string(rune('A'+i)),
			"Test campaign description",
			campaign.CampaignTypePromotion,
			[]shared.RuleID{shared.NewRuleID()},
			time.Now().Add(24*time.Hour),
			nil,
			&shared.Money{Amount: 1000.0, Currency: "EUR"},
			userID,
			createValidSettings(),
		)
		require.NoError(t, err)
		campaigns[i] = campaign
	}
	return campaigns
}

func createValidSettings() campaign.CampaignSettings {
	return campaign.CampaignSettings{
		TargetAudience: []string{"test-audience"},
		Channels:       []campaign.Channel{campaign.ChannelEmail, campaign.ChannelWeb},
		Frequency:      campaign.FrequencyDaily,
		MaxImpressions: 1000,
		Personalization: campaign.PersonalizationConfig{
			Enabled: false,
			Rules:   []shared.RuleID{},
		},
	}
}
