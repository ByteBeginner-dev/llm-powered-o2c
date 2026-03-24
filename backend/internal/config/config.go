package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	GroqApiKey  string
	DataDir     string
	Port        string
}

// Load reads all values strictly from environment variables.
func Load() *Config {
	return &Config{
		// Constructed from individual DB_ vars if DATABASE_URL not set directly
		DatabaseURL: buildDatabaseURL(),
		GroqApiKey:  mustGetEnv("GroqAPIKey"),
		DataDir:     mustGetEnv("DataDir"),
		Port:        mustGetEnv("PORT"),
	}
}

// buildDatabaseURL uses DATABASE_URL if set, otherwise builds from DB_* vars
func buildDatabaseURL() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		url = stripPsqlWrapper(url)
		return url
	}

	// Fallback:
	host := mustGetEnv("DB_HOST")
	port := mustGetEnv("DB_PORT")
	user := mustGetEnv("DB_USER")
	password := mustGetEnv("DB_PASSWORD")
	dbname := mustGetEnv("DB_NAME")
	sslmode := mustGetEnv("DB_SSL_MODE")

	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)
}

func stripPsqlWrapper(raw string) string {
	// Remove leading `psql ` prefix
	if len(raw) > 5 && raw[:5] == "psql " {
		raw = raw[5:]
	}
	// Remove surrounding quotes if present
	if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	return raw
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("REQUIRED environment variable %q is not set. Add it to your .env file.", key))
	}
	return val
}
