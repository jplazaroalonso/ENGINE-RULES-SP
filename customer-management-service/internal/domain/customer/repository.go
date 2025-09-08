package customer

import (
	"context"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/domain/shared"
)

// CustomerRepository defines the interface for customer persistence
type CustomerRepository interface {
	// Basic CRUD operations
	Save(ctx context.Context, customer *Customer) error
	FindByID(ctx context.Context, id CustomerID) (*Customer, error)
	FindByEmail(ctx context.Context, email shared.EmailAddress) (*Customer, error)
	Update(ctx context.Context, customer *Customer) error
	Delete(ctx context.Context, id CustomerID) error

	// List operations
	List(ctx context.Context, criteria ListCriteria) ([]*Customer, error)
	Count(ctx context.Context, criteria ListCriteria) (int64, error)

	// Search operations
	Search(ctx context.Context, query string, criteria ListCriteria) ([]*Customer, error)
	FindByStatus(ctx context.Context, status CustomerStatus) ([]*Customer, error)
	FindBySegment(ctx context.Context, segmentID SegmentID) ([]*Customer, error)
	FindByTags(ctx context.Context, tags []string) ([]*Customer, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*Customer, error)

	// Existence checks
	ExistsByEmail(ctx context.Context, email shared.EmailAddress) (bool, error)
	ExistsByID(ctx context.Context, id CustomerID) (bool, error)

	// Bulk operations
	BulkUpdate(ctx context.Context, updates []BulkUpdate) error
	BulkDelete(ctx context.Context, ids []CustomerID) error
}

// CustomerSegmentRepository defines the interface for customer segment persistence
type CustomerSegmentRepository interface {
	// Basic CRUD operations
	Save(ctx context.Context, segment *CustomerSegment) error
	FindByID(ctx context.Context, id SegmentID) (*CustomerSegment, error)
	FindByName(ctx context.Context, name string) (*CustomerSegment, error)
	Update(ctx context.Context, segment *CustomerSegment) error
	Delete(ctx context.Context, id SegmentID) error

	// List operations
	List(ctx context.Context, criteria SegmentListCriteria) ([]*CustomerSegment, error)
	Count(ctx context.Context, criteria SegmentListCriteria) (int64, error)

	// Search operations
	Search(ctx context.Context, query string, criteria SegmentListCriteria) ([]*CustomerSegment, error)
	FindByStatus(ctx context.Context, status SegmentStatus) ([]*CustomerSegment, error)
	FindByRuleID(ctx context.Context, ruleID shared.RuleID) ([]*CustomerSegment, error)
	FindByCreatedBy(ctx context.Context, createdBy shared.UserID) ([]*CustomerSegment, error)

	// Existence checks
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByID(ctx context.Context, id SegmentID) (bool, error)
}

// CustomerSegmentMembershipRepository defines the interface for customer segment membership persistence
type CustomerSegmentMembershipRepository interface {
	// Membership operations
	AddMembership(ctx context.Context, customerID CustomerID, segmentID SegmentID) error
	RemoveMembership(ctx context.Context, customerID CustomerID, segmentID SegmentID) error
	GetMemberships(ctx context.Context, customerID CustomerID) ([]SegmentID, error)
	GetMembers(ctx context.Context, segmentID SegmentID) ([]CustomerID, error)

	// Bulk operations
	BulkAddMemberships(ctx context.Context, memberships []Membership) error
	BulkRemoveMemberships(ctx context.Context, memberships []Membership) error

	// Existence checks
	HasMembership(ctx context.Context, customerID CustomerID, segmentID SegmentID) (bool, error)

	// History operations
	GetMembershipHistory(ctx context.Context, customerID CustomerID, segmentID SegmentID) ([]MembershipHistory, error)
}

// CustomerEventRepository defines the interface for customer event persistence
type CustomerEventRepository interface {
	// Event operations
	Save(ctx context.Context, event CustomerEvent) error
	SaveBatch(ctx context.Context, events []CustomerEvent) error

	// Query operations
	FindByCustomerID(ctx context.Context, customerID CustomerID, criteria EventCriteria) ([]CustomerEvent, error)
	FindByEventType(ctx context.Context, eventType string, criteria EventCriteria) ([]CustomerEvent, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time, criteria EventCriteria) ([]CustomerEvent, error)
	FindBySessionID(ctx context.Context, sessionID string) ([]CustomerEvent, error)

	// Analytics operations
	GetEventCount(ctx context.Context, customerID CustomerID, eventType string, dateRange *DateRange) (int64, error)
	GetEventMetrics(ctx context.Context, customerID CustomerID, dateRange *DateRange) (map[string]interface{}, error)
}

// ListCriteria represents criteria for listing customers
type ListCriteria struct {
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	SortBy     string                 `json:"sortBy"`
	SortOrder  string                 `json:"sortOrder"`
	Filters    map[string]interface{} `json:"filters"`
	Status     *CustomerStatus        `json:"status,omitempty"`
	SegmentIDs []SegmentID            `json:"segmentIds,omitempty"`
	Tags       []string               `json:"tags,omitempty"`
	DateRange  *DateRange             `json:"dateRange,omitempty"`
}

// SegmentListCriteria represents criteria for listing customer segments
type SegmentListCriteria struct {
	Page      int                    `json:"page"`
	Limit     int                    `json:"limit"`
	SortBy    string                 `json:"sortBy"`
	SortOrder string                 `json:"sortOrder"`
	Filters   map[string]interface{} `json:"filters"`
	Status    *SegmentStatus         `json:"status,omitempty"`
	RuleID    *shared.RuleID         `json:"ruleId,omitempty"`
	CreatedBy *shared.UserID         `json:"createdBy,omitempty"`
}

// EventCriteria represents criteria for querying customer events
type EventCriteria struct {
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	SortBy     string                 `json:"sortBy"`
	SortOrder  string                 `json:"sortOrder"`
	EventTypes []string               `json:"eventTypes,omitempty"`
	Channels   []string               `json:"channels,omitempty"`
	DateRange  *DateRange             `json:"dateRange,omitempty"`
	Filters    map[string]interface{} `json:"filters"`
}

// BulkUpdate represents a bulk update operation
type BulkUpdate struct {
	CustomerID CustomerID             `json:"customerId"`
	Updates    map[string]interface{} `json:"updates"`
}

// Membership represents a customer segment membership
type Membership struct {
	CustomerID CustomerID `json:"customerId"`
	SegmentID  SegmentID  `json:"segmentId"`
	JoinedAt   time.Time  `json:"joinedAt"`
}

// MembershipHistory represents the history of a customer segment membership
type MembershipHistory struct {
	ID         string     `json:"id"`
	CustomerID CustomerID `json:"customerId"`
	SegmentID  SegmentID  `json:"segmentId"`
	JoinedAt   time.Time  `json:"joinedAt"`
	LeftAt     *time.Time `json:"leftAt,omitempty"`
	IsActive   bool       `json:"isActive"`
}

// DateRange represents a date range
type DateRange struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
