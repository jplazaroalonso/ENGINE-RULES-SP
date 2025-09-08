package settings

import (
	"context"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// ConfigurationService defines the interface for configuration management
type ConfigurationService interface {
	// Configuration management
	GetConfiguration(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName) (*Configuration, error)
	SetConfiguration(ctx context.Context, configuration *Configuration) error
	UpdateConfiguration(ctx context.Context, key string, value interface{}, environment Environment, organizationID *OrganizationID, service *ServiceName, updatedBy UserID) error
	DeleteConfiguration(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName, deletedBy UserID) error

	// Bulk operations
	GetConfigurationsByService(ctx context.Context, service ServiceName, environment Environment, organizationID *OrganizationID) ([]*Configuration, error)
	GetConfigurationsByEnvironment(ctx context.Context, environment Environment, organizationID *OrganizationID) ([]*Configuration, error)
	BulkUpdateConfigurations(ctx context.Context, configurations []*Configuration) error

	// Validation
	ValidateConfiguration(ctx context.Context, configuration *Configuration) error
	ValidateConfigurationKey(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName) error
}

// FeatureFlagService defines the interface for feature flag management
type FeatureFlagService interface {
	// Feature flag management
	GetFeatureFlag(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName) (*FeatureFlag, error)
	SetFeatureFlag(ctx context.Context, featureFlag *FeatureFlag) error
	UpdateFeatureFlag(ctx context.Context, key string, updates map[string]interface{}, environment Environment, organizationID *OrganizationID, service *ServiceName, updatedBy UserID) error
	DeleteFeatureFlag(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName, deletedBy UserID) error

	// Feature flag evaluation
	EvaluateFeatureFlag(ctx context.Context, key string, context map[string]interface{}, environment Environment, organizationID *OrganizationID, service *ServiceName) (*FeatureFlagEvaluation, error)
	EvaluateFeatureFlags(ctx context.Context, keys []string, context map[string]interface{}, environment Environment, organizationID *OrganizationID, service *ServiceName) (map[string]*FeatureFlagEvaluation, error)

	// Bulk operations
	GetFeatureFlagsByService(ctx context.Context, service ServiceName, environment Environment, organizationID *OrganizationID) ([]*FeatureFlag, error)
	GetFeatureFlagsByEnvironment(ctx context.Context, environment Environment, organizationID *OrganizationID) ([]*FeatureFlag, error)
	BulkUpdateFeatureFlags(ctx context.Context, featureFlags []*FeatureFlag) error

	// Validation
	ValidateFeatureFlag(ctx context.Context, featureFlag *FeatureFlag) error
	ValidateFeatureFlagKey(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName) error
}

// UserPreferenceService defines the interface for user preference management
type UserPreferenceService interface {
	// User preference management
	GetUserPreference(ctx context.Context, userID UserID, category string, key string, organizationID *OrganizationID) (*UserPreference, error)
	SetUserPreference(ctx context.Context, preference *UserPreference) error
	UpdateUserPreference(ctx context.Context, userID UserID, category string, key string, value interface{}, organizationID *OrganizationID) error
	DeleteUserPreference(ctx context.Context, userID UserID, category string, key string, organizationID *OrganizationID) error

	// Bulk operations
	GetUserPreferences(ctx context.Context, userID UserID, organizationID *OrganizationID) ([]*UserPreference, error)
	GetUserPreferencesByCategory(ctx context.Context, userID UserID, category string, organizationID *OrganizationID) ([]*UserPreference, error)
	BulkUpdateUserPreferences(ctx context.Context, preferences []*UserPreference) error

	// Preference hierarchy
	GetEffectiveUserPreference(ctx context.Context, userID UserID, category string, key string, organizationID *OrganizationID) (*UserPreference, error)
	GetEffectiveUserPreferences(ctx context.Context, userID UserID, organizationID *OrganizationID) ([]*UserPreference, error)

	// Validation
	ValidateUserPreference(ctx context.Context, preference *UserPreference) error
}

// OrganizationSettingService defines the interface for organization setting management
type OrganizationSettingService interface {
	// Organization setting management
	GetOrganizationSetting(ctx context.Context, organizationID OrganizationID, category string, key string) (*OrganizationSetting, error)
	SetOrganizationSetting(ctx context.Context, setting *OrganizationSetting) error
	UpdateOrganizationSetting(ctx context.Context, organizationID OrganizationID, category string, key string, value interface{}, updatedBy UserID) error
	DeleteOrganizationSetting(ctx context.Context, organizationID OrganizationID, category string, key string, deletedBy UserID) error

	// Bulk operations
	GetOrganizationSettings(ctx context.Context, organizationID OrganizationID) ([]*OrganizationSetting, error)
	GetOrganizationSettingsByCategory(ctx context.Context, organizationID OrganizationID, category string) ([]*OrganizationSetting, error)
	BulkUpdateOrganizationSettings(ctx context.Context, settings []*OrganizationSetting) error

	// Setting hierarchy and inheritance
	GetEffectiveOrganizationSetting(ctx context.Context, organizationID OrganizationID, category string, key string) (*OrganizationSetting, error)
	GetEffectiveOrganizationSettings(ctx context.Context, organizationID OrganizationID) ([]*OrganizationSetting, error)
	InheritFromParent(ctx context.Context, organizationID OrganizationID, parentID OrganizationID, updatedBy UserID) error

	// Validation
	ValidateOrganizationSetting(ctx context.Context, setting *OrganizationSetting) error
}

// SettingsCacheService defines the interface for settings caching
type SettingsCacheService interface {
	// Cache operations
	GetCachedConfiguration(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName) (*Configuration, error)
	SetCachedConfiguration(ctx context.Context, configuration *Configuration, ttl time.Duration) error
	InvalidateCachedConfiguration(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName) error

	GetCachedFeatureFlag(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName) (*FeatureFlag, error)
	SetCachedFeatureFlag(ctx context.Context, featureFlag *FeatureFlag, ttl time.Duration) error
	InvalidateCachedFeatureFlag(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName) error

	GetCachedUserPreference(ctx context.Context, userID UserID, category string, key string, organizationID *OrganizationID) (*UserPreference, error)
	SetCachedUserPreference(ctx context.Context, preference *UserPreference, ttl time.Duration) error
	InvalidateCachedUserPreference(ctx context.Context, userID UserID, category string, key string, organizationID *OrganizationID) error

	GetCachedOrganizationSetting(ctx context.Context, organizationID OrganizationID, category string, key string) (*OrganizationSetting, error)
	SetCachedOrganizationSetting(ctx context.Context, setting *OrganizationSetting, ttl time.Duration) error
	InvalidateCachedOrganizationSetting(ctx context.Context, organizationID OrganizationID, category string, key string) error

	// Bulk cache operations
	InvalidateAllCachedSettings(ctx context.Context) error
	InvalidateCachedSettingsByOrganization(ctx context.Context, organizationID OrganizationID) error
	InvalidateCachedSettingsByService(ctx context.Context, service ServiceName) error
}

// SettingsAuditService defines the interface for settings audit and compliance
type SettingsAuditService interface {
	// Audit logging
	LogConfigurationChange(ctx context.Context, configuration *Configuration, action string, userID UserID, details map[string]interface{}) error
	LogFeatureFlagChange(ctx context.Context, featureFlag *FeatureFlag, action string, userID UserID, details map[string]interface{}) error
	LogUserPreferenceChange(ctx context.Context, preference *UserPreference, action string, userID UserID, details map[string]interface{}) error
	LogOrganizationSettingChange(ctx context.Context, setting *OrganizationSetting, action string, userID UserID, details map[string]interface{}) error

	// Audit queries
	GetConfigurationAuditLog(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName, criteria AuditCriteria) ([]AuditRecord, error)
	GetFeatureFlagAuditLog(ctx context.Context, key string, environment Environment, organizationID *OrganizationID, service *ServiceName, criteria AuditCriteria) ([]AuditRecord, error)
	GetUserPreferenceAuditLog(ctx context.Context, userID UserID, organizationID *OrganizationID, criteria AuditCriteria) ([]AuditRecord, error)
	GetOrganizationSettingAuditLog(ctx context.Context, organizationID OrganizationID, criteria AuditCriteria) ([]AuditRecord, error)

	// Compliance reporting
	GenerateComplianceReport(ctx context.Context, criteria ComplianceCriteria) (*ComplianceReport, error)
	ValidateCompliance(ctx context.Context, organizationID OrganizationID) (*ComplianceValidation, error)
}

// FeatureFlagEvaluation represents the result of feature flag evaluation
type FeatureFlagEvaluation struct {
	FeatureFlagID string                 `json:"featureFlagId"`
	Key           string                 `json:"key"`
	IsEnabled     bool                   `json:"isEnabled"`
	Variant       *string                `json:"variant,omitempty"`
	Value         interface{}            `json:"value,omitempty"`
	Reason        string                 `json:"reason"`
	EvaluatedAt   time.Time              `json:"evaluatedAt"`
	Context       map[string]interface{} `json:"context"`
}

// NewFeatureFlagEvaluation creates a new feature flag evaluation
func NewFeatureFlagEvaluation(featureFlagID, key string, isEnabled bool, variant *string, value interface{}, reason string, context map[string]interface{}) *FeatureFlagEvaluation {
	return &FeatureFlagEvaluation{
		FeatureFlagID: featureFlagID,
		Key:           key,
		IsEnabled:     isEnabled,
		Variant:       variant,
		Value:         value,
		Reason:        reason,
		EvaluatedAt:   time.Now(),
		Context:       context,
	}
}

// AuditCriteria represents criteria for querying audit logs
type AuditCriteria struct {
	Page      int                    `json:"page"`
	Limit     int                    `json:"limit"`
	SortBy    string                 `json:"sortBy"`
	SortOrder string                 `json:"sortOrder"`
	Actions   []string               `json:"actions,omitempty"`
	Users     []string               `json:"users,omitempty"`
	DateRange *DateRange             `json:"dateRange,omitempty"`
	Filters   map[string]interface{} `json:"filters"`
}

// DateRange represents a date range
type DateRange struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

// NewDateRange creates a new date range
func NewDateRange(startDate, endDate time.Time) (DateRange, error) {
	if startDate.After(endDate) {
		return DateRange{}, shared.NewValidationError("start date cannot be after end date", nil)
	}

	return DateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

// AuditRecord represents an audit record
type AuditRecord struct {
	ID         string                 `json:"id"`
	EntityType string                 `json:"entityType"` // CONFIGURATION, FEATURE_FLAG, USER_PREFERENCE, ORGANIZATION_SETTING
	EntityID   string                 `json:"entityId"`
	Action     string                 `json:"action"`
	UserID     string                 `json:"userId"`
	Timestamp  time.Time              `json:"timestamp"`
	Details    map[string]interface{} `json:"details"`
	IPAddress  *string                `json:"ipAddress,omitempty"`
	UserAgent  *string                `json:"userAgent,omitempty"`
}

// ComplianceCriteria represents criteria for compliance reporting
type ComplianceCriteria struct {
	DateRange       *DateRange       `json:"dateRange,omitempty"`
	OrganizationIDs []OrganizationID `json:"organizationIds,omitempty"`
	ServiceNames    []ServiceName    `json:"serviceNames,omitempty"`
	ComplianceTypes []string         `json:"complianceTypes,omitempty"`
}

// ComplianceReport represents a compliance report
type ComplianceReport struct {
	ID              string                 `json:"id"`
	Title           string                 `json:"title"`
	GeneratedAt     time.Time              `json:"generatedAt"`
	GeneratedBy     string                 `json:"generatedBy"`
	Criteria        ComplianceCriteria     `json:"criteria"`
	Summary         map[string]interface{} `json:"summary"`
	Details         map[string]interface{} `json:"details"`
	Recommendations []string               `json:"recommendations"`
}

// ComplianceValidation represents a compliance validation result
type ComplianceValidation struct {
	OrganizationID OrganizationID    `json:"organizationId"`
	IsCompliant    bool              `json:"isCompliant"`
	Issues         []ComplianceIssue `json:"issues"`
	ValidatedAt    time.Time         `json:"validatedAt"`
	ValidatedBy    string            `json:"validatedBy"`
}

// ComplianceIssue represents a compliance issue
type ComplianceIssue struct {
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"` // HIGH, MEDIUM, LOW
	Description string                 `json:"description"`
	Details     map[string]interface{} `json:"details"`
	Resolution  *string                `json:"resolution,omitempty"`
}
