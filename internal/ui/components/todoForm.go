package components

import (
	"prodBooster/internal/models"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TodoForm struct {
	todoList   *models.TodoList
	titleInput textinput.Model
	descInput  textinput.Model
	priority   models.Priority
	focusIndex int
	width      int
	height     int
	isActive   bool
	editMode   bool
	editingID  int
}

func NewTodoForm(todoList *models.TodoList) *TodoForm {
	ti := textinput.New()
	ti.Placeholder = "What needs to be done? (e.g., Buy groceries, Finish report)"
	ti.Focus()

	di := textinput.New()
	di.Placeholder = "Any extra details? (optional)"

	return &TodoForm{
		todoList:   todoList,
		titleInput: ti,
		descInput:  di,
		priority:   models.PriorityMedium,
		focusIndex: 0,
		isActive:   false,
		editMode:   false,
	}
}

func (f *TodoForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *TodoForm) Update(msg tea.Msg) (*TodoForm, tea.Cmd) {
	if !f.isActive {
		return f, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			f.isActive = false
			f.Reset()
			return f, nil

		case "tab", "shift+tab":
			f.focusIndex++
			if f.focusIndex > 2 {
				f.focusIndex = 0
			}

			f.titleInput.Blur()
			f.descInput.Blur()

			if f.focusIndex == 0 {
				cmd = f.titleInput.Focus()
			} else if f.focusIndex == 1 {
				cmd = f.descInput.Focus()
			}

			return f, cmd

		case "up":
			if f.focusIndex == 2 {
				if f.priority < models.PriorityHigh {
					f.priority++
				}
			}

		case "down":
			if f.focusIndex == 2 {
				if f.priority > models.PriorityLow {
					f.priority--
				}
			}

		case "enter":
			if f.focusIndex == 2 {
				// Submit form
				if err := f.Submit(); err != nil {
					// TODO: Show error
				}
				f.isActive = false
				f.Reset()
				return f, nil
			} else {
				// Move to next field
				f.focusIndex++
				if f.focusIndex == 1 {
					return f, f.descInput.Focus()
				} else if f.focusIndex == 2 {
					f.titleInput.Blur()
					f.descInput.Blur()
				}
			}
		}
	}

	if f.focusIndex == 0 {
		f.titleInput, cmd = f.titleInput.Update(msg)
	} else if f.focusIndex == 1 {
		f.descInput, cmd = f.descInput.Update(msg)
	}

	return f, cmd
}

func (f *TodoForm) View() string {
	if !f.isActive {
		return ""
	}

	formStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(f.width - 4).
		Height(f.height - 4)

	title := "Create New Todo"
	if f.editMode {
		title = "Edit Todo"
	}

	focusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)

	titleLabel := "📝 What's the task?"
	if f.focusIndex == 0 {
		titleLabel = focusStyle.Render("→ " + titleLabel)
	}

	descLabel := "💬 Add some details"
	if f.focusIndex == 1 {
		descLabel = focusStyle.Render("→ " + descLabel)
	}

	priorityLabel := "⭐ How urgent is this?"
	if f.focusIndex == 2 {
		priorityLabel = focusStyle.Render("→ " + priorityLabel)
	} else {
		priorityLabel = normalStyle.Render(priorityLabel)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(title),
		"",
		titleLabel,
		f.titleInput.View(),
		"",
		descLabel,
		f.descInput.View(),
		hintStyle.Render("  💡 Optional - add context, notes, or anything helpful"),
		"",
		priorityLabel+" "+f.priority.String()+" (use ↑↓ to change)",
		hintStyle.Render("  🔥 High = Do this ASAP! | 📌 Medium = Normal stuff | 💤 Low = When you have time"),
		"",
		"",
		normalStyle.Render("✨ Press Enter to save • Tab to move around • Esc to cancel"),
	)

	return formStyle.Render(content)
}

func (f *TodoForm) SetSize(width, height int) {
	f.width = width
	f.height = height
}

func (f *TodoForm) Activate() {
	f.isActive = true
	f.titleInput.Focus()
}

func (f *TodoForm) Deactivate() {
	f.isActive = false
	f.Reset()
}

func (f *TodoForm) IsActive() bool {
	return f.isActive
}

func (f *TodoForm) Reset() {
	f.titleInput.SetValue("")
	f.descInput.SetValue("")
	f.priority = models.PriorityMedium
	f.focusIndex = 0
	f.editMode = false
	f.editingID = 0
}

func (f *TodoForm) Submit() error {
	title := f.titleInput.Value()
	desc := f.descInput.Value()

	if title == "" {
		return nil // Don't submit empty todos
	}

	if f.editMode {
		return f.todoList.Update(f.editingID, title, desc, f.priority, nil)
	}

	return f.todoList.Add(title, desc, f.priority, nil)
}

func (f *TodoForm) LoadForEdit(todo *models.Todo) {
	f.editMode = true
	f.editingID = todo.ID
	f.titleInput.SetValue(todo.Title)
	f.descInput.SetValue(todo.Description)
	f.priority = todo.Priority
	f.isActive = true
	f.titleInput.Focus()
}
