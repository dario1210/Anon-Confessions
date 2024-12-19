package main

import (
	"anon-confessions/router.go"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	const file string = "test.db"

	log.Println("Starting the application...")
	log.Println("Initializing database connection...")
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Pinging the database to ensure connectivity...")
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection successful!")

	r := router.NewRouter()
	r.RegisterRoutes()

	log.Println("Starting the HTTP server on :9000...")
	if err := r.Run(":9000"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
