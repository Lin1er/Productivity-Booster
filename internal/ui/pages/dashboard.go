// Package pages contains individual page implementations
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

type dashboardFocus int

const (
	focusTodos dashboardFocus = iota
	focusEvents
	focusNotes
)

type dashboardMode int

const (
	dashboardModeNormal dashboardMode = iota
	dashboardModeQuickAddTodo
	dashboardModeQuickAddEvent
	dashboardModeQuickAddNote
)

// Dashboard list items
type dashboardTodoItem struct{ todo *models.Todo }

func (t dashboardTodoItem) Title() string       { return t.todo.Title }
func (t dashboardTodoItem) Description() string { return "" }
func (t dashboardTodoItem) FilterValue() string { return t.todo.Title }

type dashboardEventItem struct{ event *models.Event }

func (e dashboardEventItem) Title() string       { return e.event.Title }
func (e dashboardEventItem) Description() string { return "" }
func (e dashboardEventItem) FilterValue() string { return e.event.Title }

type dashboardNoteItem struct{ note *models.Note }

func (n dashboardNoteItem) Title() string       { return n.note.Title }
func (n dashboardNoteItem) Description() string { return "" }
func (n dashboardNoteItem) FilterValue() string { return n.note.Title }

// Custom delegates for dashboard items
type dashboardTodoDelegate struct{}

func (d dashboardTodoDelegate) Height() int                             { return 1 }
func (d dashboardTodoDelegate) Spacing() int                            { return 0 }
func (d dashboardTodoDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d dashboardTodoDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	todo, ok := item.(dashboardTodoItem)
	if !ok {
		return
	}

	var style lipgloss.Style
	now := time.Now()

	if todo.todo.Completed {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("70")).Strikethrough(true)
	} else if todo.todo.DueTime != nil && todo.todo.DueTime.Before(now) {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	} else if todo.todo.DueTime != nil {
		todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
		if todo.todo.DueTime.Before(todayEnd) {
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
		} else {
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("147"))
		}
	} else {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	}

	if index == m.Index() {
		style = style.Background(lipgloss.Color("238"))
	}

	icon := "○"
	if todo.todo.Completed {
		icon = "✓"
	} else if todo.todo.Priority == models.PriorityHigh {
		icon = "●"
	} else if todo.todo.Priority == models.PriorityMedium {
		icon = "◐"
	}

	fmt.Fprint(w, style.Render(fmt.Sprintf("%s %s", icon, todo.Title())))
}

type dashboardEventDelegate struct{}

func (d dashboardEventDelegate) Height() int                             { return 1 }
func (d dashboardEventDelegate) Spacing() int                            { return 0 }
func (d dashboardEventDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d dashboardEventDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	event, ok := item.(dashboardEventItem)
	if !ok {
		return
	}

	var style lipgloss.Style
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	if event.event.StartTime.Before(now) {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Strikethrough(true)
	} else if event.event.StartTime.After(todayStart) && event.event.StartTime.Before(todayEnd) {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	} else {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("147"))
	}

	if index == m.Index() {
		style = style.Background(lipgloss.Color("238"))
	}

	timeStr := event.event.StartTime.Format("Jan 2 15:04")
	fmt.Fprint(w, style.Render(fmt.Sprintf("📅 %s • %s", timeStr, event.Title())))
}

type dashboardNoteDelegate struct{}

func (d dashboardNoteDelegate) Height() int                             { return 1 }
func (d dashboardNoteDelegate) Spacing() int                            { return 0 }
func (d dashboardNoteDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d dashboardNoteDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	note, ok := item.(dashboardNoteItem)
	if !ok {
		return
	}

	var style lipgloss.Style
	age := time.Since(note.note.CreatedAt)

	if age < 24*time.Hour {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("51"))
	} else if age < 7*24*time.Hour {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("45"))
	} else {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	}

	if index == m.Index() {
		style = style.Background(lipgloss.Color("238"))
	}

	fmt.Fprint(w, style.Render(fmt.Sprintf("📝 %s", note.Title())))
}

type DashboardPage struct {
	TodoList    *models.TodoList
	NoteList    *models.NoteList
	EventList   *models.EventList
	currentPage models.PageType
	width       int
	height      int
	mode        dashboardMode
	focus       dashboardFocus
	todoList    list.Model
	eventList   list.Model
	noteList    list.Model
	todoForm    *components.TodoForm
	eventForm   *components.EventForm
	noteForm    *components.NoteForm
}

func NewDashboardPage(todoList_ *models.TodoList, noteList_ *models.NoteList, eventList_ *models.EventList) *DashboardPage {
	// Sort and create todo list
	now := time.Now()
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	sortTodos(todoList_.Todos, now, todayEnd)

	todoItems := make([]list.Item, 0)
	for _, todo := range todoList_.Todos {
		todoItems = append(todoItems, dashboardTodoItem{todo: todo})
	}
	todoListModel := list.New(todoItems, dashboardTodoDelegate{}, 0, 0)
	todoListModel.Title = "Todos"
	todoListModel.SetShowStatusBar(false)
	todoListModel.SetFilteringEnabled(false)

	// Sort and create event list
	sortEvents(eventList_.Events)

	eventItems := make([]list.Item, 0)
	for _, event := range eventList_.Events {
		eventItems = append(eventItems, dashboardEventItem{event: event})
	}
	eventListModel := list.New(eventItems, dashboardEventDelegate{}, 0, 0)
	eventListModel.Title = "Events"
	eventListModel.SetShowStatusBar(false)
	eventListModel.SetFilteringEnabled(false)

	// Sort and create note list
	sortNotes(noteList_.Notes)

	noteItems := make([]list.Item, 0)
	for _, note := range noteList_.Notes {
		noteItems = append(noteItems, dashboardNoteItem{note: note})
	}
	noteListModel := list.New(noteItems, dashboardNoteDelegate{}, 0, 0)
	noteListModel.Title = "Notes"
	noteListModel.SetShowStatusBar(false)
	noteListModel.SetFilteringEnabled(false)

	return &DashboardPage{
		TodoList:    todoList_,
		NoteList:    noteList_,
		EventList:   eventList_,
		currentPage: models.PageTypeDashboard(),
		width:       80,
		height:      24,
		mode:        dashboardModeNormal,
		focus:       focusTodos,
		todoList:    todoListModel,
		eventList:   eventListModel,
		noteList:    noteListModel,
		todoForm:    components.NewTodoForm(todoList_),
		eventForm:   components.NewEventForm(eventList_),
		noteForm:    components.NewNoteForm(noteList_),
	}
}

func (p *DashboardPage) Init() tea.Cmd {
	return nil
}

func (p *DashboardPage) Update(msg tea.Msg) (Page, tea.Cmd) {
	// Handle form modes
	switch p.mode {
	case dashboardModeQuickAddTodo:
		if p.todoForm.IsActive() {
			var cmd tea.Cmd
			p.todoForm, cmd = p.todoForm.Update(msg)
			if !p.todoForm.IsActive() {
				p.mode = dashboardModeNormal
				p.updateLists()
				p.todoList.Select(0) // Reset to first item
			}
			return p, cmd
		}
	case dashboardModeQuickAddEvent:
		if p.eventForm.IsActive() {
			var cmd tea.Cmd
			p.eventForm, cmd = p.eventForm.Update(msg)
			if !p.eventForm.IsActive() {
				p.mode = dashboardModeNormal
				p.updateLists()
				p.eventList.Select(0) // Reset to first item
			}
			return p, cmd
		}
	case dashboardModeQuickAddNote:
		if p.noteForm.IsActive() {
			var cmd tea.Cmd
			p.noteForm, cmd = p.noteForm.Update(msg)
			if !p.noteForm.IsActive() {
				p.mode = dashboardModeNormal
				p.updateLists()
				p.noteList.Select(0) // Reset to first item
			}
			return p, cmd
		}
	}

	// Normal mode navigation
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// Cycle focus: Todos -> Events -> Notes -> Todos
			p.focus = (p.focus + 1) % 3
			return p, nil

		case "a":
			// Quick add based on current focus
			switch p.focus {
			case focusTodos:
				p.mode = dashboardModeQuickAddTodo
				p.todoForm.Activate()
			case focusEvents:
				p.mode = dashboardModeQuickAddEvent
				p.eventForm.Activate()
			case focusNotes:
				p.mode = dashboardModeQuickAddNote
				p.noteForm.Activate()
			}
			return p, nil
		}
	}

	// Update focused list
	var cmd tea.Cmd
	switch p.focus {
	case focusTodos:
		p.todoList, cmd = p.todoList.Update(msg)
	case focusEvents:
		p.eventList, cmd = p.eventList.Update(msg)
	case focusNotes:
		p.noteList, cmd = p.noteList.Update(msg)
	}

	return p, cmd
}

// updateLists refreshes all lists after CRUD operations
func (p *DashboardPage) updateLists() {
	// Update todo list
	now := time.Now()
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	sortTodos(p.TodoList.Todos, now, todayEnd)

	todoItems := make([]list.Item, 0)
	for _, todo := range p.TodoList.Todos {
		todoItems = append(todoItems, dashboardTodoItem{todo: todo})
	}
	p.todoList.SetItems(todoItems)

	// Update event list
	sortEvents(p.EventList.Events)

	eventItems := make([]list.Item, 0)
	for _, event := range p.EventList.Events {
		eventItems = append(eventItems, dashboardEventItem{event: event})
	}
	p.eventList.SetItems(eventItems)

	// Update note list
	sortNotes(p.NoteList.Notes)

	noteItems := make([]list.Item, 0)
	for _, note := range p.NoteList.Notes {
		noteItems = append(noteItems, dashboardNoteItem{note: note})
	}
	p.noteList.SetItems(noteItems)
}

func (p *DashboardPage) SetSize(width, height int) {
	p.width = width
	p.height = height

	// Calculate card dimensions
	cardWidth := (width - 8) / 3
	cardHeight := height - 7

	p.todoList.SetSize(cardWidth-2, cardHeight)
	p.eventList.SetSize(cardWidth-2, cardHeight)
	p.noteList.SetSize(cardWidth-2, cardHeight)

	p.todoForm.SetSize(width, height)
	p.eventForm.SetSize(width, height)
	p.noteForm.SetSize(width, height)
}

func (p *DashboardPage) View() string {
	// Show form overlay if in add mode
	switch p.mode {
	case dashboardModeQuickAddTodo:
		return lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, p.todoForm.View())
	case dashboardModeQuickAddEvent:
		return lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, p.eventForm.View())
	case dashboardModeQuickAddNote:
		return lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, p.noteForm.View())
	}

	topBar := components.NewTopBar(models.PageTypeDashboard())
	topBar.SetSize(p.width, 1)

	// Calculate stats for quick summary
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	var todayTodos, overdueTodos, pendingTodos int
	for _, todo := range p.TodoList.Todos {
		if !todo.Completed {
			pendingTodos++
			if todo.DueTime != nil && todo.DueTime.Before(now) {
				overdueTodos++
			} else if todo.DueTime != nil && todo.DueTime.After(todayStart) && todo.DueTime.Before(todayEnd) {
				todayTodos++
			}
		}
	}

	// Hero section with key stats
	heroText := fmt.Sprintf("✨ %s • %d tasks pending • %d overdue • %d due today",
		now.Format("Monday, January 2, 2006"),
		pendingTodos,
		overdueTodos,
		todayTodos)
	heroStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("62")).
		Padding(0, 2).
		Width(p.width).
		Align(lipgloss.Center)

	// Card dimensions
	cardWidth := (p.width - 8) / 3
	cardHeight := p.height - 8

	// Create card styles
	todoCardStyle := lipgloss.NewStyle().
		Width(cardWidth).
		Height(cardHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("45"))

	eventCardStyle := lipgloss.NewStyle().
		Width(cardWidth).
		Height(cardHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("213"))

	noteCardStyle := lipgloss.NewStyle().
		Width(cardWidth).
		Height(cardHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("228"))

	// Highlight focused card
	switch p.focus {
	case focusTodos:
		todoCardStyle = todoCardStyle.BorderForeground(lipgloss.Color("51")).Bold(true)
	case focusEvents:
		eventCardStyle = eventCardStyle.BorderForeground(lipgloss.Color("201")).Bold(true)
	case focusNotes:
		noteCardStyle = noteCardStyle.BorderForeground(lipgloss.Color("226")).Bold(true)
	}

	// Render cards
	todoCard := todoCardStyle.Render(p.todoList.View())
	eventCard := eventCardStyle.Render(p.eventList.View())
	noteCard := noteCardStyle.Render(p.noteList.View())

	cardsRow := lipgloss.JoinHorizontal(lipgloss.Top, todoCard, eventCard, noteCard)

	// Help text
	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("✨ Tab: switch cards • a: quick add • ↑/↓: browse • Enter: go to page • q: quit")

	// Assemble view
	return lipgloss.JoinVertical(lipgloss.Left,
		topBar.View(),
		heroStyle.Render(heroText),
		cardsRow,
		helpText,
	)
}

func (p *DashboardPage) IsFormActive() bool {
	return p.mode != dashboardModeNormal
}
