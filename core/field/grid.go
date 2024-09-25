package field

import (
	"sync"

	"cykl/core/common"
	"cykl/core/music"
	"cykl/core/node"
	"cykl/midi"
)

// Constants for default settings
const (
	defaultTempo               = 120.            // Default tempo (BPM)
	defaultRootKey music.Key   = 60              // Default root key (MIDI note number for Middle C)
	defaultScale   music.Scale = music.CHROMATIC // Default scale (chromatic scale)
)

// Grid represents the main structure for the grid-based sequencer.
type Grid struct {
	mu sync.Mutex // Mutex to handle concurrent access to the grid

	midi   midi.Midi       // MIDI interface to send notes and control signals
	clock  *common.Clock   // Clock to manage timing and tempo
	nodes  [][]common.Node // 2D slice to store nodes (emitters, signals, etc.)
	Height int             // Height of the grid
	Width  int             // Width of the grid

	Key   music.Key   // Current root key of the grid
	Scale music.Scale // Current scale of the grid

	Playing bool   // Flag to indicate whether the grid is currently playing
	pulse   uint64 // Global pulse counter for timing events

	clipboard [][]common.Node // Clipboard to store nodes for copy-paste operations
}

// NewGrid initializes and returns a new Grid with the given dimensions and MIDI interface.
func NewGrid(width, height int, midi midi.Midi) *Grid {
	grid := &Grid{
		midi:   midi,
		nodes:  make([][]common.Node, height), // Initialize the grid with the specified height
		Height: height,
		Width:  width,
		Key:    defaultRootKey,
		Scale:  defaultScale,
	}
	for i := range grid.nodes {
		grid.nodes[i] = make([]common.Node, width) // Initialize each row with the specified width
	}

	// Create a new clock to manage timing, using the default tempo.
	grid.clock = common.NewClock(defaultTempo, func() {
		if !grid.Playing {
			return
		}
		grid.Update()
	})

	return grid
}

// TogglePlay toggles the playing state of the grid.
func (g *Grid) TogglePlay() {
	g.Playing = !g.Playing
	if !g.Playing {
		g.Reset()
		g.midi.SilenceAll()
	}
}

// SetTempo sets the tempo of the grid.
func (g *Grid) SetTempo(tempo float64) {
	g.clock.SetTempo(tempo)
}

// Tempo returns the current tempo.
func (g *Grid) Tempo() float64 {
	return g.clock.Tempo()
}

// SetKey changes the root key of the grid and transposes all notes accordingly.
func (g *Grid) SetKey(key music.Key) {
	g.Key = key
	g.Transpose()
}

// SetScale changes the scale of the grid and transposes all notes accordingly.
func (g *Grid) SetScale(scale music.Scale) {
	g.Scale = scale
	g.Transpose()
}

// MidiDevice returns the name of the currently active MIDI device.
func (g *Grid) MidiDevice() string {
	if g.midi.ActiveDevice() == nil {
		return "no midi device"
	}
	return g.midi.ActiveDevice().String()
}

// CycleMidiDevice switches to the next available MIDI device.
func (g *Grid) CycleMidiDevice() {
	g.midi.CycleMidiDevices()
}

// Pulse returns the current pulse count divided by the number of pulses per step.
func (g *Grid) Pulse() uint64 {
	return g.pulse / uint64(common.PulsesPerStep)
}

// QuarterNote checks if the current pulse aligns with a quarter note.
func (g *Grid) QuarterNote() bool {
	if !g.Playing {
		return false
	}
	return g.pulse/uint64(common.PulsesPerStep)%uint64(common.StepsPerQuarterNote) == 0
}

// CopyOrCut copies or cuts a selection of nodes from the grid to the clipboard.
func (g *Grid) CopyOrCut(startX, startY, endX, endY int, cut bool) {
	nodes := make([][]common.Node, endY-startY+1) // Initialize the clipboard with the selection size
	for i := range nodes {
		nodes[i] = make([]common.Node, endX-startX+1)
	}
	count := 0
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			_, ok := g.nodes[y][x].(common.Copyable)
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

// Paste pastes nodes from the clipboard into the grid at the specified location.
func (g *Grid) Paste(startX, startY, endX, endY int) {
	if len(g.clipboard) == 0 {
		return
	}
	h, w := len(g.clipboard), len(g.clipboard[0])
	for y := 0; y < h && startY+y <= endY; y++ {
		for x := 0; x < w && startX+x <= endX; x++ {
			if _, ok := g.clipboard[y][x].(common.Copyable); !ok {
				continue
			}
			g.nodes[startY+y][startX+x] = g.clipboard[y][x].(common.Copyable).Copy(startX+x, startY+y)
		}
	}
}

// Nodes returns the entire grid of nodes.
func (g *Grid) Nodes() [][]common.Node {
	return g.nodes
}

// Node returns a specific node from the grid at the given coordinates.
func (g *Grid) Node(x, y int) common.Node {
	return g.nodes[y][x]
}

// AddNodeFromSymbol adds a node to the grid based on a given symbol.
func (g *Grid) AddNodeFromSymbol(symbol string, x, y int) {
	switch symbol {
	case "b":
		g.AddNode(node.NewBangEmitter(g.midi, common.NONE, !g.Playing), x, y)
	case "s":
		g.AddNode(node.NewSpreadEmitter(g.midi, common.NONE), x, y)
	case "c":
		g.AddNode(node.NewCycleEmitter(g.midi, common.NONE), x, y)
	case "d":
		g.AddNode(node.NewDiceEmitter(g.midi, common.NONE), x, y)
	case "q":
		g.AddNode(node.NewQuotaEmitter(g.midi, common.NONE), x, y)
	case "e":
		g.AddNode(node.NewEuclidEmitter(g.midi, common.NONE), x, y)
	case "z":
		g.AddNode(node.NewZoneEmitter(g.midi, common.NONE), x, y)
	case "p":
		g.AddNode(node.NewPassEmitter(g.midi, common.NONE), x, y)
	case "h":
		g.AddNode(node.NewHoleEmitter(common.NONE, x, y), x, y)
	}
}

// AddNode adds a node to the grid at the specified coordinates.
func (g *Grid) AddNode(e common.Node, x, y int) {
	destinationNode, isDestBehavior := g.nodes[y][x].(common.Behavioral)
	newNode, isNewBehavior := e.(common.Behavioral)
	if isDestBehavior && isNewBehavior {
		destinationNode.SetBehavior(newNode.Behavior())
		return
	}
	g.nodes[y][x] = e
}

// RemoveNodes removes nodes from a specified rectangular region of the grid.
func (g *Grid) RemoveNodes(startX, startY, endX, endY int) {
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			g.nodes[y][x] = nil
		}
	}
}

// ToggleNodeMutes toggles the mute state for all nodes in a specified region.
func (g *Grid) ToggleNodeMutes(startX, startY, endX, endY int) {
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			if _, ok := g.nodes[y][x].(music.Audible); !ok {
				continue
			}
			g.nodes[y][x].(music.Audible).SetMute(!g.nodes[y][x].(music.Audible).Muted())
		}
	}
}

// SetAllNodeMutes sets the mute state for all nodes in the grid.
func (g *Grid) SetAllNodeMutes(mute bool) {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if _, ok := g.nodes[y][x].(music.Audible); !ok {
				continue
			}
			g.nodes[y][x].(music.Audible).SetMute(mute)
		}
	}
}

// Update advances the grid by one step, moving signals and triggering emitters.
func (g *Grid) Update() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.pulse%uint64(common.PulsesPerStep) != 0 {
		g.Tick()
		return
	}
	for y := g.Height - 1; y >= 0; y-- {
		for x := g.Width - 1; x >= 0; x-- {
			if g.nodes[y][x] == nil {
				continue
			}

			if n, ok := g.nodes[y][x].(common.Tickable); ok {
				n.Tick()
			}

			if n, ok := g.nodes[y][x].(common.Movable); ok {
				g.Move(n, x, y)
			}

			if n, ok := g.nodes[y][x].(music.Audible); ok {
				n.Trig(g.Key, g.Scale, common.NONE, g.pulse)
				g.Emit(n, x, y)
			}
		}
	}
	g.pulse++
}

// Tick updates all active notes within the grid on every pulse.
func (g *Grid) Tick() {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if n, ok := g.nodes[y][x].(common.Tickable); ok {
				n.Tick()
			}
		}
	}
	g.pulse++
}

// Transpose transposes all notes in the grid to match the current key and scale.
func (g *Grid) Transpose() {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if n, ok := g.nodes[y][x].(music.Audible); ok {
				n.Note().Transpose(g.Key, g.Scale)
			}
		}
	}
}

// Reset stops playback and resets the grid to its initial state.
func (g *Grid) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Playing = false
	g.pulse = 0
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if _, ok := g.nodes[y][x].(common.Movable); ok {
				g.nodes[y][x] = nil
			}

			if n, ok := g.nodes[y][x].(common.Tickable); ok {
				n.Reset()
			}
		}
	}
}

// Emit generates a signal at the specified coordinates and direction.
func (g *Grid) Emit(emitter music.Audible, x, y int) {
	for _, direction := range emitter.Emit(g.pulse) {
		newX, newY := direction.NextPosition(x, y)
		if (newX == x && newY == y) || g.outOfBounds(newX, newY) {
			continue
		}

		if n, ok := g.nodes[newY][newX].(common.Behavioral); ok && n.Behavior().ShouldPropagate() {
			g.PropagateZone(g.nodes[newY][newX].(*node.Emitter), direction, newX, newY)
			continue
		} else if n, ok := g.nodes[newY][newX].(music.Audible); ok {
			n.Arm()
			n.Trig(g.Key, g.Scale, direction, g.pulse)
			continue
		} else if n, ok := g.nodes[newY][newX].(*node.HoleEmitter); ok {
			g.Teleport(n, node.NewSignal(direction, g.pulse), newX, newY)
			continue
		} else if n, ok := g.nodes[newY][newX].(*node.Signal); ok {
			g.Move(n, newX, newY)
		}
		g.nodes[newY][newX] = node.NewSignal(direction, g.pulse)
	}
}

// Move moves a node in the specified direction.
func (g *Grid) Move(movable common.Movable, x, y int) {
	if !movable.MustMove(g.pulse) {
		return
	}

	direction := movable.(common.Node).Direction()
	newX, newY := direction.NextPosition(x, y)

	if g.outOfBounds(newX, newY) {
		g.nodes[y][x] = nil
		return
	}

	if g.nodes[newY][newX] == nil {
		g.nodes[newY][newX] = g.nodes[y][x]
	} else if n, ok := g.nodes[newY][newX].(common.Behavioral); ok && n.Behavior().ShouldPropagate() {
		g.PropagateZone(g.nodes[newY][newX].(*node.Emitter), direction, newX, newY)
	} else if n, ok := g.nodes[newY][newX].(music.Audible); ok {
		n.Arm()
		n.Trig(g.Key, g.Scale, direction, g.pulse)
	} else if n, ok := g.nodes[newY][newX].(*node.HoleEmitter); ok {
		g.Teleport(n, g.nodes[y][x], newX, newY)
	} else if n, ok := g.nodes[newY][newX].(*node.Signal); ok {
		g.Move(n, newX, newY)
		g.nodes[newY][newX] = g.nodes[y][x]
	}

	g.nodes[y][x] = nil
}

func (g *Grid) PropagateZone(e *node.Emitter, direction common.Direction, x, y int) {
	if e == nil {
		return
	}
	e.Arm()
	e.Trig(g.Key, g.Scale, direction, g.pulse)
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			newX, newY := x+dx, y+dy
			if newX < 0 || newX >= g.Width ||
				newY < 0 || newY >= g.Height {
				continue
			}
			if n, ok := g.nodes[newY][newX].(*node.Emitter); ok && !n.Activated() && n.Behavior().ShouldPropagate() {
				g.PropagateZone(n, direction, newX, newY)
			} else if n, ok := g.nodes[newY][newX].(*node.Emitter); ok && !n.Activated() {
				n.Arm()
				n.Trig(g.Key, g.Scale, direction, g.pulse)
			} else if n, ok := g.nodes[newY][newX].(*node.HoleEmitter); ok {
				g.Teleport(n, node.NewSignal(direction, g.pulse), newX, newY)
			}
		}
	}
}

func (g *Grid) Teleport(t *node.HoleEmitter, m common.Node, x, y int) {
	teleportX, teleportY := t.Teleport()
	if g.outOfBounds(teleportX, teleportY) {
		return
	}
	if x == teleportX && y == teleportY {
		return
	}
	if n, ok := g.nodes[teleportY][teleportX].(music.Audible); ok {
		n.Arm()
		n.Trig(g.Key, g.Scale, common.NONE, g.pulse)
	} else if n, ok := g.nodes[teleportY][teleportX].(*node.HoleEmitter); ok {
		g.Teleport(n, m, teleportX, teleportY)
	} else if g.nodes[teleportY][teleportX] == nil {
		g.nodes[teleportY][teleportX] = m
	}
}

// Resize changes the size of the grid and preserves existing nodes within the new dimensions.
func (g *Grid) Resize(newWidth, newHeight int) {
	newNodes := make([][]common.Node, newHeight)
	for i := range newNodes {
		newNodes[i] = make([]common.Node, newWidth)
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

// outOfBounds checks if the specified coordinates are outside the grid dimensions.
func (g *Grid) outOfBounds(x, y int) bool {
	return x >= g.Width || y >= g.Height || x < 0 || y < 0
}
