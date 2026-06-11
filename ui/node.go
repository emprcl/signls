package ui

import (
	"log"
	"signls/core/field"
	"signls/core/node"
	"signls/ui/param"
	"signls/ui/util"

	"github.com/charmbracelet/lipgloss"
)

var (
	gridStyle = lipgloss.NewStyle().
			Background(lipgloss.AdaptiveColor{Light: "254", Dark: "234"})
	cursorStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("190")).
			Foreground(lipgloss.Color("0"))
	teleportDestinationStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("160")).
					Foreground(lipgloss.Color("15"))
	selectionStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("238")).
			Foreground(lipgloss.Color("244"))
	emitterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15"))
	mutedEmitterStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("247")).
				Foreground(lipgloss.Color("236"))
	activeEmitterStyle = lipgloss.NewStyle().
				Background(lipgloss.AdaptiveColor{Light: "0", Dark: "15"}).
				Foreground(lipgloss.AdaptiveColor{Light: "15", Dark: "0"})
)

func (m mainModel) inSelectionRange(x, y int) bool {
	return x >= m.cursorX &&
		x <= m.selectionX &&
		y >= m.cursorY &&
		y <= m.selectionY
}

func (m mainModel) renderNode(c field.Cell, x, y int) string {
	// render cursor
	isCursor := x == m.cursorX && y == m.cursorY && m.mode != BANK

	isTeleportDestination := false
	if m.mode == EDIT && len(m.params) > 0 {
		p, ok := m.activeParam().(param.Destination)
		if ok {
			destinationX, destinationY := p.Position()
			isTeleportDestination = (destinationX == x && destinationY == y)
		}
	}

	// render grid
	teleportDestinationSymbol := node.HoleDestinationSymbol
	if c.Kind == field.CellEmpty && isCursor {
		return cursorStyle.Render("  ")
	} else if c.Kind == field.CellEmpty && isTeleportDestination && !m.blink && m.mode != BANK {
		return cursorStyle.Render(teleportDestinationSymbol)
	} else if c.Kind == field.CellEmpty && isTeleportDestination && (m.blink || m.mode == BANK) {
		return teleportDestinationStyle.Render(teleportDestinationSymbol)
	} else if c.Kind == field.CellEmpty && m.inSelectionRange(x, y) && m.mode != BANK {
		return selectionStyle.Render("..")
	} else if c.Kind == field.CellEmpty {
		if (x+y)%2 == 0 {
			return "  "
		}
		return gridStyle.Render("  ")
	}

	// render node
	switch c.Kind {
	case field.CellSignal:
		if isCursor {
			return cursorStyle.Render("  ")
		}
		return activeEmitterStyle.Render("  ")
	case field.CellAudible:
		symbol := util.Normalize(c.Symbol)

		if isCursor && m.mode != EDIT {
			return cursorStyle.Render(symbol)
		} else if isTeleportDestination && m.mode == EDIT && m.blink {
			return teleportDestinationStyle.Render(teleportDestinationSymbol)
		} else if isCursor && m.mode == EDIT && m.blink {
			return cursorStyle.Render(symbol)
		} else if c.Activated && c.Muted {
			return activeEmitterStyle.Render(symbol)
		} else if c.Muted {
			return mutedEmitterStyle.Render(symbol)
		} else if c.Activated {
			return activeEmitterStyle.
				Foreground(lipgloss.Color(c.Color)).
				Render(symbol)
		} else {
			return emitterStyle.
				Background(lipgloss.Color(c.Color)).
				Render(symbol)
		}
	case field.CellHole:
		symbol := c.Symbol

		if isCursor && m.mode != EDIT {
			return cursorStyle.Render(symbol)
		} else if isCursor && m.mode == EDIT && m.blink {
			return cursorStyle.Render(symbol)
		} else if c.Activated {
			return activeEmitterStyle.
				Foreground(lipgloss.Color(c.Color)).
				Render(symbol)
		} else {
			return emitterStyle.
				Background(lipgloss.Color(c.Color)).
				Render(symbol)
		}
	default:
		log.Fatalf("cannot render node kind: %d", c.Kind)
		return ""
	}
}
