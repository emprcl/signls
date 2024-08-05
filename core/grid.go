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

	grid.AddOnceEmitter(3, 3, 1)

	grid.clock = newClock(120., func() {
		grid.Update()
	})

	return grid
}

func (g *Grid) Nodes() [][]Node {
	return g.nodes
}

func (g *Grid) AddOnceEmitter(x, y int, direction uint8) {
	g.nodes[y][x] = &OnceEmitter{
		direction:  direction,
		shouldEmit: true,
	}
}

func (g *Grid) AddSignal(x, y int, direction uint8) {
	g.nodes[y][x] = &Signal{
		direction: direction,
	}
}

func (g *Grid) Update() {
	g.RunSignalsAndEmitters()
	g.RunTriggers()
}

func (g *Grid) RunSignalsAndEmitters() {
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			if g.nodes[y][x] == nil {
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
			}
			g.nodes[y][x].Update(g, x, y)
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

	if newX >= g.w || newY >= g.h {
		g.nodes[y][x] = nil
		return
	}

	if t, ok := g.nodes[newY][newX].(Trigger); ok {
		t.Trig()
	} else if g.nodes[newY][newX] == nil {
		g.nodes[newY][newX] = g.nodes[y][x]
	}

	g.nodes[y][x] = nil
}
