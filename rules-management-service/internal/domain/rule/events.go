package rule

import (
	"time"
)

// RuleCreatedEvent is published when a new rule is created.
type RuleCreatedEvent struct {
	RuleID    string    `json:"rule_id"`
	Name      string    `json:"name"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

func (e RuleCreatedEvent) EventType() string {
	return "RuleCreated"
}
