package ui

import (
	"fmt"
	"time"

	"signls/core/field"
	"signls/filesystem"
	"signls/ui/param"
	"signls/ui/util"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	cells         [][]field.Cell
	saver         *saver
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
	ti.SetWidth(12)
	// In bubbles v2 the input draws its own cursor; keep it inline (virtual) so
	// it renders within the control bar layout, and color it like the v1 cursor.
	ti.SetVirtualCursor(true)
	styles := ti.Styles()
	styles.Cursor.Color = lipgloss.Color("190")
	ti.SetStyles(styles)
	model := mainModel{
		bank:       bank,
		grid:       grid,
		keymap:     newKeyMap(config.KeyMap),
		help:       help.New(),
		input:      ti,
		gridParams: param.NewParamsForGrid(grid),
		saver:      newSaver(saveDebounce, func() { grid.Save(bank) }),
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

// requestWindowSize asks bubbletea to re-send the current window size. In v2
// tea.RequestWindowSize is a Msg, so it must be wrapped in a Cmd.
func requestWindowSize() tea.Cmd {
	return func() tea.Msg {
		return tea.RequestWindowSize()
	}
}

func save(m mainModel) tea.Cmd {
	return func() tea.Msg {
		m.grid.Save(m.bank)
		return saveMsg(true)
	}
}

// requestSave schedules a debounced save instead of writing on every edit, so a
// burst of edits coalesces into a single disk write. The write happens on the
// saver's timer goroutine.
func (m mainModel) requestSave() {
	m.saver.request()
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(tick(), blink())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		return m.windowResize(msg.Width, msg.Height), nil

	case tickMsg:
		return m.handleBankMetaCommand()

	case blinkMsg:
		m.blink = !m.blink
		return m, blink()

	case tea.KeyPressMsg:
		if m.input.Focused() {
			var cmd tea.Cmd
			switch {
			case key.Matches(msg, m.keymap.EditNode):
				m.input.Blur()
				m.grid.Write(func() {
					m.activeParam().SetEditValue(m.input.Value())
				})
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
			if m.mode == EDIT || m.mode == CONFIG {
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
			if m.mode == EDIT || m.mode == CONFIG {
				m.handleParamAltEdit(dir)
				m.requestSave()
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
			if m.mode == MOVE {
				m.grid.Write(func() {
					param.NewDirection(m.selectedEmitters()).SetFromKeyString(dir)
				})
				m.requestSave()
				return m, nil
			}
			m.handleParamEdit(dir)
			m.requestSave()
			return m, nil
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
			m.requestSave()
			return m, nil
		case key.Matches(msg, m.keymap.MuteNode):
			m.grid.ToggleNodeMutes(m.cursorX, m.cursorY, m.selectionX, m.selectionY)
			m.requestSave()
			return m, nil
		case key.Matches(msg, m.keymap.MuteAllNode):
			m.grid.SetAllNodeMutes(!m.mute)
			m.mute = !m.mute
			m.requestSave()
			return m, nil
		case key.Matches(msg, m.keymap.RemoveNode):
			if m.mode == BANK {
				m.bank.ClearGrid(m.selectedGrid)
				return m.loadGridFromBank(), requestWindowSize()
			}
			m.mode = MOVE
			m.grid.RemoveNodes(m.cursorX, m.cursorY, m.selectionX, m.selectionY)
			m.requestSave()
			return m, nil
		case key.Matches(msg, m.keymap.EditNode):
			if m.mode == BANK {
				m.mode = MOVE
				return m.loadGridFromBank(), requestWindowSize()
			}
			if m.mode == CONFIG {
				m.mode = MOVE
				return m, nil
			}
			if len(m.selectedEmitters()) == 0 {
				return m, nil
			}
			m.mode = m.toggleMode(EDIT)
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
			m.grid.TriggerNode(m.cursorX, m.cursorY)
			return m, nil
		case key.Matches(msg, m.keymap.Bank):
			m.selectedGrid = m.bank.ActiveIndex()
			m.mode = m.toggleMode(BANK)
			return m, nil
		case key.Matches(msg, m.keymap.RootNoteUp):
			if m.mode == EDIT {
				return m, nil
			}
			param.Get("root", m.gridParams).Up()
			m.requestSave()
			return m, nil
		case key.Matches(msg, m.keymap.RootNoteDown):
			if m.mode == EDIT {
				return m, nil
			}
			param.Get("root", m.gridParams).Down()
			m.requestSave()
			return m, nil
		case key.Matches(msg, m.keymap.ScaleUp):
			if m.mode == EDIT {
				return m, nil
			}
			param.Get("scale", m.gridParams).Up()
			m.requestSave()
			return m, nil
		case key.Matches(msg, m.keymap.ScaleDown):
			if m.mode == EDIT {
				return m, nil
			}
			param.Get("scale", m.gridParams).Down()
			m.requestSave()
			return m, nil
		case key.Matches(msg, m.keymap.TempoUp):
			m.grid.SetTempo(m.grid.Tempo() + 1)
			m.requestSave()
			return m, nil
		case key.Matches(msg, m.keymap.TempoDown):
			m.grid.SetTempo(m.grid.Tempo() - 1)
			m.requestSave()
			return m, nil
		case key.Matches(msg, m.keymap.Configuration):
			m.mode = m.toggleMode(CONFIG)
			m.params = param.NewParamsForMidi(m.grid)
			m.param = 0
			m.paramPage = 0
			return m, nil
		case key.Matches(msg, m.keymap.Copy):
			if m.mode == BANK {
				m.bankClipboard = m.bank.GridAt(m.selectedGrid)
				return m, nil
			}
			m.grid.CopyOrCut(m.cursorX, m.cursorY, m.selectionX, m.selectionY, false)
			return m, nil
		case key.Matches(msg, m.keymap.Cut):
			if m.mode == BANK {
				m.bankClipboard = m.bank.GridAt(m.selectedGrid)
				m.bank.ClearGrid(m.selectedGrid)
				if m.bank.ActiveIndex() == m.selectedGrid {
					return m.loadGridFromBank(), requestWindowSize()
				}
				return m, requestWindowSize()
			}
			m.grid.CopyOrCut(m.cursorX, m.cursorY, m.selectionX, m.selectionY, true)
			return m, nil
		case key.Matches(msg, m.keymap.Paste):
			if m.mode == BANK {
				m.bank.SetGrid(m.selectedGrid, m.bankClipboard)
			}
			m.grid.Paste(m.cursorX, m.cursorY, m.selectionX, m.selectionY)
			m.params = param.NewParamsForNodes(m.grid, m.selectedEmitters())
			m.requestSave()
			return m, nil
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
			m.requestSave()
			return m, nil
		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, tea.ClearScreen
		case key.Matches(msg, m.keymap.Quit):
			m.grid.Reset()
			// Cancel any pending debounced save; the final save below is
			// synchronous so the latest state is persisted before quitting.
			m.saver.stop()
			return m, tea.Sequence(save(m), tea.Quit)
		}
	}

	return m, nil
}

// view wraps rendered content in a tea.View. In bubbletea v2 the alt screen is
// declared on the View rather than toggled with a command in Init.
func (m mainModel) view(content string) tea.View {
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m mainModel) View() tea.View {
	// Take a consistent copy of the grid's display state, then render the grid
	// without holding any lock. The slow lipgloss render never blocks the clock
	// goroutine — only the fast Snapshot copy does.
	m.cells = m.grid.Snapshot()

	help := lipgloss.NewStyle().
		MarginLeft(2).
		Render(m.help.View(m.keymap))

	paramHelp := ""
	if m.mode == EDIT || m.mode == CONFIG {
		paramHelp = m.help.Styles.ShortDesc.
			MarginLeft(16).
			Render(m.activeParam().Help())
	}

	if m.help.ShowAll {
		return m.view(lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().
				MarginTop(1).
				MarginLeft(2).
				Render(fmt.Sprintf(helpHeader, m.version)),
			lipgloss.NewStyle().
				MarginTop(1).
				Height(m.viewport.Height+controlsHeight-1).
				Render(help),
		))
	}

	// The control bar still reads live grid/param state (tempo, pulse, selected
	// node, params). It's a single row, so render it under a short read lock.
	var control string
	m.grid.Read(func() {
		control = m.renderControl()
	})

	return m.view(lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderGrid(),
		control,
		paramHelp,
		help,
	))
}

func (m mainModel) handleParamEdit(dir string) {
	if len(m.activeParamPage()) < m.param+1 {
		return
	}

	edit := func() {
		switch dir {
		case "up":
			m.activeParam().Up()
		case "down":
			m.activeParam().Down()
		case "left":
			m.activeParam().Left()
		case "right":
			m.activeParam().Right()
		}

		// Preview only for up/down (not the alt left/right edits), and only
		// while stopped. The note copy happens here under the same lock.
		if dir == "up" || dir == "down" {
			if p, ok := m.activeParam().(*param.Key); ok && !m.grid.Playing {
				p.Preview()
			}
		}
	}

	// In EDIT mode the params mutate node state directly, so take the write
	// lock. In CONFIG mode they mutate grid state through locking Grid methods,
	// so run them directly to avoid a re-entrant lock.
	m.editParam(edit)
}

func (m mainModel) handleParamAltEdit(dir string) {
	if len(m.activeParamPage()) < m.param+1 {
		return
	}

	m.editParam(func() {
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
	})
}

// editParam runs a parameter mutation, holding the grid write lock when the
// active params mutate node state directly (EDIT mode). CONFIG-mode params
// serialize themselves through locking Grid methods, so they run unwrapped.
func (m mainModel) editParam(fn func()) {
	if m.mode == EDIT {
		m.grid.Write(fn)
		return
	}
	fn()
}

func (m mainModel) activeParam() param.Param {
	return m.params[m.paramPage][m.param]
}

func (m mainModel) activeParamPage() []param.Param {
	return m.params[m.paramPage]
}

func (m mainModel) renderGrid() string {
	lines := make([]string, 0, m.viewport.Height)
	for y := m.viewport.offsetY; y < m.viewport.offsetY+m.viewport.Height; y++ {
		nodes := make([]string, 0, m.viewport.Width)
		for x := m.viewport.offsetX; x < m.viewport.offsetX+m.viewport.Width; x++ {
			nodes = append(nodes, m.renderNode(m.cells[y][x], x, y))
		}
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left, nodes...))
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m *mainModel) moveParam(dir string) {
	if len(m.activeParamPage()) == 0 {
		return
	}
	switch dir {
	case "up":
		if m.paramPage-1 < 0 {
			return
		}
		m.param = 0
		m.paramPage--
	case "down":
		if m.paramPage+1 >= len(m.params) {
			return
		}
		m.param = 0
		m.paramPage++
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
	m.bank.SetActive(m.selectedGrid)
	isPlaying := m.grid.Playing
	m.grid.Load(m.selectedGrid, m.bank.ActiveGrid())
	m.grid.SetPlaying(isPlaying)
	m.cursorX = 1
	m.cursorY = 1
	m.selectionX = 1
	m.selectionY = 1
	m.mode = MOVE
	m.param = 0
	m.paramPage = 0
	return m.windowResize(m.viewport.Width, m.viewport.Height)
}

func (m mainModel) handleBankMetaCommand() (mainModel, tea.Cmd) {
	// BankIndex is written by the clock goroutine (meta commands), so read it
	// under the lock.
	var bankIndex int
	m.grid.Read(func() { bankIndex = m.grid.BankIndex })
	if bankIndex == m.bank.ActiveIndex() {
		return m, tick()
	}
	m.bank.SetActive(bankIndex)
	m.grid.Load(bankIndex, m.bank.ActiveGrid())
	m.grid.SetPlaying(true)
	m.mode = MOVE
	m.param = 0
	m.paramPage = 0
	return m.windowResize(m.viewport.Width, m.viewport.Height), tea.Batch(requestWindowSize(), tick())
}

func (m mainModel) windowResize(width, height int) mainModel {
	m.help.SetWidth(width)
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
	if m.selectionX > m.grid.Width-1 {
		m.selectionX = m.grid.Width - 1
	}
	if m.selectionY > m.grid.Height-1 {
		m.selectionY = m.grid.Height - 1
	}
	return m
}

func (m mainModel) toggleMode(mo mode) mode {
	if m.mode == mo {
		return MOVE
	}
	return mo
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
