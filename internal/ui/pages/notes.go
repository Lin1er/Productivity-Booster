package pages

import (
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"prodBooster/internal/models"
	"prodBooster/internal/ui/components"

	"github.com/charmbracelet/lipgloss"
)

// noteItem implements list.Item interface
type noteItem struct {
	note *models.Note
}

func (n noteItem) Title() string {
	return "📝 " + n.note.Title
}

func (n noteItem) Description() string {
	preview := n.note.Content
	if len(preview) > 100 {
		preview = preview[:100] + "..."
	}
	return preview
}

func (n noteItem) FilterValue() string {
	return n.note.Title
}

// Custom delegate for colored note items
type noteDelegate struct{}

func (d noteDelegate) Height() int                             { return 1 }
func (d noteDelegate) Spacing() int                            { return 0 }
func (d noteDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d noteDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	note, ok := item.(noteItem)
	if !ok {
		return
	}

	// Color by age - newer notes are brighter
	var style lipgloss.Style
	age := time.Since(note.note.CreatedAt)

	if age < 24*time.Hour {
		// New today - bright cyan
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("51"))
	} else if age < 7*24*time.Hour {
		// This week - cyan
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("45"))
	} else if age < 30*24*time.Hour {
		// This month - blue
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	} else {
		// Old - dim blue
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	}

	// Highlight selected item
	if index == m.Index() {
		style = style.Background(lipgloss.Color("238"))
	}

	fmt.Fprint(w, style.Render(note.Title()))
}

type NotesPage struct {
	NoteList     *models.NoteList
	form         *components.NoteForm
	searchBar    *components.SearchBar
	list         list.Model
	currentPage  models.PageType
	width        int
	height       int
	sidebarWidth int
}

func NewNotesPage(noteList_ *models.NoteList) *NotesPage {
	// Sort notes initially - newest first
	sortNotes(noteList_.Notes)

	// Create list items from notes
	items := make([]list.Item, len(noteList_.Notes))
	for i, note := range noteList_.Notes {
		items[i] = noteItem{note: note}
	}

	// Use custom delegate for colored rendering
	delegate := noteDelegate{}

	l := list.New(items, delegate, 0, 0)
	l.Title = "📒 My Notes"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false) // We use our own search

	return &NotesPage{
		NoteList:     noteList_,
		form:         components.NewNoteForm(noteList_),
		searchBar:    components.NewSearchBar(),
		list:         l,
		currentPage:  models.PageTypeNotes(),
		width:        80,
		height:       24,
		sidebarWidth: 40,
	}
}

func (p *NotesPage) Init() tea.Cmd {
	return nil
}

func (p *NotesPage) Update(msg tea.Msg) (Page, tea.Cmd) {
	// If form is active, route all input to form
	if p.form.IsActive() {
		updatedForm, cmd := p.form.Update(msg)
		p.form = updatedForm

		// Reload list after form submission
		if !p.form.IsActive() {
			p.updateListItems()
			p.list.Select(0) // Reset to first item
		}

		return p, cmd
	}

	// If search is active, route to search bar
	if p.searchBar.IsActive() {
		updatedSearch, cmd := p.searchBar.Update(msg)
		p.searchBar = updatedSearch

		// Update list when search changes
		if !p.searchBar.IsActive() {
			p.updateListItems()
		}

		return p, cmd
	}

	// Normal page navigation
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			// Create new note
			p.form.Activate()

		case "e":
			// Edit selected note
			if item, ok := p.list.SelectedItem().(noteItem); ok {
				p.form.LoadForEdit(item.note)
			}

		case "d", "delete":
			// Delete selected note
			if item, ok := p.list.SelectedItem().(noteItem); ok {
				if err := p.NoteList.Remove(item.note.ID); err != nil {
					// TODO: Handle error display
				}
				p.updateListItems()
				p.list.Select(0) // Reset to first item
			}

		case "/":
			// Activate search
			p.searchBar.Activate()
			return p, nil
		}
	}

	// Update the list for navigation and sync selection
	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)

	// Sync NoteList.Selected with list's selected index
	p.NoteList.Selected = p.list.Index()

	return p, cmd
}

// updateListItems refreshes the list with current notes and filters
func (p *NotesPage) updateListItems() {
	filteredNotes := []*models.Note{}
	for _, note := range p.NoteList.Notes {
		// Apply text search
		if !p.searchBar.Match(note.Title + " " + note.Content) {
			continue
		}
		filteredNotes = append(filteredNotes, note)
	}

	// Sort notes by creation date - newest first
	sortNotes(filteredNotes)

	// Convert to list items
	items := make([]list.Item, len(filteredNotes))
	for i, note := range filteredNotes {
		items[i] = noteItem{note: note}
	}

	p.list.SetItems(items)
}

// sortNotes sorts by creation date - newest first
func sortNotes(notes []*models.Note) {
	for i := 0; i < len(notes); i++ {
		for j := i + 1; j < len(notes); j++ {
			if notes[j].CreatedAt.After(notes[i].CreatedAt) {
				notes[i], notes[j] = notes[j], notes[i]
			}
		}
	}
}

func (p *NotesPage) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.sidebarWidth = width / 3
	if p.sidebarWidth < 30 {
		p.sidebarWidth = 30
	}
	if p.sidebarWidth > 50 {
		p.sidebarWidth = 50
	}

	p.list.SetSize(p.sidebarWidth-4, height-6) // Account for borders and padding
	p.form.SetSize(width, height)
}

func (p *NotesPage) View() string {
	// If form is active, show form overlay
	if p.form.IsActive() {
		return lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, p.form.View())
	}

	// If search is active, show search overlay
	if p.searchBar.IsActive() {
		return lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Top, p.searchBar.View())
	}

	topBar := components.NewTopBar(models.PageTypeNotes())
	topBar.SetSize(p.width, 1)

	// Sidebar with list
	sidebarStyle := lipgloss.NewStyle().
		Width(p.sidebarWidth).
		Height(p.height - 6).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	sidebar := sidebarStyle.Render(p.list.View())

	// Content pane with selected note detail
	contentWidth := p.width - p.sidebarWidth - 4
	contentStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Height(p.height-6).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	var content string
	if len(p.NoteList.Notes) == 0 {
		content = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("✨ No notes yet!\n\nPress 'n' to jot down your first thought 💭")
	} else if p.NoteList.Selected < len(p.NoteList.Notes) {
		note := p.NoteList.Notes[p.NoteList.Selected]

		// Title
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("213"))

		// Created date
		dateStr := "📅 Created " + note.CreatedAt.Format("Mon, Jan 2, 2006 at 3:04 PM")

		content = lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render(note.Title),
			"",
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render(dateStr),
			"",
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Width(contentWidth-4).
				Render(note.Content),
		)
	} else {
		content = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("👈 Pick a note to read it")
	}

	// Add help text
	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("✨ n: new note • e: edit • d: delete • /: search • ↑/↓: browse • q: quit")

	// Combine sidebar and content
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, contentStyle.Render(content))

	// Build the view
	return lipgloss.JoinVertical(lipgloss.Left,
		topBar.View(),
		mainContent,
		helpText,
	)
}

func (p *NotesPage) IsFormActive() bool {
	return p.form.IsActive() || p.searchBar.IsActive()
}
