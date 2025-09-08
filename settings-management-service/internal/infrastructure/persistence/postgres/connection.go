package postgres

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectionConfig holds the database connection configuration
type ConnectionConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConnectionConfig creates a new connection config from environment variables
func NewConnectionConfig() *ConnectionConfig {
	return &ConnectionConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "settings_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// Connect establishes a connection to PostgreSQL
func Connect(config *ConnectionConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	// Configure GORM logger
	var gormLogger logger.Interface
	if os.Getenv("GORM_LOG_LEVEL") == "debug" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

// AutoMigrate runs database migrations for all models
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Migrate all models
	err := db.AutoMigrate(
		&ConfigurationDBModel{},
		&FeatureFlagDBModel{},
		&UserPreferenceDBModel{},
		&OrganizationSettingDBModel{},
	)
	if err != nil {
		return fmt.Errorf("failed to run database migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// CreateIndexes creates additional database indexes for better performance
func CreateIndexes(db *gorm.DB) error {
	log.Println("Creating additional database indexes...")

	// Create indexes for configurations
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_configurations_created_at ON configurations(created_at)").Error; err != nil {
		return fmt.Errorf("failed to create index on configurations.created_at: %w", err)
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_configurations_updated_at ON configurations(updated_at)").Error; err != nil {
		return fmt.Errorf("failed to create index on configurations.updated_at: %w", err)
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_configurations_created_by ON configurations(created_by)").Error; err != nil {
		return fmt.Errorf("failed to create index on configurations.created_by: %w", err)
	}

	// Create indexes for feature flags
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_feature_flags_created_at ON feature_flags(created_at)").Error; err != nil {
		return fmt.Errorf("failed to create index on feature_flags.created_at: %w", err)
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_feature_flags_updated_at ON feature_flags(updated_at)").Error; err != nil {
		return fmt.Errorf("failed to create index on feature_flags.updated_at: %w", err)
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_feature_flags_created_by ON feature_flags(created_by)").Error; err != nil {
		return fmt.Errorf("failed to create index on feature_flags.created_by: %w", err)
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_feature_flags_is_enabled ON feature_flags(is_enabled)").Error; err != nil {
		return fmt.Errorf("failed to create index on feature_flags.is_enabled: %w", err)
	}

	// Create indexes for user preferences
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_user_preferences_created_at ON user_preferences(created_at)").Error; err != nil {
		return fmt.Errorf("failed to create index on user_preferences.created_at: %w", err)
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_user_preferences_updated_at ON user_preferences(updated_at)").Error; err != nil {
		return fmt.Errorf("failed to create index on user_preferences.updated_at: %w", err)
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_user_preferences_created_by ON user_preferences(created_by)").Error; err != nil {
		return fmt.Errorf("failed to create index on user_preferences.created_by: %w", err)
	}

	// Create indexes for organization settings
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_organization_settings_created_at ON organization_settings(created_at)").Error; err != nil {
		return fmt.Errorf("failed to create index on organization_settings.created_at: %w", err)
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_organization_settings_updated_at ON organization_settings(updated_at)").Error; err != nil {
		return fmt.Errorf("failed to create index on organization_settings.updated_at: %w", err)
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_organization_settings_created_by ON organization_settings(created_by)").Error; err != nil {
		return fmt.Errorf("failed to create index on organization_settings.created_by: %w", err)
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_organization_settings_parent_id ON organization_settings(parent_id)").Error; err != nil {
		return fmt.Errorf("failed to create index on organization_settings.parent_id: %w", err)
	}

	log.Println("Additional database indexes created successfully")
	return nil
}

// Close closes the database connection
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Println("Database connection closed")
	return nil
}

// HealthCheck performs a health check on the database connection
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
