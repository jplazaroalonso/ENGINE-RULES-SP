package shared

// Validator defines the interface for struct validation.
type Validator interface {
	Validate(s interface{}) error
}
