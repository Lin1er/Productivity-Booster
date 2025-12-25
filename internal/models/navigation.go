// Package models contains data structures
package models

type PageType int

const (
	PageDashboard PageType = iota // 0
	PageTodos                     // 1
	PageNotes                     // 2
	PageCalendar                  // 3
)

// String makes PageType printable for debugging
func (p PageType) String() string {
	switch p {
	case PageDashboard:
		return "Dashboard"
	case PageTodos:
		return "To-Do List"
	case PageNotes:
		return "Notes"
	case PageCalendar:
		return "Calendar"
	default:
		return "Unknown"
	}
}

// PageInfo holds metadata about each page
type PageInfo struct {
	Type  PageType
	Title string
	Key   string // Keyboard shortcut
	Icon  string // Display icon
}

// GetPages returns all available pages
func GetPages() []PageInfo {
	return []PageInfo{
		{Type: PageDashboard, Title: "Dashboard", Key: "1", Icon: "📊"},
		{Type: PageTodos, Title: "To-Dos", Key: "2", Icon: "✓"},
		{Type: PageNotes, Title: "Notes", Key: "3", Icon: "📝"},
		{Type: PageCalendar, Title: "Calendar", Key: "4", Icon: "📅"},
	}
}

func PageTypeDashboard() PageType {
	return PageDashboard
}

func PageTypeTodos() PageType {
	return PageTodos
}

func PageTypeNotes() PageType {
	return PageNotes
}

func PageTypeCalendar() PageType {
	return PageCalendar
}
