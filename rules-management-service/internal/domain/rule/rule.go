package rule

import (
	"fmt"
	"strings"
	"time"

	"rules-management-service/internal/domain/shared"

	"github.com/google/uuid"
)

// Rule represents a business rule aggregate
type Rule struct {
	id          RuleID
	name        string
	description string
	dslContent  string
	status      Status
	priority    Priority
	version     int
	createdAt   time.Time
	updatedAt   time.Time
	createdBy   string
	approvedBy  *string
	approvedAt  *time.Time
	templateID  *uuid.UUID
	category    string
	tags        []string
	events      []shared.DomainEvent
}

// RuleID is a typed identifier for rules
type RuleID struct {
	value uuid.UUID
}

func NewRuleID() RuleID {
	return RuleID{value: uuid.New()}
}

func RuleIDFromStr(id string) (RuleID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return RuleID{}, fmt.Errorf("invalid rule id: %w", err)
	}
	return RuleID{value: parsedUUID}, nil
}

func (id RuleID) String() string {
	return id.value.String()
}

// Status represents rule lifecycle status
type Status string

const (
	StatusDraft       Status = "DRAFT"
	StatusUnderReview Status = "UNDER_REVIEW"
	StatusApproved    Status = "APPROVED"
	StatusActive      Status = "ACTIVE"
	StatusInactive    Status = "INACTIVE"
	StatusDeprecated  Status = "DEPRECATED"
)

// Priority represents rule execution priority
type Priority string

const (
	PriorityLow      Priority = "LOW"
	PriorityMedium   Priority = "MEDIUM"
	PriorityHigh     Priority = "HIGH"
	PriorityCritical Priority = "CRITICAL"
)

// NewRule creates a new rule with validation
func NewRule(name, description, dslContent, createdBy string, priority Priority, category string, tags []string) (*Rule, error) {
	if err := validateRuleName(name); err != nil {
		return nil, shared.NewDomainError("invalid rule name", err)
	}

	if err := validateDSLContent(dslContent); err != nil {
		return nil, shared.NewDomainError("invalid DSL content", err)
	}

	rule := &Rule{
		id:          NewRuleID(),
		name:        name,
		description: description,
		dslContent:  dslContent,
		status:      StatusDraft,
		priority:    priority,
		version:     1,
		createdAt:   time.Now().UTC(),
		updatedAt:   time.Now().UTC(),
		createdBy:   createdBy,
		category:    category,
		tags:        tags,
		events:      make([]shared.DomainEvent, 0),
	}

	// Raise domain event
	// rule.addEvent(NewRuleCreatedEvent(rule.id, rule.name, rule.priority))

	return rule, nil
}

// Getters
func (r *Rule) ID() RuleID                   { return r.id }
func (r *Rule) Name() string                 { return r.name }
func (r *Rule) Description() string          { return r.description }
func (r *Rule) DSLContent() string           { return r.dslContent }
func (r *Rule) Status() Status               { return r.status }
func (r *Rule) Priority() Priority           { return r.priority }
func (r *Rule) Version() int                 { return r.version }
func (r *Rule) CreatedAt() time.Time         { return r.createdAt }
func (r *Rule) UpdatedAt() time.Time         { return r.updatedAt }
func (r *Rule) CreatedBy() string            { return r.createdBy }
func (r *Rule) ApprovedBy() *string          { return r.approvedBy }
func (r *Rule) ApprovedAt() *time.Time       { return r.approvedAt }
func (r *Rule) TemplateID() *uuid.UUID       { return r.templateID }
func (r *Rule) Category() string             { return r.category }
func (r *Rule) Tags() []string               { return r.tags }
func (r *Rule) Events() []shared.DomainEvent { return r.events }

func (r *Rule) ClearEvents() {
	r.events = make([]shared.DomainEvent, 0)
}

func (r *Rule) addEvent(event shared.DomainEvent) {
	r.events = append(r.events, event)
}

// ReconstructRule re-creates a rule from existing data. For repository use.
func ReconstructRule(
	id RuleID,
	name string,
	description string,
	dslContent string,
	status Status,
	priority Priority,
	version int,
	createdAt time.Time,
	updatedAt time.Time,
	createdBy string,
	approvedBy *string,
	approvedAt *time.Time,
	templateID *uuid.UUID,
	category string,
	tags []string,
) *Rule {
	return &Rule{
		id:          id,
		name:        name,
		description: description,
		dslContent:  dslContent,
		status:      status,
		priority:    priority,
		version:     version,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		createdBy:   createdBy,
		approvedBy:  approvedBy,
		approvedAt:  approvedAt,
		templateID:  templateID,
		category:    category,
		tags:        tags,
		events:      make([]shared.DomainEvent, 0), // Rehydrated entities have no new events
	}
}

// Private validation functions
func validateRuleName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("rule name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("rule name cannot exceed 100 characters")
	}
	return nil
}

func validateDSLContent(content string) error {
	if strings.TrimSpace(content) == "" {
		return fmt.Errorf("DSL content cannot be empty")
	}
	// Add more sophisticated DSL syntax validation here
	return nil
}
