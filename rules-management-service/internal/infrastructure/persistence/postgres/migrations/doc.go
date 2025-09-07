// This file is intentionally left blank.
// It is used to ensure that the Go toolchain correctly identifies this directory as a package.
package migrations

import (
	"gorm.io/gorm"
)

// ApplyMigrations applies database migrations
func ApplyMigrations(db *gorm.DB) error {
	// Auto-migrate the schema
	return db.AutoMigrate()
}
