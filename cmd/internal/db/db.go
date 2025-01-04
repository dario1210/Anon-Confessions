package db

import (
	"anon-confessions/cmd/internal/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DbConnection initializes and returns a new database connection.
func DbConnection(dbName string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to get raw database connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get *sql.DB from gorm.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	// Configure the connection pool.
	sqlDB.SetMaxOpenConns(10)               // Limit max open connections
	sqlDB.SetMaxIdleConns(5)                // Limit idle connections
	sqlDB.SetConnMaxLifetime(1 * time.Hour) // Set connection lifetime

	return db, nil
}

// RunMigrations runs the migrations using golang-migrate.
func RunMigrations(cfg *config.Migrations) error {
	if cfg == nil {
		return fmt.Errorf("invalid migration configuration: config is nil")
	}

	log.Println("Initializing migrations...")
	m, err := migrate.New(cfg.MigrationPath, cfg.DBURL)
	if err != nil {
		log.Printf("Failed to initialize migrations. MigrationPath: %s, DBURL: %s/n", cfg.MigrationPath, cfg.DBURL)
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	log.Println("Running migrations...")
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply.")
		} else {
			log.Printf("Migration failed: %v\n", err)
			return fmt.Errorf("migration failed: %w", err)
		}
	} else {
		log.Println("Migrations applied successfully!")
	}

	return nil
}
