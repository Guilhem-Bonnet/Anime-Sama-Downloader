package sqlite

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

// TestInitializeDefaults tests the default settings initialization.
func TestInitializeDefaults(t *testing.T) {
	db := getMemoryDB(t)
	defer db.Close()

	// Create tables first (normally done by migrations)
	if err := createTestSettingsTable(db); err != nil {
		t.Fatalf("Failed to create test settings table: %v", err)
	}

	// Initialize defaults
	if err := InitializeDefaults(db); err != nil {
		t.Fatalf("InitializeDefaults failed: %v", err)
	}

	// Verify all default settings exist
	expectedSettings := map[string]string{
		"download_path": "/home/user/Downloads/anime",
		"quality":       "720p",
		"hls_enabled":   "true",
	}

	for key, expectedValue := range expectedSettings {
		var value string
		err := db.QueryRow("SELECT value_json FROM settings WHERE key = ?", key).Scan(&value)
		if err == sql.ErrNoRows {
			t.Errorf("Setting %s not found", key)
		} else if err != nil {
			t.Errorf("Failed to query setting %s: %v", key, err)
		} else {
			// Remove JSON quotes from value
			jsonValue := string([]byte(value)[1 : len(value)-1])
			if jsonValue != expectedValue {
				t.Errorf("Setting %s has incorrect value: got %s, want %s", key, jsonValue, expectedValue)
			}
		}
	}
}

// TestInitializeDefaultsIdempotent tests that InitializeDefaults is idempotent.
func TestInitializeDefaultsIdempotent(t *testing.T) {
	db := getMemoryDB(t)
	defer db.Close()

	// Create tables first
	if err := createTestSettingsTable(db); err != nil {
		t.Fatalf("Failed to create test settings table: %v", err)
	}

	// First initialization
	if err := InitializeDefaults(db); err != nil {
		t.Fatalf("First InitializeDefaults failed: %v", err)
	}

	// Second initialization (should not fail)
	if err := InitializeDefaults(db); err != nil {
		t.Fatalf("Second InitializeDefaults failed: %v", err)
	}

	// Verify count of settings (should be 3, not 6)
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM settings").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count settings: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 settings, got %d", count)
	}
}

// TestVerifyDatabaseIntegrity tests database integrity check.
func TestVerifyDatabaseIntegrity(t *testing.T) {
	db := getMemoryDB(t)
	defer db.Close()

	if err := createTestSettingsTable(db); err != nil {
		t.Fatalf("Failed to create settings table: %v", err)
	}

	if err := createTestJobsTable(db); err != nil {
		t.Fatalf("Failed to create jobs table: %v", err)
	}

	// Integrity check should pass
	if err := verifyDatabaseIntegrity(db); err != nil {
		t.Fatalf("verifyDatabaseIntegrity failed: %v", err)
	}
}

// TestFullInitializationFlowWithMigrations tests complete initialization in a real database context.
func TestFullInitializationFlowWithMigrations(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "full_test.db")

	// Open database with migrations (this will create the schema)
	db, err := Open(context.Background(), dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize defaults
	if err := InitializeDefaults(db.SQL); err != nil {
		t.Fatalf("InitializeDefaults failed: %v", err)
	}

	// Verify defaults were inserted
	var count int
	err = db.SQL.QueryRow("SELECT COUNT(*) FROM settings").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count settings: %v", err)
	}
	if count < 3 {
		t.Errorf("Expected at least 3 default settings, got %d", count)
	}

	// Verify integrity check passes
	if err := verifyDatabaseIntegrity(db.SQL); err != nil {
		t.Fatalf("verifyDatabaseIntegrity failed: %v", err)
	}
}

// TestInsertDefaultSettingsChecksExisting tests that existing settings are not duplicated.
func TestInsertDefaultSettingsChecksExisting(t *testing.T) {
	db := getMemoryDB(t)
	defer db.Close()

	if err := createTestSettingsTable(db); err != nil {
		t.Fatalf("Failed to create settings table: %v", err)
	}

	// Pre-insert one setting manually
	_, err := db.Exec("INSERT INTO settings (key, value_json, updated_at) VALUES (?, json(?), datetime('now'))", "quality", `"1080p"`)
	if err != nil {
		t.Fatalf("Failed to insert pre-existing setting: %v", err)
	}

	// Now initialize defaults
	if err := insertDefaultSettings(db); err != nil {
		t.Fatalf("insertDefaultSettings failed: %v", err)
	}

	// Verify that the pre-inserted value wasn't overwritten
	var value string
	err = db.QueryRow("SELECT value_json FROM settings WHERE key = ?", "quality").Scan(&value)
	if err != nil {
		t.Fatalf("Failed to query setting: %v", err)
	}

	// The value should still be "1080p", not "720p"
	if value != `"1080p"` {
		t.Errorf("Pre-existing setting was overwritten: got %s, want %s", value, `"1080p"`)
	}
}

// Test database schema structure (tests against actual migration schema)
// TestSettingsTableSchema verifies that the settings table has the correct structure.
func TestSettingsTableSchema(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "schema_test.db")

	db, err := Open(context.Background(), dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Query table info
	rows, err := db.SQL.Query("PRAGMA table_info(settings)")
	if err != nil {
		t.Fatalf("Failed to query table info: %v", err)
	}
	defer rows.Close()

	columns := make(map[string]bool)
	for rows.Next() {
		var cid, notnull, pk interface{}
		var name, type_ string
		var dfltValue interface{}

		if err := rows.Scan(&cid, &name, &type_, &notnull, &dfltValue, &pk); err != nil {
			t.Fatalf("Failed to scan column info: %v", err)
		}
		columns[name] = true
	}

	// Verify required columns exist (from migrations)
	expectedColumns := []string{"key", "value_json", "updated_at"}
	for _, col := range expectedColumns {
		if !columns[col] {
			t.Errorf("Column %s not found in settings table", col)
		}
	}
}

// TestJobsTableSchema verifies that the jobs table has the correct structure.
func TestJobsTableSchema(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "schema_test.db")

	db, err := Open(context.Background(), dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Query table info
	rows, err := db.SQL.Query("PRAGMA table_info(jobs)")
	if err != nil {
		t.Fatalf("Failed to query table info: %v", err)
	}
	defer rows.Close()

	columns := make(map[string]bool)
	for rows.Next() {
		var cid, notnull, pk interface{}
		var name, type_ string
		var dfltValue interface{}

		if err := rows.Scan(&cid, &name, &type_, &notnull, &dfltValue, &pk); err != nil {
			t.Fatalf("Failed to scan column info: %v", err)
		}
		columns[name] = true
	}

	// Verify required columns exist (from migrations)
	expectedColumns := []string{"id", "type", "state", "created_at", "updated_at"}
	for _, col := range expectedColumns {
		if !columns[col] {
			t.Errorf("Column %s not found in jobs table", col)
		}
	}
}

// Helper functions for testing

// createTestSettingsTable creates a test settings table matching the migration schema.
func createTestSettingsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value_json BLOB NOT NULL,
		updated_at TEXT NOT NULL
	);
	`
	_, err := db.Exec(query)
	return err
}

// createTestJobsTable creates a test jobs table matching the migration schema.
func createTestJobsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS jobs (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		state TEXT NOT NULL,
		progress REAL NOT NULL DEFAULT 0,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		params_json BLOB,
		result_json BLOB,
		error_code TEXT,
		error_message TEXT
	);
	`
	_, err := db.Exec(query)
	return err
}

// getMemoryDB returns an in-memory SQLite database for testing.
func getMemoryDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create memory database: %v", err)
	}
	return db
}
