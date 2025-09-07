package rule

import (
	"time"

	"github.com/google/uuid"
)

// RuleVersion represents a version of a rule
type RuleVersion struct {
	id         uuid.UUID
	ruleID     RuleID
	version    int
	dslContent string
	changeLog  string
	createdAt  time.Time
	createdBy  string
}
