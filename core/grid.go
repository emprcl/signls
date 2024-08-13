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
	Pulse   uint64
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
	/*grid.AddSimultaneousEmitter(3, 3, RIGHT, true)
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
	grid.AddSimultaneousEmitter(6, 12, UP, false)*/

	grid.AddSimultaneousEmitter(7, 7, RIGHT, true)
	grid.AddSimultaneousEmitter(11, 7, DOWN, false)
	grid.AddSimultaneousEmitter(11, 11, LEFT, false)
	grid.AddSimultaneousEmitter(7, 11, UP, false)

	grid.AddSimultaneousEmitter(7, 2, RIGHT, true)
	grid.AddSimultaneousEmitter(12, 2, LEFT, false)
	grid.AddSimultaneousEmitter(7, 3, RIGHT, true)
	grid.AddSimultaneousEmitter(9, 3, LEFT, false)

	/*grid.AddSimultaneousEmitter(8, 8, RIGHT, true)
	grid.AddSimultaneousEmitter(10, 8, DOWN, false)
	grid.AddSimultaneousEmitter(10, 10, LEFT, false)
	grid.AddSimultaneousEmitter(8, 10, UP, false)

	grid.AddSimultaneousEmitter(3, 0, RIGHT, true)
	grid.AddSimultaneousEmitter(4, 0, DOWN, false)
	grid.AddSimultaneousEmitter(4, 1, LEFT, false)
	grid.AddSimultaneousEmitter(3, 1, UP, false)*/

	grid.clock = newClock(60., func() {
		if !grid.Playing {
			return
		}
		grid.Update()
	})

	return grid
}

func (g *Grid) TogglePlay() {
	if !g.Playing {
		g.Pulse = 0
	}
	g.Playing = !g.Playing
}

func (g *Grid) Nodes() [][]Node {
	return g.nodes
}

func (g *Grid) AddSimultaneousEmitter(x, y int, direction Direction, emitOnPlay bool) {
	g.nodes[y][x] = &SimultaneousEmitter{
		direction: direction,
		armed:     emitOnPlay,
		note: Note{
			Channel:  uint8(y),
			Key:      30 + 7*uint8(x),
			Velocity: 100,
		},
	}
}

func (g *Grid) AddSignal(x, y int, direction Direction) {
	if n, ok := g.nodes[y][x].(Emitter); ok {
		n.Arm()
		n.Trig(g, x, y)
		return
	}
	g.nodes[y][x] = &Signal{
		direction: direction,
		pulse:     g.Pulse,
	}
}

func (g *Grid) Update() {
	g.Pulse++
	g.RunNodes()
}

func (g *Grid) RunNodes() {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.nodes[y][x] == nil {
				continue
			}

			if n, ok := g.nodes[y][x].(Movable); ok {
				n.Move(g, x, y)
			}

			if n, ok := g.nodes[y][x].(Emitter); ok {
				n.Trig(g, x, y)
				n.Emit(g, x, y)
			}
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

func (g *Grid) Trig(x, y int) {
	note := g.nodes[y][x].(Emitter).Note()
	g.midi.NoteOn(0, note.Channel, note.Key, note.Velocity)
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
	} else if n, ok := g.nodes[newY][newX].(Emitter); ok {
		n.Arm()
		n.Trig(g, newX, newY)
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
