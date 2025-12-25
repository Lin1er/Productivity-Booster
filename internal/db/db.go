// Package db provides database connection and query functionalities.
package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create db directory: %w", err)
	}
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}
	return nil
}

func createTables() error {
	// Create Notes table
	notesTable := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Create Todos table
	todosTable := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		completed BOOLEAN DEFAULT 0,
		priority INTEGER DEFAULT 0,
		due_date DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Create Events table
	eventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		start_time DATETIME NOT NULL,
		end_time DATETIME NOT NULL,
		location TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Execute table creation statements
	tables := []string{notesTable, todosTable, eventsTable}
	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func Get() *sql.DB {
	if DB != nil {
		return DB
	}
	return nil
}
