package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/domain/shared"
)

// StructValidator defines the interface for struct validation
type StructValidator interface {
	Validate(s interface{}) error
}

// Validator implements the StructValidator interface using go-playground/validator
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	v := validator.New()

	// Register custom validators
	v.RegisterValidation("email", validateEmail)
	v.RegisterValidation("uuid", validateUUID)
	v.RegisterValidation("age", validateAge)
	v.RegisterValidation("gender", validateGender)
	v.RegisterValidation("status", validateStatus)

	// Register custom tag name function
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{
		validator: v,
	}
}

// Validate validates a struct
func (v *Validator) Validate(s interface{}) error {
	if err := v.validator.Struct(s); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return v.formatValidationErrors(validationErrors)
	}
	return nil
}

// formatValidationErrors formats validation errors into a user-friendly format
func (v *Validator) formatValidationErrors(errs validator.ValidationErrors) error {
	var messages []string

	for _, err := range errs {
		field := err.Field()
		tag := err.Tag()
		param := err.Param()

		var message string
		switch tag {
		case "required":
			message = fmt.Sprintf("%s is required", field)
		case "min":
			message = fmt.Sprintf("%s must be at least %s", field, param)
		case "max":
			message = fmt.Sprintf("%s must be at most %s", field, param)
		case "email":
			message = fmt.Sprintf("%s must be a valid email address", field)
		case "uuid":
			message = fmt.Sprintf("%s must be a valid UUID", field)
		case "age":
			message = fmt.Sprintf("%s must be between 0 and 150", field)
		case "gender":
			message = fmt.Sprintf("%s must be one of: MALE, FEMALE, OTHER, UNKNOWN", field)
		case "status":
			message = fmt.Sprintf("%s must be one of: ACTIVE, INACTIVE, SUSPENDED, DELETED", field)
		case "oneof":
			message = fmt.Sprintf("%s must be one of: %s", field, param)
		case "omitempty":
			// Skip omitempty errors
			continue
		default:
			message = fmt.Sprintf("%s is invalid", field)
		}

		messages = append(messages, message)
	}

	return shared.NewValidationError(strings.Join(messages, "; "), nil)
}

// Custom validators

// validateEmail validates email format
func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	if email == "" {
		return true // Let required validator handle empty values
	}

	// Basic email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// validateUUID validates UUID format
func validateUUID(fl validator.FieldLevel) bool {
	uuid := fl.Field().String()
	if uuid == "" {
		return true // Let required validator handle empty values
	}

	// Basic UUID validation (36 characters with hyphens)
	return len(uuid) == 36 && strings.Count(uuid, "-") == 4
}

// validateAge validates age range
func validateAge(fl validator.FieldLevel) bool {
	age := fl.Field().Int()
	return age >= 0 && age <= 150
}

// validateGender validates gender values
func validateGender(fl validator.FieldLevel) bool {
	gender := fl.Field().String()
	if gender == "" {
		return true // Let required validator handle empty values
	}

	validGenders := []string{"MALE", "FEMALE", "OTHER", "UNKNOWN"}
	for _, valid := range validGenders {
		if gender == valid {
			return true
		}
	}
	return false
}

// validateStatus validates status values
func validateStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return true // Let required validator handle empty values
	}

	validStatuses := []string{"ACTIVE", "INACTIVE", "SUSPENDED", "DELETED"}
	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}
