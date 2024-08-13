package ui

import (
	"cykl/core"
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	gridStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("234"))
	cursorStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("190")).
			Foreground(lipgloss.Color("0"))
	emitterStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15"))
	activeEmitterStyle = lipgloss.NewStyle().
				Bold(true).
				Background(lipgloss.Color("15")).
				Foreground(lipgloss.Color("0"))
)

func (m mainModel) renderNode(node core.Node, i, j int) string {
	// render cursor
	isCursor := false
	if j == m.cursorX && i == m.cursorY {
		isCursor = true
	}

	// render grid
	if node == nil && isCursor {
		return cursorStyle.Render("  ")
	} else if node == nil {
		if (i+j)%2 == 0 {
			return "  "
		}
		return gridStyle.Render("░░")
	}

	// render node
	switch node.(type) {
	case *core.Signal:
		if isCursor {
			return cursorStyle.Render("  ")
		}
		return activeEmitterStyle.Render("  ")
	default:
		symbol := fmt.Sprintf("%s%s", node.Symbol(), node.Direction().Symbol())

		if isCursor && !m.insert {
			return cursorStyle.Render(symbol)
		} else if isCursor && m.insert && m.blink {
			return cursorStyle.Render(symbol)
		} else if node.Activated() {
			return activeEmitterStyle.
				Foreground(lipgloss.Color(node.Color())).
				Render(symbol)
		} else {
			return emitterStyle.
				Background(lipgloss.Color(node.Color())).
				Render(symbol)
		}
	}
}
