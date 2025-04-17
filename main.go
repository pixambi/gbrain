package main

import (
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pixambi/gbrain/cmd"
	"github.com/pixambi/gbrain/internal/db"
)

func main() {
	// Create data directory if it doesn't exist
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v", err)
	}

	dataDir := filepath.Join(homeDir, ".gbrain")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Error creating data directory: %v", err)
	}

	dbPath := filepath.Join(dataDir, "gbrain.db")
	db, err := db.NewDb(dbPath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Initialize the database
	if err := db.Init(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Create and start the application
	m := cmd.NewApp(*db)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
