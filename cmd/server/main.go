package main

import (
	"anon-confessions/cmd/internal/app"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)


func main() {

	slog.Info("Initializing API server...")

	// Initialize the application
	app, err := app.NewApp()
	if err != nil {
		slog.Error("Could not initialize the application", slog.String("error", err.Error()))
		return
	}

	// Run the application
	err = app.Run()
	if err != nil {
		slog.Error("Application runtime error", slog.String("error", err.Error()))
	}
}
