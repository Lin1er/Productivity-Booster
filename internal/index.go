// Package Index Description: Main application index managing pages and state
package index

import (
	// "fmt"

	"prodBooster/internal/db"
	"prodBooster/internal/models"
	"prodBooster/internal/ui/pages"

	tea "github.com/charmbracelet/bubbletea"
)

type Instance struct {
	// models used in the application
	todoList  *models.TodoList
	noteList  *models.NoteList
	eventList *models.EventList

	currentPage models.PageType
	pages       map[models.PageType]pages.Page // Map of page type to page instance

	// Dimensions of the terminal window
	width  int
	height int
}

func NewInstance() *Instance {
	database := db.Get()

	// Use database-backed models
	todoList_ := models.NewTodoList(database)
	noteList_ := models.NewNoteList(database)
	eventList_ := models.NewEventList(database)

	pageMap := make(map[models.PageType]pages.Page)
	pageMap[models.PageDashboard] = pages.NewDashboardPage(todoList_, noteList_, eventList_)
	pageMap[models.PageTodos] = pages.NewTodosPage(todoList_)
	pageMap[models.PageNotes] = pages.NewNotesPage(noteList_)
	pageMap[models.PageCalendar] = pages.NewCalendarPage(eventList_)

	return &Instance{
		todoList:    todoList_,
		noteList:    noteList_,
		eventList:   eventList_,
		currentPage: models.PageDashboard,
		pages:       pageMap,
		width:       80,
		height:      24,
	}
}

func (i *Instance) Init() tea.Cmd {
	return nil
}

func (i *Instance) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		i.width = msg.Width
		i.height = msg.Height
		// Update size for all pages
		for _, page := range i.pages {
			page.SetSize(msg.Width, msg.Height)
		}
		return i, nil

	case tea.KeyMsg:
		// Check if current page has an active form
		currentPage := i.pages[i.currentPage]
		if currentPage.IsFormActive() {
			// If form is active, only allow 'q' for quit, route everything else to page
			// if msg.String() == "q" {
			// 	return i, tea.Quit
			// }
			updatedPage, cmd := currentPage.Update(msg)
			i.pages[i.currentPage] = updatedPage
			return i, cmd
		}

		// Normal navigation when no form is active
		switch msg.String() {
		case "q":
			return i, tea.Quit
		case "UP":
			// Handle up key
		case "DOWN":
			// Handle down key
		case "1":
			i.currentPage = models.PageDashboard
			return i, nil
		case "2":
			i.currentPage = models.PageTodos
			return i, nil
		case "3":
			i.currentPage = models.PageNotes
			return i, nil
		case "4":
			i.currentPage = models.PageCalendar
			return i, nil
		default:
			updatedPage, cmd := currentPage.Update(msg)
			i.pages[i.currentPage] = updatedPage
			return i, cmd
		}
	}
	return i, nil
}

func (i *Instance) View() string {
	currentPage := i.pages[i.currentPage]
	return currentPage.View()
}
