package ui

import (
	"fmt"

	"signls/core/common"
	"signls/filesystem"
	"signls/ui/param"
	"signls/ui/util"

	"github.com/charmbracelet/lipgloss"
)

const (
	maxGrids     = 32
	gridsPerLine = 16
)

var (
	backgroundColor = lipgloss.Color("0")
	foregroundColor = lipgloss.Color("15")
	baseStyle       = lipgloss.NewStyle().
			Background(lipgloss.Color(backgroundColor)).
			Foreground(foregroundColor)

	controlStyle = baseStyle.Copy().
			PaddingTop(1).
			PaddingLeft(2)
	cellStyle = baseStyle.Copy().
			Background(lipgloss.Color(backgroundColor)).
			PaddingRight(2)
	activeCellStyle = cellStyle.
			Foreground(lipgloss.Color("190"))
	bankStyle = baseStyle.Copy().
			Background(lipgloss.Color("79")).
			Foreground(lipgloss.Color("0"))
	bankStyleOdd = baseStyle.Copy().
			Background(lipgloss.Color("85")).
			Foreground(lipgloss.Color("0"))
	activeBankStyle = baseStyle.Copy().
			Background(lipgloss.Color("15")).
			Foreground(lipgloss.Color("0"))
)

func (m mainModel) renderControl() string {
	//lipgloss.NewStyle().Background(lipgloss.Color("190")

	if m.bankMode {
		return controlStyle.Render(m.bankSelection())
	}

	var pane string
	if m.edit && m.input.Focused() {
		pane = fmt.Sprintf(
			"%s %s",
			m.params[m.param].Name(),
			m.input.View(),
		)
	} else if m.edit {
		pane = m.nodeEdit()
	} else {
		pane = m.gridInfo()
	}

	return lipgloss.Place(m.viewport.Width, controlsHeight, lipgloss.Left, lipgloss.Top,
		controlStyle.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				lipgloss.JoinVertical(
					lipgloss.Left,
					cellStyle.Width(11).Render(m.selectedNodeName()),
					cellStyle.Width(11).Render(m.modeName()),
				),
				pane,
			),
		),
		lipgloss.WithWhitespaceBackground(backgroundColor),
	)
}

func (m mainModel) bankSelection() string {
	banks := make([]string, maxGrids)
	for i, g := range m.bank.Grids[:maxGrids] {
		label := bankGridLabel(i, g)
		if i == m.selectedGrid {
			banks[i] = cursorStyle.MarginRight(1).Render(label)
		} else if i == m.bank.Active {
			banks[i] = activeBankStyle.Render(label)
		} else if (i < gridsPerLine && i%2 == 0) || (i >= gridsPerLine && i%2 == 1) {
			banks[i] = bankStyle.Render(label)
		} else {
			banks[i] = bankStyleOdd.Render(label)
		}

	}
	pane := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			banks[:gridsPerLine]...,
		),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			banks[gridsPerLine:maxGrids]...,
		),
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.JoinVertical(
			lipgloss.Left,
			activeBankStyle.MarginRight(9).Render(bankGridLabel(m.bank.Active, m.bank.ActiveGrid())),
			cellStyle.Render(m.modeName()),
		),
		pane,
	)
}

func (m mainModel) gridInfo() string {
	root := param.Get("root", m.gridParams)
	scale := param.Get("scale", m.gridParams)
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.JoinVertical(
			lipgloss.Left,
			cellStyle.Width(8).Render(fmt.Sprintf("%d,%d", m.cursorX, m.cursorY)),
			cellStyle.Width(8).Render(fmt.Sprintf("%d,%d", m.grid.Width-1, m.grid.Height-1)),
		),
		lipgloss.JoinVertical(
			lipgloss.Left,
			cellStyle.Width(7).Render(fmt.Sprintf("%.f %s", m.grid.Tempo(), m.tempoSymbol())),
			cellStyle.Width(7).Render(m.transportSymbol()),
		),
		lipgloss.JoinVertical(
			lipgloss.Left,
			cellStyle.Width(8).Render(root.Display()),
			cellStyle.Width(8).Render(scale.Display()),
		),
		lipgloss.JoinVertical(
			lipgloss.Left,
			fmt.Sprintf(
				"%s%s",
				activeBankStyle.Render(bankGridLabel(m.bank.Active, m.bank.ActiveGrid())),
				cellStyle.PaddingLeft(1).Render(m.bank.Filename()),
			),
			cellStyle.Render(m.grid.MidiDevice()),
		),
	)
}

func (m mainModel) nodeEdit() string {
	params := []string{}
	for k, p := range m.params {
		style := cellStyle
		if k == m.param {
			style = activeCellStyle
		}
		params = append(
			params,
			lipgloss.JoinVertical(
				lipgloss.Left,
				style.Render(p.Display()),
				style.Render(p.Name()),
			),
		)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		params...,
	)
}

func (m mainModel) tempoSymbol() string {
	if m.grid.QuarterNote() {
		return "●"
	}
	return " "
}

func (m mainModel) transportSymbol() string {
	if m.grid.Playing {
		return "▶"
	}
	return "■"
}

func (m mainModel) modeName() string {
	if m.bankMode {
		return "bank"
	}
	if m.edit {
		return "edit"
	}
	return "move"
}

func (m mainModel) selectedNode() common.Node {
	return m.grid.Nodes()[m.cursorY][m.cursorX]
}

func (m mainModel) selectedEmitters() []common.Node {
	nodes := []common.Node{}
	for y := m.cursorY; y <= m.selectionY; y++ {
		for x := m.cursorX; x <= m.selectionX; x++ {
			if m.grid.Nodes()[y][x] == nil {
				continue
			} else if _, ok := m.grid.Nodes()[y][x].(common.Movable); ok {
				continue
			}
			nodes = append(nodes, m.grid.Nodes()[y][x])
		}
	}
	return nodes
}

func (m mainModel) selectedNodeName() string {
	nodes := m.selectedEmitters()
	if len(nodes) == 0 {
		return "empty"
	} else if len(nodes) > 1 {
		return fmt.Sprintf("%d nodes", len(nodes))
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		emitterStyle.
			Background(lipgloss.Color(nodes[0].Color())).
			Render(util.Normalize(nodes[0].Symbol())),
		baseStyle.PaddingLeft(1).Render(nodes[0].Name()),
	)
}

func bankGridLabel(nb int, g filesystem.Grid) string {
	label := fmt.Sprintf("%2d", nb+1)
	if !g.IsEmpty() {
		label = util.Normalize(fmt.Sprintf("%2d\u0320", nb+1))
	}
	return label
}
