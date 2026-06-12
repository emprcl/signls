package field

import (
	"signls/core/common"
	"signls/core/music"
	"signls/core/node"
)

// CellKind classifies a grid cell so the ui can pick the right rendering
// without inspecting (and racing on) live node state.
type CellKind uint8

const (
	// CellEmpty is an empty grid position.
	CellEmpty CellKind = iota
	// CellSignal is a moving signal (common.Movable).
	CellSignal
	// CellAudible is an emitter that triggers notes (music.Audible).
	CellAudible
	// CellHole is a teleport hole (*node.HoleEmitter).
	CellHole
)

// Cell is an immutable, display-oriented copy of a single grid position. It
// holds everything the ui needs to render a cell, so the render path never
// touches live node state.
type Cell struct {
	Kind      CellKind
	Symbol    string
	Color     string
	Activated bool
	Muted     bool
}

// Snapshot returns an immutable copy of the grid's display state, taken under
// the read lock. The ui renders from the returned cells without holding the
// lock, so the (slow) render never blocks the clock goroutine — only the fast
// copy does. The returned slice is owned by the caller.
func (g *Grid) Snapshot() [][]Cell {
	g.mu.RLock()
	defer g.mu.RUnlock()

	cells := make([][]Cell, g.Height)
	for y := 0; y < g.Height; y++ {
		cells[y] = make([]Cell, g.Width)
		for x := 0; x < g.Width; x++ {
			n := g.nodes[y][x]
			if n == nil {
				continue
			}
			c := Cell{
				Symbol:    n.Symbol(),
				Color:     n.Color(),
				Activated: n.Activated(),
			}
			// Mirror the type precedence used by the renderer.
			switch t := n.(type) {
			case common.Movable:
				c.Kind = CellSignal
			case music.Audible:
				c.Kind = CellAudible
				c.Muted = t.Muted()
			case *node.HoleEmitter:
				c.Kind = CellHole
			}
			cells[y][x] = c
		}
	}
	return cells
}
