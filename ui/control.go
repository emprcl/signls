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
			MarginRight(2)
	activeCellStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("190")).
			MarginRight(2)
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
			cellStyle.Render("1:16"), // TODO: implement
		),
		lipgloss.JoinVertical(
			lipgloss.Left,
			cellStyle.Render("C4"),
			cellStyle.Render("Dorien"),
		),
		lipgloss.JoinVertical(
			lipgloss.Left,
			cellStyle.Render(fmt.Sprintf("%d", m.grid.Pulse)),
			m.grid.MidiDevice(),
		),
	)
}

func (m mainModel) nodeEdit() string {
	params := []string{}
	style := cellStyle
	for k, p := range m.params {
		if k == m.param {
			style = activeCellStyle
		} else {
			style = cellStyle
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
	// return lipgloss.JoinHorizontal(
	// 	lipgloss.Left,
	// 	lipgloss.JoinVertical(
	// 		lipgloss.Left,
	// 		cellStyle.Render("c4"),
	// 		cellStyle.Render("note"),
	// 	),
	// 	lipgloss.JoinVertical(
	// 		lipgloss.Left,
	// 		cellStyle.Render("100"),
	// 		cellStyle.Render("vel"),
	// 	),
	// 	lipgloss.JoinVertical(
	// 		lipgloss.Left,
	// 		cellStyle.Render("100"),
	// 		cellStyle.Render("len"),
	// 	),
	// 	lipgloss.JoinVertical(
	// 		lipgloss.Left,
	// 		cellStyle.Render("1"),
	// 		cellStyle.Render("chan"),
	// 	),
	// )
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
