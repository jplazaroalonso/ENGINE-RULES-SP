package rule

import "context"

// Repository defines the contract for rule persistence
type Repository interface {
	Save(ctx context.Context, rule *Rule) error
	FindByID(ctx context.Context, id RuleID) (*Rule, error)
	FindByName(ctx context.Context, name string) (*Rule, error)
	// List(ctx context.Context, criteria ListCriteria) ([]*Rule, error)
	Delete(ctx context.Context, id RuleID) error
	ExistsByName(ctx context.Context, name string) (bool, error)
}

// TemplateRepository defines the contract for rule template persistence
type TemplateRepository interface {
	Save(ctx context.Context, template *RuleTemplate) error
	FindByID(ctx context.Context, id string) (*RuleTemplate, error)
	// List(ctx context.Context, criteria ListCriteria) ([]*RuleTemplate, error)
}

// VersionRepository defines the contract for rule version persistence
type VersionRepository interface {
	Save(ctx context.Context, version *RuleVersion) error
	FindByRuleID(ctx context.Context, ruleID RuleID) ([]*RuleVersion, error)
}
