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

// todoItem implements list.Item interface
type todoItem struct {
	todo *models.Todo
}

func (t todoItem) Title() string {
	if t.todo.Completed {
		return "✓ " + t.todo.Title
	}

	icon := "○"
	if t.todo.Priority == models.PriorityHigh {
		icon = "●"
	} else if t.todo.Priority == models.PriorityMedium {
		icon = "◐"
	}

	return icon + " " + t.todo.Title
}

func (t todoItem) Description() string {
	return t.todo.Description
}

func (t todoItem) FilterValue() string {
	return t.todo.Title
}

// Custom delegate for colored todo items
type todoDelegate struct{}

func (d todoDelegate) Height() int                             { return 1 }
func (d todoDelegate) Spacing() int                            { return 0 }
func (d todoDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d todoDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	todo, ok := item.(todoItem)
	if !ok {
		return
	}

	// Determine color based on status
	var style lipgloss.Style
	now := time.Now()

	if todo.todo.Completed {
		// Completed - dim green
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("70")).Strikethrough(true)
	} else if todo.todo.DueTime != nil && todo.todo.DueTime.Before(now) {
		// Overdue - bright red
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	} else if todo.todo.DueTime != nil {
		todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
		if todo.todo.DueTime.Before(todayEnd) {
			// Due today - yellow/orange
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
		} else {
			// Future - normal color by priority
			if todo.todo.Priority == models.PriorityHigh {
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("213"))
			} else if todo.todo.Priority == models.PriorityMedium {
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("147"))
			} else {
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
			}
		}
	} else {
		// No due date - normal color by priority
		if todo.todo.Priority == models.PriorityHigh {
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("213"))
		} else if todo.todo.Priority == models.PriorityMedium {
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("147"))
		} else {
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
		}
	}

	// Highlight selected item
	if index == m.Index() {
		style = style.Background(lipgloss.Color("238"))
	}

	fmt.Fprint(w, style.Render(todo.Title()))
}

type TodosPage struct {
	currentPage  models.PageType
	TodoList     *models.TodoList
	form         *components.TodoForm
	searchBar    *components.SearchBar
	list         list.Model
	width        int
	height       int
	sidebarWidth int
}

func NewTodosPage(todoList_ *models.TodoList) *TodosPage {
	// Sort todos initially
	now := time.Now()
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	sortTodos(todoList_.Todos, now, todayEnd)

	// Create list items from todos
	items := make([]list.Item, len(todoList_.Todos))
	for i, todo := range todoList_.Todos {
		items[i] = todoItem{todo: todo}
	}

	// Use custom delegate for colored rendering
	delegate := todoDelegate{}

	l := list.New(items, delegate, 0, 0)
	l.Title = "✅ My Tasks"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false) // We use our own search

	return &TodosPage{
		currentPage:  models.PageTypeTodos(),
		TodoList:     todoList_,
		form:         components.NewTodoForm(todoList_),
		searchBar:    components.NewSearchBar(),
		list:         l,
		width:        80,
		height:       24,
		sidebarWidth: 40,
	}
}

func (p *TodosPage) Init() tea.Cmd {
	return nil
}

func (p *TodosPage) Update(msg tea.Msg) (Page, tea.Cmd) {
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
		case "enter", " ":
			// Toggle completed status
			if item, ok := p.list.SelectedItem().(todoItem); ok {
				currentIndex := p.list.Index()
				if err := p.TodoList.ToggleCompleted(item.todo.ID); err != nil {
					// TODO: Handle error display in UI
				}
				p.updateListItems()
				// Try to maintain position, but clamp to valid range
				if currentIndex < len(p.list.Items()) {
					p.list.Select(currentIndex)
				} else if len(p.list.Items()) > 0 {
					p.list.Select(len(p.list.Items()) - 1)
				}
			}

		case "n":
			// Create new todo
			p.form.Activate()

		case "e":
			// Edit selected todo
			if item, ok := p.list.SelectedItem().(todoItem); ok {
				p.form.LoadForEdit(item.todo)
			}

		case "d":
			// Delete selected todo
			if item, ok := p.list.SelectedItem().(todoItem); ok {
				if err := p.TodoList.Remove(item.todo.ID); err != nil {
					// TODO: Handle error display in UI
				}
				p.updateListItems()
				p.list.Select(0) // Reset to first item
			}

		case "/":
			// Open search
			p.searchBar.Activate()
		}
	}

	// Update list component and sync selection
	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)

	// Sync TodoList.Selected with list's selected index
	p.TodoList.Selected = p.list.Index()

	return p, cmd
}

// updateListItems refreshes the list with current todos and filters
func (p *TodosPage) updateListItems() {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	filteredTodos := []*models.Todo{}
	for _, todo := range p.TodoList.Todos {
		// Apply text search
		if !p.searchBar.Match(todo.Title + " " + todo.Description) {
			continue
		}

		// Apply filter
		filter := p.searchBar.GetFilter()
		switch filter {
		case components.FilterPending:
			if todo.Completed {
				continue
			}
		case components.FilterCompleted:
			if !todo.Completed {
				continue
			}
		case components.FilterHighPriority:
			if todo.Priority != models.PriorityHigh {
				continue
			}
		case components.FilterMediumPriority:
			if todo.Priority != models.PriorityMedium {
				continue
			}
		case components.FilterLowPriority:
			if todo.Priority != models.PriorityLow {
				continue
			}
		case components.FilterToday:
			if todo.DueTime == nil || !todo.DueTime.After(todayStart) || !todo.DueTime.Before(todayEnd) {
				continue
			}
		case components.FilterOverdue:
			if todo.DueTime == nil || !todo.DueTime.Before(now) || todo.Completed {
				continue
			}
		}

		filteredTodos = append(filteredTodos, todo)
	}

	// Sort todos by priority: overdue → today → pending → completed
	sortTodos(filteredTodos, now, todayEnd)

	// Convert to list items
	items := make([]list.Item, len(filteredTodos))
	for i, todo := range filteredTodos {
		items[i] = todoItem{todo: todo}
	}

	p.list.SetItems(items)
}

// sortTodos sorts by priority: overdue > today > high priority > medium > low > completed
func sortTodos(todos []*models.Todo, now time.Time, todayEnd time.Time) {
	// Simple bubble sort with priority logic
	for i := 0; i < len(todos); i++ {
		for j := i + 1; j < len(todos); j++ {
			if shouldSwapTodos(todos[i], todos[j], now, todayEnd) {
				todos[i], todos[j] = todos[j], todos[i]
			}
		}
	}
}

func shouldSwapTodos(a, b *models.Todo, now time.Time, todayEnd time.Time) bool {
	scoreA := getTodoPriority(a, now, todayEnd)
	scoreB := getTodoPriority(b, now, todayEnd)
	return scoreB > scoreA // Higher score = higher priority
}

func getTodoPriority(todo *models.Todo, now time.Time, todayEnd time.Time) int {
	if todo.Completed {
		return 0 // Lowest priority
	}
	if todo.DueTime != nil {
		if todo.DueTime.Before(now) {
			return 100 // Overdue - highest priority
		}
		if todo.DueTime.Before(todayEnd) {
			return 90 // Due today
		}
	}
	// Not overdue, sort by priority
	switch todo.Priority {
	case models.PriorityHigh:
		return 80
	case models.PriorityMedium:
		return 70
	default:
		return 60
	}
}

func (p *TodosPage) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.sidebarWidth = width / 3
	if p.sidebarWidth < 30 {
		p.sidebarWidth = 30
	}
	if p.sidebarWidth > 50 {
		p.sidebarWidth = 50
	}

	p.list.SetSize(p.sidebarWidth-4, height-8) // Account for borders and padding
	p.form.SetSize(width, height)
}

func (p *TodosPage) View() string {
	// If form is active, show form overlay
	if p.form.IsActive() {
		return lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, p.form.View())
	}

	// If search is active, show search overlay
	if p.searchBar.IsActive() {
		return lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Top, p.searchBar.View())
	}

	topBar := components.NewTopBar(models.PageTypeTodos())
	topBar.SetSize(p.width, 1)

	// Sidebar with list
	sidebarStyle := lipgloss.NewStyle().
		Width(p.sidebarWidth).
		Height(p.height - 6).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	sidebar := sidebarStyle.Render(p.list.View())

	// Content pane with selected todo detail
	contentWidth := p.width - p.sidebarWidth - 4
	contentStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Height(p.height-6).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	var content string
	if len(p.TodoList.Todos) == 0 {
		content = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("No todos yet.\nPress 'n' to create one!")
	} else if p.TodoList.Selected < len(p.TodoList.Todos) {
		todo := p.TodoList.Todos[p.TodoList.Selected]

		// Title with status icon
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("213"))

		statusIcon := "○ Pending"
		statusColor := "214"
		if todo.Completed {
			statusIcon = "✓ Completed"
			statusColor = "120"
		}

		priorityIcon := ""
		priorityColor := "240"
		switch todo.Priority {
		case models.PriorityHigh:
			priorityIcon = "🔥 High Priority - Do this ASAP!"
			priorityColor = "196"
		case models.PriorityMedium:
			priorityIcon = "📌 Medium Priority - Normal task"
			priorityColor = "214"
		case models.PriorityLow:
			priorityIcon = "💤 Low Priority - When you have time"
			priorityColor = "240"
		}

		// Due date
		dueStr := "📅 No deadline set"
		dueColor := "240"
		if todo.DueTime != nil {
			now := time.Now()
			if todo.DueTime.Before(now) && !todo.Completed {
				dueStr = "⚠️  OVERDUE! Was due " + todo.DueTime.Format("Mon, Jan 2 at 3:04 PM")
				dueColor = "196"
			} else if todo.DueTime.After(now) && todo.DueTime.Before(now.Add(24*time.Hour)) {
				dueStr = "⏰ Due today at " + todo.DueTime.Format("3:04 PM")
				dueColor = "214"
			} else {
				dueStr = "📅 Due " + todo.DueTime.Format("Mon, Jan 2 at 3:04 PM")
				dueColor = "45"
			}
		}

		content = lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render(todo.Title),
			"",
			lipgloss.NewStyle().Foreground(lipgloss.Color(statusColor)).Render(statusIcon),
			lipgloss.NewStyle().Foreground(lipgloss.Color(priorityColor)).Render(priorityIcon),
			lipgloss.NewStyle().Foreground(lipgloss.Color(dueColor)).Render(dueStr),
			"",
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Width(contentWidth-4).
				Render(todo.Description),
		)
	} else {
		content = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("👈 Pick a task from the list to see the details")
	}

	// Add help text
	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("✨ n: new task • e: edit • d: delete • /: search • space: mark done • ↑/↓: browse • q: quit")

	// Combine sidebar and content
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, contentStyle.Render(content))

	// Build the view
	return lipgloss.JoinVertical(lipgloss.Left,
		topBar.View(),
		mainContent,
		helpText,
	)
}

func (p *TodosPage) IsFormActive() bool {
	return p.form.IsActive() || p.searchBar.IsActive()
}
