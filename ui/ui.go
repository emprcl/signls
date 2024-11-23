package ui

import (
	"fmt"
	"time"

	"signls/core/common"
	"signls/core/field"
	"signls/core/node"
	"signls/filesystem"
	"signls/ui/param"
	"signls/ui/util"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	// We don't need to refresh the ui as often as the grid.
	// It saves some cpu. Right now we run it at 30 fps.
	refreshFrequency = 33 * time.Millisecond

	blinkFrequency = 500 * time.Millisecond

	controlsHeight = 4

	helpHeader = "signls %s - docs: https://empr.cl/signls/"
)

// mode is a representation of a ui mode
type mode uint8

const (
	// MOVE mode allows moving the cursor on the grid
	MOVE mode = iota
	// EDIT mode allows node parameters edits
	EDIT
	// CONFIG mode allows global parameters edits
	CONFIG
	// BANK mode allows bank grids selection
	BANK
)

// tickMsg is a message that triggers ui rrefresh
type tickMsg time.Time

// blinkMsg is a message that triggers blinking ui elements
type blinkMsg time.Time

// saveMsg is a message that notify a successfull save
type saveMsg bool

type mainModel struct {
	bank          *filesystem.Bank
	grid          *field.Grid
	viewport      viewport
	keymap        keyMap
	help          help.Model
	input         textinput.Model
	params        [][]param.Param
	gridParams    []param.Param
	bankClipboard filesystem.Grid
	mode          mode
	version       string
	cursorX       int
	cursorY       int
	selectionX    int
	selectionY    int
	selectedGrid  int
	param         int
	paramPage     int
	blink         bool
	mute          bool
}

// New creates a new mainModel that hols the ui state. It takes a new grid.
// Check the core package.
func New(config filesystem.Configuration, grid *field.Grid, bank *filesystem.Bank) tea.Model {
	ti := textinput.New()
	ti.CharLimit = 10
	ti.Width = 12
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("190"))
	model := mainModel{
		bank:       bank,
		grid:       grid,
		keymap:     newKeyMap(config.KeyMap),
		help:       help.New(),
		input:      ti,
		gridParams: param.NewParamsForGrid(grid),
		cursorX:    1,
		cursorY:    1,
		selectionX: 1,
		selectionY: 1,

		version: config.Version(),
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

func save(m mainModel) tea.Cmd {
	return func() tea.Msg {
		m.grid.Save(m.bank)
		return saveMsg(true)
	}
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tick(), blink())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		return m.windowResize(msg.Width, msg.Height), nil

	case tickMsg:
		return m, tick()

	case blinkMsg:
		m.blink = !m.blink
		m.input.Cursor.Blink = !m.input.Cursor.Blink
		return m, blink()

	case tea.KeyMsg:
		if m.input.Focused() {
			var cmd tea.Cmd
			switch {
			case key.Matches(msg, m.keymap.EditNode):
				m.input.Blur()
				m.activeParam().SetEditValue(m.input.Value())
				return m, nil
			case key.Matches(msg, m.keymap.Cancel, m.keymap.EditInput):
				m.input.Blur()
				return m, nil
			case key.Matches(msg, m.keymap.Quit):
				break
			default:
				m.input, cmd = m.input.Update(msg)
				return m, cmd
			}
		}

		switch {
		case key.Matches(msg, m.keymap.EditInput):
			if m.mode != EDIT {
				return m, nil
			}
			m.input.Focus()
			m.input.Reset()
			return m, nil
		case key.Matches(msg, m.keymap.Play):
			m.grid.TogglePlay()
			return m, nil
		case key.Matches(msg, m.keymap.Up, m.keymap.Right, m.keymap.Down, m.keymap.Left):
			dir := m.keymap.Direction(msg)
			if m.mode == BANK {
				m.moveBankGrid(dir)
				return m, nil
			}
			if m.mode == EDIT {
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
			if m.mode == EDIT {
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
			if m.mode != EDIT {
				param.NewDirection(m.selectedEmitters()).SetFromKeyString(dir)
				return m, save(m)
			}
			m.handleParamEdit(dir)
			return m, save(m)
		case key.Matches(msg, m.keymap.AddBang, m.keymap.AddSpread, m.keymap.AddCycle, m.keymap.AddDice, m.keymap.AddToll, m.keymap.AddEuclid, m.keymap.AddZone, m.keymap.AddPass, m.keymap.AddHole):
			m.grid.AddNodeFromSymbol(m.keymap.EmitterSymbol(msg), m.cursorX, m.cursorY)
			newParams := param.NewParamsForNodes(m.grid, m.selectedEmitters())
			if len(newParams) < m.paramPage+1 {
				m.paramPage = 0
			}
			if len(newParams[m.paramPage]) < m.param+1 {
				m.param = 0
			}
			m.params = newParams
			return m, save(m)
		case key.Matches(msg, m.keymap.MuteNode):
			m.grid.ToggleNodeMutes(m.cursorX, m.cursorY, m.selectionX, m.selectionY)
			return m, save(m)
		case key.Matches(msg, m.keymap.MuteAllNode):
			m.grid.SetAllNodeMutes(!m.mute)
			m.mute = !m.mute
			return m, save(m)
		case key.Matches(msg, m.keymap.RemoveNode):
			if m.mode != BANK {
				m.bank.ClearGrid(m.selectedGrid)
				return m.loadGridFromBank(), tea.WindowSize()
			}
			m.mode = MOVE
			m.grid.RemoveNodes(m.cursorX, m.cursorY, m.selectionX, m.selectionY)
			return m, save(m)
		case key.Matches(msg, m.keymap.EditNode):
			if m.mode == BANK {
				m.mode = MOVE
				return m.loadGridFromBank(), tea.WindowSize()
			}
			if len(m.selectedEmitters()) == 0 {
				return m, nil
			}
			if m.mode == EDIT {
				m.mode = MOVE
			} else {
				m.mode = EDIT
			}
			if m.mode == EDIT {
				m.params = param.NewParamsForNodes(m.grid, m.selectedEmitters())
			}
			if len(m.params) < m.paramPage+1 {
				m.paramPage = 0
			}
			if len(m.activeParamPage()) < m.param+1 {
				m.param = 0
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
		case key.Matches(msg, m.keymap.Bank):
			m.selectedGrid = m.bank.Active
			if m.mode == BANK {
				m.mode = MOVE
			} else {
				m.mode = BANK
			}
			return m, nil
		case key.Matches(msg, m.keymap.RootNoteUp):
			if m.mode == EDIT {
				return m, nil
			}
			param.Get("root", m.gridParams).Up()
			return m, save(m)
		case key.Matches(msg, m.keymap.RootNoteDown):
			if m.mode == EDIT {
				return m, nil
			}
			param.Get("root", m.gridParams).Down()
			return m, save(m)
		case key.Matches(msg, m.keymap.ScaleUp):
			if m.mode == EDIT {
				return m, nil
			}
			param.Get("scale", m.gridParams).Up()
			return m, save(m)
		case key.Matches(msg, m.keymap.ScaleDown):
			if m.mode == EDIT {
				return m, nil
			}
			param.Get("scale", m.gridParams).Down()
			return m, save(m)
		case key.Matches(msg, m.keymap.TempoUp):
			m.grid.SetTempo(m.grid.Tempo() + 1)
			return m, save(m)
		case key.Matches(msg, m.keymap.TempoDown):
			m.grid.SetTempo(m.grid.Tempo() - 1)
			return m, save(m)
		case key.Matches(msg, m.keymap.SelectMidiDevice):
			m.grid.CycleMidiDevice()
			return m, nil
		case key.Matches(msg, m.keymap.Copy):
			if m.mode == BANK {
				m.bankClipboard = m.bank.Grids[m.selectedGrid]
				return m, nil
			}
			m.grid.CopyOrCut(m.cursorX, m.cursorY, m.selectionX, m.selectionY, false)
			return m, nil
		case key.Matches(msg, m.keymap.Cut):
			if m.mode == BANK {
				m.bankClipboard = m.bank.Grids[m.selectedGrid]
				m.bank.ClearGrid(m.selectedGrid)
				if m.bank.Active == m.selectedGrid {
					return m.loadGridFromBank(), tea.WindowSize()
				}
				return m, tea.WindowSize()
			}
			m.grid.CopyOrCut(m.cursorX, m.cursorY, m.selectionX, m.selectionY, true)
			return m, nil
		case key.Matches(msg, m.keymap.Paste):
			if m.mode == BANK {
				m.bank.Grids[m.selectedGrid] = m.bankClipboard
			}
			m.grid.Paste(m.cursorX, m.cursorY, m.selectionX, m.selectionY)
			m.params = param.NewParamsForNodes(m.grid, m.selectedEmitters())
			return m, save(m)
		case key.Matches(msg, m.keymap.Cancel):
			m.mode = MOVE
			m.selectionX = m.cursorX
			m.selectionY = m.cursorY
			m.help.ShowAll = false
			return m, nil
		case key.Matches(msg, m.keymap.FitGridToWindow):
			m.cursorX, m.cursorY = 1, 1
			m.selectionX, m.selectionY = m.cursorX, m.cursorY
			m.grid.Resize(m.viewport.Width, m.viewport.Height)
			m.viewport.Update(m.cursorX, m.cursorY, m.grid.Width, m.grid.Height)
			return m, save(m)
		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, tea.ClearScreen
		case key.Matches(msg, m.keymap.Quit):
			m.grid.Reset()
			return m, tea.Batch(save(m), tea.Quit)
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
		return lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().
				MarginTop(1).
				MarginLeft(2).
				Render(fmt.Sprintf(helpHeader, m.version)),
			lipgloss.NewStyle().
				Height(m.viewport.Height+controlsHeight-1).
				Render(help),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderGrid(),
		help,
	)
}

func (m mainModel) handleParamEdit(dir string) {
	if len(m.activeParamPage()) < m.param+1 {
		return
	}

	switch dir {
	case "up":
		m.activeParam().Up()
	case "down":
		m.activeParam().Down()
	case "left":
		m.activeParam().Left()
		return // no preview for alt param
	case "right":
		m.activeParam().Right()
		return // no preview for alt param
	}

	switch p := m.activeParam().(type) {
	case *param.Key:
		if m.grid.Playing {
			return
		}
		p.Preview()
	}
}

func (m mainModel) handleParamAltEdit(dir string) {
	if len(m.activeParamPage()) < m.param+1 {
		return
	}

	switch dir {
	case "up":
		m.activeParam().AltUp()
	case "down":
		m.activeParam().AltDown()
	case "left":
		m.activeParam().AltLeft()
	case "right":
		m.activeParam().AltRight()
	}
}

func (m mainModel) activeParam() param.Param {
	return m.params[m.paramPage][m.param]
}

func (m mainModel) activeParamPage() []param.Param {
	return m.params[m.paramPage]
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
	if len(m.activeParamPage()) == 0 {
		return
	}
	switch dir {
	case "up":
		if m.paramPage+1 >= len(m.params) {
			return
		}
		m.param = 0
		m.paramPage++
	case "down":
		if m.paramPage-1 < 0 {
			return
		}
		m.param = 0
		m.paramPage--
	case "right":
		if m.param+1 >= len(m.activeParamPage()) {
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

func (m *mainModel) moveBankGrid(dir string) {
	switch dir {
	case "up":
		if m.selectedGrid-gridsPerLine < 0 {
			return
		}
		m.selectedGrid = m.selectedGrid - gridsPerLine
	case "down":
		if m.selectedGrid+gridsPerLine >= maxGrids {
			return
		}
		m.selectedGrid = m.selectedGrid + gridsPerLine
	case "left":
		if m.selectedGrid == 0 {
			return
		}
		m.selectedGrid--
	case "right":
		if m.selectedGrid == maxGrids-1 {
			return
		}
		m.selectedGrid++
	}
}

func (m mainModel) loadGridFromBank() mainModel {
	m.bank.Active = m.selectedGrid
	isPlaying := m.grid.Playing
	m.grid.Load(m.bank.ActiveGrid())
	m.grid.Playing = isPlaying
	m.cursorX = 1
	m.cursorY = 1
	m.selectionX = 1
	m.selectionY = 1
	return m.windowResize(m.viewport.Width, m.viewport.Height)
}

func (m mainModel) windowResize(width, height int) mainModel {
	m.help.Width = width
	m.viewport.Width = width / 2
	m.viewport.Height = height - controlsHeight - 1
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
	return m
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
	return util.Clamp(newX, minX, maxX), util.Clamp(newY, minY, maxY)
}
