package settings

import (
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// PreferenceType represents the type of user preference
type PreferenceType string

const (
	PreferenceTypeString  PreferenceType = "STRING"
	PreferenceTypeNumber  PreferenceType = "NUMBER"
	PreferenceTypeBoolean PreferenceType = "BOOLEAN"
	PreferenceTypeJSON    PreferenceType = "JSON"
	PreferenceTypeArray   PreferenceType = "ARRAY"
)

// String returns the string representation of the preference type
func (t PreferenceType) String() string {
	return string(t)
}

// ParsePreferenceType parses a string to PreferenceType
func ParsePreferenceType(prefType string) (PreferenceType, error) {
	switch prefType {
	case "STRING":
		return PreferenceTypeString, nil
	case "NUMBER":
		return PreferenceTypeNumber, nil
	case "BOOLEAN":
		return PreferenceTypeBoolean, nil
	case "JSON":
		return PreferenceTypeJSON, nil
	case "ARRAY":
		return PreferenceTypeArray, nil
	default:
		return "", shared.NewValidationError("invalid preference type", nil)
	}
}

// UserPreference represents a user preference aggregate
type UserPreference struct {
	id             UserPreferenceID
	userID         UserID
	organizationID *OrganizationID
	category       string
	key            string
	value          interface{}
	prefType       PreferenceType
	isDefault      bool
	createdAt      time.Time
	updatedAt      time.Time
	version        int
	events         []shared.DomainEvent
}

// NewUserPreference creates a new user preference
func NewUserPreference(
	userID UserID,
	organizationID *OrganizationID,
	category string,
	key string,
	value interface{},
	prefType PreferenceType,
	isDefault bool,
) (*UserPreference, error) {
	if userID.IsEmpty() {
		return nil, shared.NewValidationError("user ID is required", nil)
	}

	if category == "" {
		return nil, shared.NewValidationError("preference category is required", nil)
	}

	if key == "" {
		return nil, shared.NewValidationError("preference key is required", nil)
	}

	// Validate value type
	if err := validatePreferenceValueType(value, prefType); err != nil {
		return nil, err
	}

	now := time.Now()

	preference := &UserPreference{
		id:             shared.NewUserPreferenceID(),
		userID:         userID,
		organizationID: organizationID,
		category:       category,
		key:            key,
		value:          value,
		prefType:       prefType,
		isDefault:      isDefault,
		createdAt:      now,
		updatedAt:      now,
		version:        1,
		events:         []shared.DomainEvent{},
	}

	// Add user preference created event
	preference.addEvent(NewUserPreferenceCreatedEvent(preference))

	return preference, nil
}

// GetID returns the user preference ID
func (p *UserPreference) GetID() UserPreferenceID {
	return p.id
}

// GetUserID returns the user ID
func (p *UserPreference) GetUserID() UserID {
	return p.userID
}

// GetOrganizationID returns the organization ID
func (p *UserPreference) GetOrganizationID() *OrganizationID {
	return p.organizationID
}

// GetCategory returns the preference category
func (p *UserPreference) GetCategory() string {
	return p.category
}

// GetKey returns the preference key
func (p *UserPreference) GetKey() string {
	return p.key
}

// GetValue returns the preference value
func (p *UserPreference) GetValue() interface{} {
	return p.value
}

// GetType returns the preference type
func (p *UserPreference) GetType() PreferenceType {
	return p.prefType
}

// IsDefault returns whether this is a default preference
func (p *UserPreference) IsDefault() bool {
	return p.isDefault
}

// GetCreatedAt returns the creation timestamp
func (p *UserPreference) GetCreatedAt() time.Time {
	return p.createdAt
}

// GetUpdatedAt returns the last update timestamp
func (p *UserPreference) GetUpdatedAt() time.Time {
	return p.updatedAt
}

// GetVersion returns the preference version
func (p *UserPreference) GetVersion() int {
	return p.version
}

// GetEvents returns the domain events
func (p *UserPreference) GetEvents() []shared.DomainEvent {
	return p.events
}

// ClearEvents clears the domain events
func (p *UserPreference) ClearEvents() {
	p.events = []shared.DomainEvent{}
}

// UpdateValue updates the preference value
func (p *UserPreference) UpdateValue(value interface{}) error {
	// Validate value type
	if err := validatePreferenceValueType(value, p.prefType); err != nil {
		return err
	}

	p.value = value
	p.updatedAt = time.Now()
	p.version++

	p.addEvent(NewUserPreferenceUpdatedEvent(p))

	return nil
}

// MarkAsDefault marks the preference as default
func (p *UserPreference) MarkAsDefault() error {
	if p.isDefault {
		return shared.NewValidationError("preference is already marked as default", nil)
	}

	p.isDefault = true
	p.updatedAt = time.Now()
	p.version++

	p.addEvent(NewUserPreferenceUpdatedEvent(p))

	return nil
}

// UnmarkAsDefault unmarks the preference as default
func (p *UserPreference) UnmarkAsDefault() error {
	if !p.isDefault {
		return shared.NewValidationError("preference is not marked as default", nil)
	}

	p.isDefault = false
	p.updatedAt = time.Now()
	p.version++

	p.addEvent(NewUserPreferenceUpdatedEvent(p))

	return nil
}

// Delete deletes the user preference
func (p *UserPreference) Delete() error {
	p.updatedAt = time.Now()
	p.version++

	p.addEvent(NewUserPreferenceDeletedEvent(p))

	return nil
}

// addEvent adds a domain event to the user preference
func (p *UserPreference) addEvent(event shared.DomainEvent) {
	p.events = append(p.events, event)
}

// validatePreferenceValueType validates that the value matches the preference type
func validatePreferenceValueType(value interface{}, prefType PreferenceType) error {
	switch prefType {
	case PreferenceTypeString:
		if _, ok := value.(string); !ok {
			return shared.NewValidationError("value must be a string", nil)
		}
	case PreferenceTypeNumber:
		switch value.(type) {
		case int, int8, int16, int32, int64, float32, float64:
			// Valid numeric types
		default:
			return shared.NewValidationError("value must be a number", nil)
		}
	case PreferenceTypeBoolean:
		if _, ok := value.(bool); !ok {
			return shared.NewValidationError("value must be a boolean", nil)
		}
	case PreferenceTypeJSON:
		// JSON type can be any interface{}
		// Additional validation can be added here if needed
	case PreferenceTypeArray:
		// Array type should be a slice or array
		// Additional validation can be added here if needed
	default:
		return shared.NewValidationError("unknown preference type", nil)
	}

	return nil
}
