package backup

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// BackupToMemory performs a safe SQLite backup of the database at dbPath
// to an in-memory byte slice using the SQLite serialization API.
// This ensures a consistent snapshot even if the database is in WAL mode
// and actively being written to by another process.
func BackupToMemory(dbPath string) ([]byte, error) {
	// Open source database in read-only mode
	db, err := sql.Open("sqlite", dbPath+"?mode=ro")
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	// Verify connection works
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	// Get a dedicated connection for serialization
	conn, err := db.Conn(context.Background())
	if err != nil {
		return nil, fmt.Errorf("acquiring connection: %w", err)
	}
	defer conn.Close()

	// Serialize the database to a byte slice using the driver's
	// Serialize method, which captures a consistent snapshot
	var data []byte
	err = conn.Raw(func(dc any) error {
		serializer, ok := dc.(interface {
			Serialize() ([]byte, error)
		})
		if !ok {
			return fmt.Errorf("sqlite driver does not support Serialize")
		}

		var serErr error
		data, serErr = serializer.Serialize()
		return serErr
	})

	if err != nil {
		return nil, fmt.Errorf("serializing database: %w", err)
	}

	return data, nil
}
