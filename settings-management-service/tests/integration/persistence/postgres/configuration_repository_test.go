package postgres_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
	postgres_persistence "github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/infrastructure/persistence/postgres"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&postgres_persistence.ConfigurationDBModel{},
	)
	require.NoError(t, err)

	return db
}

func TestConfigurationRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres_persistence.NewConfigurationRepository(db)

	config, err := settings.NewConfiguration(
		"database.host",
		map[string]interface{}{"host": "localhost", "port": 5432},
		"database",
		"development",
		"Database connection settings",
		false,
	)
	require.NoError(t, err)

	err = repo.Save(config)
	assert.NoError(t, err)

	// Verify the configuration was saved
	var dbModel postgres_persistence.ConfigurationDBModel
	err = db.First(&dbModel, "id = ?", config.ID()).Error
	assert.NoError(t, err)
	assert.Equal(t, config.ID(), dbModel.ID)
	assert.Equal(t, config.Key(), dbModel.Key)
	assert.Equal(t, config.Category(), dbModel.Category)
	assert.Equal(t, config.Environment(), dbModel.Environment)
	assert.Equal(t, config.Description(), dbModel.Description)
	assert.Equal(t, config.IsEncrypted(), dbModel.IsEncrypted)
}

func TestConfigurationRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres_persistence.NewConfigurationRepository(db)

	// Create and save a configuration
	config, err := settings.NewConfiguration(
		"database.host",
		map[string]interface{}{"host": "localhost", "port": 5432},
		"database",
		"development",
		"Database connection settings",
		false,
	)
	require.NoError(t, err)

	err = repo.Save(config)
	require.NoError(t, err)

	// Find the configuration by ID
	foundConfig, err := repo.FindByID(config.ID())
	assert.NoError(t, err)
	assert.NotNil(t, foundConfig)
	assert.Equal(t, config.ID(), foundConfig.ID())
	assert.Equal(t, config.Key(), foundConfig.Key())
	assert.Equal(t, config.Value(), foundConfig.Value())
	assert.Equal(t, config.Category(), foundConfig.Category())
	assert.Equal(t, config.Environment(), foundConfig.Environment())
	assert.Equal(t, config.Description(), foundConfig.Description())
	assert.Equal(t, config.IsEncrypted(), foundConfig.IsEncrypted())
}

func TestConfigurationRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres_persistence.NewConfigurationRepository(db)

	nonExistentID := uuid.New()
	foundConfig, err := repo.FindByID(nonExistentID)
	assert.Error(t, err)
	assert.Nil(t, foundConfig)
}

func TestConfigurationRepository_FindByKey(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres_persistence.NewConfigurationRepository(db)

	// Create and save a configuration
	config, err := settings.NewConfiguration(
		"database.host",
		map[string]interface{}{"host": "localhost", "port": 5432},
		"database",
		"development",
		"Database connection settings",
		false,
	)
	require.NoError(t, err)

	err = repo.Save(config)
	require.NoError(t, err)

	// Find the configuration by key
	foundConfig, err := repo.FindByKey("database.host")
	assert.NoError(t, err)
	assert.NotNil(t, foundConfig)
	assert.Equal(t, config.ID(), foundConfig.ID())
	assert.Equal(t, config.Key(), foundConfig.Key())
}

func TestConfigurationRepository_FindByKey_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres_persistence.NewConfigurationRepository(db)

	foundConfig, err := repo.FindByKey("non-existent-key")
	assert.Error(t, err)
	assert.Nil(t, foundConfig)
}

func TestConfigurationRepository_ExistsByKey(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres_persistence.NewConfigurationRepository(db)

	// Create and save a configuration
	config, err := settings.NewConfiguration(
		"database.host",
		map[string]interface{}{"host": "localhost", "port": 5432},
		"database",
		"development",
		"Database connection settings",
		false,
	)
	require.NoError(t, err)

	err = repo.Save(config)
	require.NoError(t, err)

	// Check if configuration exists by key
	exists, err := repo.ExistsByKey("database.host")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Check if non-existent configuration exists
	exists, err = repo.ExistsByKey("non-existent-key")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestConfigurationRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres_persistence.NewConfigurationRepository(db)

	// Create and save multiple configurations
	configs := []*settings.Configuration{
		func() *settings.Configuration {
			config, _ := settings.NewConfiguration(
				"database.host",
				map[string]interface{}{"host": "localhost", "port": 5432},
				"database",
				"development",
				"Database connection settings",
				false,
			)
			return config
		}(),
		func() *settings.Configuration {
			config, _ := settings.NewConfiguration(
				"redis.host",
				map[string]interface{}{"host": "localhost", "port": 6379},
				"cache",
				"development",
				"Redis connection settings",
				false,
			)
			return config
		}(),
		func() *settings.Configuration {
			config, _ := settings.NewConfiguration(
				"api.timeout",
				map[string]interface{}{"timeout": 30},
				"api",
				"production",
				"API timeout settings",
				false,
			)
			return config
		}(),
	}

	for _, config := range configs {
		err := repo.Save(config)
		require.NoError(t, err)
	}

	// Test listing all configurations
	allConfigs, err := repo.List(settings.ListFilters{}, settings.ListOptions{Page: 1, Limit: 10})
	assert.NoError(t, err)
	assert.Len(t, allConfigs, 3)

	// Test listing configurations by category
	databaseConfigs, err := repo.List(
		settings.ListFilters{Category: "database"},
		settings.ListOptions{Page: 1, Limit: 10},
	)
	assert.NoError(t, err)
	assert.Len(t, databaseConfigs, 1)
	assert.Equal(t, "database.host", databaseConfigs[0].Key())

	// Test listing configurations by environment
	devConfigs, err := repo.List(
		settings.ListFilters{Environment: "development"},
		settings.ListOptions{Page: 1, Limit: 10},
	)
	assert.NoError(t, err)
	assert.Len(t, devConfigs, 2)

	// Test pagination
	paginatedConfigs, err := repo.List(
		settings.ListFilters{},
		settings.ListOptions{Page: 1, Limit: 2},
	)
	assert.NoError(t, err)
	assert.Len(t, paginatedConfigs, 2)
}

func TestConfigurationRepository_Count(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres_persistence.NewConfigurationRepository(db)

	// Create and save multiple configurations
	configs := []*settings.Configuration{
		func() *settings.Configuration {
			config, _ := settings.NewConfiguration(
				"database.host",
				map[string]interface{}{"host": "localhost", "port": 5432},
				"database",
				"development",
				"Database connection settings",
				false,
			)
			return config
		}(),
		func() *settings.Configuration {
			config, _ := settings.NewConfiguration(
				"redis.host",
				map[string]interface{}{"host": "localhost", "port": 6379},
				"cache",
				"development",
				"Redis connection settings",
				false,
			)
			return config
		}(),
	}

	for _, config := range configs {
		err := repo.Save(config)
		require.NoError(t, err)
	}

	// Test counting all configurations
	count, err := repo.Count(settings.ListFilters{})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)

	// Test counting configurations by category
	count, err = repo.Count(settings.ListFilters{Category: "database"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Test counting configurations by environment
	count, err = repo.Count(settings.ListFilters{Environment: "development"})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestConfigurationRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres_persistence.NewConfigurationRepository(db)

	// Create and save a configuration
	config, err := settings.NewConfiguration(
		"database.host",
		map[string]interface{}{"host": "localhost", "port": 5432},
		"database",
		"development",
		"Database connection settings",
		false,
	)
	require.NoError(t, err)

	err = repo.Save(config)
	require.NoError(t, err)

	// Delete the configuration
	err = repo.Delete(config.ID())
	assert.NoError(t, err)

	// Verify the configuration was deleted
	_, err = repo.FindByID(config.ID())
	assert.Error(t, err)
}

func TestConfigurationRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres_persistence.NewConfigurationRepository(db)

	nonExistentID := uuid.New()
	err := repo.Delete(nonExistentID)
	assert.Error(t, err)
}

func TestConfigurationRepository_FindByCategory(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres_persistence.NewConfigurationRepository(db)

	// Create and save configurations with different categories
	configs := []*settings.Configuration{
		func() *settings.Configuration {
			config, _ := settings.NewConfiguration(
				"database.host",
				map[string]interface{}{"host": "localhost", "port": 5432},
				"database",
				"development",
				"Database connection settings",
				false,
			)
			return config
		}(),
		func() *settings.Configuration {
			config, _ := settings.NewConfiguration(
				"database.port",
				map[string]interface{}{"port": 5432},
				"database",
				"development",
				"Database port settings",
				false,
			)
			return config
		}(),
		func() *settings.Configuration {
			config, _ := settings.NewConfiguration(
				"redis.host",
				map[string]interface{}{"host": "localhost", "port": 6379},
				"cache",
				"development",
				"Redis connection settings",
				false,
			)
			return config
		}(),
	}

	for _, config := range configs {
		err := repo.Save(config)
		require.NoError(t, err)
	}

	// Find configurations by category
	databaseConfigs, err := repo.FindByCategory("database")
	assert.NoError(t, err)
	assert.Len(t, databaseConfigs, 2)

	// Verify all returned configurations have the correct category
	for _, config := range databaseConfigs {
		assert.Equal(t, "database", config.Category())
	}
}

func TestConfigurationRepository_FindByEnvironment(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres_persistence.NewConfigurationRepository(db)

	// Create and save configurations with different environments
	configs := []*settings.Configuration{
		func() *settings.Configuration {
			config, _ := settings.NewConfiguration(
				"database.host",
				map[string]interface{}{"host": "localhost", "port": 5432},
				"database",
				"development",
				"Database connection settings",
				false,
			)
			return config
		}(),
		func() *settings.Configuration {
			config, _ := settings.NewConfiguration(
				"redis.host",
				map[string]interface{}{"host": "localhost", "port": 6379},
				"cache",
				"development",
				"Redis connection settings",
				false,
			)
			return config
		}(),
		func() *settings.Configuration {
			config, _ := settings.NewConfiguration(
				"api.timeout",
				map[string]interface{}{"timeout": 30},
				"api",
				"production",
				"API timeout settings",
				false,
			)
			return config
		}(),
	}

	for _, config := range configs {
		err := repo.Save(config)
		require.NoError(t, err)
	}

	// Find configurations by environment
	devConfigs, err := repo.FindByEnvironment("development")
	assert.NoError(t, err)
	assert.Len(t, devConfigs, 2)

	// Verify all returned configurations have the correct environment
	for _, config := range devConfigs {
		assert.Equal(t, "development", config.Environment())
	}
}
