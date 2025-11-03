package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(256) NOT NULL,
	comment TEXT,
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

var db *sql.DB

// Init initializes the database connection.
// If the database file does not exist, it creates it and sets up the schema.
//
// Possible errors:
//   - Wrong driver name
//   - Driver initialization problems
//   - Database file is locked by another process
//   - Incorrect database file format
//   - Invalid path to database file
//   - Insufficient permissions to create database file
//   - Insufficient disk space
func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := false
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection is valid
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create schema if this is a new database
	if install {
		if _, err = db.Exec(schema); err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
	}

	return nil
}

// GetDB returns the active database connection
func GetDB() *sql.DB {
	return db
}

// Close closes the database connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
