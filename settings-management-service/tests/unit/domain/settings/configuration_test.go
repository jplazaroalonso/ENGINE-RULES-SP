package settings_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
)

func TestNewConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		value       map[string]interface{}
		category    string
		environment string
		description string
		isEncrypted bool
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid configuration",
			key:         "database.host",
			value:       map[string]interface{}{"host": "localhost", "port": 5432},
			category:    "database",
			environment: "development",
			description: "Database connection settings",
			isEncrypted: false,
			expectError: false,
		},
		{
			name:        "empty key",
			key:         "",
			value:       map[string]interface{}{"host": "localhost"},
			category:    "database",
			environment: "development",
			description: "Database connection settings",
			isEncrypted: false,
			expectError: true,
			errorMsg:    "key cannot be empty",
		},
		{
			name:        "empty category",
			key:         "database.host",
			value:       map[string]interface{}{"host": "localhost"},
			category:    "",
			environment: "development",
			description: "Database connection settings",
			isEncrypted: false,
			expectError: true,
			errorMsg:    "category cannot be empty",
		},
		{
			name:        "empty environment",
			key:         "database.host",
			value:       map[string]interface{}{"host": "localhost"},
			category:    "database",
			environment: "",
			description: "Database connection settings",
			isEncrypted: false,
			expectError: true,
			errorMsg:    "environment cannot be empty",
		},
		{
			name:        "nil value",
			key:         "database.host",
			value:       nil,
			category:    "database",
			environment: "development",
			description: "Database connection settings",
			isEncrypted: false,
			expectError: true,
			errorMsg:    "value cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := settings.NewConfiguration(
				tt.key,
				tt.value,
				tt.category,
				tt.environment,
				tt.description,
				tt.isEncrypted,
			)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, config)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
				assert.Equal(t, tt.key, config.Key())
				assert.Equal(t, tt.value, config.Value())
				assert.Equal(t, tt.category, config.Category())
				assert.Equal(t, tt.environment, config.Environment())
				assert.Equal(t, tt.description, config.Description())
				assert.Equal(t, tt.isEncrypted, config.IsEncrypted())
				assert.NotEqual(t, uuid.Nil, config.ID())
				assert.False(t, config.CreatedAt().IsZero())
				assert.False(t, config.UpdatedAt().IsZero())
			}
		})
	}
}

func TestConfiguration_UpdateValue(t *testing.T) {
	config, err := settings.NewConfiguration(
		"database.host",
		map[string]interface{}{"host": "localhost", "port": 5432},
		"database",
		"development",
		"Database connection settings",
		false,
	)
	require.NoError(t, err)

	originalUpdatedAt := config.UpdatedAt()

	// Wait a bit to ensure updated_at changes
	time.Sleep(1 * time.Millisecond)

	newValue := map[string]interface{}{"host": "newhost", "port": 3306}
	err = config.UpdateValue(newValue)
	assert.NoError(t, err)
	assert.Equal(t, newValue, config.Value())
	assert.True(t, config.UpdatedAt().After(originalUpdatedAt))
}

func TestConfiguration_UpdateValue_Invalid(t *testing.T) {
	config, err := settings.NewConfiguration(
		"database.host",
		map[string]interface{}{"host": "localhost", "port": 5432},
		"database",
		"development",
		"Database connection settings",
		false,
	)
	require.NoError(t, err)

	err = config.UpdateValue(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "value cannot be nil")
}

func TestConfiguration_UpdateCategory(t *testing.T) {
	config, err := settings.NewConfiguration(
		"database.host",
		map[string]interface{}{"host": "localhost", "port": 5432},
		"database",
		"development",
		"Database connection settings",
		false,
	)
	require.NoError(t, err)

	originalUpdatedAt := config.UpdatedAt()

	// Wait a bit to ensure updated_at changes
	time.Sleep(1 * time.Millisecond)

	err = config.UpdateCategory("newcategory")
	assert.NoError(t, err)
	assert.Equal(t, "newcategory", config.Category())
	assert.True(t, config.UpdatedAt().After(originalUpdatedAt))
}

func TestConfiguration_UpdateCategory_Invalid(t *testing.T) {
	config, err := settings.NewConfiguration(
		"database.host",
		map[string]interface{}{"host": "localhost", "port": 5432},
		"database",
		"development",
		"Database connection settings",
		false,
	)
	require.NoError(t, err)

	err = config.UpdateCategory("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "category cannot be empty")
}

func TestConfiguration_UpdateEnvironment(t *testing.T) {
	config, err := settings.NewConfiguration(
		"database.host",
		map[string]interface{}{"host": "localhost", "port": 5432},
		"database",
		"development",
		"Database connection settings",
		false,
	)
	require.NoError(t, err)

	originalUpdatedAt := config.UpdatedAt()

	// Wait a bit to ensure updated_at changes
	time.Sleep(1 * time.Millisecond)

	err = config.UpdateEnvironment("production")
	assert.NoError(t, err)
	assert.Equal(t, "production", config.Environment())
	assert.True(t, config.UpdatedAt().After(originalUpdatedAt))
}

func TestConfiguration_UpdateEnvironment_Invalid(t *testing.T) {
	config, err := settings.NewConfiguration(
		"database.host",
		map[string]interface{}{"host": "localhost", "port": 5432},
		"database",
		"development",
		"Database connection settings",
		false,
	)
	require.NoError(t, err)

	err = config.UpdateEnvironment("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "environment cannot be empty")
}

func TestConfiguration_UpdateDescription(t *testing.T) {
	config, err := settings.NewConfiguration(
		"database.host",
		map[string]interface{}{"host": "localhost", "port": 5432},
		"database",
		"development",
		"Database connection settings",
		false,
	)
	require.NoError(t, err)

	originalUpdatedAt := config.UpdatedAt()

	// Wait a bit to ensure updated_at changes
	time.Sleep(1 * time.Millisecond)

	err = config.UpdateDescription("New description")
	assert.NoError(t, err)
	assert.Equal(t, "New description", config.Description())
	assert.True(t, config.UpdatedAt().After(originalUpdatedAt))
}

func TestConfiguration_UpdateEncryption(t *testing.T) {
	config, err := settings.NewConfiguration(
		"database.host",
		map[string]interface{}{"host": "localhost", "port": 5432},
		"database",
		"development",
		"Database connection settings",
		false,
	)
	require.NoError(t, err)

	originalUpdatedAt := config.UpdatedAt()

	// Wait a bit to ensure updated_at changes
	time.Sleep(1 * time.Millisecond)

	err = config.UpdateEncryption(true)
	assert.NoError(t, err)
	assert.True(t, config.IsEncrypted())
	assert.True(t, config.UpdatedAt().After(originalUpdatedAt))
}

func TestConfiguration_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		config   *settings.Configuration
		expected bool
	}{
		{
			name: "valid configuration",
			config: func() *settings.Configuration {
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
			expected: true,
		},
		{
			name: "configuration with empty key",
			config: func() *settings.Configuration {
				config, _ := settings.NewConfiguration(
					"",
					map[string]interface{}{"host": "localhost", "port": 5432},
					"database",
					"development",
					"Database connection settings",
					false,
				)
				return config
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.config.IsValid())
		})
	}
}
