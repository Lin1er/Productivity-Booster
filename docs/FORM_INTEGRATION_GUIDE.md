# Example: How to Integrate Forms into Pages

## TodosPage with Form Integration

Berikut contoh cara mengintegrasikan form ke TodosPage:

### 1. Update struct TodosPage di `internal/ui/pages/todos.go`

```go
type TodosPage struct {
	currentPage models.PageType
	TodoList    *models.TodoList
	form        *components.TodoForm  // Add this
	width       int
	height      int
}

func NewTodosPage(todoList_ *models.TodoList) *TodosPage {
	return &TodosPage{
		currentPage: models.PageTypeTodos(),
		TodoList:    todoList_,
		form:        components.NewTodoForm(todoList_), // Add this
		width:       80,
		height:      24,
	}
}
```

### 2. Update method Update() di TodosPage

```go
func (p *TodosPage) Update(msg tea.Msg) (Page, tea.Cmd) {
	// If form is active, route all input to form
	if p.form.IsActive() {
		updatedForm, cmd := p.form.Update(msg)
		p.form = updatedForm
		return p, cmd
	}

	// Normal page navigation
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			p.TodoList.SelectPrevious()

		case "down", "j":
			p.TodoList.SelectNext()

		case "enter", " ":
			// Toggle completed status
			if selected := p.TodoList.GetSelected(); selected != nil {
				if err := p.TodoList.ToggleCompleted(selected.ID); err != nil {
					// TODO: Handle error display in UI
				}
			}

		case "n":
			// Create new todo
			p.form.Activate()

		case "e":
			// Edit selected todo
			if selected := p.TodoList.GetSelected(); selected != nil {
				p.form.LoadForEdit(selected)
			}

		case "d":
			// Delete selected todo
			if selected := p.TodoList.GetSelected(); selected != nil {
				if err := p.TodoList.Remove(selected.ID); err != nil {
					// TODO: Handle error display
				}
			}
		}
	}
	return p, nil
}
```

### 3. Update SetSize() untuk set size form juga

```go
func (p *TodosPage) SetSize(width, height int) {
	p.width = width - 2
	p.height = height - 5
	p.form.SetSize(width, height) // Add this
}
```

### 4. Update View() untuk render form sebagai overlay

```go
func (p *TodosPage) View() string {
	// If form is active, show form overlay
	if p.form.IsActive() {
		baseView := p.renderTodoList()
		formView := p.form.View()

		// Show form centered on screen
		return lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center,
			formView)
	}

	return p.renderTodoList()
}

func (p *TodosPage) renderTodoList() string {
	// ... existing View() code ...
	// Just move all existing View() code here

	// Add help text at bottom
	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("n: new | e: edit | d: delete | space: toggle | j/k: navigate")

	// Add helpText to your layout
}
```

## NotesPage Integration

Similar pattern untuk NotesPage:

```go
type NotesPage struct {
	NoteList    *models.NoteList
	form        *components.NoteForm  // Add
	// ... existing fields
}

func NewNotesPage(noteList_ *models.NoteList) *NotesPage {
	return &NotesPage{
		NoteList:    noteList_,
		form:        components.NewNoteForm(noteList_), // Add
		// ... existing fields
	}
}

func (p *NotesPage) Update(msg tea.Msg) (Page, tea.Cmd) {
	if p.form.IsActive() {
		updatedForm, cmd := p.form.Update(msg)
		p.form = updatedForm
		return p, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			p.form.Activate()
		case "e":
			if selected := p.NoteList.GetSelected(); selected != nil {
				p.form.LoadForEdit(selected)
			}
		case "d":
			if selected := p.NoteList.GetSelected(); selected != nil {
				p.NoteList.Remove(selected.ID)
			}
		// ... existing navigation
		}
	}
	return p, nil
}
```

## CalendarPage Integration

```go
type CalendarPage struct {
	EventList *models.EventList
	form      *components.EventForm  // Add
	// ... existing fields
}

func NewCalendarPage(eventList_ *models.EventList) *CalendarPage {
	return &CalendarPage{
		EventList: eventList_,
		form:      components.NewEventForm(eventList_), // Add
		// ... existing fields
	}
}

func (p *CalendarPage) Update(msg tea.Msg) (Page, tea.Cmd) {
	if p.form.IsActive() {
		updatedForm, cmd := p.form.Update(msg)
		p.form = updatedForm
		return p, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			p.form.Activate()
		case "e":
			if selected := p.EventList.GetSelected(); selected != nil {
				p.form.LoadForEdit(selected)
			}
		case "d":
			if selected := p.EventList.GetSelected(); selected != nil {
				p.EventList.Remove(selected.ID)
			}
		}
	}
	return p, nil
}
```

## Key Bindings Summary

- `n` - New item
- `e` - Edit selected item
- `d` - Delete selected item
- `space` - Toggle completed (todos only)
- `j/k` or `↓/↑` - Navigate list

### In Forms:

- `Tab` - Next field
- `Enter` - Submit (TodoForm) / Next field
- `Ctrl+S` - Save (NoteForm, EventForm)
- `Esc` - Cancel
- `↑/↓` - Change priority (TodoForm only, when on priority field)
