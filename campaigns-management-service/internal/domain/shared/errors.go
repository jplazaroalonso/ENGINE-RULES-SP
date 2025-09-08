package shared

import (
	"fmt"
)

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
	Cause   error
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewDomainError(message string, cause error) *DomainError {
	return &DomainError{
		Code:    "DOMAIN_ERROR",
		Message: message,
		Cause:   cause,
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Code    string
	Message string
	Cause   error
}

func (e *ValidationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewValidationError(message string, cause error) *ValidationError {
	return &ValidationError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Cause:   cause,
	}
}

// BusinessError represents a business rule violation
type BusinessError struct {
	Code    string
	Message string
	Cause   error
}

func (e *BusinessError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewBusinessError(message string, cause error) *BusinessError {
	return &BusinessError{
		Code:    "BUSINESS_ERROR",
		Message: message,
		Cause:   cause,
	}
}

// InfrastructureError represents an infrastructure-related error
type InfrastructureError struct {
	Code    string
	Message string
	Cause   error
}

func (e *InfrastructureError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewInfrastructureError(message string, cause error) *InfrastructureError {
	return &InfrastructureError{
		Code:    "INFRASTRUCTURE_ERROR",
		Message: message,
		Cause:   cause,
	}
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Code    string
	Message string
	Cause   error
}

func (e *NotFoundError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewNotFoundError(message string, cause error) *NotFoundError {
	return &NotFoundError{
		Code:    "NOT_FOUND_ERROR",
		Message: message,
		Cause:   cause,
	}
}

// Common validation functions
func ValidateRequired(value, fieldName string) error {
	if value == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

func ValidatePositive(value float64, fieldName string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive", fieldName)
	}
	return nil
}

func ValidateNonNegative(value int64, fieldName string) error {
	if value < 0 {
		return fmt.Errorf("%s must be non-negative", fieldName)
	}
	return nil
}

func ValidateMaxLength(value string, maxLength int, fieldName string) error {
	if len(value) > maxLength {
		return fmt.Errorf("%s cannot exceed %d characters", fieldName, maxLength)
	}
	return nil
}

func ValidateMinLength(value string, minLength int, fieldName string) error {
	if len(value) < minLength {
		return fmt.Errorf("%s must be at least %d characters", fieldName, minLength)
	}
	return nil
}
