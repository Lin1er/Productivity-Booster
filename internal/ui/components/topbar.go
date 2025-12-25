package components

import (
	tea "github.com/charmbracelet/bubbletea"

	"prodBooster/internal/models"
	"prodBooster/internal/ui/styles"

	"github.com/charmbracelet/lipgloss"
)

type TopBar struct {
	currentPage models.PageType
	navList []models.PageInfo
	width  int
	height int
}

func NewTopBar(currentPage_ models.PageType) *TopBar {
	return &TopBar{
		currentPage: currentPage_,
		navList: models.GetPages(),
		width:  80,
		height: 1,
	}
}

func (tb *TopBar) Init() tea.Cmd {
	return nil
}

func (tb *TopBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return tb, nil
}

func (tb *TopBar) SetSize(width, height int) {
	tb.width = width
	tb.height = height
}

func (tb *TopBar) View() string {

	selectedStyle := lipgloss.NewStyle().
		Foreground(styles.ColorPrimary).
		Bold(true).
		Underline(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(styles.ColorMuted)

	navItems := []string{}
	for _, page := range tb.navList {
		if page.Type == tb.currentPage {
			navItems = append(navItems, selectedStyle.Render(" "+page.Title+" "))
			continue
		}
		navItems = append(navItems, normalStyle.Render(" "+page.Title+" "))
	}

	topBarStyle := lipgloss.NewStyle().
		Padding(0,1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Width(tb.width).
		Height(tb.height)

	return topBarStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, navItems...))
}
