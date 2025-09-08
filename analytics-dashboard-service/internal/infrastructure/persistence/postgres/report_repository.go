package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/analytics"
	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
	"gorm.io/gorm"
)

// ReportModel represents the database model for Report
type ReportModel struct {
	ID            string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name          string `gorm:"not null"`
	Description   string
	Type          string          `gorm:"not null"`
	Template      json.RawMessage `gorm:"type:jsonb;not null"`
	Parameters    json.RawMessage `gorm:"type:jsonb;not null;default:'{}'"`
	Schedule      json.RawMessage `gorm:"type:jsonb"`
	OutputFormat  string          `gorm:"not null"`
	Recipients    json.RawMessage `gorm:"type:jsonb;not null;default:'[]'"`
	Status        string          `gorm:"not null;default:'ACTIVE'"`
	LastGenerated *time.Time
	NextRun       *time.Time
	OwnerID       string    `gorm:"type:uuid;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	Version       int       `gorm:"not null;default:1"`
}

// TableName returns the table name for ReportModel
func (ReportModel) TableName() string {
	return "reports"
}

// ReportRepository implements the ReportRepository interface
type ReportRepository struct {
	db *gorm.DB
}

// NewReportRepository creates a new report repository
func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// Save saves a report to the database
func (r *ReportRepository) Save(ctx context.Context, report *analytics.Report) error {
	model, err := r.domainToModel(report)
	if err != nil {
		return fmt.Errorf("failed to convert domain to model: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}

	// Update the domain object with the generated ID
	report.ID = shared.ReportID(model.ID)
	return nil
}

// FindByID finds a report by ID
func (r *ReportRepository) FindByID(ctx context.Context, id shared.ReportID) (*analytics.Report, error) {
	var model ReportModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find report: %w", err)
	}

	return r.modelToDomain(&model)
}

// FindByOwnerID finds reports by owner ID
func (r *ReportRepository) FindByOwnerID(ctx context.Context, ownerID shared.UserID) ([]*analytics.Report, error) {
	var models []ReportModel
	if err := r.db.WithContext(ctx).Where("owner_id = ?", ownerID.String()).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find reports by owner: %w", err)
	}

	reports := make([]*analytics.Report, len(models))
	for i, model := range models {
		report, err := r.modelToDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		reports[i] = report
	}

	return reports, nil
}

// FindByStatus finds reports by status
func (r *ReportRepository) FindByStatus(ctx context.Context, status analytics.ReportStatus) ([]*analytics.Report, error) {
	var models []ReportModel
	if err := r.db.WithContext(ctx).Where("status = ?", string(status)).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find reports by status: %w", err)
	}

	reports := make([]*analytics.Report, len(models))
	for i, model := range models {
		report, err := r.modelToDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		reports[i] = report
	}

	return reports, nil
}

// FindScheduled finds reports that are scheduled to run before the given time
func (r *ReportRepository) FindScheduled(ctx context.Context, before time.Time) ([]*analytics.Report, error) {
	var models []ReportModel
	if err := r.db.WithContext(ctx).Where("next_run IS NOT NULL AND next_run <= ?", before).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find scheduled reports: %w", err)
	}

	reports := make([]*analytics.Report, len(models))
	for i, model := range models {
		report, err := r.modelToDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		reports[i] = report
	}

	return reports, nil
}

// Delete deletes a report by ID
func (r *ReportRepository) Delete(ctx context.Context, id shared.ReportID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).Delete(&ReportModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete report: %w", err)
	}
	return nil
}

// Update updates a report
func (r *ReportRepository) Update(ctx context.Context, report *analytics.Report) error {
	model, err := r.domainToModel(report)
	if err != nil {
		return fmt.Errorf("failed to convert domain to model: %w", err)
	}

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update report: %w", err)
	}

	return nil
}

// domainToModel converts a domain Report to a database model
func (r *ReportRepository) domainToModel(report *analytics.Report) (*ReportModel, error) {
	templateBytes, err := json.Marshal(report.Template)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal template: %w", err)
	}

	parametersBytes, err := json.Marshal(report.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parameters: %w", err)
	}

	var scheduleBytes json.RawMessage
	if report.Schedule != nil {
		scheduleBytes, err = json.Marshal(report.Schedule)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal schedule: %w", err)
		}
	}

	recipientsBytes, err := json.Marshal(report.Recipients)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal recipients: %w", err)
	}

	return &ReportModel{
		ID:            report.ID.String(),
		Name:          report.Name,
		Description:   report.Description,
		Type:          string(report.Type),
		Template:      templateBytes,
		Parameters:    parametersBytes,
		Schedule:      scheduleBytes,
		OutputFormat:  string(report.OutputFormat),
		Recipients:    recipientsBytes,
		Status:        string(report.Status),
		LastGenerated: report.LastGenerated,
		NextRun:       report.NextRun,
		OwnerID:       report.OwnerID.String(),
		Version:       report.Version,
	}, nil
}

// modelToDomain converts a database model to a domain Report
func (r *ReportRepository) modelToDomain(model *ReportModel) (*analytics.Report, error) {
	var template analytics.ReportTemplate
	if err := json.Unmarshal(model.Template, &template); err != nil {
		return nil, fmt.Errorf("failed to unmarshal template: %w", err)
	}

	var parameters analytics.ReportParameters
	if err := json.Unmarshal(model.Parameters, &parameters); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parameters: %w", err)
	}

	var schedule *analytics.ReportSchedule
	if len(model.Schedule) > 0 {
		schedule = &analytics.ReportSchedule{}
		if err := json.Unmarshal(model.Schedule, schedule); err != nil {
			return nil, fmt.Errorf("failed to unmarshal schedule: %w", err)
		}
	}

	var recipients []string
	if err := json.Unmarshal(model.Recipients, &recipients); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recipients: %w", err)
	}

	return &analytics.Report{
		ID:            shared.ReportID(model.ID),
		Name:          model.Name,
		Description:   model.Description,
		Type:          analytics.ReportType(model.Type),
		Template:      template,
		Parameters:    parameters,
		Schedule:      schedule,
		OutputFormat:  analytics.OutputFormat(model.OutputFormat),
		Recipients:    recipients,
		Status:        analytics.ReportStatus(model.Status),
		LastGenerated: model.LastGenerated,
		NextRun:       model.NextRun,
		OwnerID:       shared.UserID(model.OwnerID),
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
		Version:       model.Version,
		Events:        []*shared.DomainEvent{},
	}, nil
}
