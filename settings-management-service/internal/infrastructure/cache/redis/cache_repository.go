package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/settings"
)

// CacheRepository implements the settings.CacheRepository interface using Redis
type CacheRepository struct {
	client *redis.Client
	prefix string
}

// NewCacheRepository creates a new Redis-based cache repository
func NewCacheRepository(client *redis.Client, prefix string) *CacheRepository {
	return &CacheRepository{
		client: client,
		prefix: prefix,
	}
}

// Configuration caching methods

// GetConfiguration retrieves a cached configuration
func (r *CacheRepository) GetConfiguration(ctx context.Context, key string, environment settings.Environment, organizationID *settings.OrganizationID, service *settings.ServiceName) (*settings.Configuration, error) {
	cacheKey := r.buildConfigurationKey(key, environment, organizationID, service)

	data, err := r.client.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get configuration from cache: %w", err)
	}

	var config settings.Configuration
	if err := json.Unmarshal([]byte(data), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration from cache: %w", err)
	}

	return &config, nil
}

// SetConfiguration caches a configuration
func (r *CacheRepository) SetConfiguration(ctx context.Context, configuration *settings.Configuration, ttl time.Duration) error {
	cacheKey := r.buildConfigurationKey(configuration.GetKey(), configuration.GetEnvironment(), configuration.GetOrganizationID(), configuration.GetService())

	data, err := json.Marshal(configuration)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration for cache: %w", err)
	}

	if err := r.client.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set configuration in cache: %w", err)
	}

	return nil
}

// DeleteConfiguration removes a configuration from cache
func (r *CacheRepository) DeleteConfiguration(ctx context.Context, key string, environment settings.Environment, organizationID *settings.OrganizationID, service *settings.ServiceName) error {
	cacheKey := r.buildConfigurationKey(key, environment, organizationID, service)

	if err := r.client.Del(ctx, cacheKey).Err(); err != nil {
		return fmt.Errorf("failed to delete configuration from cache: %w", err)
	}

	return nil
}

// Feature flag caching methods

// GetFeatureFlag retrieves a cached feature flag
func (r *CacheRepository) GetFeatureFlag(ctx context.Context, key string, environment settings.Environment, organizationID *settings.OrganizationID, service *settings.ServiceName) (*settings.FeatureFlag, error) {
	cacheKey := r.buildFeatureFlagKey(key, environment, organizationID, service)

	data, err := r.client.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get feature flag from cache: %w", err)
	}

	var featureFlag settings.FeatureFlag
	if err := json.Unmarshal([]byte(data), &featureFlag); err != nil {
		return nil, fmt.Errorf("failed to unmarshal feature flag from cache: %w", err)
	}

	return &featureFlag, nil
}

// SetFeatureFlag caches a feature flag
func (r *CacheRepository) SetFeatureFlag(ctx context.Context, featureFlag *settings.FeatureFlag, ttl time.Duration) error {
	cacheKey := r.buildFeatureFlagKey(featureFlag.GetKey(), featureFlag.GetEnvironment(), featureFlag.GetOrganizationID(), featureFlag.GetService())

	data, err := json.Marshal(featureFlag)
	if err != nil {
		return fmt.Errorf("failed to marshal feature flag for cache: %w", err)
	}

	if err := r.client.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set feature flag in cache: %w", err)
	}

	return nil
}

// DeleteFeatureFlag removes a feature flag from cache
func (r *CacheRepository) DeleteFeatureFlag(ctx context.Context, key string, environment settings.Environment, organizationID *settings.OrganizationID, service *settings.ServiceName) error {
	cacheKey := r.buildFeatureFlagKey(key, environment, organizationID, service)

	if err := r.client.Del(ctx, cacheKey).Err(); err != nil {
		return fmt.Errorf("failed to delete feature flag from cache: %w", err)
	}

	return nil
}

// User preference caching methods

// GetUserPreference retrieves a cached user preference
func (r *CacheRepository) GetUserPreference(ctx context.Context, userID settings.UserID, category string, key string, organizationID *settings.OrganizationID) (*settings.UserPreference, error) {
	cacheKey := r.buildUserPreferenceKey(userID, category, key, organizationID)

	data, err := r.client.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get user preference from cache: %w", err)
	}

	var preference settings.UserPreference
	if err := json.Unmarshal([]byte(data), &preference); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user preference from cache: %w", err)
	}

	return &preference, nil
}

// SetUserPreference caches a user preference
func (r *CacheRepository) SetUserPreference(ctx context.Context, preference *settings.UserPreference, ttl time.Duration) error {
	cacheKey := r.buildUserPreferenceKey(preference.GetUserID(), preference.GetCategory(), preference.GetKey(), preference.GetOrganizationID())

	data, err := json.Marshal(preference)
	if err != nil {
		return fmt.Errorf("failed to marshal user preference for cache: %w", err)
	}

	if err := r.client.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set user preference in cache: %w", err)
	}

	return nil
}

// DeleteUserPreference removes a user preference from cache
func (r *CacheRepository) DeleteUserPreference(ctx context.Context, userID settings.UserID, category string, key string, organizationID *settings.OrganizationID) error {
	cacheKey := r.buildUserPreferenceKey(userID, category, key, organizationID)

	if err := r.client.Del(ctx, cacheKey).Err(); err != nil {
		return fmt.Errorf("failed to delete user preference from cache: %w", err)
	}

	return nil
}

// Organization setting caching methods

// GetOrganizationSetting retrieves a cached organization setting
func (r *CacheRepository) GetOrganizationSetting(ctx context.Context, organizationID settings.OrganizationID, category string, key string) (*settings.OrganizationSetting, error) {
	cacheKey := r.buildOrganizationSettingKey(organizationID, category, key)

	data, err := r.client.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get organization setting from cache: %w", err)
	}

	var setting settings.OrganizationSetting
	if err := json.Unmarshal([]byte(data), &setting); err != nil {
		return nil, fmt.Errorf("failed to unmarshal organization setting from cache: %w", err)
	}

	return &setting, nil
}

// SetOrganizationSetting caches an organization setting
func (r *CacheRepository) SetOrganizationSetting(ctx context.Context, setting *settings.OrganizationSetting, ttl time.Duration) error {
	cacheKey := r.buildOrganizationSettingKey(setting.GetOrganizationID(), setting.GetCategory(), setting.GetKey())

	data, err := json.Marshal(setting)
	if err != nil {
		return fmt.Errorf("failed to marshal organization setting for cache: %w", err)
	}

	if err := r.client.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set organization setting in cache: %w", err)
	}

	return nil
}

// DeleteOrganizationSetting removes an organization setting from cache
func (r *CacheRepository) DeleteOrganizationSetting(ctx context.Context, organizationID settings.OrganizationID, category string, key string) error {
	cacheKey := r.buildOrganizationSettingKey(organizationID, category, key)

	if err := r.client.Del(ctx, cacheKey).Err(); err != nil {
		return fmt.Errorf("failed to delete organization setting from cache: %w", err)
	}

	return nil
}

// Bulk cache operations

// GetConfigurationsByService retrieves cached configurations by service
func (r *CacheRepository) GetConfigurationsByService(ctx context.Context, service settings.ServiceName, environment settings.Environment, organizationID *settings.OrganizationID) ([]*settings.Configuration, error) {
	pattern := r.buildConfigurationPattern(service, environment, organizationID)
	return r.getConfigurationsByPattern(ctx, pattern)
}

// SetConfigurationsByService caches configurations by service
func (r *CacheRepository) SetConfigurationsByService(ctx context.Context, configurations []*settings.Configuration, ttl time.Duration) error {
	for _, config := range configurations {
		if err := r.SetConfiguration(ctx, config, ttl); err != nil {
			return err
		}
	}
	return nil
}

// DeleteConfigurationsByService removes cached configurations by service
func (r *CacheRepository) DeleteConfigurationsByService(ctx context.Context, service settings.ServiceName, environment settings.Environment, organizationID *settings.OrganizationID) error {
	pattern := r.buildConfigurationPattern(service, environment, organizationID)
	return r.deleteByPattern(ctx, pattern)
}

// GetFeatureFlagsByService retrieves cached feature flags by service
func (r *CacheRepository) GetFeatureFlagsByService(ctx context.Context, service settings.ServiceName, environment settings.Environment, organizationID *settings.OrganizationID) ([]*settings.FeatureFlag, error) {
	pattern := r.buildFeatureFlagPattern(service, environment, organizationID)
	return r.getFeatureFlagsByPattern(ctx, pattern)
}

// SetFeatureFlagsByService caches feature flags by service
func (r *CacheRepository) SetFeatureFlagsByService(ctx context.Context, featureFlags []*settings.FeatureFlag, ttl time.Duration) error {
	for _, featureFlag := range featureFlags {
		if err := r.SetFeatureFlag(ctx, featureFlag, ttl); err != nil {
			return err
		}
	}
	return nil
}

// DeleteFeatureFlagsByService removes cached feature flags by service
func (r *CacheRepository) DeleteFeatureFlagsByService(ctx context.Context, service settings.ServiceName, environment settings.Environment, organizationID *settings.OrganizationID) error {
	pattern := r.buildFeatureFlagPattern(service, environment, organizationID)
	return r.deleteByPattern(ctx, pattern)
}

// GetUserPreferencesByUser retrieves cached user preferences by user
func (r *CacheRepository) GetUserPreferencesByUser(ctx context.Context, userID settings.UserID, organizationID *settings.OrganizationID) ([]*settings.UserPreference, error) {
	pattern := r.buildUserPreferencePattern(userID, organizationID)
	return r.getUserPreferencesByPattern(ctx, pattern)
}

// SetUserPreferencesByUser caches user preferences by user
func (r *CacheRepository) SetUserPreferencesByUser(ctx context.Context, preferences []*settings.UserPreference, ttl time.Duration) error {
	for _, preference := range preferences {
		if err := r.SetUserPreference(ctx, preference, ttl); err != nil {
			return err
		}
	}
	return nil
}

// DeleteUserPreferencesByUser removes cached user preferences by user
func (r *CacheRepository) DeleteUserPreferencesByUser(ctx context.Context, userID settings.UserID, organizationID *settings.OrganizationID) error {
	pattern := r.buildUserPreferencePattern(userID, organizationID)
	return r.deleteByPattern(ctx, pattern)
}

// GetOrganizationSettingsByOrganization retrieves cached organization settings by organization
func (r *CacheRepository) GetOrganizationSettingsByOrganization(ctx context.Context, organizationID settings.OrganizationID) ([]*settings.OrganizationSetting, error) {
	pattern := r.buildOrganizationSettingPattern(organizationID)
	return r.getOrganizationSettingsByPattern(ctx, pattern)
}

// SetOrganizationSettingsByOrganization caches organization settings by organization
func (r *CacheRepository) SetOrganizationSettingsByOrganization(ctx context.Context, settings []*settings.OrganizationSetting, ttl time.Duration) error {
	for _, setting := range settings {
		if err := r.SetOrganizationSetting(ctx, setting, ttl); err != nil {
			return err
		}
	}
	return nil
}

// DeleteOrganizationSettingsByOrganization removes cached organization settings by organization
func (r *CacheRepository) DeleteOrganizationSettingsByOrganization(ctx context.Context, organizationID settings.OrganizationID) error {
	pattern := r.buildOrganizationSettingPattern(organizationID)
	return r.deleteByPattern(ctx, pattern)
}

// Cache management methods

// ClearAll clears all cached data
func (r *CacheRepository) ClearAll(ctx context.Context) error {
	pattern := fmt.Sprintf("%s:*", r.prefix)
	return r.deleteByPattern(ctx, pattern)
}

// ClearByPattern clears cached data matching a pattern
func (r *CacheRepository) ClearByPattern(ctx context.Context, pattern string) error {
	fullPattern := fmt.Sprintf("%s:%s", r.prefix, pattern)
	return r.deleteByPattern(ctx, fullPattern)
}

// ClearByOrganization clears cached data for a specific organization
func (r *CacheRepository) ClearByOrganization(ctx context.Context, organizationID settings.OrganizationID) error {
	pattern := fmt.Sprintf("%s:*:org:%s:*", r.prefix, organizationID.String())
	return r.deleteByPattern(ctx, pattern)
}

// ClearByService clears cached data for a specific service
func (r *CacheRepository) ClearByService(ctx context.Context, service settings.ServiceName) error {
	pattern := fmt.Sprintf("%s:*:svc:%s:*", r.prefix, service.String())
	return r.deleteByPattern(ctx, pattern)
}

// Helper methods for building cache keys

func (r *CacheRepository) buildConfigurationKey(key string, environment settings.Environment, organizationID *settings.OrganizationID, service *settings.ServiceName) string {
	return fmt.Sprintf("%s:config:%s:env:%s:org:%s:svc:%s",
		r.prefix,
		key,
		environment.String(),
		r.getOrganizationIDString(organizationID),
		r.getServiceString(service),
	)
}

func (r *CacheRepository) buildFeatureFlagKey(key string, environment settings.Environment, organizationID *settings.OrganizationID, service *settings.ServiceName) string {
	return fmt.Sprintf("%s:feature-flag:%s:env:%s:org:%s:svc:%s",
		r.prefix,
		key,
		environment.String(),
		r.getOrganizationIDString(organizationID),
		r.getServiceString(service),
	)
}

func (r *CacheRepository) buildUserPreferenceKey(userID settings.UserID, category string, key string, organizationID *settings.OrganizationID) string {
	return fmt.Sprintf("%s:user-preference:%s:cat:%s:key:%s:org:%s",
		r.prefix,
		userID.String(),
		category,
		key,
		r.getOrganizationIDString(organizationID),
	)
}

func (r *CacheRepository) buildOrganizationSettingKey(organizationID settings.OrganizationID, category string, key string) string {
	return fmt.Sprintf("%s:organization-setting:%s:cat:%s:key:%s",
		r.prefix,
		organizationID.String(),
		category,
		key,
	)
}

func (r *CacheRepository) buildConfigurationPattern(service settings.ServiceName, environment settings.Environment, organizationID *settings.OrganizationID) string {
	return fmt.Sprintf("%s:config:*:env:%s:org:%s:svc:%s",
		r.prefix,
		environment.String(),
		r.getOrganizationIDString(organizationID),
		service.String(),
	)
}

func (r *CacheRepository) buildFeatureFlagPattern(service settings.ServiceName, environment settings.Environment, organizationID *settings.OrganizationID) string {
	return fmt.Sprintf("%s:feature-flag:*:env:%s:org:%s:svc:%s",
		r.prefix,
		environment.String(),
		r.getOrganizationIDString(organizationID),
		service.String(),
	)
}

func (r *CacheRepository) buildUserPreferencePattern(userID settings.UserID, organizationID *settings.OrganizationID) string {
	return fmt.Sprintf("%s:user-preference:%s:*:org:%s",
		r.prefix,
		userID.String(),
		r.getOrganizationIDString(organizationID),
	)
}

func (r *CacheRepository) buildOrganizationSettingPattern(organizationID settings.OrganizationID) string {
	return fmt.Sprintf("%s:organization-setting:%s:*",
		r.prefix,
		organizationID.String(),
	)
}

func (r *CacheRepository) getOrganizationIDString(organizationID *settings.OrganizationID) string {
	if organizationID == nil {
		return "nil"
	}
	return organizationID.String()
}

func (r *CacheRepository) getServiceString(service *settings.ServiceName) string {
	if service == nil {
		return "nil"
	}
	return service.String()
}

// Helper methods for pattern-based operations

func (r *CacheRepository) deleteByPattern(ctx context.Context, pattern string) error {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys for pattern %s: %w", pattern, err)
	}

	if len(keys) > 0 {
		if err := r.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to delete keys for pattern %s: %w", pattern, err)
		}
		log.Printf("Deleted %d keys matching pattern %s", len(keys), pattern)
	}

	return nil
}

func (r *CacheRepository) getConfigurationsByPattern(ctx context.Context, pattern string) ([]*settings.Configuration, error) {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys for pattern %s: %w", pattern, err)
	}

	configurations := make([]*settings.Configuration, 0, len(keys))
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue // Skip failed keys
		}

		var config settings.Configuration
		if err := json.Unmarshal([]byte(data), &config); err != nil {
			continue // Skip invalid data
		}

		configurations = append(configurations, &config)
	}

	return configurations, nil
}

func (r *CacheRepository) getFeatureFlagsByPattern(ctx context.Context, pattern string) ([]*settings.FeatureFlag, error) {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys for pattern %s: %w", pattern, err)
	}

	featureFlags := make([]*settings.FeatureFlag, 0, len(keys))
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue // Skip failed keys
		}

		var featureFlag settings.FeatureFlag
		if err := json.Unmarshal([]byte(data), &featureFlag); err != nil {
			continue // Skip invalid data
		}

		featureFlags = append(featureFlags, &featureFlag)
	}

	return featureFlags, nil
}

func (r *CacheRepository) getUserPreferencesByPattern(ctx context.Context, pattern string) ([]*settings.UserPreference, error) {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys for pattern %s: %w", pattern, err)
	}

	preferences := make([]*settings.UserPreference, 0, len(keys))
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue // Skip failed keys
		}

		var preference settings.UserPreference
		if err := json.Unmarshal([]byte(data), &preference); err != nil {
			continue // Skip invalid data
		}

		preferences = append(preferences, &preference)
	}

	return preferences, nil
}

func (r *CacheRepository) getOrganizationSettingsByPattern(ctx context.Context, pattern string) ([]*settings.OrganizationSetting, error) {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys for pattern %s: %w", pattern, err)
	}

	settings := make([]*settings.OrganizationSetting, 0, len(keys))
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue // Skip failed keys
		}

		var setting settings.OrganizationSetting
		if err := json.Unmarshal([]byte(data), &setting); err != nil {
			continue // Skip invalid data
		}

		settings = append(settings, &setting)
	}

	return settings, nil
}
