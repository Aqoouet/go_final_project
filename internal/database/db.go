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

// Init initialize connection to db
// if db file does not exist, create it and create table inside
func Init(dbFile string) error {

	_, err := os.Stat(dbFile)
	install := false
	if err != nil {
		install = true
	}

	// can catch:
	//	- wrong driver name
	//	- problems with driver
	db, err = sql.Open ("sqlite", dbFile)
	if err!=nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// can catch:
	//	- DBFile is blocked by another process
	//	- incorrect (wrong format) or empty dbFile
	//	- incorrect path to dbFile
	//	- not enough permissions to create dbFile
	//	- not enough disc space
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	if install {
		if _, err = db.Exec(schema); err!=nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}

func Close() error {
	if db!= nil {
		return db.Close()
	}
	return nil
}