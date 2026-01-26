package backup

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestBackupToMemory(t *testing.T) {
	// Create a temporary directory for the test database
	tmpDir, err := os.MkdirTemp("", "vaultage-test-*")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")

	// Create a test database with some data
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}

	// Create table and insert test data
	_, err = db.Exec(`
		CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT);
		INSERT INTO users (name) VALUES ('alice'), ('bob'), ('charlie');
	`)
	if err != nil {
		t.Fatalf("creating test data: %v", err)
	}

	// Close the database to ensure all data is flushed
	if err := db.Close(); err != nil {
		t.Fatalf("closing database: %v", err)
	}

	// Perform backup
	data, err := BackupToMemory(dbPath)
	if err != nil {
		t.Fatalf("BackupToMemory: %v", err)
	}

	// Verify backup is not empty
	if len(data) == 0 {
		t.Fatal("backup data is empty")
	}

	// Verify backup is a valid SQLite database by checking the header
	// SQLite databases start with "SQLite format 3\000"
	header := "SQLite format 3\x00"
	if len(data) < len(header) || string(data[:len(header)]) != header {
		t.Fatalf("backup does not have valid SQLite header")
	}

	t.Logf("Backup successful: %d bytes", len(data))
}

func TestBackupToMemory_WALMode(t *testing.T) {
	// Create a temporary directory for the test database
	tmpDir, err := os.MkdirTemp("", "vaultage-test-wal-*")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")

	// Create a test database in WAL mode
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}

	// Enable WAL mode
	_, err = db.Exec(`PRAGMA journal_mode=WAL`)
	if err != nil {
		t.Fatalf("enabling WAL mode: %v", err)
	}

	// Create table and insert test data
	_, err = db.Exec(`
		CREATE TABLE items (id INTEGER PRIMARY KEY, value TEXT);
		INSERT INTO items (value) VALUES ('one'), ('two'), ('three');
	`)
	if err != nil {
		t.Fatalf("creating test data: %v", err)
	}

	// Keep the database open to simulate active usage
	// The backup should still work correctly

	// Perform backup while database is open
	data, err := BackupToMemory(dbPath)
	if err != nil {
		t.Fatalf("BackupToMemory with WAL: %v", err)
	}

	// Close original database
	if err := db.Close(); err != nil {
		t.Fatalf("closing database: %v", err)
	}

	// Verify backup is valid
	if len(data) == 0 {
		t.Fatal("backup data is empty")
	}

	// Write backup to a new file and verify contents
	backupPath := filepath.Join(tmpDir, "backup.db")
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		t.Fatalf("writing backup file: %v", err)
	}

	// Open backup and verify data
	backupDB, err := sql.Open("sqlite", backupPath)
	if err != nil {
		t.Fatalf("opening backup database: %v", err)
	}
	defer backupDB.Close()

	var count int
	err = backupDB.QueryRow("SELECT COUNT(*) FROM items").Scan(&count)
	if err != nil {
		t.Fatalf("querying backup: %v", err)
	}

	if count != 3 {
		t.Fatalf("expected 3 items, got %d", count)
	}

	t.Logf("WAL mode backup successful: %d bytes, %d items", len(data), count)
}

func TestBackupToMemory_NonExistentFile(t *testing.T) {
	_, err := BackupToMemory("/nonexistent/path/to/database.db")
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
	t.Logf("Got expected error: %v", err)
}
