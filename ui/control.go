package ui

import (
	"fmt"

	"cykl/core"
	"cykl/ui/param"

	"github.com/charmbracelet/lipgloss"
)

var (
	controlStyle = lipgloss.NewStyle().
			MarginTop(1).
			MarginLeft(2)

	cellStyle = lipgloss.NewStyle().
			MarginRight(2)
	activeCellStyle = cellStyle.
			Foreground(lipgloss.Color("190"))
)

func (m mainModel) renderControl() string {
	var pane string
	if m.edit {
		pane = m.nodeEdit()
	} else {
		pane = m.gridInfo()
	}
	return controlStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.JoinVertical(
				lipgloss.Left,
				cellStyle.Width(9).Render(m.selectedNodeName()),
				cellStyle.Render(m.modeName()),
			),
			pane,
		),
	)
}

func (m mainModel) gridInfo() string {
	root := param.Get("root", m.gridParams)
	scale := param.Get("scale", m.gridParams)
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.JoinVertical(
			lipgloss.Left,
			cellStyle.Render(fmt.Sprintf("%d,%d", m.cursorX, m.cursorY)),
			cellStyle.Render(fmt.Sprintf("%d,%d", m.grid.Width, m.grid.Height)),
		),
		lipgloss.JoinVertical(
			lipgloss.Left,
			cellStyle.Render(fmt.Sprintf("%.f %s", m.grid.Tempo(), m.tempoSymbol())),
			cellStyle.Render(fmt.Sprintf("%s %d", m.transportSymbol(), m.grid.Pulse())),
		),
		lipgloss.JoinVertical(
			lipgloss.Left,
			cellStyle.Render(root.Display()),
			cellStyle.Render(scale.Display()),
		),
		lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			m.grid.MidiDevice(),
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
	if m.edit {
		return "edit"
	}
	return "move"
}

func (m mainModel) selectedNode() core.Node {
	return m.grid.Nodes()[m.cursorY][m.cursorX]
}

func (m mainModel) selectedEmitters() []core.Node {
	nodes := []core.Node{}
	for y := m.cursorY; y <= m.selectionY; y++ {
		for x := m.cursorX; x <= m.selectionX; x++ {
			if n, ok := m.grid.Nodes()[y][x].(*core.Emitter); ok {
				nodes = append(nodes, n)
			}
		}
	}
	return nodes
}

func (m mainModel) selectedNodeName() string {
	nodes := m.selectedEmitters()
	if len(nodes) == 0 {
		return "empty"
	} else if len(nodes) > 1 {
		return fmt.Sprintf("%d emit.", len(nodes))
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		emitterStyle.
			MarginRight(1).
			Background(lipgloss.Color(nodes[0].Color())).
			Render(nodes[0].Symbol()),
		nodes[0].Name(),
	)
}
