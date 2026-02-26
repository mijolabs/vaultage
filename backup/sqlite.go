package backup

import (
	"context"
	"database/sql"
	"fmt"

	"modernc.org/sqlite"
)

// sqliteConn defines the modernc.org/sqlite driver connection methods we need
// for performing safe database backups using the SQLite Online Backup API.
type sqliteConn interface {
	// NewRestore creates a backup that copies from srcPath into this connection.
	NewRestore(srcPath string) (*sqlite.Backup, error)
	// Serialize returns the database contents as a byte slice.
	Serialize() ([]byte, error)
}

// BackupToMemory performs a safe SQLite backup of the database at dbPath
// using the SQLite Online Backup API. This is the officially recommended
// method for backing up a live SQLite database.
//
// The backup process:
//  1. Opens an in-memory database as the destination
//  2. Uses sqlite3_backup to copy all pages from the source database
//  3. Serializes the in-memory database to a byte slice
//
// This approach safely handles WAL mode databases and provides a consistent
// snapshot even if the source database is actively being written to.
func BackupToMemory(dbPath string) ([]byte, error) {
	// Open in-memory database as the backup destination
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("opening memory database: %w", err)
	}
	defer db.Close()

	// Get a dedicated connection for the backup operation
	conn, err := db.Conn(context.Background())
	if err != nil {
		return nil, fmt.Errorf("acquiring connection: %w", err)
	}
	defer conn.Close()

	var data []byte
	err = conn.Raw(func(dc any) error {
		c, ok := dc.(sqliteConn)
		if !ok {
			return fmt.Errorf("unexpected driver type: %T (expected modernc.org/sqlite)", dc)
		}

		// NewRestore copies FROM the source path INTO this connection
		backup, err := c.NewRestore(dbPath)
		if err != nil {
			return fmt.Errorf("initializing backup: %w", err)
		}

		// Copy all pages in one step (-1 means copy everything)
		// This is safe for small databases; for very large databases,
		// you could use positive values to copy incrementally.
		if _, err := backup.Step(-1); err != nil {
			backup.Finish()
			return fmt.Errorf("backup step: %w", err)
		}

		if err := backup.Finish(); err != nil {
			return fmt.Errorf("finishing backup: %w", err)
		}

		// Serialize the in-memory database to bytes
		data, err = c.Serialize()
		if err != nil {
			return fmt.Errorf("serializing backup: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return data, nil
}
