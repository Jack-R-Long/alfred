package database

import (
	"database/sql"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	// Create a temporary test database file
	testDBFile := "./test_alfred.db"
	
	// Clean up any existing test database
	os.Remove(testDBFile)
	
	// Temporarily replace the database file path for testing
	originalDBPath := "./alfred-database.db"
	
	// Open test database
	var err error
	DB, err = sql.Open("sqlite3", testDBFile)
	if err != nil {
		t.Fatal("Failed to open test database:", err)
	}
	defer func() {
		DB.Close()
		os.Remove(testDBFile)
	}()

	if err = DB.Ping(); err != nil {
		t.Fatal("Failed to ping test database:", err)
	}

	// Test table creation
	createTables()

	// Verify users table exists and has correct structure
	rows, err := DB.Query("PRAGMA table_info(users)")
	if err != nil {
		t.Fatal("Failed to query table info:", err)
	}
	defer rows.Close()

	columnCount := 0
	expectedColumns := map[string]bool{
		"id":            false,
		"username":      false,
		"email":         false,
		"password_hash": false,
		"created_at":    false,
	}

	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, pk int
		var defaultValue sql.NullString

		err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
		if err != nil {
			t.Fatal("Failed to scan column info:", err)
		}

		if _, exists := expectedColumns[name]; exists {
			expectedColumns[name] = true
			columnCount++
		}
	}

	if columnCount != len(expectedColumns) {
		t.Errorf("Expected %d columns, found %d", len(expectedColumns), columnCount)
	}

	for col, found := range expectedColumns {
		if !found {
			t.Errorf("Expected column '%s' not found", col)
		}
	}

	// Test basic database operations
	testInsert := `INSERT INTO users (username, email, password_hash) VALUES ('testuser', 'test@example.com', 'hashedpassword')`
	_, err = DB.Exec(testInsert)
	if err != nil {
		t.Fatal("Failed to insert test data:", err)
	}

	// Test query
	var username, email string
	err = DB.QueryRow("SELECT username, email FROM users WHERE username = 'testuser'").Scan(&username, &email)
	if err != nil {
		t.Fatal("Failed to query test data:", err)
	}

	if username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", username)
	}

	if email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", email)
	}

	// Restore original database path (though this is a unit test, so it doesn't matter much)
	_ = originalDBPath
}

func TestClose(t *testing.T) {
	// Create a test database connection
	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal("Failed to create test database:", err)
	}

	// Set the global DB variable
	originalDB := DB
	DB = testDB

	// Test Close function
	Close()

	// Try to ping the closed database - this should fail
	err = DB.Ping()
	if err == nil {
		t.Error("Expected error when pinging closed database, but got none")
	}

	// Restore original DB
	DB = originalDB
}

func TestCloseWithNilDB(t *testing.T) {
	// Save original DB
	originalDB := DB
	
	// Set DB to nil
	DB = nil
	
	// Close should not panic
	Close()
	
	// Restore original DB
	DB = originalDB
}