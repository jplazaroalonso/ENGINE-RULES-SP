package shared

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// ConfigurationID represents a configuration identifier
type ConfigurationID struct {
	value string
}

// NewConfigurationID creates a new configuration ID
func NewConfigurationID() ConfigurationID {
	return ConfigurationID{value: uuid.New().String()}
}

// NewConfigurationIDFromString creates a configuration ID from a string
func NewConfigurationIDFromString(id string) (ConfigurationID, error) {
	if id == "" {
		return ConfigurationID{}, NewValidationError("configuration ID cannot be empty", nil)
	}

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return ConfigurationID{}, NewValidationError("invalid configuration ID format", err)
	}

	return ConfigurationID{value: id}, nil
}

// String returns the string representation of the configuration ID
func (id ConfigurationID) String() string {
	return id.value
}

// IsEmpty checks if the configuration ID is empty
func (id ConfigurationID) IsEmpty() bool {
	return id.value == ""
}

// FeatureFlagID represents a feature flag identifier
type FeatureFlagID struct {
	value string
}

// NewFeatureFlagID creates a new feature flag ID
func NewFeatureFlagID() FeatureFlagID {
	return FeatureFlagID{value: uuid.New().String()}
}

// NewFeatureFlagIDFromString creates a feature flag ID from a string
func NewFeatureFlagIDFromString(id string) (FeatureFlagID, error) {
	if id == "" {
		return FeatureFlagID{}, NewValidationError("feature flag ID cannot be empty", nil)
	}

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return FeatureFlagID{}, NewValidationError("invalid feature flag ID format", err)
	}

	return FeatureFlagID{value: id}, nil
}

// String returns the string representation of the feature flag ID
func (id FeatureFlagID) String() string {
	return id.value
}

// IsEmpty checks if the feature flag ID is empty
func (id FeatureFlagID) IsEmpty() bool {
	return id.value == ""
}

// UserPreferenceID represents a user preference identifier
type UserPreferenceID struct {
	value string
}

// NewUserPreferenceID creates a new user preference ID
func NewUserPreferenceID() UserPreferenceID {
	return UserPreferenceID{value: uuid.New().String()}
}

// NewUserPreferenceIDFromString creates a user preference ID from a string
func NewUserPreferenceIDFromString(id string) (UserPreferenceID, error) {
	if id == "" {
		return UserPreferenceID{}, NewValidationError("user preference ID cannot be empty", nil)
	}

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return UserPreferenceID{}, NewValidationError("invalid user preference ID format", err)
	}

	return UserPreferenceID{value: id}, nil
}

// String returns the string representation of the user preference ID
func (id UserPreferenceID) String() string {
	return id.value
}

// IsEmpty checks if the user preference ID is empty
func (id UserPreferenceID) IsEmpty() bool {
	return id.value == ""
}

// OrganizationSettingID represents an organization setting identifier
type OrganizationSettingID struct {
	value string
}

// NewOrganizationSettingID creates a new organization setting ID
func NewOrganizationSettingID() OrganizationSettingID {
	return OrganizationSettingID{value: uuid.New().String()}
}

// NewOrganizationSettingIDFromString creates an organization setting ID from a string
func NewOrganizationSettingIDFromString(id string) (OrganizationSettingID, error) {
	if id == "" {
		return OrganizationSettingID{}, NewValidationError("organization setting ID cannot be empty", nil)
	}

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return OrganizationSettingID{}, NewValidationError("invalid organization setting ID format", err)
	}

	return OrganizationSettingID{value: id}, nil
}

// String returns the string representation of the organization setting ID
func (id OrganizationSettingID) String() string {
	return id.value
}

// IsEmpty checks if the organization setting ID is empty
func (id OrganizationSettingID) IsEmpty() bool {
	return id.value == ""
}

// UserID represents a user identifier
type UserID struct {
	value string
}

// NewUserID creates a new user ID
func NewUserID() UserID {
	return UserID{value: uuid.New().String()}
}

// NewUserIDFromString creates a user ID from a string
func NewUserIDFromString(id string) (UserID, error) {
	if id == "" {
		return UserID{}, NewValidationError("user ID cannot be empty", nil)
	}

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return UserID{}, NewValidationError("invalid user ID format", err)
	}

	return UserID{value: id}, nil
}

// String returns the string representation of the user ID
func (id UserID) String() string {
	return id.value
}

// IsEmpty checks if the user ID is empty
func (id UserID) IsEmpty() bool {
	return id.value == ""
}

// OrganizationID represents an organization identifier
type OrganizationID struct {
	value string
}

// NewOrganizationID creates a new organization ID
func NewOrganizationID() OrganizationID {
	return OrganizationID{value: uuid.New().String()}
}

// NewOrganizationIDFromString creates an organization ID from a string
func NewOrganizationIDFromString(id string) (OrganizationID, error) {
	if id == "" {
		return OrganizationID{}, NewValidationError("organization ID cannot be empty", nil)
	}

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return OrganizationID{}, NewValidationError("invalid organization ID format", err)
	}

	return OrganizationID{value: id}, nil
}

// String returns the string representation of the organization ID
func (id OrganizationID) String() string {
	return id.value
}

// IsEmpty checks if the organization ID is empty
func (id OrganizationID) IsEmpty() bool {
	return id.value == ""
}

// ServiceName represents a service name
type ServiceName struct {
	value string
}

// NewServiceName creates a new service name
func NewServiceName(name string) (ServiceName, error) {
	if name == "" {
		return ServiceName{}, NewValidationError("service name cannot be empty", nil)
	}

	// Validate service name format (alphanumeric with hyphens and underscores)
	matched, err := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	if err != nil || !matched {
		return ServiceName{}, NewValidationError("invalid service name format", nil)
	}

	return ServiceName{value: strings.ToLower(name)}, nil
}

// String returns the string representation of the service name
func (s ServiceName) String() string {
	return s.value
}

// IsEmpty checks if the service name is empty
func (s ServiceName) IsEmpty() bool {
	return s.value == ""
}

// ValidationRules represents validation rules for configurations
type ValidationRules struct {
	MinLength     *int          `json:"minLength,omitempty"`
	MaxLength     *int          `json:"maxLength,omitempty"`
	MinValue      *float64      `json:"minValue,omitempty"`
	MaxValue      *float64      `json:"maxValue,omitempty"`
	Pattern       *string       `json:"pattern,omitempty"`
	AllowedValues []interface{} `json:"allowedValues,omitempty"`
	Required      bool          `json:"required"`
}

// NewValidationRules creates new validation rules
func NewValidationRules(minLength, maxLength *int, minValue, maxValue *float64, pattern *string, allowedValues []interface{}, required bool) ValidationRules {
	return ValidationRules{
		MinLength:     minLength,
		MaxLength:     maxLength,
		MinValue:      minValue,
		MaxValue:      maxValue,
		Pattern:       pattern,
		AllowedValues: allowedValues,
		Required:      required,
	}
}

// Validate validates a value against the rules
func (vr ValidationRules) Validate(value interface{}) error {
	// Convert value to string for length validation
	valueStr := fmt.Sprintf("%v", value)

	// Check min length
	if vr.MinLength != nil && len(valueStr) < *vr.MinLength {
		return NewValidationError(fmt.Sprintf("value length %d is less than minimum %d", len(valueStr), *vr.MinLength), nil)
	}

	// Check max length
	if vr.MaxLength != nil && len(valueStr) > *vr.MaxLength {
		return NewValidationError(fmt.Sprintf("value length %d is greater than maximum %d", len(valueStr), *vr.MaxLength), nil)
	}

	// Check numeric value range
	if vr.MinValue != nil || vr.MaxValue != nil {
		if numValue, ok := value.(float64); ok {
			if vr.MinValue != nil && numValue < *vr.MinValue {
				return NewValidationError(fmt.Sprintf("value %f is less than minimum %f", numValue, *vr.MinValue), nil)
			}
			if vr.MaxValue != nil && numValue > *vr.MaxValue {
				return NewValidationError(fmt.Sprintf("value %f is greater than maximum %f", numValue, *vr.MaxValue), nil)
			}
		}
	}

	// Check pattern
	if vr.Pattern != nil {
		matched, err := regexp.MatchString(*vr.Pattern, valueStr)
		if err != nil || !matched {
			return NewValidationError(fmt.Sprintf("value does not match pattern %s", *vr.Pattern), nil)
		}
	}

	// Check allowed values
	if len(vr.AllowedValues) > 0 {
		found := false
		for _, allowed := range vr.AllowedValues {
			if fmt.Sprintf("%v", allowed) == fmt.Sprintf("%v", value) {
				found = true
				break
			}
		}
		if !found {
			return NewValidationError(fmt.Sprintf("value %v is not in allowed values", value), nil)
		}
	}

	return nil
}
