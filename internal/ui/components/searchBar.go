package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FilterType int

const (
	FilterNone FilterType = iota
	FilterPending
	FilterCompleted
	FilterHighPriority
	FilterMediumPriority
	FilterLowPriority
	FilterToday
	FilterOverdue
)

type SearchBar struct {
	input  textinput.Model
	active bool
	query  string
	filter FilterType
	width  int
}

func NewSearchBar() *SearchBar {
	input := textinput.New()
	input.Placeholder = "Search... (Esc to close, Tab for filters)"
	input.CharLimit = 100
	input.Width = 60

	return &SearchBar{
		input:  input,
		active: false,
		filter: FilterNone,
		width:  80,
	}
}

func (s *SearchBar) Activate() {
	s.active = true
	s.input.Focus()
}

func (s *SearchBar) Deactivate() {
	s.active = false
	s.input.Blur()
	s.query = ""
	s.input.SetValue("")
	s.filter = FilterNone
}

func (s *SearchBar) IsActive() bool {
	return s.active
}

func (s *SearchBar) GetQuery() string {
	return s.query
}

func (s *SearchBar) GetFilter() FilterType {
	return s.filter
}

func (s *SearchBar) SetWidth(width int) {
	s.width = width
	s.input.Width = width - 20
}

func (s *SearchBar) Update(msg tea.Msg) (*SearchBar, tea.Cmd) {
	if !s.active {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			s.Deactivate()
			return s, nil
		case "enter":
			// Apply search and close
			s.query = s.input.Value()
			s.active = false
			s.input.Blur()
			return s, nil
		case "tab":
			// Cycle through filters
			s.filter = (s.filter + 1) % 8
			return s, nil
		}
	}

	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	s.query = s.input.Value()

	return s, cmd
}

func (s *SearchBar) View() string {
	if !s.active {
		return ""
	}

	filterLabel := ""
	filterColor := "240"

	switch s.filter {
	case FilterPending:
		filterLabel = "Pending"
		filterColor = "214"
	case FilterCompleted:
		filterLabel = "Completed"
		filterColor = "46"
	case FilterHighPriority:
		filterLabel = "High Priority"
		filterColor = "196"
	case FilterMediumPriority:
		filterLabel = "Medium Priority"
		filterColor = "214"
	case FilterLowPriority:
		filterLabel = "Low Priority"
		filterColor = "45"
	case FilterToday:
		filterLabel = "Due Today"
		filterColor = "45"
	case FilterOverdue:
		filterLabel = "Overdue"
		filterColor = "196"
	default:
		filterLabel = "No Filter"
		filterColor = "240"
	}

	filterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(filterColor)).
		Bold(true).
		Padding(0, 1)

	searchBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(s.width - 4)

	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render("🔍 Search & Filter"),
		s.input.View(),
		lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Filter: "),
			filterStyle.Render(filterLabel),
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(" (Tab to change)"),
		),
	)

	return searchBox.Render(content)
}

// Match checks if a string matches the search query
func (s *SearchBar) Match(text string) bool {
	if s.query == "" {
		return true
	}
	return strings.Contains(strings.ToLower(text), strings.ToLower(s.query))
}
