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

// eventItem implements list.Item interface
type eventItem struct {
	event *models.Event
}

func (e eventItem) Title() string {
	return "📅 " + e.event.Title
}

func (e eventItem) Description() string {
	location := ""
	if e.event.Location != "" {
		location = " @ " + e.event.Location
	}
	return e.event.StartTime.Format("Mon, Jan 2 15:04") + location
}

func (e eventItem) FilterValue() string {
	return e.event.Title
}

// Custom delegate for colored event items
type eventDelegate struct{}

func (d eventDelegate) Height() int                             { return 1 }
func (d eventDelegate) Spacing() int                            { return 0 }
func (d eventDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d eventDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	event, ok := item.(eventItem)
	if !ok {
		return
	}

	// Color by time status
	var style lipgloss.Style
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	if event.event.StartTime.Before(now) {
		// Past event - dim gray
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Strikethrough(true)
	} else if event.event.StartTime.After(todayStart) && event.event.StartTime.Before(todayEnd) {
		// Today - bright yellow/orange
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	} else if event.event.StartTime.Before(now.Add(7 * 24 * time.Hour)) {
		// This week - cyan
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("45"))
	} else {
		// Future - normal blue
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("147"))
	}

	// Highlight selected item
	if index == m.Index() {
		style = style.Background(lipgloss.Color("238"))
	}

	fmt.Fprint(w, style.Render(event.Title()))
}

type CalendarPage struct {
	EventList    *models.EventList
	form         *components.EventForm
	searchBar    *components.SearchBar
	list         list.Model
	width        int
	height       int
	sidebarWidth int
}

func NewCalendarPage(eventList_ *models.EventList) *CalendarPage {
	// Sort events initially - today > this week > future > past
	sortEvents(eventList_.Events)

	// Create list items from events
	items := make([]list.Item, len(eventList_.Events))
	for i, event := range eventList_.Events {
		items[i] = eventItem{event: event}
	}

	// Use custom delegate for colored rendering
	delegate := eventDelegate{}

	l := list.New(items, delegate, 0, 0)
	l.Title = "📅 My Calendar"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false) // We use our own search

	return &CalendarPage{
		EventList:    eventList_,
		form:         components.NewEventForm(eventList_),
		searchBar:    components.NewSearchBar(),
		list:         l,
		width:        80,
		height:       24,
		sidebarWidth: 40,
	}
}

func (p *CalendarPage) Init() tea.Cmd {
	return nil
}

func (p *CalendarPage) Update(msg tea.Msg) (Page, tea.Cmd) {
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
			// Create new event
			p.form.Activate()

		case "e":
			// Edit selected event
			if item, ok := p.list.SelectedItem().(eventItem); ok {
				p.form.LoadForEdit(item.event)
			}

		case "d", "delete":
			// Delete selected event
			if item, ok := p.list.SelectedItem().(eventItem); ok {
				if err := p.EventList.Remove(item.event.ID); err != nil {
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

	// Sync EventList.Selected with list's selected index
	p.EventList.Selected = p.list.Index()

	return p, cmd
}

// updateListItems refreshes the list with current events and filters
func (p *CalendarPage) updateListItems() {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	filteredEvents := []*models.Event{}
	for _, event := range p.EventList.Events {
		// Apply text search
		if !p.searchBar.Match(event.Title + " " + event.Content + " " + event.Location) {
			continue
		}

		// Apply filter (reuse filter types for calendar context)
		filter := p.searchBar.GetFilter()
		switch filter {
		case components.FilterToday:
			// Events happening today
			if !event.StartTime.After(todayStart) || !event.StartTime.Before(todayEnd) {
				continue
			}
		case components.FilterOverdue:
			// Past events
			if !event.StartTime.Before(now) {
				continue
			}
		case components.FilterPending:
			// Future events
			if !event.StartTime.After(now) {
				continue
			}
		}

		filteredEvents = append(filteredEvents, event)
	}

	// Sort events by time: today > this week > future > past
	sortEvents(filteredEvents)

	// Convert to list items
	items := make([]list.Item, len(filteredEvents))
	for i, event := range filteredEvents {
		items[i] = eventItem{event: event}
	}

	p.list.SetItems(items)
}

// sortEvents sorts by priority: today > this week > future > past
func sortEvents(events []*models.Event) {
	now := time.Now()
	for i := 0; i < len(events); i++ {
		for j := i + 1; j < len(events); j++ {
			if shouldSwapEvents(events[i], events[j], now) {
				events[i], events[j] = events[j], events[i]
			}
		}
	}
}

func shouldSwapEvents(a, b *models.Event, now time.Time) bool {
	scoreA := getEventPriority(a, now)
	scoreB := getEventPriority(b, now)

	if scoreA == scoreB {
		// Same priority, sort by time (earlier first for future, later first for past)
		if scoreA >= 90 { // Future events
			return a.StartTime.After(b.StartTime)
		}
		return b.StartTime.After(a.StartTime)
	}

	return scoreB > scoreA
}

func getEventPriority(event *models.Event, now time.Time) int {
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	if event.StartTime.After(todayStart) && event.StartTime.Before(todayEnd) {
		return 100 // Today - highest priority
	}
	if event.StartTime.After(now) && event.StartTime.Before(now.Add(7*24*time.Hour)) {
		return 90 // This week
	}
	if event.StartTime.After(now) {
		return 80 // Future
	}
	return 50 // Past - lower priority
}

func (p *CalendarPage) SetSize(width, height int) {
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

func (p *CalendarPage) View() string {
	// If form is active, show form overlay
	if p.form.IsActive() {
		return lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, p.form.View())
	}

	// If search is active, show search overlay
	if p.searchBar.IsActive() {
		return lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Top, p.searchBar.View())
	}

	topBar := components.NewTopBar(models.PageTypeCalendar())
	topBar.SetSize(p.width, 1)

	// Sidebar with list
	sidebarStyle := lipgloss.NewStyle().
		Width(p.sidebarWidth).
		Height(p.height - 6).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	sidebar := sidebarStyle.Render(p.list.View())

	// Content pane with selected event detail
	contentWidth := p.width - p.sidebarWidth - 4
	contentStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Height(p.height-6).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	var content string
	if len(p.EventList.Events) == 0 {
		content = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("✨ No events scheduled!\n\nPress 'n' to plan something 📅")
	} else if p.EventList.Selected < len(p.EventList.Events) {
		event := p.EventList.Events[p.EventList.Selected]

		// Title
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("213"))

		// Time information with colors
		now := time.Now()
		var timeColor string
		var timePrefix string
		if event.StartTime.Before(now) {
			timeColor = "240"
			timePrefix = "⏳ Past Event"
		} else {
			todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			todayEnd := todayStart.Add(24 * time.Hour)
			if event.StartTime.After(todayStart) && event.StartTime.Before(todayEnd) {
				timeColor = "45"
				timePrefix = "🎯 Happening Today!"
			} else {
				timeColor = "213"
				timePrefix = "📍 Coming Up"
			}
		}

		startTimeStr := "🕐 " + event.StartTime.Format("Monday, Jan 2, 2006 at 3:04 PM")
		endTimeStr := ""
		if !event.EndTime.IsZero() {
			endTimeStr = "🕑 Ends " + event.EndTime.Format("Monday, Jan 2 at 3:04 PM")
		}

		locationStr := ""
		if event.Location != "" {
			locationStr = "📍 " + event.Location
		}

		contentParts := []string{
			titleStyle.Render(event.Title),
			"",
			lipgloss.NewStyle().Foreground(lipgloss.Color(timeColor)).Render(timePrefix),
			lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(startTimeStr),
		}

		if endTimeStr != "" {
			contentParts = append(contentParts,
				lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(endTimeStr))
		}

		if locationStr != "" {
			contentParts = append(contentParts, "",
				lipgloss.NewStyle().Foreground(lipgloss.Color("45")).Render(locationStr))
		}

		if event.Content != "" {
			contentParts = append(contentParts, "",
				lipgloss.NewStyle().
					Foreground(lipgloss.Color("252")).
					Width(contentWidth-4).
					Render(event.Content))
		}

		content = lipgloss.JoinVertical(lipgloss.Left, contentParts...)
	} else {
		content = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("👈 Pick an event to see what's happening")
	}

	// Add help text
	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("✨ n: new event • e: edit • d: delete • /: search • ↑/↓: browse • q: quit")

	// Combine sidebar and content
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, contentStyle.Render(content))

	// Build the view
	return lipgloss.JoinVertical(lipgloss.Left,
		topBar.View(),
		mainContent,
		helpText,
	)
}

func (p *CalendarPage) IsFormActive() bool {
	return p.form.IsActive() || p.searchBar.IsActive()
}
