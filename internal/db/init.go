package db

import (
	"context"
	"database/sql"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	"traveler/pkg/log"
)

// Init opens (or creates) an SQLite database at dbPath and applies the schema
// from schemaPath. It returns an opened *sql.DB ready for use.
func Init(ctx context.Context, dbPath, schemaPath string) (*sql.DB, error) {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, err
	}

	// Open a database using modernc.org/sqlite (pure Go)
	dsn := "file:" + dbPath + "?_busy_timeout=5000&_pragma=journal_mode(WAL)"
	dbComm, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	// Set reasonable connection parameters for embedded SQLite
	dbComm.SetConnMaxLifetime(0)
	dbComm.SetMaxOpenConns(1)
	dbComm.SetMaxIdleConns(1)

	// Verify connection
	ctxPing, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := dbComm.PingContext(ctxPing); err != nil {
		_ = dbComm.Close()
		return nil, err
	}

	// Read and apply schema
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		// If a schema file is missing, close and return the error
		_ = dbComm.Close()
		// Enrich error with path existence info
		if _, statErr := os.Stat(schemaPath); errors.Is(statErr, fs.ErrNotExist) {
			return nil, statErr
		}
		return nil, err
	}

	if len(schemaSQL) > 0 {
		tx, err := dbComm.BeginTx(ctx, nil)
		if err != nil {
			_ = dbComm.Close()
			return nil, err
		}
		if _, err := tx.ExecContext(ctx, string(schemaSQL)); err != nil {
			_ = tx.Rollback()
			_ = dbComm.Close()
			return nil, err
		}
		if err := tx.Commit(); err != nil {
			_ = dbComm.Close()
			return nil, err
		}
		log.Info("database schema applied", "path", schemaPath)
	}

	return dbComm, nil
}
