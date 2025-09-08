package customer

import (
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/domain/shared"
)

// SegmentID represents a segment identifier
type SegmentID = shared.SegmentID

// RuleID represents a rule identifier
type RuleID = shared.RuleID

// UserID represents a user identifier
type UserID = shared.UserID

// AgeRange represents an age range
type AgeRange = shared.AgeRange

// ValueRange represents a value range
type ValueRange = shared.ValueRange

// FrequencyRange represents a frequency range
type FrequencyRange = shared.FrequencyRange

// DaysRange represents a days range
type DaysRange = shared.DaysRange

// LocationRadius represents a location radius
type LocationRadius = shared.LocationRadius

// SegmentStatus represents the status of a segment
type SegmentStatus string

const (
	SegmentStatusActive      SegmentStatus = "ACTIVE"
	SegmentStatusInactive    SegmentStatus = "INACTIVE"
	SegmentStatusCalculating SegmentStatus = "CALCULATING"
	SegmentStatusError       SegmentStatus = "ERROR"
)

// String returns the string representation of the segment status
func (s SegmentStatus) String() string {
	return string(s)
}

// ParseSegmentStatus parses a string to SegmentStatus
func ParseSegmentStatus(status string) (SegmentStatus, error) {
	switch status {
	case "ACTIVE":
		return SegmentStatusActive, nil
	case "INACTIVE":
		return SegmentStatusInactive, nil
	case "CALCULATING":
		return SegmentStatusCalculating, nil
	case "ERROR":
		return SegmentStatusError, nil
	default:
		return "", shared.NewValidationError("invalid segment status", nil)
	}
}

// DemographicCriteria represents demographic criteria for segmentation
type DemographicCriteria struct {
	AgeRange    *AgeRange   `json:"ageRange,omitempty"`
	Gender      []Gender    `json:"gender,omitempty"`
	IncomeRange *ValueRange `json:"incomeRange,omitempty"`
	Education   []string    `json:"education,omitempty"`
	Occupation  []string    `json:"occupation,omitempty"`
}

// BehavioralCriteria represents behavioral criteria for segmentation
type BehavioralCriteria struct {
	PurchaseFrequency *FrequencyRange `json:"purchaseFrequency,omitempty"`
	AverageOrderValue *ValueRange     `json:"averageOrderValue,omitempty"`
	LastPurchaseDays  *DaysRange      `json:"lastPurchaseDays,omitempty"`
	ProductCategories []string        `json:"productCategories,omitempty"`
	BrandPreferences  []string        `json:"brandPreferences,omitempty"`
}

// GeographicCriteria represents geographic criteria for segmentation
type GeographicCriteria struct {
	Countries    []string        `json:"countries,omitempty"`
	Cities       []string        `json:"cities,omitempty"`
	Regions      []string        `json:"regions,omitempty"`
	PostalCodes  []string        `json:"postalCodes,omitempty"`
	DistanceFrom *LocationRadius `json:"distanceFrom,omitempty"`
}

// CustomRule represents a custom rule for segmentation
type CustomRule struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// SegmentCriteria represents the criteria for a customer segment
type SegmentCriteria struct {
	Demographics    *DemographicCriteria `json:"demographics,omitempty"`
	Behavioral      *BehavioralCriteria  `json:"behavioral,omitempty"`
	Geographic      *GeographicCriteria  `json:"geographic,omitempty"`
	PurchaseHistory *PurchaseCriteria    `json:"purchaseHistory,omitempty"`
	CustomRules     []CustomRule         `json:"customRules"`
}

// PurchaseCriteria represents purchase-related criteria
type PurchaseCriteria struct {
	MinTotalSpent       *ValueRange `json:"minTotalSpent,omitempty"`
	MaxTotalSpent       *ValueRange `json:"maxTotalSpent,omitempty"`
	MinOrderCount       *int        `json:"minOrderCount,omitempty"`
	MaxOrderCount       *int        `json:"maxOrderCount,omitempty"`
	PreferredCategories []string    `json:"preferredCategories,omitempty"`
	ExcludedCategories  []string    `json:"excludedCategories,omitempty"`
}

// CustomerSegment represents a customer segment aggregate
type CustomerSegment struct {
	id             SegmentID
	name           string
	description    string
	ruleID         RuleID
	customerCount  int
	criteria       SegmentCriteria
	status         SegmentStatus
	createdBy      UserID
	createdAt      time.Time
	updatedAt      time.Time
	lastCalculated *time.Time
	version        int
	events         []shared.DomainEvent
}

// NewCustomerSegment creates a new customer segment
func NewCustomerSegment(name, description string, ruleID RuleID, criteria SegmentCriteria, createdBy UserID) (*CustomerSegment, error) {
	if name == "" {
		return nil, shared.NewValidationError("segment name is required", nil)
	}

	if description == "" {
		return nil, shared.NewValidationError("segment description is required", nil)
	}

	if ruleID.IsEmpty() {
		return nil, shared.NewValidationError("rule ID is required", nil)
	}

	if createdBy.IsEmpty() {
		return nil, shared.NewValidationError("created by user ID is required", nil)
	}

	// Validate that at least one criteria is specified
	if criteria.Demographics == nil && criteria.Behavioral == nil &&
		criteria.Geographic == nil && criteria.PurchaseHistory == nil &&
		len(criteria.CustomRules) == 0 {
		return nil, shared.NewValidationError("at least one criteria must be specified", nil)
	}

	now := time.Now()

	segment := &CustomerSegment{
		id:            shared.NewSegmentID(),
		name:          name,
		description:   description,
		ruleID:        ruleID,
		customerCount: 0,
		criteria:      criteria,
		status:        SegmentStatusActive,
		createdBy:     createdBy,
		createdAt:     now,
		updatedAt:     now,
		version:       1,
		events:        []shared.DomainEvent{},
	}

	// Add segment created event
	segment.addEvent(NewCustomerSegmentCreatedEvent(segment))

	return segment, nil
}

// GetID returns the segment ID
func (s *CustomerSegment) GetID() SegmentID {
	return s.id
}

// GetName returns the segment name
func (s *CustomerSegment) GetName() string {
	return s.name
}

// GetDescription returns the segment description
func (s *CustomerSegment) GetDescription() string {
	return s.description
}

// GetRuleID returns the rule ID
func (s *CustomerSegment) GetRuleID() RuleID {
	return s.ruleID
}

// GetCustomerCount returns the customer count
func (s *CustomerSegment) GetCustomerCount() int {
	return s.customerCount
}

// GetCriteria returns the segment criteria
func (s *CustomerSegment) GetCriteria() SegmentCriteria {
	return s.criteria
}

// GetStatus returns the segment status
func (s *CustomerSegment) GetStatus() SegmentStatus {
	return s.status
}

// GetCreatedBy returns the user who created the segment
func (s *CustomerSegment) GetCreatedBy() UserID {
	return s.createdBy
}

// GetCreatedAt returns the creation timestamp
func (s *CustomerSegment) GetCreatedAt() time.Time {
	return s.createdAt
}

// GetUpdatedAt returns the last update timestamp
func (s *CustomerSegment) GetUpdatedAt() time.Time {
	return s.updatedAt
}

// GetLastCalculated returns the last calculated timestamp
func (s *CustomerSegment) GetLastCalculated() *time.Time {
	return s.lastCalculated
}

// GetVersion returns the segment version
func (s *CustomerSegment) GetVersion() int {
	return s.version
}

// GetEvents returns the domain events
func (s *CustomerSegment) GetEvents() []shared.DomainEvent {
	return s.events
}

// ClearEvents clears the domain events
func (s *CustomerSegment) ClearEvents() {
	s.events = []shared.DomainEvent{}
}

// UpdateName updates the segment name
func (s *CustomerSegment) UpdateName(name string) error {
	if name == "" {
		return shared.NewValidationError("segment name cannot be empty", nil)
	}

	s.name = name
	s.updatedAt = time.Now()
	s.version++

	s.addEvent(NewCustomerSegmentUpdatedEvent(s))

	return nil
}

// UpdateDescription updates the segment description
func (s *CustomerSegment) UpdateDescription(description string) error {
	if description == "" {
		return shared.NewValidationError("segment description cannot be empty", nil)
	}

	s.description = description
	s.updatedAt = time.Now()
	s.version++

	s.addEvent(NewCustomerSegmentUpdatedEvent(s))

	return nil
}

// UpdateCriteria updates the segment criteria
func (s *CustomerSegment) UpdateCriteria(criteria SegmentCriteria) error {
	// Validate that at least one criteria is specified
	if criteria.Demographics == nil && criteria.Behavioral == nil &&
		criteria.Geographic == nil && criteria.PurchaseHistory == nil &&
		len(criteria.CustomRules) == 0 {
		return shared.NewValidationError("at least one criteria must be specified", nil)
	}

	s.criteria = criteria
	s.updatedAt = time.Now()
	s.version++

	s.addEvent(NewCustomerSegmentCriteriaUpdatedEvent(s))

	return nil
}

// Activate activates the segment
func (s *CustomerSegment) Activate() error {
	if s.status == SegmentStatusActive {
		return shared.NewValidationError("segment is already active", nil)
	}

	s.status = SegmentStatusActive
	s.updatedAt = time.Now()
	s.version++

	s.addEvent(NewCustomerSegmentActivatedEvent(s))

	return nil
}

// Deactivate deactivates the segment
func (s *CustomerSegment) Deactivate() error {
	if s.status == SegmentStatusInactive {
		return shared.NewValidationError("segment is already inactive", nil)
	}

	s.status = SegmentStatusInactive
	s.updatedAt = time.Now()
	s.version++

	s.addEvent(NewCustomerSegmentDeactivatedEvent(s))

	return nil
}

// StartCalculation starts the segment calculation process
func (s *CustomerSegment) StartCalculation() error {
	if s.status == SegmentStatusCalculating {
		return shared.NewValidationError("segment calculation is already in progress", nil)
	}

	s.status = SegmentStatusCalculating
	s.updatedAt = time.Now()
	s.version++

	s.addEvent(NewCustomerSegmentCalculationStartedEvent(s))

	return nil
}

// CompleteCalculation completes the segment calculation process
func (s *CustomerSegment) CompleteCalculation(customerCount int) error {
	if s.status != SegmentStatusCalculating {
		return shared.NewValidationError("segment is not in calculating status", nil)
	}

	s.status = SegmentStatusActive
	s.customerCount = customerCount
	now := time.Now()
	s.lastCalculated = &now
	s.updatedAt = now
	s.version++

	s.addEvent(NewCustomerSegmentCalculationCompletedEvent(s, customerCount))

	return nil
}

// FailCalculation marks the segment calculation as failed
func (s *CustomerSegment) FailCalculation(errorMessage string) error {
	if s.status != SegmentStatusCalculating {
		return shared.NewValidationError("segment is not in calculating status", nil)
	}

	s.status = SegmentStatusError
	s.updatedAt = time.Now()
	s.version++

	s.addEvent(NewCustomerSegmentCalculationFailedEvent(s, errorMessage))

	return nil
}

// UpdateCustomerCount updates the customer count
func (s *CustomerSegment) UpdateCustomerCount(customerCount int) error {
	if customerCount < 0 {
		return shared.NewValidationError("customer count cannot be negative", nil)
	}

	s.customerCount = customerCount
	s.updatedAt = time.Now()
	s.version++

	s.addEvent(NewCustomerSegmentUpdatedEvent(s))

	return nil
}

// Delete deletes the segment
func (s *CustomerSegment) Delete() error {
	if s.status == SegmentStatusCalculating {
		return shared.NewValidationError("cannot delete segment while calculation is in progress", nil)
	}

	s.status = SegmentStatusInactive
	s.updatedAt = time.Now()
	s.version++

	s.addEvent(NewCustomerSegmentDeletedEvent(s))

	return nil
}

// addEvent adds a domain event to the segment
func (s *CustomerSegment) addEvent(event shared.DomainEvent) {
	s.events = append(s.events, event)
}
