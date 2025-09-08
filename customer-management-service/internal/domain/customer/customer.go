package customer

import (
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/domain/shared"
)

// CustomerID represents a customer identifier
type CustomerID = shared.CustomerID

// EmailAddress represents an email address
type EmailAddress = shared.EmailAddress

// UserID represents a user identifier
type UserID = shared.UserID

// Money represents a monetary value
type Money = shared.Money

// CustomerStatus represents the status of a customer
type CustomerStatus string

const (
	CustomerStatusActive    CustomerStatus = "ACTIVE"
	CustomerStatusInactive  CustomerStatus = "INACTIVE"
	CustomerStatusSuspended CustomerStatus = "SUSPENDED"
	CustomerStatusDeleted   CustomerStatus = "DELETED"
)

// String returns the string representation of the customer status
func (s CustomerStatus) String() string {
	return string(s)
}

// ParseCustomerStatus parses a string to CustomerStatus
func ParseCustomerStatus(status string) (CustomerStatus, error) {
	switch status {
	case "ACTIVE":
		return CustomerStatusActive, nil
	case "INACTIVE":
		return CustomerStatusInactive, nil
	case "SUSPENDED":
		return CustomerStatusSuspended, nil
	case "DELETED":
		return CustomerStatusDeleted, nil
	default:
		return "", shared.NewValidationError("invalid customer status", nil)
	}
}

// Gender represents the gender of a customer
type Gender string

const (
	GenderMale    Gender = "MALE"
	GenderFemale  Gender = "FEMALE"
	GenderOther   Gender = "OTHER"
	GenderUnknown Gender = "UNKNOWN"
)

// String returns the string representation of the gender
func (g Gender) String() string {
	return string(g)
}

// ParseGender parses a string to Gender
func ParseGender(gender string) (Gender, error) {
	switch gender {
	case "MALE":
		return GenderMale, nil
	case "FEMALE":
		return GenderFemale, nil
	case "OTHER":
		return GenderOther, nil
	case "UNKNOWN":
		return GenderUnknown, nil
	default:
		return "", shared.NewValidationError("invalid gender", nil)
	}
}

// CustomerLocation represents the location of a customer
type CustomerLocation struct {
	Country    string   `json:"country"`
	City       string   `json:"city"`
	Region     string   `json:"region"`
	PostalCode *string  `json:"postalCode,omitempty"`
	Timezone   string   `json:"timezone"`
	Latitude   *float64 `json:"latitude,omitempty"`
	Longitude  *float64 `json:"longitude,omitempty"`
}

// NewCustomerLocation creates a new customer location
func NewCustomerLocation(country, city, region, timezone string, postalCode *string, latitude, longitude *float64) (CustomerLocation, error) {
	if country == "" {
		return CustomerLocation{}, shared.NewValidationError("country is required", nil)
	}

	if city == "" {
		return CustomerLocation{}, shared.NewValidationError("city is required", nil)
	}

	if timezone == "" {
		return CustomerLocation{}, shared.NewValidationError("timezone is required", nil)
	}

	return CustomerLocation{
		Country:    country,
		City:       city,
		Region:     region,
		PostalCode: postalCode,
		Timezone:   timezone,
		Latitude:   latitude,
		Longitude:  longitude,
	}, nil
}

// NotificationSettings represents notification preferences
type NotificationSettings struct {
	EmailNotifications bool `json:"emailNotifications"`
	SMSNotifications   bool `json:"smsNotifications"`
	PushNotifications  bool `json:"pushNotifications"`
	MarketingEmails    bool `json:"marketingEmails"`
	SystemAlerts       bool `json:"systemAlerts"`
}

// PrivacySettings represents privacy preferences
type PrivacySettings struct {
	DataSharing       bool `json:"dataSharing"`
	AnalyticsTracking bool `json:"analyticsTracking"`
	Personalization   bool `json:"personalization"`
	ThirdPartySharing bool `json:"thirdPartySharing"`
}

// CustomerPreferences represents customer preferences
type CustomerPreferences struct {
	Language              string                 `json:"language"`
	Currency              string                 `json:"currency"`
	Timezone              string                 `json:"timezone"`
	NotificationSettings  NotificationSettings   `json:"notificationSettings"`
	PrivacySettings       PrivacySettings        `json:"privacySettings"`
	MarketingConsent      bool                   `json:"marketingConsent"`
	DataProcessingConsent bool                   `json:"dataProcessingConsent"`
	CustomPreferences     map[string]interface{} `json:"customPreferences"`
}

// NewCustomerPreferences creates new customer preferences
func NewCustomerPreferences(language, currency, timezone string, notificationSettings NotificationSettings, privacySettings PrivacySettings, marketingConsent, dataProcessingConsent bool, customPreferences map[string]interface{}) (CustomerPreferences, error) {
	if language == "" {
		return CustomerPreferences{}, shared.NewValidationError("language is required", nil)
	}

	if currency == "" {
		return CustomerPreferences{}, shared.NewValidationError("currency is required", nil)
	}

	if timezone == "" {
		return CustomerPreferences{}, shared.NewValidationError("timezone is required", nil)
	}

	return CustomerPreferences{
		Language:              language,
		Currency:              currency,
		Timezone:              timezone,
		NotificationSettings:  notificationSettings,
		PrivacySettings:       privacySettings,
		MarketingConsent:      marketingConsent,
		DataProcessingConsent: dataProcessingConsent,
		CustomPreferences:     customPreferences,
	}, nil
}

// PurchaseRecord represents a purchase record
type PurchaseRecord struct {
	ID           string    `json:"id"`
	Amount       Money     `json:"amount"`
	Product      string    `json:"product"`
	Category     string    `json:"category"`
	PurchaseDate time.Time `json:"purchaseDate"`
	Channel      string    `json:"channel"`
}

// InteractionRecord represents an interaction record
type InteractionRecord struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Channel   string    `json:"channel"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
	Duration  *int      `json:"duration,omitempty"`
	Outcome   string    `json:"outcome"`
}

// DeviceInfo represents device information
type DeviceInfo struct {
	Type      string    `json:"type"`
	OS        string    `json:"os"`
	Browser   string    `json:"browser"`
	UserAgent string    `json:"userAgent"`
	FirstSeen time.Time `json:"firstSeen"`
	LastSeen  time.Time `json:"lastSeen"`
	IsActive  bool      `json:"isActive"`
}

// CustomerMetadata represents customer metadata
type CustomerMetadata struct {
	Source             string              `json:"source"`
	AcquisitionDate    time.Time           `json:"acquisitionDate"`
	LifetimeValue      Money               `json:"lifetimeValue"`
	PurchaseHistory    []PurchaseRecord    `json:"purchaseHistory"`
	InteractionHistory []InteractionRecord `json:"interactionHistory"`
	DeviceInfo         []DeviceInfo        `json:"deviceInfo"`
	ReferralSource     *string             `json:"referralSource,omitempty"`
	LastLogin          *time.Time          `json:"lastLogin,omitempty"`
	LoginCount         int                 `json:"loginCount"`
}

// NewCustomerMetadata creates new customer metadata
func NewCustomerMetadata(source string, acquisitionDate time.Time, lifetimeValue Money, purchaseHistory []PurchaseRecord, interactionHistory []InteractionRecord, deviceInfo []DeviceInfo, referralSource *string, lastLogin *time.Time, loginCount int) CustomerMetadata {
	return CustomerMetadata{
		Source:             source,
		AcquisitionDate:    acquisitionDate,
		LifetimeValue:      lifetimeValue,
		PurchaseHistory:    purchaseHistory,
		InteractionHistory: interactionHistory,
		DeviceInfo:         deviceInfo,
		ReferralSource:     referralSource,
		LastLogin:          lastLogin,
		LoginCount:         loginCount,
	}
}

// Customer represents a customer aggregate
type Customer struct {
	id           CustomerID
	email        EmailAddress
	name         string
	age          *int
	gender       *Gender
	location     *CustomerLocation
	preferences  CustomerPreferences
	segments     []SegmentID
	tags         []string
	status       CustomerStatus
	createdAt    time.Time
	updatedAt    time.Time
	lastActivity time.Time
	metadata     CustomerMetadata
	version      int
	events       []shared.DomainEvent
}

// NewCustomer creates a new customer
func NewCustomer(email EmailAddress, name string, age *int, gender *Gender, location *CustomerLocation, preferences CustomerPreferences, tags []string, metadata CustomerMetadata) (*Customer, error) {
	if name == "" {
		return nil, shared.NewValidationError("customer name is required", nil)
	}

	if age != nil && (*age < 0 || *age > 150) {
		return nil, shared.NewValidationError("age must be between 0 and 150", nil)
	}

	now := time.Now()

	customer := &Customer{
		id:           shared.NewCustomerID(),
		email:        email,
		name:         name,
		age:          age,
		gender:       gender,
		location:     location,
		preferences:  preferences,
		segments:     []SegmentID{},
		tags:         tags,
		status:       CustomerStatusActive,
		createdAt:    now,
		updatedAt:    now,
		lastActivity: now,
		metadata:     metadata,
		version:      1,
		events:       []shared.DomainEvent{},
	}

	// Add customer created event
	customer.addEvent(NewCustomerCreatedEvent(customer))

	return customer, nil
}

// GetID returns the customer ID
func (c *Customer) GetID() CustomerID {
	return c.id
}

// GetEmail returns the customer email
func (c *Customer) GetEmail() EmailAddress {
	return c.email
}

// GetName returns the customer name
func (c *Customer) GetName() string {
	return c.name
}

// GetAge returns the customer age
func (c *Customer) GetAge() *int {
	return c.age
}

// GetGender returns the customer gender
func (c *Customer) GetGender() *Gender {
	return c.gender
}

// GetLocation returns the customer location
func (c *Customer) GetLocation() *CustomerLocation {
	return c.location
}

// GetPreferences returns the customer preferences
func (c *Customer) GetPreferences() CustomerPreferences {
	return c.preferences
}

// GetSegments returns the customer segments
func (c *Customer) GetSegments() []SegmentID {
	return c.segments
}

// GetTags returns the customer tags
func (c *Customer) GetTags() []string {
	return c.tags
}

// GetStatus returns the customer status
func (c *Customer) GetStatus() CustomerStatus {
	return c.status
}

// GetCreatedAt returns the creation timestamp
func (c *Customer) GetCreatedAt() time.Time {
	return c.createdAt
}

// GetUpdatedAt returns the last update timestamp
func (c *Customer) GetUpdatedAt() time.Time {
	return c.updatedAt
}

// GetLastActivity returns the last activity timestamp
func (c *Customer) GetLastActivity() time.Time {
	return c.lastActivity
}

// GetMetadata returns the customer metadata
func (c *Customer) GetMetadata() CustomerMetadata {
	return c.metadata
}

// GetVersion returns the customer version
func (c *Customer) GetVersion() int {
	return c.version
}

// GetEvents returns the domain events
func (c *Customer) GetEvents() []shared.DomainEvent {
	return c.events
}

// ClearEvents clears the domain events
func (c *Customer) ClearEvents() {
	c.events = []shared.DomainEvent{}
}

// UpdateName updates the customer name
func (c *Customer) UpdateName(name string) error {
	if name == "" {
		return shared.NewValidationError("customer name cannot be empty", nil)
	}

	c.name = name
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewCustomerUpdatedEvent(c))

	return nil
}

// UpdateEmail updates the customer email
func (c *Customer) UpdateEmail(email EmailAddress) error {
	c.email = email
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewCustomerEmailUpdatedEvent(c))

	return nil
}

// UpdateAge updates the customer age
func (c *Customer) UpdateAge(age *int) error {
	if age != nil && (*age < 0 || *age > 150) {
		return shared.NewValidationError("age must be between 0 and 150", nil)
	}

	c.age = age
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewCustomerUpdatedEvent(c))

	return nil
}

// UpdateGender updates the customer gender
func (c *Customer) UpdateGender(gender *Gender) error {
	c.gender = gender
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewCustomerUpdatedEvent(c))

	return nil
}

// UpdateLocation updates the customer location
func (c *Customer) UpdateLocation(location *CustomerLocation) error {
	c.location = location
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewCustomerUpdatedEvent(c))

	return nil
}

// UpdatePreferences updates the customer preferences
func (c *Customer) UpdatePreferences(preferences CustomerPreferences) error {
	c.preferences = preferences
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewCustomerPreferencesUpdatedEvent(c))

	return nil
}

// AddTag adds a tag to the customer
func (c *Customer) AddTag(tag string) error {
	if tag == "" {
		return shared.NewValidationError("tag cannot be empty", nil)
	}

	// Check if tag already exists
	for _, existingTag := range c.tags {
		if existingTag == tag {
			return shared.NewValidationError("tag already exists", nil)
		}
	}

	c.tags = append(c.tags, tag)
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewCustomerUpdatedEvent(c))

	return nil
}

// RemoveTag removes a tag from the customer
func (c *Customer) RemoveTag(tag string) error {
	if tag == "" {
		return shared.NewValidationError("tag cannot be empty", nil)
	}

	for i, existingTag := range c.tags {
		if existingTag == tag {
			c.tags = append(c.tags[:i], c.tags[i+1:]...)
			c.updatedAt = time.Now()
			c.version++

			c.addEvent(NewCustomerUpdatedEvent(c))
			return nil
		}
	}

	return shared.NewValidationError("tag not found", nil)
}

// AddSegment adds a segment to the customer
func (c *Customer) AddSegment(segmentID SegmentID) error {
	// Check if segment already exists
	for _, existingSegment := range c.segments {
		if existingSegment == segmentID {
			return shared.NewValidationError("customer already belongs to this segment", nil)
		}
	}

	c.segments = append(c.segments, segmentID)
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewCustomerSegmentJoinedEvent(c, segmentID))

	return nil
}

// RemoveSegment removes a segment from the customer
func (c *Customer) RemoveSegment(segmentID SegmentID) error {
	for i, existingSegment := range c.segments {
		if existingSegment == segmentID {
			c.segments = append(c.segments[:i], c.segments[i+1:]...)
			c.updatedAt = time.Now()
			c.version++

			c.addEvent(NewCustomerSegmentLeftEvent(c, segmentID))
			return nil
		}
	}

	return shared.NewValidationError("customer does not belong to this segment", nil)
}

// Activate activates the customer
func (c *Customer) Activate() error {
	if c.status == CustomerStatusActive {
		return shared.NewValidationError("customer is already active", nil)
	}

	c.status = CustomerStatusActive
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewCustomerActivatedEvent(c))

	return nil
}

// Deactivate deactivates the customer
func (c *Customer) Deactivate() error {
	if c.status == CustomerStatusInactive {
		return shared.NewValidationError("customer is already inactive", nil)
	}

	c.status = CustomerStatusInactive
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewCustomerDeactivatedEvent(c))

	return nil
}

// Suspend suspends the customer
func (c *Customer) Suspend() error {
	if c.status == CustomerStatusSuspended {
		return shared.NewValidationError("customer is already suspended", nil)
	}

	c.status = CustomerStatusSuspended
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewCustomerSuspendedEvent(c))

	return nil
}

// Delete marks the customer as deleted (GDPR compliance)
func (c *Customer) Delete() error {
	if c.status == CustomerStatusDeleted {
		return shared.NewValidationError("customer is already deleted", nil)
	}

	c.status = CustomerStatusDeleted
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewCustomerDeletedEvent(c))

	return nil
}

// UpdateLastActivity updates the last activity timestamp
func (c *Customer) UpdateLastActivity() {
	c.lastActivity = time.Now()
	c.updatedAt = time.Now()
}

// UpdateLastLogin updates the last login timestamp and increments login count
func (c *Customer) UpdateLastLogin() {
	now := time.Now()
	c.metadata.LastLogin = &now
	c.metadata.LoginCount++
	c.lastActivity = now
	c.updatedAt = now
	c.version++

	c.addEvent(NewCustomerLoggedInEvent(c))
}

// AddPurchase adds a purchase to the customer's history
func (c *Customer) AddPurchase(purchase PurchaseRecord) error {
	c.metadata.PurchaseHistory = append(c.metadata.PurchaseHistory, purchase)

	// Update lifetime value
	newLifetimeValue, err := c.metadata.LifetimeValue.Add(purchase.Amount)
	if err != nil {
		return err
	}
	c.metadata.LifetimeValue = newLifetimeValue

	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewCustomerPurchaseEvent(c, purchase))

	return nil
}

// AddInteraction adds an interaction to the customer's history
func (c *Customer) AddInteraction(interaction InteractionRecord) error {
	c.metadata.InteractionHistory = append(c.metadata.InteractionHistory, interaction)
	c.lastActivity = interaction.Timestamp
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewCustomerInteractionEvent(c, interaction))

	return nil
}

// AddDevice adds device information
func (c *Customer) AddDevice(device DeviceInfo) error {
	// Check if device already exists
	for i, existingDevice := range c.metadata.DeviceInfo {
		if existingDevice.Type == device.Type && existingDevice.OS == device.OS && existingDevice.Browser == device.Browser {
			// Update existing device
			c.metadata.DeviceInfo[i] = device
			c.updatedAt = time.Now()
			c.version++

			c.addEvent(NewCustomerDeviceUpdatedEvent(c, device))
			return nil
		}
	}

	// Add new device
	c.metadata.DeviceInfo = append(c.metadata.DeviceInfo, device)
	c.updatedAt = time.Now()
	c.version++

	c.addEvent(NewCustomerDeviceAddedEvent(c, device))

	return nil
}

// addEvent adds a domain event to the customer
func (c *Customer) addEvent(event shared.DomainEvent) {
	c.events = append(c.events, event)
}
