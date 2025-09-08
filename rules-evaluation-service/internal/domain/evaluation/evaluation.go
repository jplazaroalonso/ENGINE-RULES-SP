package evaluation

// Evaluation represents the aggregate root for a single rule evaluation.
type Evaluation struct {
}

// Context holds the input data for a rule evaluation.
// It's a map to remain generic, as different rule categories will require different data.
type Context map[string]interface{}

// Result holds the outcome of a rule evaluation.
// It's also a map to accommodate different kinds of results.
type Result map[string]interface{}

// Status represents the status of an evaluation.
type Status string

const (
	StatusPending Status = "PENDING"
	StatusSuccess Status = "SUCCESS"
	StatusFailure Status = "FAILURE"
)
