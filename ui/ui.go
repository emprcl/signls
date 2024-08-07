package ui

import (
	"time"

	"cykl/core"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	// We don't need to refresh the ui as often as the grid.
	// It saves some cpu. Right now we run it at 30 fps.
	refreshFrequency = 33 * time.Millisecond
)

// tickMsg is a message that triggers ui rrefresh
type tickMsg time.Time

type mainModel struct {
	grid   *core.Grid
	width  int
	height int
}

// New creates a new mainModel that hols the ui state. It takes a new grid.
// Check the core package.
func New(grid *core.Grid) tea.Model {
	model := mainModel{
		grid: grid,
	}
	return model
}

func tick() tea.Cmd {
	return tea.Tick(refreshFrequency, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tick())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.grid.Resize(m.width/2, m.height)
		return m, nil

	case tickMsg:
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			m.grid.Playing = !m.grid.Playing
			return m, nil
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m mainModel) View() string {
	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderGrid(),
	)

	// Cleanup gibber
	cleanup := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height - lipgloss.Height(mainView)).
		Render("")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainView,
		cleanup,
	)
}

func (m mainModel) renderGrid() string {
	var lines []string
	for _, line := range m.grid.Nodes() {
		var nodes []string
		for _, node := range line {
			nodes = append(nodes, renderNode(node))
		}
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left, nodes...))
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}
