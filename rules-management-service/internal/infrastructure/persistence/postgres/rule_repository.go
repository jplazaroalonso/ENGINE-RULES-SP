package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/domain/rule"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/domain/shared"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/rules-management-service/internal/infrastructure/telemetry"
)

type RuleRepository struct {
	db *gorm.DB
}

func NewRuleRepository(db *gorm.DB) *RuleRepository {
	return &RuleRepository{db: db}
}

func (r *RuleRepository) Save(ctx context.Context, rule *rule.Rule) error {
	defer telemetry.DBQueryDuration.WithLabelValues("Save").Observe(time.Since(time.Now()).Seconds())
	// This is a simplified implementation. A full implementation would handle created vs updated records.
	if err := r.db.WithContext(ctx).Save(toDBModel(rule)).Error; err != nil {
		return shared.NewInfrastructureError("failed to save rule", err)
	}
	return nil
}

func (r *RuleRepository) FindByID(ctx context.Context, id rule.RuleID) (*rule.Rule, error) {
	defer telemetry.DBQueryDuration.WithLabelValues("FindByID").Observe(time.Since(time.Now()).Seconds())
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
	defer telemetry.DBQueryDuration.WithLabelValues("FindByName").Observe(time.Since(time.Now()).Seconds())
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
	defer telemetry.DBQueryDuration.WithLabelValues("ExistsByName").Observe(time.Since(time.Now()).Seconds())
	var count int64
	if err := r.db.WithContext(ctx).Model(&RuleDBModel{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, shared.NewInfrastructureError("failed to check rule existence by name", err)
	}
	return count > 0, nil
}

func (r *RuleRepository) List(ctx context.Context, options rule.ListOptions) ([]rule.Rule, error) {
	defer telemetry.DBQueryDuration.WithLabelValues("List").Observe(time.Since(time.Now()).Seconds())
	
	var rulesDB []RuleDBModel
	query := r.db.WithContext(ctx).Model(&RuleDBModel{})
	
	// Apply filters
	if options.Filters.Status != "" {
		query = query.Where("status = ?", options.Filters.Status)
	}
	if options.Filters.Category != "" {
		query = query.Where("category = ?", options.Filters.Category)
	}
	if options.Filters.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+options.Filters.Search+"%", "%"+options.Filters.Search+"%")
	}
	
	// Apply sorting
	sortOrder := "ASC"
	if options.SortOrder == "desc" {
		sortOrder = "DESC"
	}
	query = query.Order(options.SortBy + " " + sortOrder)
	
	// Apply pagination
	offset := (options.Page - 1) * options.Limit
	query = query.Offset(offset).Limit(options.Limit)
	
	if err := query.Find(&rulesDB).Error; err != nil {
		return nil, shared.NewInfrastructureError("failed to list rules", err)
	}
	
	// Convert to domain entities
	rules := make([]rule.Rule, len(rulesDB))
	for i, ruleDB := range rulesDB {
		rules[i] = *toDomainEntity(&ruleDB)
	}
	
	return rules, nil
}

func (r *RuleRepository) Count(ctx context.Context, filters rule.ListFilters) (int, error) {
	defer telemetry.DBQueryDuration.WithLabelValues("Count").Observe(time.Since(time.Now()).Seconds())
	
	var count int64
	query := r.db.WithContext(ctx).Model(&RuleDBModel{})
	
	// Apply filters
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Category != "" {
		query = query.Where("category = ?", filters.Category)
	}
	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}
	
	if err := query.Count(&count).Error; err != nil {
		return 0, shared.NewInfrastructureError("failed to count rules", err)
	}
	
	return int(count), nil
}

func (r *RuleRepository) Delete(ctx context.Context, id rule.RuleID) error {
	defer telemetry.DBQueryDuration.WithLabelValues("Delete").Observe(time.Since(time.Now()).Seconds())
	if err := r.db.WithContext(ctx).Delete(&RuleDBModel{}, "id = ?", id.String()).Error; err != nil {
		return shared.NewInfrastructureError("failed to delete rule", err)
	}
	return nil
}
