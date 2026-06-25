package config

import (
	"os"
)

// Config holds global environmental values for Lordis
type Config struct {
	Port         string
	BrevoAPIKey  string
	SessionSecret string
}

// LoadConfig fetches settings from the operating host environment
func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	brevoKey := os.Getenv("BREVO_API_KEY")

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		// Fallback safe value for local environment execution
		sessionSecret = "lordis-default-super-secret-key-string-54321"
	}

	return &Config{
		Port:          port,
		BrevoAPIKey:   brevoKey,
		SessionSecret: sessionSecret,
	}
}