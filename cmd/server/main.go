package main

import (
	"anon-confessions/cmd/internal/app"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

// @title           Anonymous Confessions API
// @version         1.0
// @description     A privacy-focused backend service that allows users to:
// @description     • Post and manage anonymous confessions.
// @description     • React to posts with likes and comments.
// @description     • Leave comments on confessions.
// @description     • Receive real-time updates through WebSocket.
// @description
// @description     The API is designed with RESTful principles, uses SQLite for data storage, and ensures anonymity without storing personal information.

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:9000
// @BasePath        /api/v1

// @securityDefinitions.apikey AccountNumberAuth
// @in header
// @name X-Account-Number
// @description A unique account number for user authentication.
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
