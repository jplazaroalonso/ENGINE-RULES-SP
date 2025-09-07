package evaluation

// EvaluationStrategy defines the interface for different rule evaluation algorithms.
type EvaluationStrategy interface {
	// Evaluate applies the rule logic to the given context.
	Evaluate(dslContent string, context Context) (Result, error)
}
