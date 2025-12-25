package components

import (
	"prodBooster/internal/models"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EventForm struct {
	eventList     *models.EventList
	titleInput    textinput.Model
	descInput     textinput.Model
	locationInput textinput.Model
	startInput    textinput.Model
	endInput      textinput.Model
	focusIndex    int
	width         int
	height        int
	isActive      bool
	editMode      bool
	editingID     int
}

func NewEventForm(eventList *models.EventList) *EventForm {
	ti := textinput.New()
	ti.Placeholder = "What's happening? (e.g., Team meeting, Lunch with Sarah)"
	ti.Focus()

	di := textinput.New()
	di.Placeholder = "Any details worth noting? (optional)"

	li := textinput.New()
	li.Placeholder = "Where? (e.g., Office, Zoom, Coffee shop)"

	si := textinput.New()
	si.Placeholder = "When does it start? → 2025-12-25 14:30"

	ei := textinput.New()
	ei.Placeholder = "When does it end? → 2025-12-25 16:00 (or leave empty)"

	return &EventForm{
		eventList:     eventList,
		titleInput:    ti,
		descInput:     di,
		locationInput: li,
		startInput:    si,
		endInput:      ei,
		focusIndex:    0,
		isActive:      false,
		editMode:      false,
	}
}

func (f *EventForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *EventForm) Update(msg tea.Msg) (*EventForm, tea.Cmd) {
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
			f.focusIndex++
			if f.focusIndex > 4 {
				f.focusIndex = 0
			}

			f.titleInput.Blur()
			f.descInput.Blur()
			f.locationInput.Blur()
			f.startInput.Blur()
			f.endInput.Blur()

			switch f.focusIndex {
			case 0:
				cmd = f.titleInput.Focus()
			case 1:
				cmd = f.descInput.Focus()
			case 2:
				cmd = f.locationInput.Focus()
			case 3:
				cmd = f.startInput.Focus()
			case 4:
				cmd = f.endInput.Focus()
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

	switch f.focusIndex {
	case 0:
		f.titleInput, cmd = f.titleInput.Update(msg)
	case 1:
		f.descInput, cmd = f.descInput.Update(msg)
	case 2:
		f.locationInput, cmd = f.locationInput.Update(msg)
	case 3:
		f.startInput, cmd = f.startInput.Update(msg)
	case 4:
		f.endInput, cmd = f.endInput.Update(msg)
	}

	return f, cmd
}

func (f *EventForm) View() string {
	if !f.isActive {
		return ""
	}

	formStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(f.width - 4).
		Height(f.height - 4)

	title := "📅 New Event"
	if f.editMode {
		title = "✏️  Edit Event"
	}

	focusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)

	fields := []struct {
		label string
		input string
		hint  string
	}{
		{"🎯 What's happening?", f.titleInput.View(), ""},
		{"💬 Details", f.descInput.View(), "💡 optional - add context if needed"},
		{"📍 Where?", f.locationInput.View(), "💡 optional - place, room, or link"},
		{"🕐 Start time", f.startInput.View(), "Format: 2025-12-25 14:30 (year-month-day hour:minute)"},
		{"🕑 End time", f.endInput.View(), "💡 optional - leave empty for no end time"},
	}

	var lines []string
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render(title))
	lines = append(lines, "")

	for i, field := range fields {
		label := field.label
		if i == f.focusIndex {
			label = focusStyle.Render("→ " + label)
		} else {
			label = normalStyle.Render(label)
		}
		lines = append(lines, label)
		lines = append(lines, field.input)
		if field.hint != "" {
			lines = append(lines, hintStyle.Render("  "+field.hint))
		}
		lines = append(lines, "")
	}

	lines = append(lines, "")
	lines = append(lines, hintStyle.Render("💡 Quick tips:"))
	lines = append(lines, hintStyle.Render("  • Time uses 24-hour format (14:30 = 2:30 PM, 09:00 = 9:00 AM)"))
	lines = append(lines, hintStyle.Render("  • No end time? No problem - just leave it blank!"))
	lines = append(lines, "")
	lines = append(lines, normalStyle.Render("✨ Ctrl+S to save • Tab to move around • Esc to cancel"))

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	return formStyle.Render(content)
}

func (f *EventForm) SetSize(width, height int) {
	f.width = width
	f.height = height
}

func (f *EventForm) Activate() {
	f.isActive = true
	f.titleInput.Focus()
}

func (f *EventForm) Deactivate() {
	f.isActive = false
	f.Reset()
}

func (f *EventForm) IsActive() bool {
	return f.isActive
}

func (f *EventForm) Reset() {
	f.titleInput.SetValue("")
	f.descInput.SetValue("")
	f.locationInput.SetValue("")
	f.startInput.SetValue("")
	f.endInput.SetValue("")
	f.focusIndex = 0
	f.editMode = false
	f.editingID = 0
}

func (f *EventForm) Submit() error {
	title := f.titleInput.Value()
	desc := f.descInput.Value()
	location := f.locationInput.Value()
	startStr := f.startInput.Value()
	endStr := f.endInput.Value()

	if title == "" || startStr == "" || endStr == "" {
		return nil // Don't submit incomplete events
	}

	// Parse times (simple parsing, you might want better validation)
	layout := "2006-01-02 15:04"
	startTime, err := time.Parse(layout, startStr)
	if err != nil {
		return err
	}

	endTime, err := time.Parse(layout, endStr)
	if err != nil {
		return err
	}

	if f.editMode {
		return f.eventList.Update(f.editingID, title, desc, location, startTime, endTime)
	}

	return f.eventList.Add(title, desc, location, startTime, endTime)
}

func (f *EventForm) LoadForEdit(event *models.Event) {
	f.editMode = true
	f.editingID = event.ID
	f.titleInput.SetValue(event.Title)
	f.descInput.SetValue(event.Content)
	f.locationInput.SetValue(event.Location)

	layout := "2006-01-02 15:04"
	f.startInput.SetValue(event.StartTime.Format(layout))
	f.endInput.SetValue(event.EndTime.Format(layout))

	f.isActive = true
	f.titleInput.Focus()
}
