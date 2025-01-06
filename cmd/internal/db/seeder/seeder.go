package main

import (
	"anon-confessions/cmd/internal/config"
	"anon-confessions/cmd/internal/db"
	"log/slog"
	"os"
)

// The seeder will run a sql file to seed the database with some initial data.
// The seeder is optional; if the SQL file is not present, the seeder will log an error and terminate.
func main() {

	cfg := config.LoadConfig()

	c, ioErr := os.ReadFile("cmd/internal/db/seeder/seeder.sql")
	if ioErr != nil {
		slog.Error("Failed to read seeder file", slog.String("error", ioErr.Error()))
		os.Exit(1)
	}

	sqlStmt := string(c)
	dbConn, err := db.DbConnection(cfg.DB.File)

	if err != nil {
		slog.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	result := dbConn.Exec(sqlStmt)
	if result.Error != nil {
		slog.Error("Failed to seed database", slog.String("error", result.Error.Error()))
		os.Exit(1)
	}
}
