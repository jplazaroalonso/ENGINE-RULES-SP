package dsl

import (
	"strings"
)

// SimpleValidator provides a basic implementation of the rule.ValidationService.
type SimpleValidator struct{}

// NewSimpleValidator creates a new SimpleValidator.
func NewSimpleValidator() *SimpleValidator {
	return &SimpleValidator{}
}

// Validate performs a basic validation of the DSL content.
func (v *SimpleValidator) Validate(dslContent string) (bool, []string) {
	var errors []string

	if !strings.Contains(dslContent, "IF") {
		errors = append(errors, "DSL is missing 'IF' keyword")
	}

	if !strings.Contains(dslContent, "THEN") {
		errors = append(errors, "DSL is missing 'THEN' keyword")
	}

	// Example of a slightly more complex check
	parts := strings.Split(dslContent, "THEN")
	if len(parts) == 2 {
		condition := strings.TrimSpace(strings.TrimPrefix(parts[0], "IF"))
		if len(condition) == 0 {
			errors = append(errors, "Condition in 'IF' statement cannot be empty")
		}

		action := strings.TrimSpace(parts[1])
		if len(action) == 0 {
			errors = append(errors, "Action in 'THEN' statement cannot be empty")
		}
	} else if len(parts) > 2 {
		errors = append(errors, "DSL contains multiple 'THEN' keywords")
	}

	if len(errors) > 0 {
		return false, errors
	}

	return true, nil
}
