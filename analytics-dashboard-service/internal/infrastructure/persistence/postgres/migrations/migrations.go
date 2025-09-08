package migrations

import (
	"embed"
	"fmt"
	"sort"
	"strings"

	"gorm.io/gorm"
)

//go:embed *.sql
var migrationFiles embed.FS

// ApplyMigrations applies all database migrations
func ApplyMigrations(db *gorm.DB) error {
	// Get all migration files
	files, err := migrationFiles.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	// Sort files by name to ensure correct order
	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Apply each migration
	for _, fileName := range migrationFiles {
		if err := applyMigration(db, fileName); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", fileName, err)
		}
	}

	return nil
}

// createMigrationsTable creates the migrations tracking table
func createMigrationsTable(db *gorm.DB) error {
	sql := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		id SERIAL PRIMARY KEY,
		version VARCHAR(255) NOT NULL UNIQUE,
		applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	`
	return db.Exec(sql).Error
}

// applyMigration applies a single migration if it hasn't been applied yet
func applyMigration(db *gorm.DB, fileName string) error {
	// Check if migration has already been applied
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", fileName).Scan(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		fmt.Printf("Migration %s already applied, skipping\n", fileName)
		return nil
	}

	// Read migration file
	content, err := migrationFiles.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
	}

	// Execute migration
	if err := db.Exec(string(content)).Error; err != nil {
		return fmt.Errorf("failed to execute migration %s: %w", fileName, err)
	}

	// Record migration as applied
	if err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", fileName).Error; err != nil {
		return fmt.Errorf("failed to record migration %s: %w", fileName, err)
	}

	fmt.Printf("Successfully applied migration: %s\n", fileName)
	return nil
}

// RollbackMigration rolls back a specific migration (for development)
func RollbackMigration(db *gorm.DB, fileName string) error {
	// This is a simplified rollback - in production, you'd want more sophisticated rollback logic
	if err := db.Exec("DELETE FROM schema_migrations WHERE version = ?", fileName).Error; err != nil {
		return fmt.Errorf("failed to rollback migration %s: %w", fileName, err)
	}

	fmt.Printf("Rolled back migration: %s\n", fileName)
	return nil
}

// GetAppliedMigrations returns a list of applied migrations
func GetAppliedMigrations(db *gorm.DB) ([]string, error) {
	var migrations []string
	if err := db.Raw("SELECT version FROM schema_migrations ORDER BY applied_at").Scan(&migrations).Error; err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}
	return migrations, nil
}
