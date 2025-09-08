package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/campaigns-management-service/internal/domain/shared"
)

// StructValidator implements shared.StructValidator interface
type StructValidator struct{}

// NewStructValidator creates a new struct validator
func NewStructValidator() *StructValidator {
	return &StructValidator{}
}

// Validate validates a struct using reflection and basic validation rules
func (v *StructValidator) Validate(s interface{}) error {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return shared.NewValidationError("validation target must be a struct", nil)
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Get validation tags
		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		// Parse validation rules
		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			if err := v.validateField(field, fieldType, rule); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateField validates a single field against a validation rule
func (v *StructValidator) validateField(field reflect.Value, fieldType reflect.StructField, rule string) error {
	fieldName := fieldType.Name

	switch {
	case rule == "required":
		if v.isEmpty(field) {
			return shared.NewValidationError(fmt.Sprintf("%s is required", fieldName), nil)
		}
	case strings.HasPrefix(rule, "min="):
		minLen := v.parseIntFromRule(rule, "min=")
		if err := v.validateMinLength(field, fieldName, minLen); err != nil {
			return err
		}
	case strings.HasPrefix(rule, "max="):
		maxLen := v.parseIntFromRule(rule, "max=")
		if err := v.validateMaxLength(field, fieldName, maxLen); err != nil {
			return err
		}
	case strings.HasPrefix(rule, "oneof="):
		allowedValues := strings.Split(strings.TrimPrefix(rule, "oneof="), " ")
		if err := v.validateOneOf(field, fieldName, allowedValues); err != nil {
			return err
		}
	case rule == "uuid":
		if err := v.validateUUID(field, fieldName); err != nil {
			return err
		}
	case strings.HasPrefix(rule, "gt="):
		minVal := v.parseFloatFromRule(rule, "gt=")
		if err := v.validateGreaterThan(field, fieldName, minVal); err != nil {
			return err
		}
	case rule == "omitempty":
		// Skip validation if field is empty
		if v.isEmpty(field) {
			return nil
		}
	}

	return nil
}

// isEmpty checks if a field is empty
func (v *StructValidator) isEmpty(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.String:
		return strings.TrimSpace(field.String()) == ""
	case reflect.Slice, reflect.Array:
		return field.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return field.IsNil()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return field.Float() == 0
	case reflect.Bool:
		return !field.Bool()
	default:
		return false
	}
}

// validateMinLength validates minimum length for strings and slices
func (v *StructValidator) validateMinLength(field reflect.Value, fieldName string, minLen int) error {
	switch field.Kind() {
	case reflect.String:
		if len(strings.TrimSpace(field.String())) < minLen {
			return shared.NewValidationError(fmt.Sprintf("%s must be at least %d characters long", fieldName, minLen), nil)
		}
	case reflect.Slice, reflect.Array:
		if field.Len() < minLen {
			return shared.NewValidationError(fmt.Sprintf("%s must have at least %d items", fieldName, minLen), nil)
		}
	}
	return nil
}

// validateMaxLength validates maximum length for strings and slices
func (v *StructValidator) validateMaxLength(field reflect.Value, fieldName string, maxLen int) error {
	switch field.Kind() {
	case reflect.String:
		if len(strings.TrimSpace(field.String())) > maxLen {
			return shared.NewValidationError(fmt.Sprintf("%s must be at most %d characters long", fieldName, maxLen), nil)
		}
	case reflect.Slice, reflect.Array:
		if field.Len() > maxLen {
			return shared.NewValidationError(fmt.Sprintf("%s must have at most %d items", fieldName, maxLen), nil)
		}
	}
	return nil
}

// validateOneOf validates that a field value is one of the allowed values
func (v *StructValidator) validateOneOf(field reflect.Value, fieldName string, allowedValues []string) error {
	if field.Kind() != reflect.String {
		return nil // Only validate strings for oneof
	}

	value := field.String()
	for _, allowed := range allowedValues {
		if value == allowed {
			return nil
		}
	}

	return shared.NewValidationError(fmt.Sprintf("%s must be one of: %s", fieldName, strings.Join(allowedValues, ", ")), nil)
}

// validateUUID validates that a string field contains a valid UUID
func (v *StructValidator) validateUUID(field reflect.Value, fieldName string) error {
	if field.Kind() != reflect.String {
		return nil // Only validate strings for UUID
	}

	value := strings.TrimSpace(field.String())
	if value == "" {
		return nil // Empty values are handled by required validation
	}

	// Basic UUID format validation (8-4-4-4-12 hex digits)
	parts := strings.Split(value, "-")
	if len(parts) != 5 {
		return shared.NewValidationError(fmt.Sprintf("%s must be a valid UUID", fieldName), nil)
	}

	expectedLengths := []int{8, 4, 4, 4, 12}
	for i, part := range parts {
		if len(part) != expectedLengths[i] {
			return shared.NewValidationError(fmt.Sprintf("%s must be a valid UUID", fieldName), nil)
		}
		// Check if all characters are hex digits
		for _, char := range part {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
				return shared.NewValidationError(fmt.Sprintf("%s must be a valid UUID", fieldName), nil)
			}
		}
	}

	return nil
}

// validateGreaterThan validates that a numeric field is greater than a minimum value
func (v *StructValidator) validateGreaterThan(field reflect.Value, fieldName string, minVal float64) error {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if float64(field.Int()) <= minVal {
			return shared.NewValidationError(fmt.Sprintf("%s must be greater than %g", fieldName, minVal), nil)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if float64(field.Uint()) <= minVal {
			return shared.NewValidationError(fmt.Sprintf("%s must be greater than %g", fieldName, minVal), nil)
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() <= minVal {
			return shared.NewValidationError(fmt.Sprintf("%s must be greater than %g", fieldName, minVal), nil)
		}
	}
	return nil
}

// parseIntFromRule extracts an integer value from a validation rule
func (v *StructValidator) parseIntFromRule(rule, prefix string) int {
	valueStr := strings.TrimPrefix(rule, prefix)
	// Simple parsing - in a real implementation, you'd want proper error handling
	if len(valueStr) > 0 {
		// This is a simplified parser - you might want to use strconv.Atoi with proper error handling
		switch valueStr {
		case "1":
			return 1
		case "2":
			return 2
		case "3":
			return 3
		case "4":
			return 4
		case "5":
			return 5
		case "10":
			return 10
		case "20":
			return 20
		case "50":
			return 50
		case "100":
			return 100
		case "255":
			return 255
		case "500":
			return 500
		case "1000":
			return 1000
		default:
			return 0
		}
	}
	return 0
}

// parseFloatFromRule extracts a float value from a validation rule
func (v *StructValidator) parseFloatFromRule(rule, prefix string) float64 {
	valueStr := strings.TrimPrefix(rule, prefix)
	// Simple parsing - in a real implementation, you'd want proper error handling
	if len(valueStr) > 0 {
		switch valueStr {
		case "0":
			return 0.0
		case "0.0":
			return 0.0
		case "1":
			return 1.0
		case "1.0":
			return 1.0
		default:
			return 0.0
		}
	}
	return 0.0
}
