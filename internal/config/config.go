package config

import (
	"os"
	"strconv"
)

const (
	// DefaultPort is the default server port
	DefaultPort = 7540
	
	// DefaultWebDir is the default directory for web assets
	DefaultWebDir = "./web"
	
	// DefaultDBFile is the default database file name
	DefaultDBFile = "scheduler.db"
)

// Config holds application configuration
type Config struct {
	Port   int    // Server port
	WebDir string // Directory for web assets
	DBFile string // Path to database file
}

// LoadConfig loads configuration from environment variables
// Supported environment variables:
//   - TODO_PORT: server port (default: 7540)
//   - TODO_DBFILE: database file path (default: scheduler.db)
func LoadConfig() *Config {
	cfg := &Config{
		Port:   DefaultPort,
		WebDir: DefaultWebDir,
		DBFile: DefaultDBFile,
	}

	// Load port from environment
	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.Port = port
		}
	}

	// Load database file from environment
	if dbFile := os.Getenv("TODO_DBFILE"); dbFile != "" {
		cfg.DBFile = dbFile
	}

	return cfg
}
