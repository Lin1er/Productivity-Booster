package main

import (
	"fmt"
	"os"
	"path/filepath"

	index "prodBooster/internal"
	"prodBooster/internal/db"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize database
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	dbPath := filepath.Join(homeDir, ".prodbooster", "data.db")
	if err := db.Init(dbPath); err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	p := tea.NewProgram(index.NewInstance())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
