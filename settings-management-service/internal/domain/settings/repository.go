package settings

import (
	"context"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// ConfigurationRepository defines the interface for configuration persistence
type ConfigurationRepository interface {
	// Basic CRUD operations
	Save(ctx context.Context, configuration *Configuration) error
	FindByID(ctx context.Context, id ConfigurationID) (*Configuration, error)
	FindByKey(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName) (*Configuration, error)
	Update(ctx context.Context, configuration *Configuration) error
	Delete(ctx context.Context, id ConfigurationID) error

	// Query operations
	List(ctx context.Context, options ListOptions) ([]*Configuration, error)
	Count(ctx context.Context, filters ListFilters) (int, error)

	// Advanced queries
	FindByService(ctx context.Context, service ServiceName, environment Environment, organizationID *OrganizationID) ([]*Configuration, error)
	FindByEnvironment(ctx context.Context, environment Environment, organizationID *OrganizationID) ([]*Configuration, error)
	FindByOrganization(ctx context.Context, organizationID OrganizationID) ([]*Configuration, error)
	FindByCategory(ctx context.Context, category string, environment Environment, organizationID *OrganizationID) ([]*Configuration, error)

	// Existence checks
	ExistsByKey(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName) (bool, error)
	ExistsByID(ctx context.Context, id ConfigurationID) (bool, error)

	// Bulk operations
	BulkSave(ctx context.Context, configurations []*Configuration) error
	BulkUpdate(ctx context.Context, configurations []*Configuration) error
	BulkDelete(ctx context.Context, ids []ConfigurationID) error
}

// FeatureFlagRepository defines the interface for feature flag persistence
type FeatureFlagRepository interface {
	// Basic CRUD operations
	Save(ctx context.Context, featureFlag *FeatureFlag) error
	FindByID(ctx context.Context, id FeatureFlagID) (*FeatureFlag, error)
	FindByKey(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName) (*FeatureFlag, error)
	Update(ctx context.Context, featureFlag *FeatureFlag) error
	Delete(ctx context.Context, id FeatureFlagID) error

	// Query operations
	List(ctx context.Context, options ListOptions) ([]*FeatureFlag, error)
	Count(ctx context.Context, filters ListFilters) (int, error)

	// Advanced queries
	FindByService(ctx context.Context, service ServiceName, environment Environment, organizationID *OrganizationID) ([]*FeatureFlag, error)
	FindByEnvironment(ctx context.Context, environment Environment, organizationID *OrganizationID) ([]*FeatureFlag, error)
	FindByOrganization(ctx context.Context, organizationID OrganizationID) ([]*FeatureFlag, error)
	FindByCategory(ctx context.Context, category string, environment Environment, organizationID *OrganizationID) ([]*FeatureFlag, error)
	FindEnabled(ctx context.Context, environment Environment, organizationID *OrganizationID) ([]*FeatureFlag, error)

	// Existence checks
	ExistsByKey(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName) (bool, error)
	ExistsByID(ctx context.Context, id FeatureFlagID) (bool, error)

	// Bulk operations
	BulkSave(ctx context.Context, featureFlags []*FeatureFlag) error
	BulkUpdate(ctx context.Context, featureFlags []*FeatureFlag) error
	BulkDelete(ctx context.Context, ids []FeatureFlagID) error
}

// UserPreferenceRepository defines the interface for user preference persistence
type UserPreferenceRepository interface {
	// Basic CRUD operations
	Save(ctx context.Context, preference *UserPreference) error
	FindByID(ctx context.Context, id UserPreferenceID) (*UserPreference, error)
	FindByKey(ctx context.Context, userID UserID, category string, key string, organizationID *OrganizationID) (*UserPreference, error)
	Update(ctx context.Context, preference *UserPreference) error
	Delete(ctx context.Context, id UserPreferenceID) error

	// Query operations
	List(ctx context.Context, options ListOptions) ([]*UserPreference, error)
	Count(ctx context.Context, filters ListFilters) (int, error)

	// Advanced queries
	FindByUser(ctx context.Context, userID UserID, organizationID *OrganizationID) ([]*UserPreference, error)
	FindByUserAndCategory(ctx context.Context, userID UserID, category string, organizationID *OrganizationID) ([]*UserPreference, error)
	FindByOrganization(ctx context.Context, organizationID OrganizationID) ([]*UserPreference, error)
	FindByCategory(ctx context.Context, category string, organizationID *OrganizationID) ([]*UserPreference, error)

	// Existence checks
	ExistsByKey(ctx context.Context, userID UserID, category string, key string, organizationID *OrganizationID) (bool, error)
	ExistsByID(ctx context.Context, id UserPreferenceID) (bool, error)

	// Bulk operations
	BulkSave(ctx context.Context, preferences []*UserPreference) error
	BulkUpdate(ctx context.Context, preferences []*UserPreference) error
	BulkDelete(ctx context.Context, ids []UserPreferenceID) error
}

// OrganizationSettingRepository defines the interface for organization setting persistence
type OrganizationSettingRepository interface {
	// Basic CRUD operations
	Save(ctx context.Context, setting *OrganizationSetting) error
	FindByID(ctx context.Context, id OrganizationSettingID) (*OrganizationSetting, error)
	FindByKey(ctx context.Context, organizationID OrganizationID, category string, key string) (*OrganizationSetting, error)
	Update(ctx context.Context, setting *OrganizationSetting) error
	Delete(ctx context.Context, id OrganizationSettingID) error

	// Query operations
	List(ctx context.Context, options ListOptions) ([]*OrganizationSetting, error)
	Count(ctx context.Context, filters ListFilters) (int, error)

	// Advanced queries
	FindByOrganization(ctx context.Context, organizationID OrganizationID) ([]*OrganizationSetting, error)
	FindByOrganizationAndCategory(ctx context.Context, organizationID OrganizationID, category string) ([]*OrganizationSetting, error)
	FindByCategory(ctx context.Context, category string) ([]*OrganizationSetting, error)
	FindByParentOrganization(ctx context.Context, parentID OrganizationID) ([]*OrganizationSetting, error)

	// Existence checks
	ExistsByKey(ctx context.Context, organizationID OrganizationID, category string, key string) (bool, error)
	ExistsByID(ctx context.Context, id OrganizationSettingID) (bool, error)

	// Bulk operations
	BulkSave(ctx context.Context, settings []*OrganizationSetting) error
	BulkUpdate(ctx context.Context, settings []*OrganizationSetting) error
	BulkDelete(ctx context.Context, ids []OrganizationSettingID) error
}

// AuditRepository defines the interface for audit log persistence
type AuditRepository interface {
	// Audit log operations
	SaveAuditRecord(ctx context.Context, record *AuditRecord) error
	FindAuditRecords(ctx context.Context, criteria AuditCriteria) ([]*AuditRecord, error)
	CountAuditRecords(ctx context.Context, criteria AuditCriteria) (int, error)

	// Audit queries by entity
	FindConfigurationAuditLog(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName, criteria AuditCriteria) ([]*AuditRecord, error)
	FindFeatureFlagAuditLog(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName, criteria AuditCriteria) ([]*AuditRecord, error)
	FindUserPreferenceAuditLog(ctx context.Context, userID UserID, organizationID *OrganizationID, criteria AuditCriteria) ([]*AuditRecord, error)
	FindOrganizationSettingAuditLog(ctx context.Context, organizationID OrganizationID, criteria AuditCriteria) ([]*AuditRecord, error)

	// Compliance operations
	SaveComplianceReport(ctx context.Context, report *ComplianceReport) error
	FindComplianceReport(ctx context.Context, id string) (*ComplianceReport, error)
	FindComplianceReports(ctx context.Context, criteria ComplianceCriteria) ([]*ComplianceReport, error)

	SaveComplianceValidation(ctx context.Context, validation *ComplianceValidation) error
	FindComplianceValidation(ctx context.Context, organizationID OrganizationID) (*ComplianceValidation, error)
	FindComplianceValidations(ctx context.Context, criteria ComplianceCriteria) ([]*ComplianceValidation, error)
}

// CacheRepository defines the interface for settings caching
type CacheRepository interface {
	// Configuration caching
	GetConfiguration(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName) (*Configuration, error)
	SetConfiguration(ctx context.Context, configuration *Configuration, ttl time.Duration) error
	DeleteConfiguration(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName) error

	// Feature flag caching
	GetFeatureFlag(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName) (*FeatureFlag, error)
	SetFeatureFlag(ctx context.Context, featureFlag *FeatureFlag, ttl time.Duration) error
	DeleteFeatureFlag(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName) error

	// User preference caching
	GetUserPreference(ctx context.Context, userID UserID, category string, key string, organizationID *OrganizationID) (*UserPreference, error)
	SetUserPreference(ctx context.Context, preference *UserPreference, ttl time.Duration) error
	DeleteUserPreference(ctx context.Context, userID UserID, category string, key string, organizationID *OrganizationID) error

	// Organization setting caching
	GetOrganizationSetting(ctx context.Context, organizationID OrganizationID, category string, key string) (*OrganizationSetting, error)
	SetOrganizationSetting(ctx context.Context, setting *OrganizationSetting, ttl time.Duration) error
	DeleteOrganizationSetting(ctx context.Context, organizationID OrganizationID, category string, key string) error

	// Bulk cache operations
	GetConfigurationsByService(ctx context.Context, service ServiceName, environment Environment, organizationID *OrganizationID) ([]*Configuration, error)
	SetConfigurationsByService(ctx context.Context, configurations []*Configuration, ttl time.Duration) error
	DeleteConfigurationsByService(ctx context.Context, service ServiceName, environment Environment, organizationID *OrganizationID) error

	GetFeatureFlagsByService(ctx context.Context, service ServiceName, environment Environment, organizationID *OrganizationID) ([]*FeatureFlag, error)
	SetFeatureFlagsByService(ctx context.Context, featureFlags []*FeatureFlag, ttl time.Duration) error
	DeleteFeatureFlagsByService(ctx context.Context, service ServiceName, environment Environment, organizationID *OrganizationID) error

	GetUserPreferencesByUser(ctx context.Context, userID UserID, organizationID *OrganizationID) ([]*UserPreference, error)
	SetUserPreferencesByUser(ctx context.Context, preferences []*UserPreference, ttl time.Duration) error
	DeleteUserPreferencesByUser(ctx context.Context, userID UserID, organizationID *OrganizationID) error

	GetOrganizationSettingsByOrganization(ctx context.Context, organizationID OrganizationID) ([]*OrganizationSetting, error)
	SetOrganizationSettingsByOrganization(ctx context.Context, settings []*OrganizationSetting, ttl time.Duration) error
	DeleteOrganizationSettingsByOrganization(ctx context.Context, organizationID OrganizationID) error

	// Cache management
	ClearAll(ctx context.Context) error
	ClearByPattern(ctx context.Context, pattern string) error
	ClearByOrganization(ctx context.Context, organizationID OrganizationID) error
	ClearByService(ctx context.Context, service ServiceName) error
}

// ListOptions defines common options for listing entities
type ListOptions struct {
	Page      int         `json:"page"`
	Limit     int         `json:"limit"`
	SortBy    string      `json:"sortBy"`
	SortOrder string      `json:"sortOrder"`
	Filters   ListFilters `json:"filters"`
}

// ListFilters defines common filters for listing entities
type ListFilters map[string]interface{}

// NewListOptions creates a new ListOptions with defaults
func NewListOptions(page, limit int, sortBy, sortOrder string, filters ListFilters) ListOptions {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}
	if filters == nil {
		filters = make(ListFilters)
	}

	return ListOptions{
		Page:      page,
		Limit:     limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Filters:   filters,
	}
}

// GetOffset calculates the offset for pagination
func (o ListOptions) GetOffset() int {
	return (o.Page - 1) * o.Limit
}

// Validate validates the list options
func (o ListOptions) Validate() error {
	if o.Page < 1 {
		return shared.NewValidationError("page must be greater than 0", nil)
	}
	if o.Limit < 1 || o.Limit > 1000 {
		return shared.NewValidationError("limit must be between 1 and 1000", nil)
	}
	if o.SortOrder != "asc" && o.SortOrder != "desc" {
		return shared.NewValidationError("sort order must be 'asc' or 'desc'", nil)
	}
	return nil
}
