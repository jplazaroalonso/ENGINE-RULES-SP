package calculation

import (
	"time"

	"github.com/google/uuid"
)

// CalculationID represents the unique identifier for a Calculation.
type CalculationID uuid.UUID

// String returns the string representation of a CalculationID.
func (id CalculationID) String() string {
	return uuid.UUID(id).String()
}

// Result represents the outcome of a calculation.
type Result struct {
	Value     float64            `json:"value"`
	Breakdown map[string]float64 `json:"breakdown"`
}

// Status represents the status of a calculation.
type Status string

const (
	StatusPending   Status = "PENDING"
	StatusCompleted Status = "COMPLETED"
	StatusFailed    Status = "FAILED"
)

// Calculation represents the core aggregate for a calculation process.
type Calculation struct {
	id          CalculationID
	ruleIDs     []string
	context     map[string]interface{}
	result      *Result
	status      Status
	createdAt   time.Time
	completedAt *time.Time
}

// NewCalculation creates a new Calculation.
func NewCalculation(ruleIDs []string, context map[string]interface{}) (*Calculation, error) {
	return &Calculation{
		id:        CalculationID(uuid.New()),
		ruleIDs:   ruleIDs,
		context:   context,
		status:    StatusPending,
		createdAt: time.Now().UTC(),
	}, nil
}

// ID returns the calculation's ID.
func (c *Calculation) ID() CalculationID {
	return c.id
}

// Status returns the calculation's status.
func (c *Calculation) Status() Status {
	return c.status
}

// Result returns the calculation's result.
func (c *Calculation) Result() *Result {
	return c.result
}

// CompletedAt returns the calculation's completion time.
func (c *Calculation) CompletedAt() *time.Time {
	return c.completedAt
}

// Complete marks the calculation as completed.
func (c *Calculation) Complete(result Result) {
	c.status = StatusCompleted
	c.result = &result
	now := time.Now().UTC()
	c.completedAt = &now
}

// Fail marks the calculation as failed.
func (c *Calculation) Fail() {
	c.status = StatusFailed
	now := time.Now().UTC()
	c.completedAt = &now
}
