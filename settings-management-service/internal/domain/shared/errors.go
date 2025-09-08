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
	Type      ErrorType              `json:"type"`
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Cause     error                  `json:"-"`
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

// Predefined error instances for common errors
var (
	ErrValidation      = &DomainError{Type: ErrorTypeValidation, Code: "VALIDATION_ERROR", Message: "validation failed"}
	ErrNotFound        = &DomainError{Type: ErrorTypeNotFound, Code: "NOT_FOUND", Message: "resource not found"}
	ErrAlreadyExists   = &DomainError{Type: ErrorTypeConflict, Code: "ALREADY_EXISTS", Message: "resource already exists"}
	ErrInternalService = &DomainError{Type: ErrorTypeInfrastructure, Code: "INTERNAL_SERVICE_ERROR", Message: "internal service error"}
)

// Error type checking functions
func IsValidationError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ErrorTypeValidation
	}
	return false
}

func IsNotFoundError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ErrorTypeNotFound
	}
	return false
}

func IsAlreadyExistsError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ErrorTypeConflict
	}
	return false
}

func IsBusinessError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ErrorTypeBusiness
	}
	return false
}

func IsInfrastructureError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ErrorTypeInfrastructure
	}
	return false
}

// Validator defines the interface for validation
type Validator interface {
	Validate(interface{}) error
}

// StructValidator defines the interface for struct validation
type StructValidator interface {
	Validator
	ValidateStruct(interface{}) error
}
