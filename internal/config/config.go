// Package config loads application configuration from environment variables with sensible defaults.
package config

import (
	"os"
	"strconv"
)

const (
	DefaultPort = 7540
	DefaultWebDir = "./web"
	DefaultDBFile = "scheduler.db"
)

type Config struct {
    Port     int
    WebDir   string
    DBFile   string
    Password string
}

func LoadConfig() *Config {
	cfg := &Config{
		Port:   DefaultPort,
		WebDir: DefaultWebDir,
		DBFile: DefaultDBFile,
	}

	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.Port = port
		}
	}

	if dbFile := os.Getenv("TODO_DBFILE"); dbFile != "" {
		cfg.DBFile = dbFile
	}

	if password := os.Getenv("TODO_PASSWORD"); password != "" {
		cfg.Password = password
	}

	return cfg
}
