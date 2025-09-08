package analytics

import (
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
)

// Report represents a report aggregate
type Report struct {
	ID            shared.ReportID       `json:"id"`
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	Type          ReportType            `json:"type"`
	Template      ReportTemplate        `json:"template"`
	Parameters    ReportParameters      `json:"parameters"`
	Schedule      *ReportSchedule       `json:"schedule,omitempty"`
	OutputFormat  OutputFormat          `json:"outputFormat"`
	Recipients    []string              `json:"recipients"`
	Status        ReportStatus          `json:"status"`
	LastGenerated *time.Time            `json:"lastGenerated,omitempty"`
	NextRun       *time.Time            `json:"nextRun,omitempty"`
	OwnerID       shared.UserID         `json:"ownerId"`
	CreatedAt     time.Time             `json:"createdAt"`
	UpdatedAt     time.Time             `json:"updatedAt"`
	Version       int                   `json:"version"`
	Events        []*shared.DomainEvent `json:"-"`
}

// ReportType represents the type of report
type ReportType string

const (
	ReportTypePerformance ReportType = "PERFORMANCE"
	ReportTypeCompliance  ReportType = "COMPLIANCE"
	ReportTypeBusiness    ReportType = "BUSINESS"
	ReportTypeCustom      ReportType = "CUSTOM"
)

// OutputFormat represents the output format of a report
type OutputFormat string

const (
	OutputFormatPDF   OutputFormat = "PDF"
	OutputFormatExcel OutputFormat = "EXCEL"
	OutputFormatCSV   OutputFormat = "CSV"
	OutputFormatJSON  OutputFormat = "JSON"
	OutputFormatHTML  OutputFormat = "HTML"
)

// ReportStatus represents the status of a report
type ReportStatus string

const (
	ReportStatusActive     ReportStatus = "ACTIVE"
	ReportStatusInactive   ReportStatus = "INACTIVE"
	ReportStatusGenerating ReportStatus = "GENERATING"
	ReportStatusError      ReportStatus = "ERROR"
)

// ReportTemplate represents the template configuration for a report
type ReportTemplate struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Layout      ReportLayout           `json:"layout"`
	Sections    []ReportSection        `json:"sections"`
	Styles      map[string]interface{} `json:"styles"`
}

// ReportLayout represents the layout of a report
type ReportLayout struct {
	Orientation string `json:"orientation"` // "portrait" or "landscape"
	PageSize    string `json:"pageSize"`    // "A4", "Letter", etc.
	Margins     struct {
		Top    float64 `json:"top"`
		Bottom float64 `json:"bottom"`
		Left   float64 `json:"left"`
		Right  float64 `json:"right"`
	} `json:"margins"`
}

// ReportSection represents a section in a report
type ReportSection struct {
	ID      string                 `json:"id"`
	Type    string                 `json:"type"` // "text", "chart", "table", "image"
	Title   string                 `json:"title"`
	Content map[string]interface{} `json:"content"`
	Order   int                    `json:"order"`
	Visible bool                   `json:"visible"`
}

// ReportParameters represents parameters for report generation
type ReportParameters map[string]interface{}

// ReportSchedule represents the schedule for report generation
type ReportSchedule struct {
	Type     ScheduleType `json:"type"`
	Interval int          `json:"interval"` // in minutes
	Days     []int        `json:"days"`     // 0-6 (Sunday-Saturday)
	Time     string       `json:"time"`     // HH:MM format
	Timezone string       `json:"timezone"`
}

// ScheduleType represents the type of schedule
type ScheduleType string

const (
	ScheduleTypeOnce    ScheduleType = "ONCE"
	ScheduleTypeHourly  ScheduleType = "HOURLY"
	ScheduleTypeDaily   ScheduleType = "DAILY"
	ScheduleTypeWeekly  ScheduleType = "WEEKLY"
	ScheduleTypeMonthly ScheduleType = "MONTHLY"
)

// NewReport creates a new report
func NewReport(name, description string, reportType ReportType, ownerID shared.UserID) *Report {
	now := time.Now()
	return &Report{
		ID:           shared.NewReportID(),
		Name:         name,
		Description:  description,
		Type:         reportType,
		Template:     ReportTemplate{},
		Parameters:   make(ReportParameters),
		OutputFormat: OutputFormatPDF,
		Recipients:   []string{},
		Status:       ReportStatusActive,
		OwnerID:      ownerID,
		CreatedAt:    now,
		UpdatedAt:    now,
		Version:      1,
		Events:       []*shared.DomainEvent{},
	}
}

// SetTemplate sets the report template
func (r *Report) SetTemplate(template ReportTemplate) {
	r.Template = template
	r.UpdatedAt = time.Now()
	r.Version++

	event := shared.NewDomainEvent("TemplateUpdated", r.ID.String(), r.Version, map[string]interface{}{
		"templateId": template.ID,
	})
	r.Events = append(r.Events, event)
}

// SetSchedule sets the report schedule
func (r *Report) SetSchedule(schedule *ReportSchedule) error {
	if schedule != nil {
		if err := r.validateSchedule(schedule); err != nil {
			return err
		}
		r.calculateNextRun(schedule)
	}

	r.Schedule = schedule
	r.UpdatedAt = time.Now()
	r.Version++

	event := shared.NewDomainEvent("ScheduleUpdated", r.ID.String(), r.Version, map[string]interface{}{
		"schedule": schedule,
	})
	r.Events = append(r.Events, event)

	return nil
}

// SetOutputFormat sets the output format
func (r *Report) SetOutputFormat(format OutputFormat) error {
	if !r.isValidOutputFormat(format) {
		return shared.NewDomainError("INVALID_OUTPUT_FORMAT", "Invalid output format", "")
	}

	r.OutputFormat = format
	r.UpdatedAt = time.Now()
	r.Version++

	event := shared.NewDomainEvent("OutputFormatChanged", r.ID.String(), r.Version, map[string]interface{}{
		"outputFormat": format,
	})
	r.Events = append(r.Events, event)

	return nil
}

// AddRecipient adds a recipient to the report
func (r *Report) AddRecipient(email string) error {
	if !r.isValidEmail(email) {
		return shared.NewDomainError("INVALID_EMAIL", "Invalid email address", "")
	}

	for _, recipient := range r.Recipients {
		if recipient == email {
			return shared.NewDomainError("DUPLICATE_RECIPIENT", "Recipient already exists", "")
		}
	}

	r.Recipients = append(r.Recipients, email)
	r.UpdatedAt = time.Now()
	r.Version++

	event := shared.NewDomainEvent("RecipientAdded", r.ID.String(), r.Version, map[string]interface{}{
		"email": email,
	})
	r.Events = append(r.Events, event)

	return nil
}

// RemoveRecipient removes a recipient from the report
func (r *Report) RemoveRecipient(email string) error {
	for i, recipient := range r.Recipients {
		if recipient == email {
			r.Recipients = append(r.Recipients[:i], r.Recipients[i+1:]...)
			r.UpdatedAt = time.Now()
			r.Version++

			event := shared.NewDomainEvent("RecipientRemoved", r.ID.String(), r.Version, map[string]interface{}{
				"email": email,
			})
			r.Events = append(r.Events, event)

			return nil
		}
	}
	return shared.ErrNotFound
}

// SetStatus sets the report status
func (r *Report) SetStatus(status ReportStatus) {
	r.Status = status
	r.UpdatedAt = time.Now()
	r.Version++

	event := shared.NewDomainEvent("StatusChanged", r.ID.String(), r.Version, map[string]interface{}{
		"status": status,
	})
	r.Events = append(r.Events, event)
}

// MarkAsGenerated marks the report as generated
func (r *Report) MarkAsGenerated() {
	now := time.Now()
	r.LastGenerated = &now
	r.Status = ReportStatusActive

	if r.Schedule != nil {
		r.calculateNextRun(r.Schedule)
	}

	r.UpdatedAt = now
	r.Version++

	event := shared.NewDomainEvent("ReportGenerated", r.ID.String(), r.Version, map[string]interface{}{
		"generatedAt": now,
	})
	r.Events = append(r.Events, event)
}

// validateSchedule validates the report schedule
func (r *Report) validateSchedule(schedule *ReportSchedule) error {
	switch schedule.Type {
	case ScheduleTypeHourly:
		if schedule.Interval < 1 || schedule.Interval > 24 {
			return shared.NewDomainError("INVALID_HOURLY_INTERVAL", "Hourly interval must be between 1 and 24 hours", "")
		}
	case ScheduleTypeDaily:
		if schedule.Time == "" {
			return shared.NewDomainError("INVALID_DAILY_SCHEDULE", "Daily schedule must have a time", "")
		}
	case ScheduleTypeWeekly:
		if len(schedule.Days) == 0 {
			return shared.NewDomainError("INVALID_WEEKLY_SCHEDULE", "Weekly schedule must have at least one day", "")
		}
	}
	return nil
}

// calculateNextRun calculates the next run time based on the schedule
func (r *Report) calculateNextRun(schedule *ReportSchedule) {
	now := time.Now()

	switch schedule.Type {
	case ScheduleTypeHourly:
		nextRun := now.Add(time.Duration(schedule.Interval) * time.Hour)
		r.NextRun = &nextRun
	case ScheduleTypeDaily:
		// Parse time and set for next day
		nextRun := now.Add(24 * time.Hour)
		r.NextRun = &nextRun
	case ScheduleTypeWeekly:
		// Calculate next occurrence based on days
		nextRun := now.Add(7 * 24 * time.Hour)
		r.NextRun = &nextRun
	}
}

// isValidOutputFormat checks if the output format is valid
func (r *Report) isValidOutputFormat(format OutputFormat) bool {
	validFormats := []OutputFormat{
		OutputFormatPDF, OutputFormatExcel, OutputFormatCSV,
		OutputFormatJSON, OutputFormatHTML,
	}

	for _, validFormat := range validFormats {
		if format == validFormat {
			return true
		}
	}
	return false
}

// isValidEmail checks if the email is valid (basic validation)
func (r *Report) isValidEmail(email string) bool {
	return len(email) > 0 && len(email) < 255
}

// ClearEvents clears the domain events
func (r *Report) ClearEvents() {
	r.Events = []*shared.DomainEvent{}
}
