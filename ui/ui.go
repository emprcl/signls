package ui

import (
	"strings"
	"time"

	"cykl/core"
	"cykl/ui/param"

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
	params     []param.Param
	cursorX    int
	cursorY    int
	selectionX int
	selectionY int
	width      int
	height     int
	param      int
	edit       bool
	blink      bool
	mute       bool
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
		case "tab":
			if m.edit {
				m.moveParam(msg)
			}
			return m, nil
		case "up", "right", "down", "left":
			if m.edit {
				m.moveParam(msg)
				return m, nil
			}
			m.blink = true
			m.cursorX, m.cursorY = moveCursor(
				msg.String(), m.cursorX, m.cursorY,
				0, m.grid.Width-1, 0, m.grid.Height-1,
			)
			m.selectionX, m.selectionY = moveCursor(
				msg.String(), m.selectionX, m.selectionY,
				m.cursorX, m.grid.Width-1, m.cursorY, m.grid.Height-1,
			)
			return m, nil
		case "shift+up", "shift+right", "shift+down", "shift+left":
			dir := strings.Replace(msg.String(), "shift+", "", 1)
			if m.edit {
				m.handleParamEdit(dir)
				return m, nil
			}
			m.selectionX, m.selectionY = moveCursor(
				dir, m.selectionX, m.selectionY,
				m.cursorX, m.grid.Width-1, m.cursorY, m.grid.Height-1,
			)
			return m, nil
		case "b", "s":
			m.grid.AddNodeFromSymbol(msg.String(), m.cursorX, m.cursorY)
			m.params = param.NewParamsForNode(m.selectedNode())
			m.edit = true
			return m, nil
		case "m":
			m.grid.ToggleNodeMutes(m.cursorX, m.cursorY, m.selectionX, m.selectionY)
			return m, nil
		case "M":
			m.grid.SetAllNodeMutes(!m.mute)
			m.mute = !m.mute
			return m, nil
		case "backspace":
			m.edit = false
			m.grid.RemoveNodes(m.cursorX, m.cursorY, m.selectionX, m.selectionY)
			return m, nil
		case "enter":
			if m.selectedNode() == nil {
				return m, nil
			}
			m.edit = !m.edit
			if m.edit {
				m.params = param.NewParamsForNode(m.selectedNode())
				m.param = 0
			}
			return m, nil
		case "ctrl+up":
			m.grid.SetTempo(m.grid.Tempo() + 1)
			return m, nil
		case "ctrl+down":
			m.grid.SetTempo(m.grid.Tempo() - 1)
			return m, nil
		case "=":
			m.grid.CycleMidiDevice()
			return m, nil
		case "ctrl+c":
			m.grid.CopyOrCut(m.cursorX, m.cursorY, m.selectionX, m.selectionY, false)
			return m, nil
		case "ctrl+x":
			m.grid.CopyOrCut(m.cursorX, m.cursorY, m.selectionX, m.selectionY, true)
			return m, nil
		case "ctrl+v":
			m.grid.Paste(m.cursorX, m.cursorY, m.selectionX, m.selectionY)
			return m, nil
		case "esc":
			m.edit = false
			m.param = 0
			m.selectionX = m.cursorX
			m.selectionY = m.cursorY
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
func (m mainModel) handleParamEdit(key string) {
	if len(m.params) < m.param+1 {
		return
	}

	switch p := m.params[m.param].(type) {
	case param.Direction:
		p.SetFromKeyString(key)
	}

	switch key {
	case "up", "right":
		m.params[m.param].Increment()
	case "down", "left":
		m.params[m.param].Decrement()
	}
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

func (m *mainModel) moveParam(msg tea.KeyMsg) {
	if len(m.params) == 0 {
		return
	}
	switch msg.String() {
	case "right":
		if m.param+1 >= len(m.params) {
			return
		}
		m.param++
	case "left":
		if m.param-1 < 0 {
			return
		}
		m.param--
	case "tab":
		m.param = (m.param + 1) % len(m.params)
	}
}

func moveCursor(dir string, x, y, minX, maxX, minY, maxY int) (int, int) {
	var newX, newY int
	switch dir {
	case "up":
		newX, newY = x, y-1
	case "right":
		newX, newY = x+1, y
	case "down":
		newX, newY = x, y+1
	case "left":
		newX, newY = x-1, y
	default:
		newX, newY = 0, 0
	}
	return clamp(newX, minX, maxX), clamp(newY, minY, maxY)
}
