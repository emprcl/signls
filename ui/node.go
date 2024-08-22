package ui

import (
	"cykl/core"
	"log"

	"github.com/charmbracelet/lipgloss"
)

var (
	gridStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("234"))
	cursorStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("190")).
			Foreground(lipgloss.Color("0"))
	selectionStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("238")).
			Foreground(lipgloss.Color("244"))
	emitterStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15"))
	mutedEmitterStyle = lipgloss.NewStyle().
				Bold(true).
				Background(lipgloss.Color("243"))
	activeEmitterStyle = lipgloss.NewStyle().
				Bold(true).
				Background(lipgloss.Color("15")).
				Foreground(lipgloss.Color("0"))
)

func (m mainModel) inSelectionRange(x, y int) bool {
	return x >= m.cursorX &&
		x <= m.selectionX &&
		y >= m.cursorY &&
		y <= m.selectionY
}

func (m mainModel) renderNode(node core.Node, x, y int) string {
	// render cursor
	isCursor := false
	if x == m.cursorX && y == m.cursorY {
		isCursor = true
	}

	// render grid
	if node == nil && isCursor {
		return cursorStyle.Render("  ")
	} else if node == nil && m.inSelectionRange(x, y) {
		return selectionStyle.Render("..")
	} else if node == nil {
		if (x+y)%2 == 0 {
			return "  "
		}
		return gridStyle.Render("░░")
	}

	// render node
	switch t := node.(type) {
	case core.Movable:
		if isCursor {
			return cursorStyle.Render("  ")
		}
		return activeEmitterStyle.Render("  ")
	case *core.Emitter:
		symbol := node.Symbol()

		if isCursor && !m.edit {
			return cursorStyle.Render(symbol)
		} else if isCursor && m.edit && m.blink {
			return cursorStyle.Render(symbol)
		} else if node.Activated() && node.(*core.Emitter).Muted() {
			return activeEmitterStyle.Render(symbol)
		} else if t.Muted() {
			return mutedEmitterStyle.Render(symbol)
		} else if node.Activated() {
			return activeEmitterStyle.
				Foreground(lipgloss.Color(node.Color())).
				Render(symbol)
		} else {
			return emitterStyle.
				Background(lipgloss.Color(node.Color())).
				Render(symbol)
		}
	default:
		log.Fatalf("cannot render node: %+v", t)
		return ""
	}
}
