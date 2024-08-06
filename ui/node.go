package ui

import (
	"cykl/core"

	"github.com/charmbracelet/lipgloss"
)

var (
	gridStyle    = lipgloss.NewStyle().Foreground(gridColor)
	signalStyle  = lipgloss.NewStyle().Foreground(signalColor)
	emitterStyle = lipgloss.NewStyle().Background(secondaryColor).Foreground(primaryColor)
)

func renderNode(node core.Node) string {
	switch node.(type) {
	case *core.BasicEmitter:
		switch node.(*core.BasicEmitter).Direction {
		case 0:
			return emitterStyle.Render("▀▀")
		case 1:
			return emitterStyle.Render(" █")
		case 2:
			return emitterStyle.Render("▄▄")
		case 3:
			return emitterStyle.Render("█ ")
		}
	case *core.Signal:
		switch node.(*core.Signal).Direction {
		case 0, 2:
			return signalStyle.Render("██")
		case 1, 3:
			return signalStyle.Render("██")
		}

	}
	return gridStyle.Render("├┤")
}
