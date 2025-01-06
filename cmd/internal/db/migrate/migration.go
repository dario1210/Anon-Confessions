package main

import (
	"anon-confessions/cmd/internal/config"
	"anon-confessions/cmd/internal/db"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
)

// This file loads the configuration, initializes a database connection, and runs the migrations.
// It is designed to be run as a standalone command to migrate the database.
// All migrations in the specified path are applied, ensuring the database schema is up-to-date.
// If any errors occur during the process, the program will terminate and log the error.
func main() {

	cfg := config.LoadConfig()

	_, err := db.DbConnection(cfg.DB.File)
	if err != nil {
		slog.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = RunMigrations(&cfg.Migrations)
	if err != nil {
		slog.Error("Failed to run migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

// RunMigrations runs the migrations using golang-migrate.
func RunMigrations(cfg *config.Migrations) error {
	if cfg == nil {
		slog.Error("Invalid migration configuration: config is nil")
		return fmt.Errorf("invalid migration configuration: config is nil")

	}

	slog.Info("Initializing migrations...", "MigrationPath", cfg.MigrationPath, "DBURL", cfg.DBURL)
	m, err := migrate.New(cfg.MigrationPath, cfg.DBURL)
	if err != nil {
		slog.Error("Failed to initialize migrations", "MigrationPath", cfg.MigrationPath, "DBURL", cfg.DBURL, "error", err)
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	slog.Info("Running migrations...")
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			slog.Info("No new migrations to apply.")
		} else {
			slog.Error("Migration failed", "error", err)
			return fmt.Errorf("migration failed: %w", err)
		}
	}
	slog.Info("Migrations applied successfully!")

	return nil
}
