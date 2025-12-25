// Package styles contains all styling definitions for the UI
package styles

import "github.com/charmbracelet/lipgloss"

const (
	SidebarWidth = 30
	StatusHeight = 1
)

const (
	ColorPrimary = lipgloss.Color("170") // purple

	ColorSuccess = lipgloss.Color("34") // Green

	ColorBorder   = lipgloss.Color("62")  // Blue
	ColorMuted    = lipgloss.Color("246") // Gray
	ColorStatusBg = lipgloss.Color("62")  // Status bar background
	ColorStatusFg = lipgloss.Color("0")   // Status bar foreground (black)

	// Priority colors
	ColorHighPriority   = lipgloss.Color("196") // Red
	ColorMediumPriority = lipgloss.Color("214") // Orange
	ColorLowPriority    = lipgloss.Color("34")  // Green
)

// === GO LESSON: Variables vs Constants ===
// 'var' declares variables that can be changed
// These are package-level variables (shared across package)
var (
	// Sidebar styles
	SidebarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 1)

	// Selected item in sidebar
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Bold(true).
				Underline(true)

	// Normal item in sidebar
	NormalItemStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	HighPriorItemStyle = lipgloss.NewStyle().
				Foreground(ColorHighPriority)

	MediumPriorItemStyle = lipgloss.NewStyle().
				Foreground(ColorMediumPriority)

	LowPriorItemStyle = lipgloss.NewStyle().
				Foreground(ColorLowPriority)

	CompletedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Italic(true).
				Strikethrough(true)

	// Content pane styles
	ContentStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2)

	// Content title style
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Underline(true).
			MarginBottom(1)

	// Status bar style
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(ColorStatusFg).
			Background(ColorStatusBg).
			Padding(0, 1)

	// Text styles
	TextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")) // White

	TextSecondaryStyle = lipgloss.NewStyle().
				Foreground(ColorMuted)

	// Success and Error styles
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")) // Green

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")) // Red

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")) // Yellow
)

// === GO LESSON: Functions that return configured styles ===
// These functions are useful when styles need dynamic sizing

// GetSidebarStyle returns a sidebar style with given dimensions
func GetSidebarStyle(height int) lipgloss.Style {
	return SidebarStyle.
		Width(SidebarWidth).
		Height(height)
}

// GetContentStyle returns a content style with given dimensions
func GetContentStyle(width, height int) lipgloss.Style {
	return ContentStyle.
		Width(width).
		Height(height)
}

// GetStatusBarStyle returns a status bar style with given width
func GetStatusBarStyle(width int) lipgloss.Style {
	return StatusBarStyle.Width(width - 2)
}
