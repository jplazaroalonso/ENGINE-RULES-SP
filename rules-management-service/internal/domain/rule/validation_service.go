package rule

// ValidationService defines the interface for validating rule DSL.
type ValidationService interface {
	Validate(dslContent string) (bool, []string)
}
