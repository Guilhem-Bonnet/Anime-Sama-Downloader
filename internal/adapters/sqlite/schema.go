package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// InitializeDefaults ensures the database has default settings initialized.
// This function is idempotent and safe to call multiple times.
// Note: Table creation is handled by migrations (internal/adapters/sqlite/migrations/)
func InitializeDefaults(db *sql.DB) error {
	// Verify database integrity
	if err := verifyDatabaseIntegrity(db); err != nil {
		return err
	}

	// Insert default settings if not already present
	if err := insertDefaultSettings(db); err != nil {
		return err
	}

	return nil
}

// insertDefaultSettings inserts default application settings if they don't already exist.
// Default settings include download path, video quality, and HLS enablement.
func insertDefaultSettings(db *sql.DB) error {
	defaults := map[string]string{
		"download_path": "/home/user/Downloads/anime",
		"quality":       "720p",
		"hls_enabled":   "true",
	}

	for key, value := range defaults {
		// Check if setting already exists
		var exists bool
		err := db.QueryRow("SELECT 1 FROM settings WHERE key = ?", key).Scan(&exists)
		if err == sql.ErrNoRows {
			// Setting doesn't exist, insert it
			// Settings table stores value as JSON BLOB according to migrations
			_, insertErr := db.Exec("INSERT INTO settings (key, value_json, updated_at) VALUES (?, json(?), datetime('now'))", key, `"`+value+`"`)
			if insertErr != nil {
				return fmt.Errorf("failed to insert default setting %s: %w", key, insertErr)
			}
		} else if err != nil {
			return fmt.Errorf("failed to check setting %s: %w", key, err)
		}
		// If exists, do nothing (idempotent)
	}

	return nil
}

// verifyDatabaseIntegrity checks the database integrity using PRAGMA integrity_check.
// Returns nil if database is valid, error otherwise.
func verifyDatabaseIntegrity(db *sql.DB) error {
	var result string
	err := db.QueryRow("PRAGMA integrity_check;").Scan(&result)
	if err != nil {
		return fmt.Errorf("failed to run integrity check: %w", err)
	}

	if result != "ok" {
		return fmt.Errorf("database integrity check failed: %s", result)
	}

	return nil
}
