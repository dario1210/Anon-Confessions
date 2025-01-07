package testutils

import (
	"anon-confessions/cmd/internal/models"
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupMockDB initializes an in-memory SQLite database and applies migrations
func SetupMockDB() *gorm.DB {
	// Open an in-memory SQLite database with shared cache
	dsn := "file::memory:?cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to in-memory SQLite database: %v", err)
	}

	// Configure the connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get raw SQL DB from GORM: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Dynamically determine the path to the migration files.
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	migrationPath := filepath.Join(basepath, "../../db/migrations_files")
	migrationPath = "file://" + strings.ReplaceAll(migrationPath, "\\", "/")

	// Run migrations
	err = runMigrations(migrationPath, "sqlite3://"+dsn)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

// runMigrations applies migration files to the given database
func runMigrations(migrationPath, dbURL string) error {
	m, err := migrate.New(migrationPath, dbURL)
	if err != nil {
		return err
	}

	// Apply migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// HTTPTestRequest gets a method and URL and a request body and returns a recorder and request.
func HTTPTestRequest(method string, url string, body []byte) (*httptest.ResponseRecorder, *http.Request) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	return w, req
}

// SeedPost creates a new post in the memory database
func SeedPost() {
	db := SetupMockDB()
	db.Create(&models.PostDBModel{
		Content: "Seeded Post",
		UserId:  1,
	})
}

// SeedComment creates a new comment in the memory database.
func SeedComment(postID int, content string) {
	db := SetupMockDB()
	db.Create(&models.CommentsDbModel{
		Content: content,
		UserId:  1,
		PostId:  postID,
	})
}
