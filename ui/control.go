package ui

import (
	"cykl/core"
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	controlStyle = lipgloss.NewStyle().
			MarginTop(1).
			MarginLeft(2)

	cellStyle = lipgloss.NewStyle().
			Width(8).
			MarginRight(0)
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
				cellStyle.Width(12).Render(m.selectedNodeName()),
				cellStyle.Render(m.modeName()),
			),
			pane,
		),
	)
}

func (m mainModel) gridInfo() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.JoinVertical(
			lipgloss.Left,
			cellStyle.Render(fmt.Sprintf("%d,%d", m.cursorX, m.cursorY)),
			cellStyle.Render(fmt.Sprintf("%d,%d", m.grid.Width, m.grid.Height)),
		),
		lipgloss.JoinVertical(
			lipgloss.Left,
			cellStyle.Render(fmt.Sprintf("%.f %s", m.grid.Tempo(), m.transportSymbol())),
			"8:8",
		),
		lipgloss.JoinVertical(
			lipgloss.Left,
			cellStyle.Render(fmt.Sprintf("%d", m.grid.Pulse)),
			m.grid.MidiDevice(),
		),
	)
}

func (m mainModel) nodeEdit() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.JoinVertical(
			lipgloss.Left,
			cellStyle.Render("c4"),
			cellStyle.Render("note"),
		),
		lipgloss.JoinVertical(
			lipgloss.Left,
			cellStyle.Render("100"),
			cellStyle.Render("vel"),
		),
		lipgloss.JoinVertical(
			lipgloss.Left,
			cellStyle.Render("100"),
			cellStyle.Render("len"),
		),
		lipgloss.JoinVertical(
			lipgloss.Left,
			cellStyle.Render("1"),
			cellStyle.Render("chan"),
		),
	)
}

func (m mainModel) transportSymbol() string {
	if m.grid.QuarterNote() {
		return "●"
	}
	return " "
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

func (m mainModel) selectedNodeName() string {
	if m.selectedNode() == nil {
		return "empty"
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		emitterStyle.
			MarginRight(1).
			Background(lipgloss.Color(m.selectedNode().Color())).
			Render(m.selectedNode().Symbol()),
		m.selectedNode().Name(),
	)
}
