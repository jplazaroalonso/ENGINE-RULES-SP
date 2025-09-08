package shared

import (
	"time"

	"github.com/google/uuid"
)

// UserID represents a user identifier
type UserID string

// NewUserID creates a new user ID
func NewUserID() UserID {
	return UserID(uuid.New().String())
}

// String returns the string representation of UserID
func (u UserID) String() string {
	return string(u)
}

// DashboardID represents a dashboard identifier
type DashboardID string

// NewDashboardID creates a new dashboard ID
func NewDashboardID() DashboardID {
	return DashboardID(uuid.New().String())
}

// String returns the string representation of DashboardID
func (d DashboardID) String() string {
	return string(d)
}

// ReportID represents a report identifier
type ReportID string

// NewReportID creates a new report ID
func NewReportID() ReportID {
	return ReportID(uuid.New().String())
}

// String returns the string representation of ReportID
func (r ReportID) String() string {
	return string(r)
}

// MetricID represents a metric identifier
type MetricID string

// NewMetricID creates a new metric ID
func NewMetricID() MetricID {
	return MetricID(uuid.New().String())
}

// String returns the string representation of MetricID
func (m MetricID) String() string {
	return string(m)
}

// WidgetID represents a widget identifier
type WidgetID string

// NewWidgetID creates a new widget ID
func NewWidgetID() WidgetID {
	return WidgetID(uuid.New().String())
}

// String returns the string representation of WidgetID
func (w WidgetID) String() string {
	return string(w)
}

// TimeRange represents a time range for analytics queries
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// NewTimeRange creates a new time range
func NewTimeRange(start, end time.Time) *TimeRange {
	return &TimeRange{
		Start: start,
		End:   end,
	}
}

// IsValid checks if the time range is valid
func (tr *TimeRange) IsValid() bool {
	return tr.Start.Before(tr.End)
}

// Duration returns the duration of the time range
func (tr *TimeRange) Duration() time.Duration {
	return tr.End.Sub(tr.Start)
}
