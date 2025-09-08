package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// StructValidator implements the shared.StructValidator interface
type StructValidator struct {
	validator *validator.Validate
}

// NewStructValidator creates a new struct validator
func NewStructValidator() *StructValidator {
	v := validator.New()

	// Register custom validators
	v.RegisterValidation("uuid", validateUUID)
	v.RegisterValidation("environment", validateEnvironment)
	v.RegisterValidation("service_name", validateServiceName)
	v.RegisterValidation("category", validateCategory)
	v.RegisterValidation("key", validateKey)
	v.RegisterValidation("description", validateDescription)
	v.RegisterValidation("tags", validateTags)

	// Register custom tag name function
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &StructValidator{
		validator: v,
	}
}

// Struct validates a struct
func (v *StructValidator) Struct(s interface{}) error {
	if err := v.validator.Struct(s); err != nil {
		return v.formatValidationError(err)
	}
	return nil
}

// Var validates a single variable
func (v *StructValidator) Var(field interface{}, tag string) error {
	if err := v.validator.Var(field, tag); err != nil {
		return v.formatValidationError(err)
	}
	return nil
}

// formatValidationError formats validation errors into a more readable format
func (v *StructValidator) formatValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errors []string
		for _, e := range validationErrors {
			errors = append(errors, v.formatFieldError(e))
		}
		return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}
	return err
}

// formatFieldError formats a single field validation error
func (v *StructValidator) formatFieldError(e validator.FieldError) string {
	field := e.Field()
	if field == "" {
		field = e.StructField()
	}

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, e.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, e.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	case "environment":
		return fmt.Sprintf("%s must be one of: DEVELOPMENT, STAGING, PRODUCTION", field)
	case "service_name":
		return fmt.Sprintf("%s must be a valid service name (1-100 characters, alphanumeric and hyphens only)", field)
	case "category":
		return fmt.Sprintf("%s must be a valid category name (1-100 characters, alphanumeric and hyphens only)", field)
	case "key":
		return fmt.Sprintf("%s must be a valid key (1-255 characters, alphanumeric and hyphens only)", field)
	case "description":
		return fmt.Sprintf("%s must be a valid description (max 1000 characters)", field)
	case "tags":
		return fmt.Sprintf("%s must be a valid array of tags", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// Custom validators

// validateUUID validates that a string is a valid UUID
func validateUUID(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Let required validator handle empty values
	}

	// Basic UUID format validation (8-4-4-4-12)
	parts := strings.Split(value, "-")
	if len(parts) != 5 {
		return false
	}

	expectedLengths := []int{8, 4, 4, 4, 12}
	for i, part := range parts {
		if len(part) != expectedLengths[i] {
			return false
		}
		// Check if all characters are hexadecimal
		for _, char := range part {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
				return false
			}
		}
	}

	return true
}

// validateEnvironment validates that a string is a valid environment
func validateEnvironment(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	validEnvironments := []string{"DEVELOPMENT", "STAGING", "PRODUCTION"}

	for _, env := range validEnvironments {
		if value == env {
			return true
		}
	}

	return false
}

// validateServiceName validates that a string is a valid service name
func validateServiceName(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Let required validator handle empty values
	}

	// Service name should be 1-100 characters, alphanumeric and hyphens only
	if len(value) < 1 || len(value) > 100 {
		return false
	}

	for _, char := range value {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-') {
			return false
		}
	}

	return true
}

// validateCategory validates that a string is a valid category name
func validateCategory(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Let required validator handle empty values
	}

	// Category should be 1-100 characters, alphanumeric and hyphens only
	if len(value) < 1 || len(value) > 100 {
		return false
	}

	for _, char := range value {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-') {
			return false
		}
	}

	return true
}

// validateKey validates that a string is a valid key
func validateKey(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Let required validator handle empty values
	}

	// Key should be 1-255 characters, alphanumeric and hyphens only
	if len(value) < 1 || len(value) > 255 {
		return false
	}

	for _, char := range value {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-') {
			return false
		}
	}

	return true
}

// validateDescription validates that a string is a valid description
func validateDescription(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Let required validator handle empty values
	}

	// Description should be max 1000 characters
	return len(value) <= 1000
}

// validateTags validates that a slice is a valid array of tags
func validateTags(fl validator.FieldLevel) bool {
	value := fl.Field()
	if value.Kind() != reflect.Slice {
		return false
	}

	// Check if it's a slice of strings
	if value.Type().Elem().Kind() != reflect.String {
		return false
	}

	// Check each tag
	for i := 0; i < value.Len(); i++ {
		tag := value.Index(i).String()
		if len(tag) == 0 || len(tag) > 100 {
			return false
		}

		// Tag should be alphanumeric and hyphens only
		for _, char := range tag {
			if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-') {
				return false
			}
		}
	}

	return true
}
