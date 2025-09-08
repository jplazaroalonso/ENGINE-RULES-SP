package postgres

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectToPostgreSQL connects to PostgreSQL database
func ConnectToPostgreSQL(dsn string) (*gorm.DB, error) {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	// Migrate customer-related tables
	if err := db.AutoMigrate(
		&CustomerDBModel{},
		&CustomerSegmentDBModel{},
		&CustomerSegmentMembershipDBModel{},
		&CustomerEventDBModel{},
	); err != nil {
		return fmt.Errorf("failed to run auto-migration: %w", err)
	}

	// Create indexes
	if err := createIndexes(db); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// createIndexes creates database indexes for better performance
func createIndexes(db *gorm.DB) error {
	// Customer indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email)").Error; err != nil {
		return fmt.Errorf("failed to create email index: %w", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customers_status ON customers(status)").Error; err != nil {
		return fmt.Errorf("failed to create status index: %w", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customers_created_at ON customers(created_at)").Error; err != nil {
		return fmt.Errorf("failed to create created_at index: %w", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customers_last_activity ON customers(last_activity)").Error; err != nil {
		return fmt.Errorf("failed to create last_activity index: %w", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customers_segments ON customers USING GIN(segments)").Error; err != nil {
		return fmt.Errorf("failed to create segments index: %w", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customers_tags ON customers USING GIN(tags)").Error; err != nil {
		return fmt.Errorf("failed to create tags index: %w", err)
	}

	// Customer segment indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customer_segments_rule_id ON customer_segments(rule_id)").Error; err != nil {
		return fmt.Errorf("failed to create rule_id index: %w", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customer_segments_status ON customer_segments(status)").Error; err != nil {
		return fmt.Errorf("failed to create segment status index: %w", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customer_segments_created_by ON customer_segments(created_by)").Error; err != nil {
		return fmt.Errorf("failed to create created_by index: %w", err)
	}

	// Customer segment membership indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customer_segment_membership_customer_id ON customer_segment_membership(customer_id)").Error; err != nil {
		return fmt.Errorf("failed to create customer_id index: %w", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customer_segment_membership_segment_id ON customer_segment_membership(segment_id)").Error; err != nil {
		return fmt.Errorf("failed to create segment_id index: %w", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customer_segment_membership_active ON customer_segment_membership(is_active)").Error; err != nil {
		return fmt.Errorf("failed to create is_active index: %w", err)
	}

	// Customer event indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customer_events_customer_id ON customer_events(customer_id)").Error; err != nil {
		return fmt.Errorf("failed to create customer_id index: %w", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customer_events_type ON customer_events(event_type)").Error; err != nil {
		return fmt.Errorf("failed to create event_type index: %w", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customer_events_occurred_at ON customer_events(occurred_at)").Error; err != nil {
		return fmt.Errorf("failed to create occurred_at index: %w", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customer_events_session_id ON customer_events(session_id)").Error; err != nil {
		return fmt.Errorf("failed to create session_id index: %w", err)
	}

	return nil
}

// CustomerSegmentDBModel represents the customer segment database model
type CustomerSegmentDBModel struct {
	ID             string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name           string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	Description    string     `gorm:"type:text"`
	RuleID         string     `gorm:"type:uuid;not null"`
	CustomerCount  int        `gorm:"type:integer;not null;default:0"`
	Criteria       JSONB      `gorm:"type:jsonb;not null"`
	Status         string     `gorm:"type:varchar(20);not null;default:'ACTIVE'"`
	CreatedBy      string     `gorm:"type:uuid;not null"`
	CreatedAt      time.Time  `gorm:"type:timestamp with time zone;default:now()"`
	UpdatedAt      time.Time  `gorm:"type:timestamp with time zone;default:now()"`
	LastCalculated *time.Time `gorm:"type:timestamp with time zone"`
	Version        int        `gorm:"type:integer;not null;default:1"`
}

// CustomerSegmentMembershipDBModel represents the customer segment membership database model
type CustomerSegmentMembershipDBModel struct {
	ID         string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerID string     `gorm:"type:uuid;not null;index"`
	SegmentID  string     `gorm:"type:uuid;not null;index"`
	JoinedAt   time.Time  `gorm:"type:timestamp with time zone;default:now()"`
	LeftAt     *time.Time `gorm:"type:timestamp with time zone"`
	IsActive   bool       `gorm:"type:boolean;not null;default:true;index"`
}

// CustomerEventDBModel represents the customer event database model
type CustomerEventDBModel struct {
	ID         string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerID string    `gorm:"type:uuid;not null;index"`
	EventType  string    `gorm:"type:varchar(50);not null;index"`
	EventData  JSONB     `gorm:"type:jsonb"`
	OccurredAt time.Time `gorm:"type:timestamp with time zone;default:now();index"`
	SessionID  *string   `gorm:"type:uuid;index"`
	DeviceInfo JSONB     `gorm:"type:jsonb"`
}
