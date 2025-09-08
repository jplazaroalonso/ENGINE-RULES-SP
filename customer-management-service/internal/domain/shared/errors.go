package shared

import (
	"fmt"
	"time"
)

// ErrorType represents the type of error
type ErrorType string

const (
	ErrorTypeValidation     ErrorType = "VALIDATION_ERROR"
	ErrorTypeBusiness       ErrorType = "BUSINESS_ERROR"
	ErrorTypeNotFound       ErrorType = "NOT_FOUND_ERROR"
	ErrorTypeConflict       ErrorType = "CONFLICT_ERROR"
	ErrorTypeInfrastructure ErrorType = "INFRASTRUCTURE_ERROR"
	ErrorTypeUnauthorized   ErrorType = "UNAUTHORIZED_ERROR"
	ErrorTypeForbidden      ErrorType = "FORBIDDEN_ERROR"
	ErrorTypeRateLimit      ErrorType = "RATE_LIMIT_ERROR"
)

// DomainError represents a domain-specific error
type DomainError struct {
	Type      ErrorType `json:"type"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Cause     error     `json:"-"`
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying cause error
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// NewValidationError creates a new validation error
func NewValidationError(message string, cause error) *DomainError {
	return &DomainError{
		Type:      ErrorTypeValidation,
		Code:      "VALIDATION_ERROR",
		Message:   message,
		Timestamp: time.Now(),
		Cause:     cause,
	}
}

// NewBusinessError creates a new business error
func NewBusinessError(message string, code string) *DomainError {
	return &DomainError{
		Type:      ErrorTypeBusiness,
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string, cause error) *DomainError {
	return &DomainError{
		Type:      ErrorTypeNotFound,
		Code:      "NOT_FOUND",
		Message:   message,
		Timestamp: time.Now(),
		Cause:     cause,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(message string, code string) *DomainError {
	return &DomainError{
		Type:      ErrorTypeConflict,
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// NewInfrastructureError creates a new infrastructure error
func NewInfrastructureError(message string, cause error) *DomainError {
	return &DomainError{
		Type:      ErrorTypeInfrastructure,
		Code:      "INFRASTRUCTURE_ERROR",
		Message:   message,
		Timestamp: time.Now(),
		Cause:     cause,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) *DomainError {
	return &DomainError{
		Type:      ErrorTypeUnauthorized,
		Code:      "UNAUTHORIZED",
		Message:   message,
		Timestamp: time.Now(),
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string) *DomainError {
	return &DomainError{
		Type:      ErrorTypeForbidden,
		Code:      "FORBIDDEN",
		Message:   message,
		Timestamp: time.Now(),
	}
}

// NewRateLimitError creates a new rate limit error
func NewRateLimitError(message string) *DomainError {
	return &DomainError{
		Type:      ErrorTypeRateLimit,
		Code:      "RATE_LIMIT_EXCEEDED",
		Message:   message,
		Timestamp: time.Now(),
	}
}

// WithDetails adds details to the error
func (e *DomainError) WithDetails(details map[string]interface{}) *DomainError {
	e.Details = details
	return e
}

// WithCode sets the error code
func (e *DomainError) WithCode(code string) *DomainError {
	e.Code = code
	return e
}
