package settings

import (
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/settings-management-service/internal/domain/shared"
)

// SettingType represents the type of organization setting
type SettingType string

const (
	SettingTypeString  SettingType = "STRING"
	SettingTypeNumber  SettingType = "NUMBER"
	SettingTypeBoolean SettingType = "BOOLEAN"
	SettingTypeJSON    SettingType = "JSON"
	SettingTypeArray   SettingType = "ARRAY"
)

// String returns the string representation of the setting type
func (t SettingType) String() string {
	return string(t)
}

// ParseSettingType parses a string to SettingType
func ParseSettingType(settingType string) (SettingType, error) {
	switch settingType {
	case "STRING":
		return SettingTypeString, nil
	case "NUMBER":
		return SettingTypeNumber, nil
	case "BOOLEAN":
		return SettingTypeBoolean, nil
	case "JSON":
		return SettingTypeJSON, nil
	case "ARRAY":
		return SettingTypeArray, nil
	default:
		return "", shared.NewValidationError("invalid setting type", nil)
	}
}

// OrganizationSetting represents an organization setting aggregate
type OrganizationSetting struct {
	id             OrganizationSettingID
	organizationID OrganizationID
	category       string
	key            string
	value          interface{}
	settingType    SettingType
	isInherited    bool
	parentID       *OrganizationID
	createdBy      UserID
	createdAt      time.Time
	updatedAt      time.Time
	version        int
	events         []shared.DomainEvent
}

// NewOrganizationSetting creates a new organization setting
func NewOrganizationSetting(
	organizationID OrganizationID,
	category string,
	key string,
	value interface{},
	settingType SettingType,
	isInherited bool,
	parentID *OrganizationID,
	createdBy UserID,
) (*OrganizationSetting, error) {
	if organizationID.IsEmpty() {
		return nil, shared.NewValidationError("organization ID is required", nil)
	}

	if category == "" {
		return nil, shared.NewValidationError("setting category is required", nil)
	}

	if key == "" {
		return nil, shared.NewValidationError("setting key is required", nil)
	}

	if createdBy.IsEmpty() {
		return nil, shared.NewValidationError("created by user ID is required", nil)
	}

	// Validate value type
	if err := validateSettingValueType(value, settingType); err != nil {
		return nil, err
	}

	// If inherited, parent ID must be provided
	if isInherited && (parentID == nil || parentID.IsEmpty()) {
		return nil, shared.NewValidationError("parent organization ID is required for inherited settings", nil)
	}

	now := time.Now()

	setting := &OrganizationSetting{
		id:             shared.NewOrganizationSettingID(),
		organizationID: organizationID,
		category:       category,
		key:            key,
		value:          value,
		settingType:    settingType,
		isInherited:    isInherited,
		parentID:       parentID,
		createdBy:      createdBy,
		createdAt:      now,
		updatedAt:      now,
		version:        1,
		events:         []shared.DomainEvent{},
	}

	// Add organization setting created event
	setting.addEvent(NewOrganizationSettingCreatedEvent(setting))

	return setting, nil
}

// GetID returns the organization setting ID
func (s *OrganizationSetting) GetID() OrganizationSettingID {
	return s.id
}

// GetOrganizationID returns the organization ID
func (s *OrganizationSetting) GetOrganizationID() OrganizationID {
	return s.organizationID
}

// GetCategory returns the setting category
func (s *OrganizationSetting) GetCategory() string {
	return s.category
}

// GetKey returns the setting key
func (s *OrganizationSetting) GetKey() string {
	return s.key
}

// GetValue returns the setting value
func (s *OrganizationSetting) GetValue() interface{} {
	return s.value
}

// GetType returns the setting type
func (s *OrganizationSetting) GetType() SettingType {
	return s.settingType
}

// IsInherited returns whether this setting is inherited
func (s *OrganizationSetting) IsInherited() bool {
	return s.isInherited
}

// GetParentID returns the parent organization ID
func (s *OrganizationSetting) GetParentID() *OrganizationID {
	return s.parentID
}

// GetCreatedBy returns the user who created the setting
func (s *OrganizationSetting) GetCreatedBy() UserID {
	return s.createdBy
}

// GetCreatedAt returns the creation timestamp
func (s *OrganizationSetting) GetCreatedAt() time.Time {
	return s.createdAt
}

// GetUpdatedAt returns the last update timestamp
func (s *OrganizationSetting) GetUpdatedAt() time.Time {
	return s.updatedAt
}

// GetVersion returns the setting version
func (s *OrganizationSetting) GetVersion() int {
	return s.version
}

// GetEvents returns the domain events
func (s *OrganizationSetting) GetEvents() []shared.DomainEvent {
	return s.events
}

// ClearEvents clears the domain events
func (s *OrganizationSetting) ClearEvents() {
	s.events = []shared.DomainEvent{}
}

// UpdateValue updates the setting value
func (s *OrganizationSetting) UpdateValue(value interface{}, updatedBy UserID) error {
	if updatedBy.IsEmpty() {
		return shared.NewValidationError("updated by user ID is required", nil)
	}

	// Validate value type
	if err := validateSettingValueType(value, s.settingType); err != nil {
		return err
	}

	s.value = value
	s.updatedAt = time.Now()
	s.version++

	s.addEvent(NewOrganizationSettingUpdatedEvent(s))

	return nil
}

// MarkAsInherited marks the setting as inherited
func (s *OrganizationSetting) MarkAsInherited(parentID OrganizationID, updatedBy UserID) error {
	if updatedBy.IsEmpty() {
		return shared.NewValidationError("updated by user ID is required", nil)
	}

	if parentID.IsEmpty() {
		return shared.NewValidationError("parent organization ID is required", nil)
	}

	if s.isInherited {
		return shared.NewValidationError("setting is already marked as inherited", nil)
	}

	s.isInherited = true
	s.parentID = &parentID
	s.updatedAt = time.Now()
	s.version++

	s.addEvent(NewOrganizationSettingUpdatedEvent(s))

	return nil
}

// UnmarkAsInherited unmarks the setting as inherited
func (s *OrganizationSetting) UnmarkAsInherited(updatedBy UserID) error {
	if updatedBy.IsEmpty() {
		return shared.NewValidationError("updated by user ID is required", nil)
	}

	if !s.isInherited {
		return shared.NewValidationError("setting is not marked as inherited", nil)
	}

	s.isInherited = false
	s.parentID = nil
	s.updatedAt = time.Now()
	s.version++

	s.addEvent(NewOrganizationSettingUpdatedEvent(s))

	return nil
}

// Delete deletes the organization setting
func (s *OrganizationSetting) Delete(deletedBy UserID) error {
	if deletedBy.IsEmpty() {
		return shared.NewValidationError("deleted by user ID is required", nil)
	}

	s.updatedAt = time.Now()
	s.version++

	s.addEvent(NewOrganizationSettingDeletedEvent(s))

	return nil
}

// addEvent adds a domain event to the organization setting
func (s *OrganizationSetting) addEvent(event shared.DomainEvent) {
	s.events = append(s.events, event)
}

// validateSettingValueType validates that the value matches the setting type
func validateSettingValueType(value interface{}, settingType SettingType) error {
	switch settingType {
	case SettingTypeString:
		if _, ok := value.(string); !ok {
			return shared.NewValidationError("value must be a string", nil)
		}
	case SettingTypeNumber:
		switch value.(type) {
		case int, int8, int16, int32, int64, float32, float64:
			// Valid numeric types
		default:
			return shared.NewValidationError("value must be a number", nil)
		}
	case SettingTypeBoolean:
		if _, ok := value.(bool); !ok {
			return shared.NewValidationError("value must be a boolean", nil)
		}
	case SettingTypeJSON:
		// JSON type can be any interface{}
		// Additional validation can be added here if needed
	case SettingTypeArray:
		// Array type should be a slice or array
		// Additional validation can be added here if needed
	default:
		return shared.NewValidationError("unknown setting type", nil)
	}

	return nil
}
