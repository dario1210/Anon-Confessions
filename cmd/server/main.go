package main

import (
	"anon-confessions/cmd/internal/app"
	"anon-confessions/cmd/internal/config"
	"log/slog"
)

func main() {

	slog.Info("Initializing API server...")

	slog.Info("Loading configuration...")
	cfg := config.LoadConfig()

	// Initialize the application
	app, err := app.NewApp(cfg)
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
