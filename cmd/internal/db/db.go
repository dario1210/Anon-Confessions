package db

import (
	"anon-confessions/cmd/internal/config"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DbConnection initializes and returns a new database connection.
func DbConnection(dbName string) (*gorm.DB, error) {
	slog.Info("Initializing database connection...", "dbName", dbName)
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to get raw database connection", "error", err)
		return nil, fmt.Errorf("failed to get raw database connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get *sql.DB from gorm.DB", "error", err)
		return nil, fmt.Errorf("failed to get *sql.DB from gorm.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		slog.Error("Database ping failed", "error", err)
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	// Configure the connection pool.
	sqlDB.SetMaxOpenConns(10)               // Limit max open connections
	sqlDB.SetMaxIdleConns(5)                // Limit idle connections
	sqlDB.SetConnMaxLifetime(1 * time.Hour) // Set connection lifetime

	slog.Info("Database connection established successfully", "dbName", dbName)
	return db, nil
}

// RunMigrations runs the migrations using golang-migrate.
func RunMigrations(cfg *config.Migrations) error {
	if cfg == nil {
		slog.Error("Invalid migration configuration: config is nil")
		return fmt.Errorf("invalid migration configuration: config is nil")
	}

	slog.Info("Initializing migrations...", "MigrationPath", cfg.MigrationPath, "DBURL", cfg.DBURL)
	m, err := migrate.New("file://cmd/internal/db/migrations", cfg.DBURL)
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
