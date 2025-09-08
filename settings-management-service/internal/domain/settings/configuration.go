package settings

import (
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// ConfigurationType represents the type of configuration
type ConfigurationType string

const (
	ConfigurationTypeString  ConfigurationType = "STRING"
	ConfigurationTypeNumber  ConfigurationType = "NUMBER"
	ConfigurationTypeBoolean ConfigurationType = "BOOLEAN"
	ConfigurationTypeJSON    ConfigurationType = "JSON"
	ConfigurationTypeArray   ConfigurationType = "ARRAY"
	ConfigurationTypeObject  ConfigurationType = "OBJECT"
)

// String returns the string representation of the configuration type
func (t ConfigurationType) String() string {
	return string(t)
}

// ParseConfigurationType parses a string to ConfigurationType
func ParseConfigurationType(configType string) (ConfigurationType, error) {
	switch configType {
	case "STRING":
		return ConfigurationTypeString, nil
	case "NUMBER":
		return ConfigurationTypeNumber, nil
	case "BOOLEAN":
		return ConfigurationTypeBoolean, nil
	case "JSON":
		return ConfigurationTypeJSON, nil
	case "ARRAY":
		return ConfigurationTypeArray, nil
	case "OBJECT":
		return ConfigurationTypeObject, nil
	default:
		return "", shared.NewValidationError("invalid configuration type", nil)
	}
}

// Environment represents the environment
type Environment string

const (
	EnvironmentDevelopment Environment = "DEVELOPMENT"
	EnvironmentStaging     Environment = "STAGING"
	EnvironmentProduction  Environment = "PRODUCTION"
	EnvironmentTesting     Environment = "TESTING"
)

// String returns the string representation of the environment
func (e Environment) String() string {
	return string(e)
}

// ParseEnvironment parses a string to Environment
func ParseEnvironment(env string) (Environment, error) {
	switch env {
	case "DEVELOPMENT":
		return EnvironmentDevelopment, nil
	case "STAGING":
		return EnvironmentStaging, nil
	case "PRODUCTION":
		return EnvironmentProduction, nil
	case "TESTING":
		return EnvironmentTesting, nil
	default:
		return "", shared.NewValidationError("invalid environment", nil)
	}
}

// Configuration represents a configuration aggregate
type Configuration struct {
	id              ConfigurationID
	key             string
	value           interface{}
	configType      ConfigurationType
	category        string
	description     string
	environment     Environment
	organizationID  *OrganizationID
	service         *ServiceName
	isEncrypted     bool
	isSensitive     bool
	validationRules ValidationRules
	defaultValue    interface{}
	createdBy       UserID
	createdAt       time.Time
	updatedAt       time.Time
	version         int
	events          []shared.DomainEvent
}

// NewConfiguration creates a new configuration
func NewConfiguration(
	key string,
	value interface{},
	configType ConfigurationType,
	category string,
	description string,
	environment Environment,
	organizationID *OrganizationID,
	service *ServiceName,
	isEncrypted bool,
	isSensitive bool,
	validationRules ValidationRules,
	defaultValue interface{},
	createdBy UserID,
) (*Configuration, error) {
	if key == "" {
		return nil, shared.NewValidationError("configuration key is required", nil)
	}

	if category == "" {
		return nil, shared.NewValidationError("configuration category is required", nil)
	}

	if description == "" {
		return nil, shared.NewValidationError("configuration description is required", nil)
	}

	if createdBy.IsEmpty() {
		return nil, shared.NewValidationError("created by user ID is required", nil)
	}

	// Validate key format (should follow service.category.key pattern)
	if service != nil && !service.IsEmpty() {
		expectedPrefix := service.String() + "."
		if !startsWith(key, expectedPrefix) {
			return nil, shared.NewValidationError("configuration key must start with service name", nil)
		}
	}

	// Validate value type
	if err := validateValueType(value, configType); err != nil {
		return nil, err
	}

	// Validate value against validation rules
	if err := validationRules.Validate(value); err != nil {
		return nil, err
	}

	now := time.Now()

	configuration := &Configuration{
		id:              shared.NewConfigurationID(),
		key:             key,
		value:           value,
		configType:      configType,
		category:        category,
		description:     description,
		environment:     environment,
		organizationID:  organizationID,
		service:         service,
		isEncrypted:     isEncrypted,
		isSensitive:     isSensitive,
		validationRules: validationRules,
		defaultValue:    defaultValue,
		createdBy:       createdBy,
		createdAt:       now,
		updatedAt:       now,
		version:         1,
		events:          []shared.DomainEvent{},
	}

	// Add configuration created event
	configuration.addEvent(NewConfigurationCreatedEvent(configuration))

	return configuration, nil
}

// GetID returns the configuration ID
func (c *Configuration) GetID() ConfigurationID {
	return c.id
}

// GetKey returns the configuration key
func (c *Configuration) GetKey() string {
	return c.key
}

// GetValue returns the configuration value
func (c *Configuration) GetValue() interface{} {
	return c.value
}

// GetType returns the configuration type
func (c *Configuration) GetType() ConfigurationType {
	return c.configType
}

// GetCategory returns the configuration category
func (c *Configuration) GetCategory() string {
	return c.category
}

// GetDescription returns the configuration description
func (c *Configuration) GetDescription() string {
	return c.description
}

// GetEnvironment returns the environment
func (c *Configuration) GetEnvironment() Environment {
	return c.environment
}

// GetOrganizationID returns the organization ID
func (c *Configuration) GetOrganizationID() *OrganizationID {
	return c.organizationID
}

// GetService returns the service name
func (c *Configuration) GetService() *ServiceName {
	return c.service
}

// IsEncrypted returns whether the configuration is encrypted
func (c *Configuration) IsEncrypted() bool {
	return c.isEncrypted
}

// IsSensitive returns whether the configuration is sensitive
func (c *Configuration) IsSensitive() bool {
	return c.isSensitive
}

// GetValidationRules returns the validation rules
func (c *Configuration) GetValidationRules() ValidationRules {
	return c.validationRules
}

// GetDefaultValue returns the default value
func (c *Configuration) GetDefaultValue() interface{} {
	return c.defaultValue
}

// GetCreatedBy returns the user who created the configuration
func (c *Configuration) GetCreatedBy() UserID {
	return c.createdBy
}

// GetCreatedAt returns the creation timestamp
func (c *Configuration) GetCreatedAt() time.Time {
	return c.createdAt
}

// GetUpdatedAt returns the last update timestamp
func (c *Configuration) GetUpdatedAt() time.Time {
	return c.updatedAt
}

// GetVersion returns the configuration version
func (c *Configuration) GetVersion() int {
	return c.version
}

// GetEvents returns the domain events
func (c *Configuration) GetEvents() []shared.DomainEvent {
	return c.events
}

// ClearEvents clears the domain events
func (c *Configuration) ClearEvents() {
	c.events = []shared.DomainEvent{}
}

// UpdateValue updates the configuration value
func (c *Configuration) UpdateValue(value interface{}, updatedBy UserID) error {
	if updatedBy.IsEmpty() {
		return shared.NewValidationError("updated by user ID is required", nil)
	}

	// Validate value type
	if err := validateValueType(value, c.configType); err != nil {
		return err
	}

	// Validate value against validation rules
	if err := c.validationRules.Validate(value); err != nil {
		return err
	}

	c.value = value
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewConfigurationUpdatedEvent(c))

	return nil
}

// UpdateDescription updates the configuration description
func (c *Configuration) UpdateDescription(description string, updatedBy UserID) error {
	if description == "" {
		return shared.NewValidationError("configuration description cannot be empty", nil)
	}

	if updatedBy.IsEmpty() {
		return shared.NewValidationError("updated by user ID is required", nil)
	}

	c.description = description
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewConfigurationUpdatedEvent(c))

	return nil
}

// UpdateValidationRules updates the validation rules
func (c *Configuration) UpdateValidationRules(validationRules ValidationRules, updatedBy UserID) error {
	if updatedBy.IsEmpty() {
		return shared.NewValidationError("updated by user ID is required", nil)
	}

	// Validate current value against new rules
	if err := validationRules.Validate(c.value); err != nil {
		return err
	}

	c.validationRules = validationRules
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewConfigurationUpdatedEvent(c))

	return nil
}

// MarkAsSensitive marks the configuration as sensitive
func (c *Configuration) MarkAsSensitive(updatedBy UserID) error {
	if updatedBy.IsEmpty() {
		return shared.NewValidationError("updated by user ID is required", nil)
	}

	c.isSensitive = true
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewConfigurationUpdatedEvent(c))

	return nil
}

// MarkAsEncrypted marks the configuration as encrypted
func (c *Configuration) MarkAsEncrypted(updatedBy UserID) error {
	if updatedBy.IsEmpty() {
		return shared.NewValidationError("updated by user ID is required", nil)
	}

	c.isEncrypted = true
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewConfigurationUpdatedEvent(c))

	return nil
}

// Delete deletes the configuration
func (c *Configuration) Delete(deletedBy UserID) error {
	if deletedBy.IsEmpty() {
		return shared.NewValidationError("deleted by user ID is required", nil)
	}

	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewConfigurationDeletedEvent(c))

	return nil
}

// addEvent adds a domain event to the configuration
func (c *Configuration) addEvent(event shared.DomainEvent) {
	c.events = append(c.events, event)
}

// validateValueType validates that the value matches the configuration type
func validateValueType(value interface{}, configType ConfigurationType) error {
	switch configType {
	case ConfigurationTypeString:
		if _, ok := value.(string); !ok {
			return shared.NewValidationError("value must be a string", nil)
		}
	case ConfigurationTypeNumber:
		switch value.(type) {
		case int, int8, int16, int32, int64, float32, float64:
			// Valid numeric types
		default:
			return shared.NewValidationError("value must be a number", nil)
		}
	case ConfigurationTypeBoolean:
		if _, ok := value.(bool); !ok {
			return shared.NewValidationError("value must be a boolean", nil)
		}
	case ConfigurationTypeJSON, ConfigurationTypeObject:
		// JSON and Object types can be any interface{}
		// Additional validation can be added here if needed
	case ConfigurationTypeArray:
		// Array type should be a slice or array
		// Additional validation can be added here if needed
	default:
		return shared.NewValidationError("unknown configuration type", nil)
	}

	return nil
}

// startsWith checks if a string starts with a prefix
func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
