package config

import (
	"os"
)

type Config struct {
	MongoURI string
	DBName   string
}

func NewConfigFromEnv() *Config {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "lordis"
	}

	return &Config{
		MongoURI: uri,
		DBName:   dbName,
	}
}
