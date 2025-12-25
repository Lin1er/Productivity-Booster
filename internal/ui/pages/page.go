// Package pages contains individual page implementations
package pages

import (
	tea "github.com/charmbracelet/bubbletea"
)

// === GO LESSON: Interfaces ===
// An interface defines a contract - a set of methods
// Any type that has these methods automatically implements this interface
// This is IMPLICIT interface implementation (no 'implements' keyword)

// Page interface that all pages must implement
// This allows us to treat different pages uniformly
type Page interface {
	// Init is called when the page is first loaded
	Init() tea.Cmd

	// Update handles messages/events for this page
	Update(msg tea.Msg) (Page, tea.Cmd)

	// View renders the page content
	View() string

	// SetSize updates the page dimensions
	SetSize(width, height int)

	// IsFormActive returns true if a form is currently active
	IsFormActive() bool
}
