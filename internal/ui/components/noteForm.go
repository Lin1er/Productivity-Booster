package components

import (
	"prodBooster/internal/models"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type NoteForm struct {
	noteList     *models.NoteList
	titleInput   textinput.Model
	contentInput textarea.Model
	focusIndex   int
	width        int
	height       int
	isActive     bool
	editMode     bool
	editingID    int
}

func NewNoteForm(noteList *models.NoteList) *NoteForm {
	ti := textinput.New()
	ti.Placeholder = "Give your note a title (e.g., Meeting notes, Ideas, Reminders)"
	ti.Focus()

	ta := textarea.New()
	ta.Blur()
	ta.Placeholder = "Write anything here... your thoughts, plans, random ideas 💭\n\nCtrl+S to save • Esc to cancel"

	return &NoteForm{
		noteList:     noteList,
		titleInput:   ti,
		contentInput: ta,
		focusIndex:   0,
		isActive:     false,
		editMode:     false,
	}
}

func (f *NoteForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *NoteForm) Update(msg tea.Msg) (*NoteForm, tea.Cmd) {
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

		case "tab":
			if f.focusIndex == 0 {
				f.focusIndex = 1
				f.titleInput.Blur()
				cmd = f.contentInput.Focus()
			} else {
				f.focusIndex = 0
				f.contentInput.Blur()
				cmd = f.titleInput.Focus()
			}
			return f, cmd

		case "ctrl+s":
			// Submit form
			if err := f.Submit(); err != nil {
				// TODO: Show error
			}
			f.isActive = false
			f.Reset()
			return f, nil
		}
	}

	if f.focusIndex == 0 {
		f.titleInput, cmd = f.titleInput.Update(msg)
	} else {
		f.contentInput, cmd = f.contentInput.Update(msg)
	}

	return f, cmd
}

func (f *NoteForm) View() string {
	if !f.isActive {
		return ""
	}

	formStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(f.width - 4).
		Height(f.height - 4)

	title := "✍️  New Note"
	if f.editMode {
		title = "✏️  Edit Note"
	}

	focusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	titleLabel := "📌 Note title"
	if f.focusIndex == 0 {
		titleLabel = focusStyle.Render("→ " + titleLabel)
	} else {
		titleLabel = normalStyle.Render(titleLabel)
	}

	contentLabel := "📝 Your thoughts"
	if f.focusIndex == 1 {
		contentLabel = focusStyle.Render("→ " + contentLabel)
	} else {
		contentLabel = normalStyle.Render(contentLabel)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(title),
		"",
		titleLabel,
		f.titleInput.View(),
		"",
		contentLabel,
		f.contentInput.View(),
		"",
		normalStyle.Render("✨ Ctrl+S to save • Tab to switch • Esc to cancel"),
	)

	return formStyle.Render(content)
}

func (f *NoteForm) SetSize(width, height int) {
	f.width = width
	f.height = height
	f.contentInput.SetWidth(width - 10)
	f.contentInput.SetHeight(height - 15)
}

func (f *NoteForm) Activate() {
	f.isActive = true
	f.titleInput.Focus()
}

func (f *NoteForm) Deactivate() {
	f.isActive = false
	f.Reset()
}

func (f *NoteForm) IsActive() bool {
	return f.isActive
}

func (f *NoteForm) Reset() {
	f.titleInput.SetValue("")
	f.contentInput.SetValue("")
	f.focusIndex = 0
	f.editMode = false
	f.editingID = 0
}

func (f *NoteForm) Submit() error {
	title := f.titleInput.Value()
	content := f.contentInput.Value()

	if title == "" {
		return nil // Don't submit empty notes
	}

	if f.editMode {
		return f.noteList.Update(f.editingID, title, content)
	}

	return f.noteList.Add(title, content)
}

func (f *NoteForm) LoadForEdit(note *models.Note) {
	f.editMode = true
	f.editingID = note.ID
	f.titleInput.SetValue(note.Title)
	f.contentInput.SetValue(note.Content)
	f.isActive = true
	f.titleInput.Focus()
}
