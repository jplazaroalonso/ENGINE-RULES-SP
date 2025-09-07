package postgres

import (
	"context"

	"gorm.io/gorm"

	"rules-management-service/internal/domain/rule"
	"rules-management-service/internal/domain/shared"
)

type RuleRepository struct {
	db *gorm.DB
}

func NewRuleRepository(db *gorm.DB) *RuleRepository {
	return &RuleRepository{db: db}
}

func (r *RuleRepository) Save(ctx context.Context, rule *rule.Rule) error {
	// This is a simplified implementation. A full implementation would handle created vs updated records.
	if err := r.db.WithContext(ctx).Save(toDBModel(rule)).Error; err != nil {
		return shared.NewInfrastructureError("failed to save rule", err)
	}
	return nil
}

func (r *RuleRepository) FindByID(ctx context.Context, id rule.RuleID) (*rule.Rule, error) {
	var ruleDB RuleDBModel
	if err := r.db.WithContext(ctx).First(&ruleDB, "id = ?", id.String()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.NewNotFoundError("rule not found", err)
		}
		return nil, shared.NewInfrastructureError("failed to find rule by id", err)
	}
	return toDomainEntity(&ruleDB), nil
}

func (r *RuleRepository) FindByName(ctx context.Context, name string) (*rule.Rule, error) {
	var ruleDB RuleDBModel
	if err := r.db.WithContext(ctx).First(&ruleDB, "name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not an error if not found, service layer decides
		}
		return nil, shared.NewInfrastructureError("failed to find rule by name", err)
	}
	return toDomainEntity(&ruleDB), nil
}

func (r *RuleRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&RuleDBModel{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, shared.NewInfrastructureError("failed to check rule existence by name", err)
	}
	return count > 0, nil
}

func (r *RuleRepository) Delete(ctx context.Context, id rule.RuleID) error {
	if err := r.db.WithContext(ctx).Delete(&RuleDBModel{}, "id = ?", id.String()).Error; err != nil {
		return shared.NewInfrastructureError("failed to delete rule", err)
	}
	return nil
}
