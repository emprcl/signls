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
	gridParams []param.Param
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
		grid:       grid,
		gridParams: param.NewParamsForGrid(grid),
		cursorX:    1,
		cursorY:    1,
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
				msg.String(), 1, m.cursorX, m.cursorY,
				0, m.grid.Width-1, 0, m.grid.Height-1,
			)
			m.selectionX, m.selectionY = moveCursor(
				msg.String(), 1, m.selectionX, m.selectionY,
				m.cursorX, m.grid.Width-1, m.cursorY, m.grid.Height-1,
			)
			m.params = param.NewParamsForNode(m.grid, m.selectedNode())
			return m, nil
		case "shift+up", "shift+right", "shift+down", "shift+left":
			dir := strings.Replace(msg.String(), "shift+", "", 1)
			m.selectionX, m.selectionY = moveCursor(
				dir, 1, m.selectionX, m.selectionY,
				m.cursorX, m.grid.Width-1, m.cursorY, m.grid.Height-1,
			)
			return m, nil
		case "ctrl+up", "ctrl+right", "ctrl+down", "ctrl+left":
			dir := strings.Replace(msg.String(), "ctrl+", "", 1)
			if !m.edit && m.selectedNode() != nil {
				param.NewDirection(m.selectedNode()).SetFromKeyString(dir)
				return m, nil
			}
			m.handleParamEdit(dir)
			return m, nil
		case "b", "s", "c":
			m.grid.AddNodeFromSymbol(msg.String(), m.cursorX, m.cursorY)
			m.params = param.NewParamsForNode(m.grid, m.selectedNode())
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
			return m, nil
		case "!":
			if !m.grid.Playing {
				return m, nil
			}
			if _, ok := m.selectedNode().(*core.Emitter); !ok {
				return m, nil
			}
			m.selectedNode().(*core.Emitter).Arm()
			m.selectedNode().(*core.Emitter).Trig(m.grid.Key, m.grid.Scale, m.grid.Pulse())
			return m, nil
		case "*":
			if m.edit {
				return m, nil
			}
			param.Get("root", m.gridParams).Increment()
			return m, nil
		case "ù":
			if m.edit {
				return m, nil
			}
			param.Get("root", m.gridParams).Decrement()
			return m, nil
		case "µ":
			if m.edit {
				return m, nil
			}
			param.Get("scale", m.gridParams).Increment()
			return m, nil
		case "%":
			if m.edit {
				return m, nil
			}
			param.Get("scale", m.gridParams).Decrement()
			return m, nil
		case "=":
			m.grid.SetTempo(m.grid.Tempo() + 1)
			return m, nil
		case ")":
			m.grid.SetTempo(m.grid.Tempo() - 1)
			return m, nil
		case "f2":
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
			m.params = param.NewParamsForNode(m.grid, m.selectedNode())
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
			m.grid.Reset()
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
		return
	}

	switch key {
	case "up":
		m.params[m.param].Increment()
	case "down":
		m.params[m.param].Decrement()
	case "left":
		m.params[m.param].Left()
	case "right":
		m.params[m.param].Right()
	}

	switch p := m.params[m.param].(type) {
	case param.Key:
		if m.grid.Playing {
			return
		}
		p.Preview()
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

func moveCursor(dir string, speed, x, y, minX, maxX, minY, maxY int) (int, int) {
	var newX, newY int
	switch dir {
	case "up":
		newX, newY = x, y-speed
	case "right":
		newX, newY = x+speed, y
	case "down":
		newX, newY = x, y+speed
	case "left":
		newX, newY = x-speed, y
	default:
		newX, newY = 0, 0
	}
	return clamp(newX, minX, maxX), clamp(newY, minY, maxY)
}
