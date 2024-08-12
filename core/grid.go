package core

import (
	"cykl/midi"
)

type Grid struct {
	midi   midi.Midi
	clock  *clock
	nodes  [][]Node
	Height int
	Width  int

	Playing bool
}

func NewGrid(width, height int, midi midi.Midi) *Grid {
	grid := &Grid{
		midi:   midi,
		nodes:  make([][]Node, height),
		Height: height,
		Width:  width,
	}
	for i := range grid.nodes {
		grid.nodes[i] = make([]Node, width)
	}

	// Basic test case
	grid.AddSimultaneousEmitter(3, 3, RIGHT, true)
	grid.AddSimultaneousEmitter(15, 3, DOWN, false)
	grid.AddSimultaneousEmitter(15, 15, LEFT, false)
	grid.AddSimultaneousEmitter(3, 15, UP, false)

	grid.AddSimultaneousEmitter(4, 4, RIGHT, true)
	grid.AddSimultaneousEmitter(14, 4, DOWN, false)
	grid.AddSimultaneousEmitter(14, 14, LEFT, false)
	grid.AddSimultaneousEmitter(4, 14, UP, false)

	grid.AddSimultaneousEmitter(5, 5, RIGHT, true)
	grid.AddSimultaneousEmitter(13, 5, DOWN, false)
	grid.AddSimultaneousEmitter(13, 13, LEFT, false)
	grid.AddSimultaneousEmitter(5, 13, UP, false)

	grid.AddSimultaneousEmitter(6, 6, RIGHT, true)
	grid.AddSimultaneousEmitter(12, 6, DOWN, false)
	grid.AddSimultaneousEmitter(12, 12, LEFT, false)
	grid.AddSimultaneousEmitter(6, 12, UP, false)

	grid.AddSimultaneousEmitter(7, 7, RIGHT, true)
	grid.AddSimultaneousEmitter(11, 7, DOWN, false)
	grid.AddSimultaneousEmitter(11, 11, LEFT, false)
	grid.AddSimultaneousEmitter(7, 11, UP, false)

	grid.AddSimultaneousEmitter(8, 8, RIGHT, true)
	grid.AddSimultaneousEmitter(10, 8, DOWN, false)
	grid.AddSimultaneousEmitter(10, 10, LEFT, false)
	grid.AddSimultaneousEmitter(8, 10, UP, false)

	grid.AddSimultaneousEmitter(20, 3, RIGHT, true)
	grid.AddSimultaneousEmitter(21, 3, DOWN, false)
	grid.AddSimultaneousEmitter(21, 4, LEFT, false)
	grid.AddSimultaneousEmitter(20, 4, UP, false)

	grid.clock = newClock(60., func() {
		grid.Update()
	})

	return grid
}

func (g *Grid) Nodes() [][]Node {
	return g.nodes
}

func (g *Grid) AddSimultaneousEmitter(x, y int, direction Direction, emitOnPlay bool) {
	g.nodes[y][x] = &SimultaneousEmitter{
		direction: direction,
		activated: emitOnPlay,
	}
}

func (g *Grid) AddSignal(x, y int, direction Direction) {
	if g.nodes[y][x] != nil {
		if t, ok := g.nodes[y][x].(Emitter); ok {
			t.Emit()
		}
		// TODO: mutualise code with move
		return
	}
	g.nodes[y][x] = &Signal{
		direction: direction,
		updated:   true,
	}
}

func (g *Grid) Update() {
	if !g.Playing {
		return
	}
	g.RunSignalsAndEmitters()
	g.RunTriggers()
	g.RunResets()
}

func (g *Grid) RunSignalsAndEmitters() {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.nodes[y][x] == nil {
				continue
			} else if _, ok := g.nodes[y][x].(Trigger); ok {
				continue
			}
			g.nodes[y][x].Update(g, x, y)
		}
	}
}

func (g *Grid) RunTriggers() {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.nodes[y][x] == nil {
				continue
			} else if _, ok := g.nodes[y][x].(Trigger); !ok {
				continue
			}
			g.nodes[y][x].Update(g, x, y)
		}
	}
}

func (g *Grid) RunResets() {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.nodes[y][x] == nil {
				continue
			}
			g.nodes[y][x].Reset()
		}
	}
}

func (g *Grid) Emit(x, y int, direction Direction) {
	switch direction {
	case 0:
		g.AddSignal(x, y-1, direction)
	case 1:
		g.AddSignal(x+1, y, direction)
	case 2:
		g.AddSignal(x, y+1, direction)
	case 3:
		g.AddSignal(x-1, y, direction)
	}
}

func (g *Grid) Move(x, y int, direction Direction) {
	newX, newY := x, y
	switch direction {
	case UP:
		newY -= 1
	case RIGHT:
		newX += 1
	case DOWN:
		newY += 1
	case LEFT:
		newX -= 1
	}

	if newX >= g.Width || newY >= g.Height ||
		newX < 0 || newY < 0 {
		g.nodes[y][x] = nil
		return
	}

	if g.nodes[newY][newX] == nil {
		g.nodes[newY][newX] = g.nodes[y][x]
	} else if t, ok := g.nodes[newY][newX].(Trigger); ok {
		t.Trig()
	} else if t, ok := g.nodes[newY][newX].(Emitter); ok {
		t.Emit()
	}

	g.nodes[y][x] = nil
}

func (g *Grid) Resize(newWidth, newHeight int) {
	newNodes := make([][]Node, newHeight)
	for i := range newNodes {
		newNodes[i] = make([]Node, newWidth)
	}

	minWidth := g.Width
	if newWidth < g.Width {
		minWidth = newWidth
	}

	minHeight := g.Height
	if newHeight < g.Height {
		minHeight = newHeight
	}

	for y := 0; y < minHeight; y++ {
		for x := 0; x < minWidth; x++ {
			newNodes[y][x] = g.nodes[y][x]
		}
	}

	g.Width = newWidth
	g.Height = newHeight
	g.nodes = newNodes
}
