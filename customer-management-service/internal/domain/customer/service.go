package customer

import (
	"context"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/customer-management-service/internal/domain/shared"
)

// CustomerSegmentationService defines the interface for customer segmentation
type CustomerSegmentationService interface {
	// Segment evaluation
	EvaluateCustomer(ctx context.Context, customerID CustomerID, segmentID SegmentID) (bool, error)
	CalculateSegment(ctx context.Context, segmentID SegmentID) ([]CustomerID, error)
	UpdateSegmentMembership(ctx context.Context, customerID CustomerID) error
	GetCustomerSegments(ctx context.Context, customerID CustomerID) ([]SegmentID, error)

	// Segment management
	RecalculateSegment(ctx context.Context, segmentID SegmentID) error
	RecalculateAllSegments(ctx context.Context) error

	// Segment validation
	ValidateSegmentCriteria(ctx context.Context, criteria SegmentCriteria) error
}

// CustomerAnalyticsService defines the interface for customer analytics
type CustomerAnalyticsService interface {
	// Metrics calculation
	CalculateCustomerMetrics(ctx context.Context, customerID CustomerID) (*CustomerMetrics, error)
	GetCustomerInsights(ctx context.Context, customerID CustomerID) (*CustomerInsights, error)

	// Event tracking
	TrackCustomerEvent(ctx context.Context, customerID CustomerID, event CustomerEvent) error
	TrackCustomerEvents(ctx context.Context, events []CustomerEvent) error

	// Reporting
	GenerateCustomerReport(ctx context.Context, criteria ReportCriteria) (*CustomerReport, error)
	GenerateSegmentReport(ctx context.Context, segmentID SegmentID, criteria ReportCriteria) (*CustomerReport, error)

	// Analytics queries
	GetCustomerBehaviorPatterns(ctx context.Context, customerID CustomerID, dateRange *DateRange) ([]BehaviorPattern, error)
	GetCustomerEngagementScore(ctx context.Context, customerID CustomerID) (float64, error)
	GetCustomerChurnRisk(ctx context.Context, customerID CustomerID) (float64, error)
}

// CustomerPrivacyService defines the interface for customer privacy and GDPR compliance
type CustomerPrivacyService interface {
	// Data management
	AnonymizeCustomerData(ctx context.Context, customerID CustomerID) error
	ExportCustomerData(ctx context.Context, customerID CustomerID) (*CustomerDataExport, error)
	DeleteCustomerData(ctx context.Context, customerID CustomerID) error

	// Consent management
	UpdatePrivacyConsent(ctx context.Context, customerID CustomerID, consent PrivacyConsent) error
	GetPrivacyConsent(ctx context.Context, customerID CustomerID) (*PrivacyConsent, error)

	// Data processing
	ProcessDataRequest(ctx context.Context, request DataRequest) (*DataRequestResponse, error)
	ValidateDataProcessing(ctx context.Context, customerID CustomerID, purpose string) (bool, error)
}

// CustomerValidationService defines the interface for customer data validation
type CustomerValidationService interface {
	// Email validation
	ValidateEmail(ctx context.Context, email shared.EmailAddress) error
	IsEmailAvailable(ctx context.Context, email shared.EmailAddress) (bool, error)

	// Customer validation
	ValidateCustomer(ctx context.Context, customer *Customer) error
	ValidateCustomerUpdate(ctx context.Context, customer *Customer, updates map[string]interface{}) error

	// Segment validation
	ValidateSegment(ctx context.Context, segment *CustomerSegment) error
	ValidateSegmentCriteria(ctx context.Context, criteria SegmentCriteria) error

	// Business rules validation
	ValidateBusinessRules(ctx context.Context, customer *Customer) error
	ValidateSegmentBusinessRules(ctx context.Context, segment *CustomerSegment) error
}

// CustomerNotificationService defines the interface for customer notifications
type CustomerNotificationService interface {
	// Notification sending
	SendWelcomeNotification(ctx context.Context, customerID CustomerID) error
	SendSegmentJoinedNotification(ctx context.Context, customerID CustomerID, segmentID SegmentID) error
	SendSegmentLeftNotification(ctx context.Context, customerID CustomerID, segmentID SegmentID) error
	SendPrivacyUpdateNotification(ctx context.Context, customerID CustomerID) error

	// Notification preferences
	GetNotificationPreferences(ctx context.Context, customerID CustomerID) (*NotificationSettings, error)
	UpdateNotificationPreferences(ctx context.Context, customerID CustomerID, preferences NotificationSettings) error

	// Notification history
	GetNotificationHistory(ctx context.Context, customerID CustomerID, criteria NotificationCriteria) ([]NotificationRecord, error)
}

// DataRequest represents a data request (GDPR)
type DataRequest struct {
	ID          string                 `json:"id"`
	CustomerID  CustomerID             `json:"customerId"`
	RequestType string                 `json:"requestType"` // EXPORT, DELETE, ANONYMIZE
	Purpose     string                 `json:"purpose"`
	RequestedAt time.Time              `json:"requestedAt"`
	RequestedBy string                 `json:"requestedBy"`
	Data        map[string]interface{} `json:"data"`
}

// DataRequestResponse represents a response to a data request
type DataRequestResponse struct {
	RequestID   string                 `json:"requestId"`
	Status      string                 `json:"status"` // PENDING, PROCESSING, COMPLETED, FAILED
	Result      map[string]interface{} `json:"result,omitempty"`
	ProcessedAt *time.Time             `json:"processedAt,omitempty"`
	Error       *string                `json:"error,omitempty"`
}

// NotificationCriteria represents criteria for querying notifications
type NotificationCriteria struct {
	Page      int                    `json:"page"`
	Limit     int                    `json:"limit"`
	SortBy    string                 `json:"sortBy"`
	SortOrder string                 `json:"sortOrder"`
	Types     []string               `json:"types,omitempty"`
	Channels  []string               `json:"channels,omitempty"`
	DateRange *DateRange             `json:"dateRange,omitempty"`
	Status    []string               `json:"status,omitempty"`
	Filters   map[string]interface{} `json:"filters"`
}

// NotificationRecord represents a notification record
type NotificationRecord struct {
	ID          string                 `json:"id"`
	CustomerID  CustomerID             `json:"customerId"`
	Type        string                 `json:"type"`
	Channel     string                 `json:"channel"`
	Title       string                 `json:"title"`
	Content     string                 `json:"content"`
	Status      string                 `json:"status"` // SENT, DELIVERED, READ, FAILED
	SentAt      time.Time              `json:"sentAt"`
	DeliveredAt *time.Time             `json:"deliveredAt,omitempty"`
	ReadAt      *time.Time             `json:"readAt,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// CustomerIntegrationService defines the interface for customer integrations
type CustomerIntegrationService interface {
	// External system integration
	SyncWithExternalSystem(ctx context.Context, customerID CustomerID, system string) error
	ImportFromExternalSystem(ctx context.Context, system string, data map[string]interface{}) error
	ExportToExternalSystem(ctx context.Context, customerID CustomerID, system string) error

	// Data synchronization
	SyncCustomerData(ctx context.Context, customerID CustomerID) error
	SyncSegmentData(ctx context.Context, segmentID SegmentID) error

	// Integration validation
	ValidateIntegrationData(ctx context.Context, system string, data map[string]interface{}) error
	TestIntegration(ctx context.Context, system string) error
}

// CustomerAuditService defines the interface for customer audit and compliance
type CustomerAuditService interface {
	// Audit logging
	LogCustomerAction(ctx context.Context, customerID CustomerID, action string, details map[string]interface{}) error
	LogSegmentAction(ctx context.Context, segmentID SegmentID, action string, details map[string]interface{}) error

	// Audit queries
	GetCustomerAuditLog(ctx context.Context, customerID CustomerID, criteria AuditCriteria) ([]AuditRecord, error)
	GetSegmentAuditLog(ctx context.Context, segmentID SegmentID, criteria AuditCriteria) ([]AuditRecord, error)

	// Compliance reporting
	GenerateComplianceReport(ctx context.Context, criteria ComplianceCriteria) (*ComplianceReport, error)
	ValidateCompliance(ctx context.Context, customerID CustomerID) (*ComplianceValidation, error)
}

// AuditCriteria represents criteria for querying audit logs
type AuditCriteria struct {
	Page      int                    `json:"page"`
	Limit     int                    `json:"limit"`
	SortBy    string                 `json:"sortBy"`
	SortOrder string                 `json:"sortOrder"`
	Actions   []string               `json:"actions,omitempty"`
	Users     []string               `json:"users,omitempty"`
	DateRange *DateRange             `json:"dateRange,omitempty"`
	Filters   map[string]interface{} `json:"filters"`
}

// AuditRecord represents an audit record
type AuditRecord struct {
	ID         string                 `json:"id"`
	EntityType string                 `json:"entityType"` // CUSTOMER, SEGMENT
	EntityID   string                 `json:"entityId"`
	Action     string                 `json:"action"`
	UserID     string                 `json:"userId"`
	Timestamp  time.Time              `json:"timestamp"`
	Details    map[string]interface{} `json:"details"`
	IPAddress  *string                `json:"ipAddress,omitempty"`
	UserAgent  *string                `json:"userAgent,omitempty"`
}

// ComplianceCriteria represents criteria for compliance reporting
type ComplianceCriteria struct {
	DateRange       *DateRange   `json:"dateRange,omitempty"`
	CustomerIDs     []CustomerID `json:"customerIds,omitempty"`
	SegmentIDs      []SegmentID  `json:"segmentIds,omitempty"`
	ComplianceTypes []string     `json:"complianceTypes,omitempty"`
}

// ComplianceReport represents a compliance report
type ComplianceReport struct {
	ID              string                 `json:"id"`
	Title           string                 `json:"title"`
	GeneratedAt     time.Time              `json:"generatedAt"`
	GeneratedBy     string                 `json:"generatedBy"`
	Criteria        ComplianceCriteria     `json:"criteria"`
	Summary         map[string]interface{} `json:"summary"`
	Details         map[string]interface{} `json:"details"`
	Recommendations []string               `json:"recommendations"`
}

// ComplianceValidation represents a compliance validation result
type ComplianceValidation struct {
	CustomerID  CustomerID        `json:"customerId"`
	IsCompliant bool              `json:"isCompliant"`
	Issues      []ComplianceIssue `json:"issues"`
	ValidatedAt time.Time         `json:"validatedAt"`
	ValidatedBy string            `json:"validatedBy"`
}

// ComplianceIssue represents a compliance issue
type ComplianceIssue struct {
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"` // HIGH, MEDIUM, LOW
	Description string                 `json:"description"`
	Details     map[string]interface{} `json:"details"`
	Resolution  *string                `json:"resolution,omitempty"`
}
