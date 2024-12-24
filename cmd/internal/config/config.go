package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type SQLiteConfig struct {
	File string
}

type Migrations struct {
	MigrationPath string
	DBURL         string
}

type Config struct {
	Port       string
	DB         SQLiteConfig
	Migrations Migrations
}

// LoadConfig loads the application configuration from environment variables.
// It initializes default values for configuration parameters if they are not
// present in the environment.
//
// Returns:
// - *Config: A pointer to the populated Config struct.
func LoadConfig() *Config {
	cfg := &Config{
		Port: getEnv("PORT", "9000"),
		DB: SQLiteConfig{
			File: getEnv("DBFILE", "anon_confessions.db"),
		},
		Migrations: Migrations{
			MigrationPath: getEnv("MIGRATIONS_PATH", "file://cmd/internal/db/migrations"),
			DBURL:         getEnv("DB_URL", "sqlite3://../../anon_confessions.db"),
		},
	}

	return cfg
}

// getEnv retrieves the value of the environment variable named by the key.
// If the variable is not present, it returns the default value provided.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
