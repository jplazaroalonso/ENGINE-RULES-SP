package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"rules-management-service/internal/domain/rule"
)

// RuleDBModel is the GORM model for the Rule entity
type RuleDBModel struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"not null;uniqueIndex"`
	Description string
	DSLContent  string `gorm:"type:text"`
	Status      string
	Priority    string
	Version     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   string
	ApprovedBy  *string
	ApprovedAt  *time.Time
	TemplateID  *string
	Category    string
	Tags        pq.StringArray `gorm:"type:text[]"`
}

func (RuleDBModel) TableName() string {
	return "rules"
}

// toDBModel converts a domain Rule entity to a GORM model
func toDBModel(r *rule.Rule) *RuleDBModel {
	var templateID *string
	if r.TemplateID() != nil {
		s := r.TemplateID().String()
		templateID = &s
	}

	return &RuleDBModel{
		ID:          r.ID().String(),
		Name:        r.Name(),
		Description: r.Description(),
		DSLContent:  r.DSLContent(),
		Status:      string(r.Status()),
		Priority:    string(r.Priority()),
		Version:     r.Version(),
		CreatedAt:   r.CreatedAt(),
		UpdatedAt:   r.UpdatedAt(),
		CreatedBy:   r.CreatedBy(),
		ApprovedBy:  r.ApprovedBy(),
		ApprovedAt:  r.ApprovedAt(),
		TemplateID:  templateID,
		Category:    r.Category(),
		Tags:        r.Tags(),
	}
}

// toDomainEntity converts a GORM model to a domain Rule entity
func toDomainEntity(dbm *RuleDBModel) *rule.Rule {
	ruleID, _ := rule.RuleIDFromStr(dbm.ID)

	var templateUUID *uuid.UUID
	if dbm.TemplateID != nil {
		parsed, err := uuid.Parse(*dbm.TemplateID)
		if err == nil {
			templateUUID = &parsed
		}
	}

	return rule.ReconstructRule(
		ruleID,
		dbm.Name,
		dbm.Description,
		dbm.DSLContent,
		rule.Status(dbm.Status),
		rule.Priority(dbm.Priority),
		dbm.Version,
		dbm.CreatedAt,
		dbm.UpdatedAt,
		dbm.CreatedBy,
		dbm.ApprovedBy,
		dbm.ApprovedAt,
		templateUUID,
		dbm.Category,
		dbm.Tags,
	)
}
