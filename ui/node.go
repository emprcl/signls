package ui

import (
	"cykl/core"

	"github.com/charmbracelet/lipgloss"
)

var (
	gridStyle    = lipgloss.NewStyle().Foreground(gridColor)
	cursorStyle  = lipgloss.NewStyle().Foreground(signalColor)
	signalStyle  = lipgloss.NewStyle().Foreground(signalColor)
	emitterStyle = lipgloss.NewStyle().Foreground(primaryColor)
)

func (m mainModel) renderNode(node core.Node, i, j int) string {
	if j == m.cursorX && i == m.cursorY {
		return signalStyle.Render("██")
	}
	switch node.(type) {
	case *core.BasicEmitter:
		var emitter string
		switch node.(*core.BasicEmitter).Direction {
		case 0:
			emitter = "▀▀"
		case 1:
			emitter = " █"
		case 2:
			emitter = "▄▄"
		case 3:
			emitter = "█ "
		}
		if node.(*core.BasicEmitter).Activated {
			return emitterStyle.Background(signalColor).Render(emitter)
		} else {
			return emitterStyle.Background(secondaryColor).Render(emitter)
		}
	case *core.Signal:
		switch node.(*core.Signal).Direction {
		case 0, 2:
			return signalStyle.Render("██")
		case 1, 3:
			return signalStyle.Render("██")
		}

	}
	if (i+j)%2 == 0 {
		return "  "
	}
	return gridStyle.Render("░░")
}
