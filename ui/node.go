package ui

import (
	"cykl/core"

	"github.com/charmbracelet/lipgloss"
)

var (
	gridStyle    = lipgloss.NewStyle().Foreground(secondaryColor)
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
		return signalStyle.Render("≡≡")
	}
	return gridStyle.Render("··")
}
