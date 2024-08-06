package core

import (
	"cykl/midi"
)

type Grid struct {
	midi  midi.Midi
	clock *clock
	nodes [][]Node
	h     int
	w     int

	Playing bool
}

func NewGrid(width, height int, midi midi.Midi) *Grid {
	grid := &Grid{
		midi:  midi,
		nodes: make([][]Node, height),
		h:     height,
		w:     width,
	}
	for i := range grid.nodes {
		grid.nodes[i] = make([]Node, width)
	}

	// Basic test case
	grid.AddBasicEmitter(3, 3, 1, true)
	grid.AddBasicEmitter(15, 3, 2, false)
	grid.AddBasicEmitter(15, 15, 3, false)
	grid.AddBasicEmitter(3, 15, 0, false)

	grid.AddBasicEmitter(4, 4, 1, true)
	grid.AddBasicEmitter(14, 4, 2, false)
	grid.AddBasicEmitter(14, 14, 3, false)
	grid.AddBasicEmitter(4, 14, 0, false)

	grid.AddBasicEmitter(5, 5, 1, true)
	grid.AddBasicEmitter(13, 5, 2, false)
	grid.AddBasicEmitter(13, 13, 3, false)
	grid.AddBasicEmitter(5, 13, 0, false)

	grid.AddBasicEmitter(6, 6, 1, true)
	grid.AddBasicEmitter(12, 6, 2, false)
	grid.AddBasicEmitter(12, 12, 3, false)
	grid.AddBasicEmitter(6, 12, 0, false)

	grid.AddBasicEmitter(7, 7, 1, true)
	grid.AddBasicEmitter(11, 7, 2, false)
	grid.AddBasicEmitter(11, 11, 3, false)
	grid.AddBasicEmitter(7, 11, 0, false)

	grid.AddBasicEmitter(8, 8, 1, true)
	grid.AddBasicEmitter(10, 8, 2, false)
	grid.AddBasicEmitter(10, 10, 3, false)
	grid.AddBasicEmitter(8, 10, 0, false)

	grid.clock = newClock(60., func() {
		grid.Update()
	})

	return grid
}

func (g *Grid) Nodes() [][]Node {
	return g.nodes
}

func (g *Grid) AddBasicEmitter(x, y int, direction uint8, emitOnPlay bool) {
	g.nodes[y][x] = &BasicEmitter{
		Direction: direction,
		Activated: emitOnPlay,
	}
}

func (g *Grid) AddSignal(x, y int, direction uint8) {
	g.nodes[y][x] = &Signal{
		Direction: direction,
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
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
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
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
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
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			if g.nodes[y][x] == nil {
				continue
			}
			g.nodes[y][x].Reset()
		}
	}
}

func (g *Grid) Emit(x, y int, direction uint8) {
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

func (g *Grid) Move(x, y int, direction uint8) {
	newX, newY := x, y
	switch direction {
	case 0:
		newY -= 1
	case 1:
		newX += 1
	case 2:
		newY += 1
	case 3:
		newX -= 1
	}

	if newX >= g.w || newY >= g.h ||
		newX < 0 || newY < 0 {
		g.nodes[y][x] = nil
		return
	}

	if g.nodes[newY][newX] == nil {
		g.nodes[newY][newX] = g.nodes[y][x]
	} else if t, ok := g.nodes[newY][newX].(Trigger); ok {
		t.Trig()
	} else if t, ok := g.nodes[newY][newX].(*BasicEmitter); ok {
		t.Emit()
	}

	g.nodes[y][x] = nil
}
