package shared

import (
	"fmt"

	"github.com/google/uuid"
)

// UserID represents a user identifier
type UserID struct {
	value uuid.UUID
}

func NewUserID() UserID {
	return UserID{value: uuid.New()}
}

func NewUserIDFromString(s string) (UserID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return UserID{}, NewValidationError("invalid user ID format", err)
	}
	return UserID{value: id}, nil
}

func (u UserID) String() string {
	return u.value.String()
}

func (u UserID) IsEmpty() bool {
	return u.value == uuid.Nil
}

// RuleID represents a rule identifier
type RuleID struct {
	value uuid.UUID
}

func NewRuleID() RuleID {
	return RuleID{value: uuid.New()}
}

func NewRuleIDFromString(s string) (RuleID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return RuleID{}, NewValidationError("invalid rule ID format", err)
	}
	return RuleID{value: id}, nil
}

func (r RuleID) String() string {
	return r.value.String()
}

func (r RuleID) IsEmpty() bool {
	return r.value == uuid.Nil
}

// CustomerID represents a customer identifier
type CustomerID struct {
	value uuid.UUID
}

func NewCustomerID() CustomerID {
	return CustomerID{value: uuid.New()}
}

func NewCustomerIDFromString(s string) (CustomerID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return CustomerID{}, NewValidationError("invalid customer ID format", err)
	}
	return CustomerID{value: id}, nil
}

func (c CustomerID) String() string {
	return c.value.String()
}

func (c CustomerID) IsEmpty() bool {
	return c.value == uuid.Nil
}

// Money represents a monetary value
type Money struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

func NewMoney(amount float64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, NewValidationError("amount cannot be negative", nil)
	}
	if currency == "" {
		return Money{}, NewValidationError("currency cannot be empty", nil)
	}
	return Money{Amount: amount, Currency: currency}, nil
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", m.Amount, m.Currency)
}

func (m Money) IsZero() bool {
	return m.Amount == 0
}

func (m Money) IsPositive() bool {
	return m.Amount > 0
}

func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, NewValidationError("cannot add money with different currencies", nil)
	}
	return NewMoney(m.Amount+other.Amount, m.Currency)
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, NewValidationError("cannot subtract money with different currencies", nil)
	}
	return NewMoney(m.Amount-other.Amount, m.Currency)
}

func (m Money) Multiply(factor float64) (Money, error) {
	return NewMoney(m.Amount*factor, m.Currency)
}

func (m Money) Divide(divisor float64) (Money, error) {
	if divisor == 0 {
		return Money{}, NewValidationError("cannot divide by zero", nil)
	}
	return NewMoney(m.Amount/divisor, m.Currency)
}

func (m Money) Equals(other Money) bool {
	return m.Amount == other.Amount && m.Currency == other.Currency
}

func (m Money) GreaterThan(other Money) (bool, error) {
	if m.Currency != other.Currency {
		return false, NewValidationError("cannot compare money with different currencies", nil)
	}
	return m.Amount > other.Amount, nil
}

func (m Money) LessThan(other Money) (bool, error) {
	if m.Currency != other.Currency {
		return false, NewValidationError("cannot compare money with different currencies", nil)
	}
	return m.Amount < other.Amount, nil
}

// StructValidator defines the interface for struct validation
type StructValidator interface {
	Validate(s interface{}) error
}
