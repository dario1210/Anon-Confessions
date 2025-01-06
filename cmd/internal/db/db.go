package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DbConnection initializes and returns a new database connection.
func DbConnection(dbName string) (*gorm.DB, error) {
	slog.Info("Initializing database connection...", "dbName", dbName)

	// Since we are using sqlite, we do not need to create the database explicitly.
	// If it does not exist, sqllite will create the db.
	// Careful any other type of database we need to be sure that it exits before attempting to connect.
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to open database connection", "error", err)
		return nil, fmt.Errorf("failed to open database connection: %w", err)
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
