package settings_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
)

func TestNewFeatureFlag(t *testing.T) {
	tests := []struct {
		name              string
		nameValue         string
		key               string
		description       string
		status            settings.FeatureFlagStatus
		rolloutPercentage int
		targetAudience    map[string]interface{}
		conditions        map[string]interface{}
		expectError       bool
		errorMsg          string
	}{
		{
			name:              "valid feature flag",
			nameValue:         "New Feature",
			key:               "new-feature",
			description:       "A new feature for testing",
			status:            settings.FeatureFlagStatusActive,
			rolloutPercentage: 50,
			targetAudience:    map[string]interface{}{"users": []string{"user1", "user2"}},
			conditions:        map[string]interface{}{"country": "US"},
			expectError:       false,
		},
		{
			name:              "empty name",
			nameValue:         "",
			key:               "new-feature",
			description:       "A new feature for testing",
			status:            settings.FeatureFlagStatusActive,
			rolloutPercentage: 50,
			targetAudience:    map[string]interface{}{"users": []string{"user1", "user2"}},
			conditions:        map[string]interface{}{"country": "US"},
			expectError:       true,
			errorMsg:          "name cannot be empty",
		},
		{
			name:              "empty key",
			nameValue:         "New Feature",
			key:               "",
			description:       "A new feature for testing",
			status:            settings.FeatureFlagStatusActive,
			rolloutPercentage: 50,
			targetAudience:    map[string]interface{}{"users": []string{"user1", "user2"}},
			conditions:        map[string]interface{}{"country": "US"},
			expectError:       true,
			errorMsg:          "key cannot be empty",
		},
		{
			name:              "invalid rollout percentage - negative",
			nameValue:         "New Feature",
			key:               "new-feature",
			description:       "A new feature for testing",
			status:            settings.FeatureFlagStatusActive,
			rolloutPercentage: -1,
			targetAudience:    map[string]interface{}{"users": []string{"user1", "user2"}},
			conditions:        map[string]interface{}{"country": "US"},
			expectError:       true,
			errorMsg:          "rollout percentage must be between 0 and 100",
		},
		{
			name:              "invalid rollout percentage - over 100",
			nameValue:         "New Feature",
			key:               "new-feature",
			description:       "A new feature for testing",
			status:            settings.FeatureFlagStatusActive,
			rolloutPercentage: 101,
			targetAudience:    map[string]interface{}{"users": []string{"user1", "user2"}},
			conditions:        map[string]interface{}{"country": "US"},
			expectError:       true,
			errorMsg:          "rollout percentage must be between 0 and 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag, err := settings.NewFeatureFlag(
				tt.nameValue,
				tt.key,
				tt.description,
				tt.status,
				tt.rolloutPercentage,
				tt.targetAudience,
				tt.conditions,
			)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, flag)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, flag)
				assert.Equal(t, tt.nameValue, flag.Name())
				assert.Equal(t, tt.key, flag.Key())
				assert.Equal(t, tt.description, flag.Description())
				assert.Equal(t, tt.status, flag.Status())
				assert.Equal(t, tt.rolloutPercentage, flag.RolloutPercentage())
				assert.Equal(t, tt.targetAudience, flag.TargetAudience())
				assert.Equal(t, tt.conditions, flag.Conditions())
				assert.NotEqual(t, uuid.Nil, flag.ID())
				assert.False(t, flag.CreatedAt().IsZero())
				assert.False(t, flag.UpdatedAt().IsZero())
			}
		})
	}
}

func TestFeatureFlag_UpdateName(t *testing.T) {
	flag, err := settings.NewFeatureFlag(
		"New Feature",
		"new-feature",
		"A new feature for testing",
		settings.FeatureFlagStatusActive,
		50,
		map[string]interface{}{"users": []string{"user1", "user2"}},
		map[string]interface{}{"country": "US"},
	)
	require.NoError(t, err)

	originalUpdatedAt := flag.UpdatedAt()

	// Wait a bit to ensure updated_at changes
	time.Sleep(1 * time.Millisecond)

	err = flag.UpdateName("Updated Feature")
	assert.NoError(t, err)
	assert.Equal(t, "Updated Feature", flag.Name())
	assert.True(t, flag.UpdatedAt().After(originalUpdatedAt))
}

func TestFeatureFlag_UpdateName_Invalid(t *testing.T) {
	flag, err := settings.NewFeatureFlag(
		"New Feature",
		"new-feature",
		"A new feature for testing",
		settings.FeatureFlagStatusActive,
		50,
		map[string]interface{}{"users": []string{"user1", "user2"}},
		map[string]interface{}{"country": "US"},
	)
	require.NoError(t, err)

	err = flag.UpdateName("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name cannot be empty")
}

func TestFeatureFlag_UpdateDescription(t *testing.T) {
	flag, err := settings.NewFeatureFlag(
		"New Feature",
		"new-feature",
		"A new feature for testing",
		settings.FeatureFlagStatusActive,
		50,
		map[string]interface{}{"users": []string{"user1", "user2"}},
		map[string]interface{}{"country": "US"},
	)
	require.NoError(t, err)

	originalUpdatedAt := flag.UpdatedAt()

	// Wait a bit to ensure updated_at changes
	time.Sleep(1 * time.Millisecond)

	err = flag.UpdateDescription("Updated description")
	assert.NoError(t, err)
	assert.Equal(t, "Updated description", flag.Description())
	assert.True(t, flag.UpdatedAt().After(originalUpdatedAt))
}

func TestFeatureFlag_UpdateStatus(t *testing.T) {
	flag, err := settings.NewFeatureFlag(
		"New Feature",
		"new-feature",
		"A new feature for testing",
		settings.FeatureFlagStatusActive,
		50,
		map[string]interface{}{"users": []string{"user1", "user2"}},
		map[string]interface{}{"country": "US"},
	)
	require.NoError(t, err)

	originalUpdatedAt := flag.UpdatedAt()

	// Wait a bit to ensure updated_at changes
	time.Sleep(1 * time.Millisecond)

	err = flag.UpdateStatus(settings.FeatureFlagStatusInactive)
	assert.NoError(t, err)
	assert.Equal(t, settings.FeatureFlagStatusInactive, flag.Status())
	assert.True(t, flag.UpdatedAt().After(originalUpdatedAt))
}

func TestFeatureFlag_UpdateRolloutPercentage(t *testing.T) {
	flag, err := settings.NewFeatureFlag(
		"New Feature",
		"new-feature",
		"A new feature for testing",
		settings.FeatureFlagStatusActive,
		50,
		map[string]interface{}{"users": []string{"user1", "user2"}},
		map[string]interface{}{"country": "US"},
	)
	require.NoError(t, err)

	originalUpdatedAt := flag.UpdatedAt()

	// Wait a bit to ensure updated_at changes
	time.Sleep(1 * time.Millisecond)

	err = flag.UpdateRolloutPercentage(75)
	assert.NoError(t, err)
	assert.Equal(t, 75, flag.RolloutPercentage())
	assert.True(t, flag.UpdatedAt().After(originalUpdatedAt))
}

func TestFeatureFlag_UpdateRolloutPercentage_Invalid(t *testing.T) {
	flag, err := settings.NewFeatureFlag(
		"New Feature",
		"new-feature",
		"A new feature for testing",
		settings.FeatureFlagStatusActive,
		50,
		map[string]interface{}{"users": []string{"user1", "user2"}},
		map[string]interface{}{"country": "US"},
	)
	require.NoError(t, err)

	err = flag.UpdateRolloutPercentage(-1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rollout percentage must be between 0 and 100")

	err = flag.UpdateRolloutPercentage(101)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rollout percentage must be between 0 and 100")
}

func TestFeatureFlag_UpdateTargetAudience(t *testing.T) {
	flag, err := settings.NewFeatureFlag(
		"New Feature",
		"new-feature",
		"A new feature for testing",
		settings.FeatureFlagStatusActive,
		50,
		map[string]interface{}{"users": []string{"user1", "user2"}},
		map[string]interface{}{"country": "US"},
	)
	require.NoError(t, err)

	originalUpdatedAt := flag.UpdatedAt()

	// Wait a bit to ensure updated_at changes
	time.Sleep(1 * time.Millisecond)

	newTargetAudience := map[string]interface{}{"users": []string{"user3", "user4"}}
	err = flag.UpdateTargetAudience(newTargetAudience)
	assert.NoError(t, err)
	assert.Equal(t, newTargetAudience, flag.TargetAudience())
	assert.True(t, flag.UpdatedAt().After(originalUpdatedAt))
}

func TestFeatureFlag_UpdateConditions(t *testing.T) {
	flag, err := settings.NewFeatureFlag(
		"New Feature",
		"new-feature",
		"A new feature for testing",
		settings.FeatureFlagStatusActive,
		50,
		map[string]interface{}{"users": []string{"user1", "user2"}},
		map[string]interface{}{"country": "US"},
	)
	require.NoError(t, err)

	originalUpdatedAt := flag.UpdatedAt()

	// Wait a bit to ensure updated_at changes
	time.Sleep(1 * time.Millisecond)

	newConditions := map[string]interface{}{"country": "CA"}
	err = flag.UpdateConditions(newConditions)
	assert.NoError(t, err)
	assert.Equal(t, newConditions, flag.Conditions())
	assert.True(t, flag.UpdatedAt().After(originalUpdatedAt))
}

func TestFeatureFlag_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   settings.FeatureFlagStatus
		expected bool
	}{
		{
			name:     "active status",
			status:   settings.FeatureFlagStatusActive,
			expected: true,
		},
		{
			name:     "inactive status",
			status:   settings.FeatureFlagStatusInactive,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag, err := settings.NewFeatureFlag(
				"New Feature",
				"new-feature",
				"A new feature for testing",
				tt.status,
				50,
				map[string]interface{}{"users": []string{"user1", "user2"}},
				map[string]interface{}{"country": "US"},
			)
			require.NoError(t, err)

			assert.Equal(t, tt.expected, flag.IsActive())
		})
	}
}

func TestFeatureFlag_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		flag     *settings.FeatureFlag
		expected bool
	}{
		{
			name: "valid feature flag",
			flag: func() *settings.FeatureFlag {
				flag, _ := settings.NewFeatureFlag(
					"New Feature",
					"new-feature",
					"A new feature for testing",
					settings.FeatureFlagStatusActive,
					50,
					map[string]interface{}{"users": []string{"user1", "user2"}},
					map[string]interface{}{"country": "US"},
				)
				return flag
			}(),
			expected: true,
		},
		{
			name: "feature flag with empty name",
			flag: func() *settings.FeatureFlag {
				flag, _ := settings.NewFeatureFlag(
					"",
					"new-feature",
					"A new feature for testing",
					settings.FeatureFlagStatusActive,
					50,
					map[string]interface{}{"users": []string{"user1", "user2"}},
					map[string]interface{}{"country": "US"},
				)
				return flag
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.flag.IsValid())
		})
	}
}
