package db

import (
	"log"

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
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}

	return db, nil
}

// RunMigrations runs the migrations using golang-migrate.
// TODO: MAKE DBURL AND MIGRATIONSPATH BE READ FROM ENV, MAKE DB URL RELATIVE PATH.
func RunMigrations() {
	// Path to migrations folder
	migrationsPath := "file://cmd/internal/db/migrations"
	// SQLite database connection URL
	databaseURL := "sqlite3://C:/Users/User/Desktop/Anon-Repository/test.db"

	// Initialize migrate
	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")
}
