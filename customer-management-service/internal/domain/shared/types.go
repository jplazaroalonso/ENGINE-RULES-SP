package shared

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// CustomerID represents a customer identifier
type CustomerID struct {
	value string
}

// NewCustomerID creates a new customer ID
func NewCustomerID() CustomerID {
	return CustomerID{value: uuid.New().String()}
}

// NewCustomerIDFromString creates a customer ID from a string
func NewCustomerIDFromString(id string) (CustomerID, error) {
	if id == "" {
		return CustomerID{}, NewValidationError("customer ID cannot be empty", nil)
	}
	
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return CustomerID{}, NewValidationError("invalid customer ID format", err)
	}
	
	return CustomerID{value: id}, nil
}

// String returns the string representation of the customer ID
func (id CustomerID) String() string {
	return id.value
}

// IsEmpty checks if the customer ID is empty
func (id CustomerID) IsEmpty() bool {
	return id.value == ""
}

// SegmentID represents a segment identifier
type SegmentID struct {
	value string
}

// NewSegmentID creates a new segment ID
func NewSegmentID() SegmentID {
	return SegmentID{value: uuid.New().String()}
}

// NewSegmentIDFromString creates a segment ID from a string
func NewSegmentIDFromString(id string) (SegmentID, error) {
	if id == "" {
		return SegmentID{}, NewValidationError("segment ID cannot be empty", nil)
	}
	
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return SegmentID{}, NewValidationError("invalid segment ID format", err)
	}
	
	return SegmentID{value: id}, nil
}

// String returns the string representation of the segment ID
func (id SegmentID) String() string {
	return id.value
}

// IsEmpty checks if the segment ID is empty
func (id SegmentID) IsEmpty() bool {
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

// RuleID represents a rule identifier
type RuleID struct {
	value string
}

// NewRuleID creates a new rule ID
func NewRuleID() RuleID {
	return RuleID{value: uuid.New().String()}
}

// NewRuleIDFromString creates a rule ID from a string
func NewRuleIDFromString(id string) (RuleID, error) {
	if id == "" {
		return RuleID{}, NewValidationError("rule ID cannot be empty", nil)
	}
	
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return RuleID{}, NewValidationError("invalid rule ID format", err)
	}
	
	return RuleID{value: id}, nil
}

// String returns the string representation of the rule ID
func (id RuleID) String() string {
	return id.value
}

// IsEmpty checks if the rule ID is empty
func (id RuleID) IsEmpty() bool {
	return id.value == ""
}

// Money represents a monetary value
type Money struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// NewMoney creates a new money value
func NewMoney(amount float64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, NewValidationError("money amount cannot be negative", nil)
	}
	
	if currency == "" {
		return Money{}, NewValidationError("currency cannot be empty", nil)
	}
	
	// Validate currency format (3-letter ISO code)
	if len(currency) != 3 {
		return Money{}, NewValidationError("currency must be a 3-letter ISO code", nil)
	}
	
	return Money{
		Amount:   amount,
		Currency: strings.ToUpper(currency),
	}, nil
}

// Add adds another money value to this one
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, NewValidationError("cannot add money with different currencies", nil)
	}
	
	return Money{
		Amount:   m.Amount + other.Amount,
		Currency: m.Currency,
	}, nil
}

// Subtract subtracts another money value from this one
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, NewValidationError("cannot subtract money with different currencies", nil)
	}
	
	return Money{
		Amount:   m.Amount - other.Amount,
		Currency: m.Currency,
	}, nil
}

// Multiply multiplies the money value by a factor
func (m Money) Multiply(factor float64) Money {
	return Money{
		Amount:   m.Amount * factor,
		Currency: m.Currency,
	}
}

// IsZero checks if the money amount is zero
func (m Money) IsZero() bool {
	return m.Amount == 0
}

// IsPositive checks if the money amount is positive
func (m Money) IsPositive() bool {
	return m.Amount > 0
}

// IsNegative checks if the money amount is negative
func (m Money) IsNegative() bool {
	return m.Amount < 0
}

// EmailAddress represents an email address value object
type EmailAddress struct {
	value string
}

// NewEmailAddress creates a new email address
func NewEmailAddress(email string) (EmailAddress, error) {
	if email == "" {
		return EmailAddress{}, NewValidationError("email address cannot be empty", nil)
	}
	
	// Basic email validation regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return EmailAddress{}, NewValidationError("invalid email address format", nil)
	}
	
	return EmailAddress{value: strings.ToLower(strings.TrimSpace(email))}, nil
}

// String returns the string representation of the email address
func (e EmailAddress) String() string {
	return e.value
}

// Domain returns the domain part of the email address
func (e EmailAddress) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// LocalPart returns the local part of the email address
func (e EmailAddress) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) == 2 {
		return parts[0]
	}
	return ""
}

// EventID represents an event identifier
type EventID struct {
	value string
}

// NewEventID creates a new event ID
func NewEventID() EventID {
	return EventID{value: uuid.New().String()}
}

// String returns the string representation of the event ID
func (id EventID) String() string {
	return id.value
}

// AgeRange represents an age range value object
type AgeRange struct {
	Min *int `json:"min,omitempty"`
	Max *int `json:"max,omitempty"`
}

// NewAgeRange creates a new age range
func NewAgeRange(min, max *int) (AgeRange, error) {
	if min != nil && *min < 0 {
		return AgeRange{}, NewValidationError("minimum age cannot be negative", nil)
	}
	
	if max != nil && *max < 0 {
		return AgeRange{}, NewValidationError("maximum age cannot be negative", nil)
	}
	
	if min != nil && max != nil && *min > *max {
		return AgeRange{}, NewValidationError("minimum age cannot be greater than maximum age", nil)
	}
	
	return AgeRange{Min: min, Max: max}, nil
}

// Contains checks if an age is within the range
func (ar AgeRange) Contains(age int) bool {
	if ar.Min != nil && age < *ar.Min {
		return false
	}
	
	if ar.Max != nil && age > *ar.Max {
		return false
	}
	
	return true
}

// ValueRange represents a value range
type ValueRange struct {
	Min       *float64 `json:"min,omitempty"`
	Max       *float64 `json:"max,omitempty"`
	Currency  string   `json:"currency,omitempty"`
}

// NewValueRange creates a new value range
func NewValueRange(min, max *float64, currency string) (ValueRange, error) {
	if min != nil && *min < 0 {
		return ValueRange{}, NewValidationError("minimum value cannot be negative", nil)
	}
	
	if max != nil && *max < 0 {
		return ValueRange{}, NewValidationError("maximum value cannot be negative", nil)
	}
	
	if min != nil && max != nil && *min > *max {
		return ValueRange{}, NewValidationError("minimum value cannot be greater than maximum value", nil)
	}
	
	return ValueRange{Min: min, Max: max, Currency: currency}, nil
}

// Contains checks if a value is within the range
func (vr ValueRange) Contains(value float64) bool {
	if vr.Min != nil && value < *vr.Min {
		return false
	}
	
	if vr.Max != nil && value > *vr.Max {
		return false
	}
	
	return true
}

// FrequencyRange represents a frequency range
type FrequencyRange struct {
	Min    *int    `json:"min,omitempty"`
	Max    *int    `json:"max,omitempty"`
	Period string  `json:"period,omitempty"`
}

// NewFrequencyRange creates a new frequency range
func NewFrequencyRange(min, max *int, period string) (FrequencyRange, error) {
	if min != nil && *min < 0 {
		return FrequencyRange{}, NewValidationError("minimum frequency cannot be negative", nil)
	}
	
	if max != nil && *max < 0 {
		return FrequencyRange{}, NewValidationError("maximum frequency cannot be negative", nil)
	}
	
	if min != nil && max != nil && *min > *max {
		return FrequencyRange{}, NewValidationError("minimum frequency cannot be greater than maximum frequency", nil)
	}
	
	validPeriods := []string{"DAILY", "WEEKLY", "MONTHLY", "YEARLY"}
	if period != "" {
		valid := false
		for _, p := range validPeriods {
			if p == period {
				valid = true
				break
			}
		}
		if !valid {
			return FrequencyRange{}, NewValidationError(fmt.Sprintf("invalid period: %s", period), nil)
		}
	}
	
	return FrequencyRange{Min: min, Max: max, Period: period}, nil
}

// DaysRange represents a days range
type DaysRange struct {
	Min *int `json:"min,omitempty"`
	Max *int `json:"max,omitempty"`
}

// NewDaysRange creates a new days range
func NewDaysRange(min, max *int) (DaysRange, error) {
	if min != nil && *min < 0 {
		return DaysRange{}, NewValidationError("minimum days cannot be negative", nil)
	}
	
	if max != nil && *max < 0 {
		return DaysRange{}, NewValidationError("maximum days cannot be negative", nil)
	}
	
	if min != nil && max != nil && *min > *max {
		return DaysRange{}, NewValidationError("minimum days cannot be greater than maximum days", nil)
	}
	
	return DaysRange{Min: min, Max: max}, nil
}

// Contains checks if a number of days is within the range
func (dr DaysRange) Contains(days int) bool {
	if dr.Min != nil && days < *dr.Min {
		return false
	}
	
	if dr.Max != nil && days > *dr.Max {
		return false
	}
	
	return true
}

// LocationRadius represents a location radius
type LocationRadius struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"` // in kilometers
}

// NewLocationRadius creates a new location radius
func NewLocationRadius(latitude, longitude, radius float64) (LocationRadius, error) {
	if latitude < -90 || latitude > 90 {
		return LocationRadius{}, NewValidationError("latitude must be between -90 and 90", nil)
	}
	
	if longitude < -180 || longitude > 180 {
		return LocationRadius{}, NewValidationError("longitude must be between -180 and 180", nil)
	}
	
	if radius < 0 {
		return LocationRadius{}, NewValidationError("radius cannot be negative", nil)
	}
	
	return LocationRadius{
		Latitude:  latitude,
		Longitude: longitude,
		Radius:    radius,
	}, nil
}
