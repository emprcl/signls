package field

import (
	"signls/core/common"
	"signls/core/music"
	"signls/core/music/meta"
	"signls/core/node"
	"signls/core/theory"
	"signls/midi"
	"sync"
)

const (
	defaultTempo                = 120.
	defaultRootKey theory.Key   = 60
	defaultScale   theory.Scale = theory.CHROMATIC
)

// Grid represents the main structure for the grid-based sequencer.
//
// The grid is accessed concurrently: the clock goroutine drives Update on
// every pulse, while the ui goroutine reads node state to render and mutates
// it in response to user input. All access to the shared state below (nodes,
// dimensions, pulse, key/scale, playing state, device) must be guarded by mu.
// Public methods acquire the lock themselves; the ui render path takes a
// consistent copy with Snapshot and wraps any remaining live reads with Read.
type Grid struct {
	mu sync.RWMutex

	midi      midi.Midi
	device    midi.Device
	clock     *common.Clock
	nodes     [][]common.Node
	Height    int
	Width     int
	BankIndex int

	Key   theory.Key
	Scale theory.Scale

	Playing bool

	SendClock     bool
	SendTransport bool

	pulse uint64 // Global pulse counter for timing events

	clipboard [][]common.Node
}

// NewGrid initializes and returns a new Grid with the given dimensions and MIDI interface.
func NewGrid(width, height int, midi midi.Midi, device string) *Grid {
	d := midi.NewDevice(device, "")
	grid := &Grid{
		midi:   midi,
		device: d,
		nodes:  make([][]common.Node, height),
		Height: height,
		Width:  width,
		Key:    defaultRootKey,
		Scale:  defaultScale,
	}
	for i := range grid.nodes {
		grid.nodes[i] = make([]common.Node, width)
	}

	grid.clock = common.NewClock(defaultTempo, func() {
		grid.Update()
	})

	return grid
}

// Read runs fn while holding the read lock. The ui uses it to wrap reads that
// span several grid/node accesses (such as rendering the control bar) so they
// observe a consistent state without exposing the lock itself. For rendering
// the grid, prefer Snapshot, which copies the display state and releases the
// lock immediately.
func (g *Grid) Read(fn func()) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	fn()
}

// TriggerNode manually arms and triggers the emitter at the given coordinates,
// if any. Used by the ui to trigger a node from a key press.
func (g *Grid) TriggerNode(x, y int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	e, ok := g.nodes[y][x].(*node.Emitter)
	if !ok {
		return
	}
	e.Arm()
	e.Trig(g.Key, g.Scale, common.NONE, g.pulse/uint64(common.PulsesPerStep))
}

// TogglePlay toggles the playing state of the grid.
func (g *Grid) TogglePlay() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Playing = !g.Playing
	if !g.Playing {
		g.reset()
		g.midi.SilenceAll()
	}

	if !g.SendTransport {
		return
	}

	if g.Playing {
		g.midi.TransportStart(g.device.ID)
	} else {
		g.midi.TransportStop(g.device.ID)
	}
}

// SetPlaying sets the playing state of the grid. Used by the ui when switching
// banks to preserve the previous playing state.
func (g *Grid) SetPlaying(playing bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Playing = playing
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
func (g *Grid) SetKey(key theory.Key) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Key = key
	g.transpose()
}

// SetScale changes the scale of the grid and transposes all notes accordingly.
func (g *Grid) SetScale(scale theory.Scale) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Scale = scale
	g.transpose()
}

// MidiDevice returns the name of the currently active MIDI device.
func (g *Grid) MidiDevice() midi.Device {
	return g.device
}

// SetMidiDevice sets the midi device.
func (g *Grid) SetMidiDevice(device midi.Device) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.device = device
}

// Midi returns the Midi interface.
func (g *Grid) Midi() midi.Midi {
	return g.midi
}

// Pulse returns the current pulse step.
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
	g.mu.Lock()
	defer g.mu.Unlock()
	nodes := make([][]common.Node, endY-startY+1)
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
	g.mu.Lock()
	defer g.mu.Unlock()
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
	g.mu.Lock()
	defer g.mu.Unlock()
	switch symbol {
	case "b":
		g.addNode(node.NewBangEmitter(g.midi, &g.device, common.NONE, !g.Playing), x, y)
	case "s":
		g.addNode(node.NewSpreadEmitter(g.midi, &g.device, common.NONE), x, y)
	case "c":
		g.addNode(node.NewCycleEmitter(g.midi, &g.device, common.NONE), x, y)
	case "d":
		g.addNode(node.NewDiceEmitter(g.midi, &g.device, common.NONE), x, y)
	case "t":
		g.addNode(node.NewTollEmitter(g.midi, &g.device, common.NONE), x, y)
	case "e":
		g.addNode(node.NewEuclidEmitter(g.midi, &g.device, common.NONE), x, y)
	case "z":
		g.addNode(node.NewZoneEmitter(g.midi, &g.device, common.NONE), x, y)
	case "p":
		g.addNode(node.NewPassEmitter(g.midi, &g.device, common.NONE), x, y)
	case "h":
		g.addNode(node.NewHoleEmitter(common.NONE, x, y, g.Width, g.Height), x, y)
	}
}

// AddNode adds a node to the grid at the specified coordinates.
func (g *Grid) AddNode(e common.Node, x, y int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.addNode(e, x, y)
}

// addNode adds a node to the grid at the specified coordinates. The caller must
// hold the lock.
func (g *Grid) addNode(e common.Node, x, y int) {
	destinationNode, isDestBehavior := g.nodes[y][x].(common.Behavioral)
	newNode, isNewBehavior := e.(common.Behavioral)
	if isDestBehavior && isNewBehavior {
		destinationNode.SetBehavior(newNode.Behavior())
		return
	}
	g.nodes[y][x] = e
}

// RemoveNodes removes nodes from a specified region of the grid.
func (g *Grid) RemoveNodes(startX, startY, endX, endY int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			g.nodes[y][x] = nil
		}
	}
}

// ToggleNodeMutes toggles the mute state for all nodes in a specified region.
func (g *Grid) ToggleNodeMutes(startX, startY, endX, endY int) {
	g.mu.Lock()
	defer g.mu.Unlock()
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
	g.mu.Lock()
	defer g.mu.Unlock()
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
// It is called by the clock goroutine on every pulse.
func (g *Grid) Update() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if !g.Playing {
		return
	}
	if g.SendClock {
		g.midi.SendClock(g.device.ID)
	}
	if g.pulse%uint64(common.PulsesPerStep) != 0 {
		g.tick()
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
				g.ExecuteMetaCommands(n)
				g.Emit(n, x, y)
			}
		}
	}
	g.pulse++
}

// tick updates all active notes within the grid on every pulse. The caller
// must hold the lock.
func (g *Grid) tick() {
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
	g.mu.Lock()
	defer g.mu.Unlock()
	g.transpose()
}

// transpose transposes all notes in the grid to match the current key and
// scale. The caller must hold the lock.
func (g *Grid) transpose() {
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
	g.reset()
}

// reset stops playback and resets the grid to its initial state. The caller
// must hold the lock.
func (g *Grid) reset() {
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

// Emit makes specified emitter generates signals.
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

// ExecuteMetaCommands executes meta commands from the node.
func (g *Grid) ExecuteMetaCommands(node music.Audible) {
	for _, cmd := range node.Note().MetaCommands {
		if !cmd.Executed() {
			continue
		}
		switch c := cmd.(type) {
		case *meta.RootCommand:
			g.Key = theory.Key(c.Value().Computed())
		case *meta.ScaleCommand:
			g.Scale = theory.AllScales()[c.Value().Computed()]
		case *meta.TempoCommand:
			g.SetTempo(float64(c.Value().Computed()))
		case *meta.BankCommand:
			g.BankIndex = c.Value().Computed()
		}

		cmd.Reset()
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

// PropagateZone propagates a trigger to neighboring nodes.
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

// Teleport moves a node through a Hole emitter.
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
	g.mu.Lock()
	defer g.mu.Unlock()
	g.resize(newWidth, newHeight)
}

// resize changes the size of the grid and preserves existing nodes within the
// new dimensions. The caller must hold the lock.
func (g *Grid) resize(newWidth, newHeight int) {
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
