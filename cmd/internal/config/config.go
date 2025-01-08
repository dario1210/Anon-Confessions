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

var (
	defaultFileName       = "anon_confessions142.db?_foreign_keys=1"
	defaultDBURL          = "sqlite3://" + defaultFileName
	defaultMigrationsPath = "file://./cmd/internal/db/migrations_files"
	defaultPort           = "9000"
)

// LoadConfig loads the application configuration from environment variables.
// It initializes default values for configuration parameters if they are not
// present in the environment.
func LoadConfig() *Config {

	cfg := &Config{
		Port: getEnv("PORT", defaultPort),
		DB: SQLiteConfig{
			File: getEnv("DB_FILE", defaultFileName),
		},
		Migrations: Migrations{
			MigrationPath: getEnv("MIGRATIONS_PATH", defaultMigrationsPath),
			DBURL:         getEnv("DB_URL", defaultDBURL),
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
