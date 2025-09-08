package analytics

import (
	"time"

	"github.com/juanpablolazaro/ENGINE-RULES-SP/analytics-dashboard-service/internal/domain/shared"
)

// Dashboard represents the main dashboard aggregate
type Dashboard struct {
	ID              shared.DashboardID    `json:"id"`
	Name            string                `json:"name"`
	Description     string                `json:"description"`
	Layout          DashboardLayout       `json:"layout"`
	Widgets         []Widget              `json:"widgets"`
	Filters         DashboardFilters      `json:"filters"`
	RefreshInterval int                   `json:"refreshInterval"` // seconds
	IsPublic        bool                  `json:"isPublic"`
	OwnerID         shared.UserID         `json:"ownerId"`
	CreatedAt       time.Time             `json:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	Version         int                   `json:"version"`
	Events          []*shared.DomainEvent `json:"-"`
}

// DashboardLayout represents the layout configuration of a dashboard
type DashboardLayout struct {
	Columns    int  `json:"columns"`
	Rows       int  `json:"rows"`
	GridSize   int  `json:"gridSize"`
	Responsive bool `json:"responsive"`
}

// Widget represents a widget in a dashboard
type Widget struct {
	ID              shared.WidgetID `json:"id"`
	Type            WidgetType      `json:"type"`
	Title           string          `json:"title"`
	Position        WidgetPosition  `json:"position"`
	Size            WidgetSize      `json:"size"`
	Configuration   WidgetConfig    `json:"configuration"`
	DataSource      DataSource      `json:"dataSource"`
	RefreshInterval int             `json:"refreshInterval"`
}

// WidgetType represents the type of widget
type WidgetType string

const (
	WidgetTypeChart   WidgetType = "CHART"
	WidgetTypeTable   WidgetType = "TABLE"
	WidgetTypeKPI     WidgetType = "KPI"
	WidgetTypeGauge   WidgetType = "GAUGE"
	WidgetTypeHeatmap WidgetType = "HEATMAP"
	WidgetTypeMap     WidgetType = "MAP"
	WidgetTypeText    WidgetType = "TEXT"
	WidgetTypeImage   WidgetType = "IMAGE"
)

// WidgetPosition represents the position of a widget
type WidgetPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// WidgetSize represents the size of a widget
type WidgetSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// WidgetConfig represents the configuration of a widget
type WidgetConfig map[string]interface{}

// DataSource represents a data source for widgets
type DataSource struct {
	Type       string                 `json:"type"`
	Endpoint   string                 `json:"endpoint"`
	Query      string                 `json:"query"`
	Parameters map[string]interface{} `json:"parameters"`
}

// DashboardFilters represents filters applied to a dashboard
type DashboardFilters map[string]interface{}

// NewDashboard creates a new dashboard
func NewDashboard(name, description string, ownerID shared.UserID) *Dashboard {
	now := time.Now()
	return &Dashboard{
		ID:              shared.NewDashboardID(),
		Name:            name,
		Description:     description,
		Layout:          DashboardLayout{Columns: 4, Rows: 3, GridSize: 12, Responsive: true},
		Widgets:         []Widget{},
		Filters:         make(DashboardFilters),
		RefreshInterval: 300, // 5 minutes default
		IsPublic:        false,
		OwnerID:         ownerID,
		CreatedAt:       now,
		UpdatedAt:       now,
		Version:         1,
		Events:          []*shared.DomainEvent{},
	}
}

// AddWidget adds a widget to the dashboard
func (d *Dashboard) AddWidget(widget Widget) error {
	if !d.isValidWidgetPosition(widget.Position, widget.Size) {
		return shared.NewDomainError("INVALID_WIDGET_POSITION", "Widget position is outside dashboard bounds", "")
	}

	d.Widgets = append(d.Widgets, widget)
	d.UpdatedAt = time.Now()
	d.Version++

	event := shared.NewDomainEvent("WidgetAdded", d.ID.String(), d.Version, map[string]interface{}{
		"widgetId": widget.ID.String(),
		"type":     widget.Type,
	})
	d.Events = append(d.Events, event)

	return nil
}

// RemoveWidget removes a widget from the dashboard
func (d *Dashboard) RemoveWidget(widgetID shared.WidgetID) error {
	for i, widget := range d.Widgets {
		if widget.ID == widgetID {
			d.Widgets = append(d.Widgets[:i], d.Widgets[i+1:]...)
			d.UpdatedAt = time.Now()
			d.Version++

			event := shared.NewDomainEvent("WidgetRemoved", d.ID.String(), d.Version, map[string]interface{}{
				"widgetId": widgetID.String(),
			})
			d.Events = append(d.Events, event)

			return nil
		}
	}
	return shared.ErrNotFound
}

// UpdateLayout updates the dashboard layout
func (d *Dashboard) UpdateLayout(layout DashboardLayout) error {
	if layout.Columns <= 0 || layout.Rows <= 0 {
		return shared.NewDomainError("INVALID_LAYOUT", "Layout dimensions must be positive", "")
	}

	d.Layout = layout
	d.UpdatedAt = time.Now()
	d.Version++

	event := shared.NewDomainEvent("LayoutUpdated", d.ID.String(), d.Version, map[string]interface{}{
		"layout": layout,
	})
	d.Events = append(d.Events, event)

	return nil
}

// SetPublic sets the dashboard as public or private
func (d *Dashboard) SetPublic(isPublic bool) {
	d.IsPublic = isPublic
	d.UpdatedAt = time.Now()
	d.Version++

	event := shared.NewDomainEvent("VisibilityChanged", d.ID.String(), d.Version, map[string]interface{}{
		"isPublic": isPublic,
	})
	d.Events = append(d.Events, event)
}

// SetRefreshInterval sets the refresh interval for the dashboard
func (d *Dashboard) SetRefreshInterval(interval int) error {
	if interval < 30 || interval > 3600 {
		return shared.NewDomainError("INVALID_REFRESH_INTERVAL", "Refresh interval must be between 30 seconds and 1 hour", "")
	}

	d.RefreshInterval = interval
	d.UpdatedAt = time.Now()
	d.Version++

	event := shared.NewDomainEvent("RefreshIntervalChanged", d.ID.String(), d.Version, map[string]interface{}{
		"refreshInterval": interval,
	})
	d.Events = append(d.Events, event)

	return nil
}

// isValidWidgetPosition checks if the widget position is valid within the dashboard bounds
func (d *Dashboard) isValidWidgetPosition(position WidgetPosition, size WidgetSize) bool {
	return position.X >= 0 && position.Y >= 0 &&
		position.X+size.Width <= d.Layout.Columns &&
		position.Y+size.Height <= d.Layout.Rows
}

// ClearEvents clears the domain events
func (d *Dashboard) ClearEvents() {
	d.Events = []*shared.DomainEvent{}
}
