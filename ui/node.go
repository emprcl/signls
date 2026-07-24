package ui

import (
	"log"

	"signls/core/field"
	"signls/core/node"
	"signls/ui/param"
	"signls/ui/util"
)

// gridBlankCell and gridBlankCellAlt cache the two rendered empty grid cells
// (the alternating checkerboard positions). On a sparse grid these are by far
// the most common cells, and their output is constant. They are computed lazily
// on first render (rather than at init) so lipgloss has already detected the
// terminal background for its adaptive colors. renderNode runs only on the
// bubbletea (View) goroutine, so the lazy writes need no synchronization.
var (
	gridBlankCell    string
	gridBlankCellAlt string
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
		return m.styles.cursor.Render("  ")
	} else if c.Kind == field.CellEmpty && isTeleportDestination && !m.blink && m.mode != BANK {
		return m.styles.cursor.Render(teleportDestinationSymbol)
	} else if c.Kind == field.CellEmpty && isTeleportDestination && (m.blink || m.mode == BANK) {
		return m.styles.holeDestination.Render(teleportDestinationSymbol)
	} else if c.Kind == field.CellEmpty && m.inSelectionRange(x, y) && m.mode != BANK {
		return m.styles.selection.Render("..")
	} else if c.Kind == field.CellEmpty {
		if (x+y)%2 == 0 {
			// The alternating cells: their own background when the theme sets
			// one, otherwise the bare terminal background.
			if !m.styles.gridAltSet {
				return "  "
			}
			if gridBlankCellAlt == "" {
				gridBlankCellAlt = m.styles.gridAlt.Render("  ")
			}
			return gridBlankCellAlt
		}
		if gridBlankCell == "" {
			gridBlankCell = m.styles.grid.Render("  ")
		}
		return gridBlankCell
	}

	// render node
	switch c.Kind {
	case field.CellSignal:
		if isCursor {
			return m.styles.cursor.Render("  ")
		}
		return m.styles.activeEmitter.Render("  ")
	case field.CellAudible:
		symbol := util.Normalize(c.Symbol)

		if isCursor && m.mode != EDIT {
			return m.styles.cursor.Render(symbol)
		} else if isTeleportDestination && m.mode == EDIT && m.blink {
			return m.styles.holeDestination.Render(teleportDestinationSymbol)
		} else if isCursor && m.mode == EDIT && m.blink {
			return m.styles.cursor.Render(symbol)
		} else if c.Activated && c.Muted {
			return m.styles.activeEmitter.Render(symbol)
		} else if c.Muted {
			return m.styles.mutedEmitter.Render(symbol)
		} else if c.Activated {
			return m.styles.activeEmitter.
				Foreground(m.styles.nodeColor(c.ColorKey)).
				Render(symbol)
		} else {
			return m.styles.emitter.
				Background(m.styles.nodeColor(c.ColorKey)).
				Render(symbol)
		}
	case field.CellHole:
		symbol := c.Symbol

		if isCursor && m.mode != EDIT {
			return m.styles.cursor.Render(symbol)
		} else if isCursor && m.mode == EDIT && m.blink {
			return m.styles.cursor.Render(symbol)
		} else if c.Activated {
			return m.styles.activeEmitter.
				Foreground(m.styles.nodeColor(c.ColorKey)).
				Render(symbol)
		} else {
			return m.styles.emitter.
				Background(m.styles.nodeColor(c.ColorKey)).
				Render(symbol)
		}
	default:
		log.Fatalf("cannot render node kind: %d", c.Kind)
		return ""
	}
}
