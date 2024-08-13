package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	controlStyle = lipgloss.NewStyle().
			MarginTop(1).
			MarginLeft(1)

	cellStyle = lipgloss.NewStyle().MarginRight(2)
)

func (m mainModel) renderControl() string {
	return controlStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				cellStyle.Render("emit."),
				"",
			),
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				cellStyle.Render("120.0"),
				cellStyle.Render(fmt.Sprintf("%d,%d", m.cursorX, m.cursorY)),
			),
		),
	)
}
