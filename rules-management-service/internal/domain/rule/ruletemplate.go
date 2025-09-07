package rule

import (
	"time"

	"github.com/google/uuid"
)

// RuleTemplate represents a reusable rule template
type RuleTemplate struct {
	id          uuid.UUID
	name        string
	description string
	category    string
	dslTemplate string
	parameters  []TemplateParameter
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
	createdBy   string
}

// TemplateParameter represents a parameter in a rule template
type TemplateParameter struct {
	id           uuid.UUID
	templateID   uuid.UUID
	name         string
	paramType    string
	description  string
	required     bool
	defaultValue *string
}
