package shared

import "fmt"

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
