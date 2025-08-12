// internal/models/tab.go
package models

import (
	"encoding/json"
	"time"
)

type Tab struct {
	ID            int       `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Artist        string    `json:"artist" db:"artist"`
	Content       [6]string `json:"content" db:"content"` // 6 strings
	Tuning        [6]string `json:"tuning" db:"tuning"`   // E A D G B e
	Tempo         int       `json:"tempo" db:"tempo"`
	TimeSignature string    `json:"time_signature" db:"time_signature"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// main.go
package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/guitar-tab-tui/internal/storage"
	"github.com/guitar-tab-tui/internal/ui"
)

func main() {
	// Initialize storage
	storage, err := storage.NewSQLiteStorage("tabs.db")
	if err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}

	// Create the main application model
	m := ui.NewModel(storage)

	// Start the Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}

// internal/models/tab.go
package models

import (
	"encoding/json"
	"time"
)

type Tab struct {
	ID            int       `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Artist        string    `json:"artist" db:"artist"`
	Content       [6]string `json:"content" db:"content"` // 6 strings
	Tuning        [6]string `json:"tuning" db:"tuning"`   // E A D G B e
	Tempo         int       `json:"tempo" db:"tempo"`
	TimeSignature string    `json:"time_signature" db:"time_signature"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

func NewEmptyTab(name string) *Tab {
	emptyLine := "----------------"
	return &Tab{
		Name:          name,
		Artist:        "",
		Content:       [6]string{emptyLine, emptyLine, emptyLine, emptyLine, emptyLine, emptyLine},
		Tuning:        [6]string{"e", "B", "G", "D", "A", "E"},
		Tempo:         120,
		TimeSignature: "4/4",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

type Position struct {
	String   int
	Position int
}

type SessionState struct {
	CurrentTab    *Tab
	CursorPos     Position
	ViewMode      ViewMode
	PlaybackState PlaybackState
	EditMode      EditMode
}

type ViewMode int

const (
	ViewHome ViewMode = iota
	ViewEditor
	ViewBrowser
	ViewSettings
	ViewHelp
)

type EditMode int

const (
	EditNormal EditMode = iota
	EditInsert
	EditSelect
)

type PlaybackState struct {
	IsPlaying     bool
	Position      int
	Highlighted   []Position
	Tempo         int
}
