package testutils

import (
	"testing"
)

func TestSetupMockDB(t *testing.T) {
	// Initialize the mock database
	db := SetupMockDB()

	// Query the list of tables in the database
	var tableNames []string
	err := db.Raw("SELECT name FROM sqlite_master WHERE type='table'").Scan(&tableNames).Error
	if err != nil {
		t.Fatalf("Failed to query tables: %v", err)
	}

	// Ensure the tables created by migrations exist
	if len(tableNames) == 0 {
		t.Fatalf("Expected tables to be created, but none were found")
	}

	t.Logf("Tables in the database: %v", tableNames)
}
