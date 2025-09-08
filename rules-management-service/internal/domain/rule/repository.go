package rule

import "context"

// ListOptions represents options for listing rules
type ListOptions struct {
	Page      int
	Limit     int
	SortBy    string
	SortOrder string
	Filters   ListFilters
}

// ListFilters represents filters for listing rules
type ListFilters struct {
	Status   string
	Category string
	Search   string
}

// Repository defines the contract for rule persistence
type Repository interface {
	Save(ctx context.Context, rule *Rule) error
	FindByID(ctx context.Context, id RuleID) (*Rule, error)
	FindByName(ctx context.Context, name string) (*Rule, error)
	List(ctx context.Context, options ListOptions) ([]Rule, error)
	Count(ctx context.Context, filters ListFilters) (int, error)
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
