package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/domain/shared"
)

// StructValidator implements the shared.Validator interface using go-playground/validator.
type StructValidator struct {
	validator *validator.Validate
}

// NewStructValidator creates a new StructValidator.
func NewStructValidator() *StructValidator {
	return &StructValidator{
		validator: validator.New(),
	}
}

// Validate validates a struct using go-playground/validator.
func (v *StructValidator) Validate(s interface{}) error {
	if err := v.validator.Struct(s); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}

// Ensure StructValidator implements shared.Validator interface.
var _ shared.Validator = (*StructValidator)(nil)
