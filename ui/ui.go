package ui

import (
	"strings"
	"time"

	"cykl/core"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	// We don't need to refresh the ui as often as the grid.
	// It saves some cpu. Right now we run it at 30 fps.
	refreshFrequency = 33 * time.Millisecond

	blinkFrequency = 500 * time.Millisecond

	controlsHeight = 4
)

// tickMsg is a message that triggers ui rrefresh
type tickMsg time.Time

// blinkMsg is a message that triggers blinking ui elements
type blinkMsg time.Time

type mainModel struct {
	grid       *core.Grid
	cursorX    int
	cursorY    int
	selectionX int
	selectionY int
	width      int
	height     int
	insert     bool
	blink      bool
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

func blink() tea.Cmd {
	return tea.Tick(blinkFrequency, func(t time.Time) tea.Msg {
		return blinkMsg(t)
	})
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tick(), blink())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.grid.Resize(m.width/2, m.height-controlsHeight)
		if m.cursorX > m.grid.Width-1 {
			m.cursorX = m.grid.Width - 1
		}
		if m.cursorY > m.grid.Height-1 {
			m.cursorY = m.grid.Height - 1
		}
		return m, nil

	case tickMsg:
		return m, tick()

	case blinkMsg:
		m.blink = !m.blink
		return m, blink()

	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			m.grid.TogglePlay()
			return m, nil
		case "ctrl+d":
			m.grid.Reset()
			return m, nil
		case "up", "right", "down", "left":
			if m.insert {
				m.grid.Node(m.cursorX, m.cursorY).
					SetDirection(core.DirectionFromString(msg.String()))
			} else {
				m.blink = true
				m.cursorX, m.cursorY = moveCursor(
					msg.String(), m.cursorX, m.cursorY,
					0, m.grid.Width-1, 0, m.grid.Height-1,
				)
				m.selectionX, m.selectionY = moveCursor(
					msg.String(), m.selectionX, m.selectionY,
					m.cursorX, m.grid.Width-1, m.cursorY, m.grid.Height-1,
				)
			}
			return m, nil
		case "shift+up", "shift+right", "shift+down", "shift+left":
			dir := strings.Replace(msg.String(), "shift+", "", 1)
			m.selectionX, m.selectionY = moveCursor(
				dir, m.selectionX, m.selectionY,
				m.cursorX, m.grid.Width-1, m.cursorY, m.grid.Height-1,
			)
			return m, nil
		case "i", "s":
			m.grid.AddNodeFromSymbol(msg.String(), m.cursorX, m.cursorY)
			m.insert = true
			return m, nil
		case "backspace":
			m.insert = false
			m.grid.RemoveNodes(m.cursorX, m.cursorY, m.selectionX, m.selectionY)
			return m, nil
		case "enter":
			m.insert = !m.insert
			return m, nil
		case "n":
			m.grid.Update()
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
	for y, line := range m.grid.Nodes() {
		var nodes []string
		for x, node := range line {
			nodes = append(nodes, m.renderNode(node, x, y))
		}
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left, nodes...))
	}
	lines = append(lines, m.renderControl())
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func moveCursor(dir string, x, y, minX, maxX, minY, maxY int) (int, int) {
	newX, newY := core.DirectionFromString(dir).
		NextPosition(x, y)
	return clamp(newX, minX, maxX), clamp(newY, minY, maxY)
}
