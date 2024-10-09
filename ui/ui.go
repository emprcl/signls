package ui

import (
	"time"

	"cykl/core/common"
	"cykl/core/field"
	"cykl/core/node"
	"cykl/filesystem"
	"cykl/ui/param"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	grid       *field.Grid
	viewport   viewport
	keymap     keyMap
	help       help.Model
	params     []param.Param
	gridParams []param.Param
	cursorX    int
	cursorY    int
	selectionX int
	selectionY int
	param      int
	edit       bool
	blink      bool
	mute       bool
}

// New creates a new mainModel that hols the ui state. It takes a new grid.
// Check the core package.
func New(config filesystem.Configuration, grid *field.Grid) tea.Model {
	model := mainModel{
		grid:       grid,
		keymap:     newKeyMap(config.KeyMap),
		help:       help.New(),
		gridParams: param.NewParamsForGrid(grid),
		cursorX:    1,
		cursorY:    1,
		selectionX: 1,
		selectionY: 1,
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
		m.help.Width = msg.Width
		m.viewport.Width = msg.Width / 2
		m.viewport.Height = msg.Height - controlsHeight - 1
		if m.viewport.Width > m.grid.Width || m.viewport.Height > m.grid.Height {
			m.grid.Resize(m.viewport.Width, m.viewport.Height)
		}
		m.viewport.Update(m.cursorX, m.cursorY, m.grid.Width, m.grid.Height)
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
		switch {
		case key.Matches(msg, m.keymap.Play):
			m.grid.TogglePlay()
			return m, nil
		case key.Matches(msg, m.keymap.Up, m.keymap.Right, m.keymap.Down, m.keymap.Left):
			dir := m.keymap.Direction(msg)
			if m.edit {
				m.moveParam(dir)
				return m, nil
			}
			m.blink = true
			m.cursorX, m.cursorY = moveCursor(
				dir, 1, m.cursorX, m.cursorY,
				0, m.grid.Width-1, 0, m.grid.Height-1,
			)
			m.selectionX, m.selectionY = moveCursor(
				dir, 1, m.selectionX, m.selectionY,
				m.cursorX, m.grid.Width-1, m.cursorY, m.grid.Height-1,
			)
			m.params = param.NewParamsForNodes(m.grid, m.selectedEmitters())
			m.viewport.Update(m.cursorX, m.cursorY, m.grid.Width, m.grid.Height)
			return m, nil
		case key.Matches(msg, m.keymap.SelectionUp, m.keymap.SelectionRight, m.keymap.SelectionDown, m.keymap.SelectionLeft):
			dir := m.keymap.Direction(msg)
			if m.edit {
				m.handleParamAltEdit(dir)
				return m, nil
			}
			m.selectionX, m.selectionY = moveCursor(
				dir, 1, m.selectionX, m.selectionY,
				m.cursorX, m.grid.Width-1, m.cursorY, m.grid.Height-1,
			)
			m.params = param.NewParamsForNodes(m.grid, m.selectedEmitters())
			return m, nil
		case key.Matches(msg, m.keymap.EditUp, m.keymap.EditRight, m.keymap.EditDown, m.keymap.EditLeft):
			dir := m.keymap.Direction(msg)
			if !m.edit {
				param.NewDirection(m.selectedEmitters()).SetFromKeyString(dir)
				return m, nil
			}
			m.handleParamEdit(dir)
			return m, nil
		case key.Matches(msg, m.keymap.AddBang, m.keymap.AddRelay, m.keymap.AddCycle, m.keymap.AddDice, m.keymap.AddToll, m.keymap.AddEuclid, m.keymap.AddZone, m.keymap.AddPass, m.keymap.AddHole):
			m.grid.AddNodeFromSymbol(m.keymap.EmitterSymbol(msg), m.cursorX, m.cursorY)
			m.params = param.NewParamsForNodes(m.grid, m.selectedEmitters())
			return m, nil
		case key.Matches(msg, m.keymap.MuteNode):
			m.grid.ToggleNodeMutes(m.cursorX, m.cursorY, m.selectionX, m.selectionY)
			return m, nil
		case key.Matches(msg, m.keymap.MuteAllNode):
			m.grid.SetAllNodeMutes(!m.mute)
			m.mute = !m.mute
			return m, nil
		case key.Matches(msg, m.keymap.RemoveNode):
			m.edit = false
			m.grid.RemoveNodes(m.cursorX, m.cursorY, m.selectionX, m.selectionY)
			return m, nil
		case key.Matches(msg, m.keymap.EditNode):
			if len(m.selectedEmitters()) == 0 {
				return m, nil
			}
			if len(m.params) < m.param+1 {
				m.param = 0
			}
			m.edit = !m.edit
			if m.edit {
				m.params = param.NewParamsForNodes(m.grid, m.selectedEmitters())
			}
			return m, nil
		case key.Matches(msg, m.keymap.TriggerNode):
			if !m.grid.Playing {
				return m, nil
			}
			if _, ok := m.selectedNode().(*node.Emitter); !ok {
				return m, nil
			}
			m.selectedNode().(*node.Emitter).Arm()
			m.selectedNode().(*node.Emitter).Trig(m.grid.Key, m.grid.Scale, common.NONE, m.grid.Pulse())
			return m, nil
		case key.Matches(msg, m.keymap.RootNoteUp):
			if m.edit {
				return m, nil
			}
			param.Get("root", m.gridParams).Up()
			return m, nil
		case key.Matches(msg, m.keymap.RootNoteDown):
			if m.edit {
				return m, nil
			}
			param.Get("root", m.gridParams).Down()
			return m, nil
		case key.Matches(msg, m.keymap.ScaleUp):
			if m.edit {
				return m, nil
			}
			param.Get("scale", m.gridParams).Up()
			return m, nil
		case key.Matches(msg, m.keymap.ScaleDown):
			if m.edit {
				return m, nil
			}
			param.Get("scale", m.gridParams).Down()
			return m, nil
		case key.Matches(msg, m.keymap.TempoUp):
			m.grid.SetTempo(m.grid.Tempo() + 1)
			return m, nil
		case key.Matches(msg, m.keymap.TempoDown):
			m.grid.SetTempo(m.grid.Tempo() - 1)
			return m, nil
		case key.Matches(msg, m.keymap.SelectMidiDevice):
			m.grid.CycleMidiDevice()
			return m, nil
		case key.Matches(msg, m.keymap.Copy):
			m.grid.CopyOrCut(m.cursorX, m.cursorY, m.selectionX, m.selectionY, false)
			return m, nil
		case key.Matches(msg, m.keymap.Cut):
			m.grid.CopyOrCut(m.cursorX, m.cursorY, m.selectionX, m.selectionY, true)
			return m, nil
		case key.Matches(msg, m.keymap.Paste):
			m.grid.Paste(m.cursorX, m.cursorY, m.selectionX, m.selectionY)
			m.params = param.NewParamsForNodes(m.grid, m.selectedEmitters())
			return m, nil
		case key.Matches(msg, m.keymap.Cancel):
			m.edit = false
			m.param = 0
			m.selectionX = m.cursorX
			m.selectionY = m.cursorY
			m.help.ShowAll = false
			return m, nil
		case key.Matches(msg, m.keymap.FitGridToWindow):
			m.cursorX, m.cursorY = 1, 1
			m.selectionX, m.selectionY = m.cursorX, m.cursorY
			m.grid.Resize(m.viewport.Width, m.viewport.Height)
			m.viewport.Update(m.cursorX, m.cursorY, m.grid.Width, m.grid.Height)
			return m, nil
		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, tea.ClearScreen
		case key.Matches(msg, m.keymap.Quit):
			m.grid.Reset()
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m mainModel) View() string {
	help := lipgloss.NewStyle().
		MarginLeft(2).
		MarginTop(1).
		Render(m.help.View(m.keymap))

	if m.help.ShowAll {
		return help
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderGrid(),
		help,
	)
}

func (m mainModel) handleParamEdit(dir string) {
	if len(m.params) < m.param+1 {
		return
	}

	switch dir {
	case "up":
		m.params[m.param].Up()
	case "down":
		m.params[m.param].Down()
	case "left":
		m.params[m.param].Left()
	case "right":
		m.params[m.param].Right()
	}

	switch p := m.params[m.param].(type) {
	case *param.Key:
		if m.grid.Playing {
			return
		}
		p.Preview()
	}
}

func (m mainModel) handleParamAltEdit(dir string) {
	if len(m.params) < m.param+1 {
		return
	}

	switch dir {
	case "up":
		m.params[m.param].AltUp()
	case "down":
		m.params[m.param].AltDown()
	case "left":
		m.params[m.param].AltLeft()
	case "right":
		m.params[m.param].AltRight()
	}
}

func (m mainModel) renderGrid() string {
	var lines []string
	for y := m.viewport.offsetY; y < m.viewport.offsetY+m.viewport.Height; y++ {
		var nodes []string
		for x := m.viewport.offsetX; x < m.viewport.offsetX+m.viewport.Width; x++ {
			nodes = append(nodes, m.renderNode(m.grid.Nodes()[y][x], x, y))
		}
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left, nodes...))
	}
	lines = append(lines, m.renderControl())
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m *mainModel) moveParam(dir string) {
	if len(m.params) == 0 {
		return
	}
	switch dir {
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
