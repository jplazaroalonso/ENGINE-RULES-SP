package shared

import "fmt"

// DomainError represents an error that occurred in the domain logic.
type DomainError struct {
	message string
	cause   error
}

func NewDomainError(message string, cause error) *DomainError {
	return &DomainError{message: message, cause: cause}
}

func (e *DomainError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

// ValidationError represents an error in input validation.
type ValidationError struct {
	message string
	cause   error
}

func NewValidationError(message string, cause error) *ValidationError {
	return &ValidationError{message: message, cause: cause}
}

func (e *ValidationError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

// InfrastructureError represents an error from an infrastructure component.
type InfrastructureError struct {
	message string
	cause   error
}

func NewInfrastructureError(message string, cause error) *InfrastructureError {
	return &InfrastructureError{message: message, cause: cause}
}

func (e *InfrastructureError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

// BusinessError represents a business rule violation.
type BusinessError struct {
	message string
	cause   error
}

func NewBusinessError(message string, cause error) *BusinessError {
	return &BusinessError{message: message, cause: cause}
}

func (e *BusinessError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

// NotFoundError represents a resource not found error.
type NotFoundError struct {
	message string
	cause   error
}

func NewNotFoundError(message string, cause error) *NotFoundError {
	return &NotFoundError{message: message, cause: cause}
}

func (e *NotFoundError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}
