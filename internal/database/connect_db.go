package database

import "os"

type Config struct {
	DatabaseURL   string
	MigrationPath string
}

func NewConfigFromEnv() Config {
	migrationPath := os.Getenv("MIGRATION_PATH")
	if migrationPath == "" {
		migrationPath = "migration/migrations.sql"
	}

	return Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		MigrationPath: migrationPath,
	}
}

func (c Config) HasDatabaseURL() bool {
	return c.DatabaseURL != ""
}
