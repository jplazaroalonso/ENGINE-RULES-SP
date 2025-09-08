package shared

import "errors"

// Domain errors
var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrNotFound         = errors.New("resource not found")
	ErrUnauthorized     = errors.New("unauthorized access")
	ErrForbidden        = errors.New("forbidden access")
	ErrConflict         = errors.New("resource conflict")
	ErrValidationFailed = errors.New("validation failed")
	ErrInternalError    = errors.New("internal error")
)

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e DomainError) Error() string {
	return e.Message
}

// NewDomainError creates a new domain error
func NewDomainError(code, message, details string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}
