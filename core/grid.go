package core

import (
	"cykl/midi"
	"sync"
)

const (
	defaultTempo = 120.
)

type Grid struct {
	mu sync.Mutex

	midi   midi.Midi
	clock  *clock
	nodes  [][]Node
	Height int
	Width  int

	Key   uint8
	Scale uint8

	Playing bool
	Pulse   uint64

	clipboard [][]Node
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

	grid.clock = newClock(defaultTempo, func() {
		if !grid.Playing {
			return
		}
		grid.Update()
	})

	return grid
}

func (g *Grid) TogglePlay() {
	g.Playing = !g.Playing
	if !g.Playing {
		g.Reset()
		g.midi.SilenceAll()
	}
}

func (g *Grid) SetTempo(tempo float64) {
	g.clock.setTempo(tempo)
}

func (g *Grid) Tempo() float64 {
	return g.clock.tempo
}

func (g *Grid) MidiDevice() string {
	if g.midi.ActiveDevice() == nil {
		return "no midi device"
	}
	return g.midi.ActiveDevice().String()
}

func (g *Grid) CycleMidiDevice() {
	g.midi.CycleMidiDevices()
}

func (g *Grid) QuarterNote() bool {
	if !g.Playing {
		return false
	}
	return g.Pulse/uint64(pulsesPerStep)%uint64(stepsPerQuarterNote) == 1
}

func (g *Grid) CopyOrCut(startX, startY, endX, endY int, cut bool) {
	nodes := make([][]Node, endY-startY+1)
	for i := range nodes {
		nodes[i] = make([]Node, endX-startX+1)
	}
	count := 0
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			_, ok := g.nodes[y][x].(Emitter)
			if ok {
				nodes[y-startY][x-startX] = g.nodes[y][x]
				count++
			}
			if ok && cut {
				g.nodes[y][x] = nil
			}
		}
	}
	if count == 0 {
		return
	}
	g.clipboard = nodes
}

func (g *Grid) Paste(startX, startY, endX, endY int) {
	h, w := len(g.clipboard), len(g.clipboard[0])
	for y := 0; y < h && startY+y <= endY; y++ {
		for x := 0; x < w && startX+x <= endX; x++ {
			if _, ok := g.clipboard[y][x].(Emitter); !ok {
				continue
			}
			g.nodes[startY+y][startX+x] = g.clipboard[y][x].(Emitter).Copy()
		}
	}
}

func (g *Grid) Nodes() [][]Node {
	return g.nodes
}

func (g *Grid) Node(x, y int) Node {
	return g.nodes[y][x]
}

func (g *Grid) AddNodeFromSymbol(symbol string, x, y int) {
	switch symbol {
	case "b":
		g.AddNode(NewInitEmitter(g.midi, NONE), x, y)
	case "s":
		g.AddNode(NewSimpleEmitter(g.midi, NONE), x, y)
	}
}

func (g *Grid) AddNode(node Node, x, y int) {
	g.nodes[y][x] = node
}

func (g *Grid) RemoveNodes(startX, startY, endX, endY int) {
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			g.nodes[y][x] = nil
		}
	}
}

func (g *Grid) ToggleNodeMutes(startX, startY, endX, endY int) {
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			if _, ok := g.nodes[y][x].(Emitter); !ok {
				continue
			}
			g.nodes[y][x].(Emitter).SetMute(!g.nodes[y][x].(Emitter).Muted())
		}
	}
}

func (g *Grid) SetAllNodeMutes(mute bool) {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if _, ok := g.nodes[y][x].(Emitter); !ok {
				continue
			}
			g.nodes[y][x].(Emitter).SetMute(mute)
		}
	}
}

func (g *Grid) Update() {
	g.mu.Lock()
	defer g.mu.Unlock()
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.nodes[y][x] == nil {
				continue
			}

			if n, ok := g.nodes[y][x].(Movable); ok {
				n.Move(g, x, y)
			}

			if n, ok := g.nodes[y][x].(Emitter); ok {
				n.Trig(g.Pulse)
				n.Emit(g, x, y)
			}
		}
	}
	g.Pulse++
}

func (g *Grid) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Playing = false
	g.Pulse = 0
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if _, ok := g.nodes[y][x].(Movable); ok {
				g.nodes[y][x] = nil
			}

			if n, ok := g.nodes[y][x].(Emitter); ok {
				n.Reset()
			}
		}
	}
}

func (g *Grid) Emit(x, y int, direction Direction) {
	newX, newY := direction.NextPosition(x, y)
	if newX == x && newY == y {
		return
	}
	if n, ok := g.nodes[newY][newX].(Emitter); ok {
		n.Arm()
		n.Trig(g.Pulse)
		return
	}
	g.nodes[newY][newX] = NewSignal(direction, g.Pulse)
}

func (g *Grid) Move(x, y int, direction Direction) {
	newX, newY := direction.NextPosition(x, y)

	if newX >= g.Width || newY >= g.Height ||
		newX < 0 || newY < 0 {
		g.nodes[y][x] = nil
		return
	}

	if g.nodes[newY][newX] == nil {
		g.nodes[newY][newX] = g.nodes[y][x]
	} else if n, ok := g.nodes[newY][newX].(Emitter); ok {
		n.Arm()
		n.Trig(g.Pulse)
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
