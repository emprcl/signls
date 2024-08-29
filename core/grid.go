package core

import (
	"cykl/midi"
	"sync"
)

// Constants for default settings
const (
	defaultTempo         = 120.      // Default tempo (BPM)
	defaultRootKey Key   = 60        // Default root key (MIDI note number for Middle C)
	defaultScale   Scale = CHROMATIC // Default scale (chromatic scale)
)

// Grid represents the main structure for the grid-based sequencer.
type Grid struct {
	mu sync.Mutex // Mutex to handle concurrent access to the grid

	midi   midi.Midi // MIDI interface to send notes and control signals
	clock  *clock    // Clock to manage timing and tempo
	nodes  [][]Node  // 2D slice to store nodes (emitters, signals, etc.)
	Height int       // Height of the grid
	Width  int       // Width of the grid

	Key   Key   // Current root key of the grid
	Scale Scale // Current scale of the grid

	Playing bool   // Flag to indicate whether the grid is currently playing
	pulse   uint64 // Global pulse counter for timing events

	clipboard [][]Node // Clipboard to store nodes for copy-paste operations
}

// NewGrid initializes and returns a new Grid with the given dimensions and MIDI interface.
func NewGrid(width, height int, midi midi.Midi) *Grid {
	grid := &Grid{
		midi:   midi,
		nodes:  make([][]Node, height), // Initialize the grid with the specified height
		Height: height,
		Width:  width,
		Key:    defaultRootKey,
		Scale:  defaultScale,
	}
	for i := range grid.nodes {
		grid.nodes[i] = make([]Node, width) // Initialize each row with the specified width
	}

	// Create a new clock to manage timing, using the default tempo.
	grid.clock = newClock(defaultTempo, func() {
		if !grid.Playing {
			return
		}
		grid.Update() // Update the grid on each clock tick
	})

	return grid
}

// TogglePlay toggles the playing state of the grid.
func (g *Grid) TogglePlay() {
	g.Playing = !g.Playing
	if !g.Playing {
		g.Reset()           // Reset the grid when stopping
		g.midi.SilenceAll() // Silence all MIDI notes
	}
}

// SetTempo sets the tempo of the grid.
func (g *Grid) SetTempo(tempo float64) {
	g.clock.setTempo(tempo)
}

// Tempo returns the current tempo.
func (g *Grid) Tempo() float64 {
	return g.clock.tempo
}

// SetKey changes the root key of the grid and transposes all notes accordingly.
func (g *Grid) SetKey(key Key) {
	g.Key = key
	g.Transpose() // Transpose all notes to the new key
}

// SetScale changes the scale of the grid and transposes all notes accordingly.
func (g *Grid) SetScale(scale Scale) {
	g.Scale = scale
	g.Transpose() // Transpose all notes to the new scale
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
	return g.pulse / uint64(pulsesPerStep)
}

// QuarterNote checks if the current pulse aligns with a quarter note.
func (g *Grid) QuarterNote() bool {
	if !g.Playing {
		return false
	}
	return g.pulse/uint64(pulsesPerStep)%uint64(stepsPerQuarterNote) == 0
}

// CopyOrCut copies or cuts a selection of nodes from the grid to the clipboard.
func (g *Grid) CopyOrCut(startX, startY, endX, endY int, cut bool) {
	nodes := make([][]Node, endY-startY+1) // Initialize the clipboard with the selection size
	for i := range nodes {
		nodes[i] = make([]Node, endX-startX+1)
	}
	count := 0
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			_, ok := g.nodes[y][x].(*Emitter)
			if ok {
				nodes[y-startY][x-startX] = g.nodes[y][x] // Copy node to clipboard
				count++
			}
			if ok && cut {
				g.nodes[y][x] = nil // Remove node from grid if cutting
			}
		}
	}
	if count == 0 {
		return // If no nodes were copied, exit
	}
	g.clipboard = nodes // Store the copied nodes in the clipboard
}

// Paste pastes nodes from the clipboard into the grid at the specified location.
func (g *Grid) Paste(startX, startY, endX, endY int) {
	h, w := len(g.clipboard), len(g.clipboard[0])
	for y := 0; y < h && startY+y <= endY; y++ {
		for x := 0; x < w && startX+x <= endX; x++ {
			if _, ok := g.clipboard[y][x].(*Emitter); !ok {
				continue // Skip if clipboard does not contain an Emitter
			}
			g.nodes[startY+y][startX+x] = g.clipboard[y][x].(*Emitter).Copy() // Paste a copy of the node
		}
	}
}

// Nodes returns the entire grid of nodes.
func (g *Grid) Nodes() [][]Node {
	return g.nodes
}

// Node returns a specific node from the grid at the given coordinates.
func (g *Grid) Node(x, y int) Node {
	return g.nodes[y][x]
}

// AddNodeFromSymbol adds a node to the grid based on a given symbol.
func (g *Grid) AddNodeFromSymbol(symbol string, x, y int) {
	switch symbol {
	case "b":
		g.AddNode(NewBangEmitter(g.midi, NONE, !g.Playing), x, y)
	case "c":
		g.AddNode(NewCycleEmitter(g.midi, NONE), x, y)
	case "s":
		g.AddNode(NewSpreadEmitter(g.midi, NONE), x, y)
	}
}

// AddNode adds a node to the grid at the specified coordinates.
func (g *Grid) AddNode(node *Emitter, x, y int) {
	if n, ok := g.nodes[y][x].(*Emitter); g.nodes[y][x] != nil && ok {
		n.behavior = node.behavior // Update existing node's behavior
		return
	}
	g.nodes[y][x] = node // Add the new node to the grid
}

// RemoveNodes removes nodes from a specified rectangular region of the grid.
func (g *Grid) RemoveNodes(startX, startY, endX, endY int) {
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			g.nodes[y][x] = nil // Remove node from grid
		}
	}
}

// ToggleNodeMutes toggles the mute state for all nodes in a specified region.
func (g *Grid) ToggleNodeMutes(startX, startY, endX, endY int) {
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			if _, ok := g.nodes[y][x].(*Emitter); !ok {
				continue // Skip if not an Emitter
			}
			g.nodes[y][x].(*Emitter).SetMute(!g.nodes[y][x].(*Emitter).Muted()) // Toggle mute
		}
	}
}

// SetAllNodeMutes sets the mute state for all nodes in the grid.
func (g *Grid) SetAllNodeMutes(mute bool) {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if _, ok := g.nodes[y][x].(*Emitter); !ok {
				continue // Skip if not an Emitter
			}
			g.nodes[y][x].(*Emitter).SetMute(mute) // Set mute state
		}
	}
}

// Update advances the grid by one step, moving signals and triggering emitters.
func (g *Grid) Update() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.pulse%uint64(pulsesPerStep) != 0 {
		g.Tick() // Handle partial pulse updates
		return
	}
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.nodes[y][x] == nil {
				continue // Skip empty nodes
			}

			if n, ok := g.nodes[y][x].(Movable); ok {
				n.Move(g, x, y) // Move the node if it implements the Movable interface
			}

			if n, ok := g.nodes[y][x].(*Emitter); ok {
				n.Trig(g.Key, g.Scale, g.pulse) // Trigger the emitter
				n.Emit(g, x, y)                 // Emit signals from the emitter
			}
		}
	}
	g.pulse++ // Increment the pulse counter
}

// Tick updates all active notes within the grid on every pulse.
func (g *Grid) Tick() {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if n, ok := g.nodes[y][x].(*Emitter); ok {
				n.Note().Tick() // Advance the note's tick counter
			}
		}
	}
	g.pulse++
}

// Transpose transposes all notes in the grid to match the current key and scale.
func (g *Grid) Transpose() {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if n, ok := g.nodes[y][x].(*Emitter); ok {
				n.Note().Transpose(g.Key, g.Scale) // Transpose the note to the current key and scale
			}
		}
	}
}

// Reset stops playback and resets the grid to its initial state.
func (g *Grid) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Playing = false
	g.pulse = 0 // Reset the pulse counter
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if _, ok := g.nodes[y][x].(Movable); ok {
				g.nodes[y][x] = nil // Remove movable nodes from the grid
			}

			if n, ok := g.nodes[y][x].(*Emitter); ok {
				n.Reset() // Reset the emitter
			}
		}
	}
}

// Emit generates a signal at the specified coordinates and direction.
func (g *Grid) Emit(x, y int, direction Direction) {
	newX, newY := direction.NextPosition(x, y)
	if (newX == x && newY == y) || g.outOfBounds(newX, newY) {
		return // Skip if the new position is out of bounds or unchanged
	}
	if n, ok := g.nodes[newY][newX].(*Emitter); ok {
		n.Arm()
		n.Trig(g.Key, g.Scale, g.pulse) // Trigger the emitter if it exists
		return
	}
	g.nodes[newY][newX] = NewSignal(direction, g.pulse) // Create a new signal
}

// Move moves a node in the specified direction.
func (g *Grid) Move(x, y int, direction Direction) {
	newX, newY := direction.NextPosition(x, y)

	if g.outOfBounds(newX, newY) {
		g.nodes[y][x] = nil // Remove the node if it moves out of bounds
		return
	}

	if g.nodes[newY][newX] == nil {
		g.nodes[newY][newX] = g.nodes[y][x] // Move the node to the new position
	} else if n, ok := g.nodes[newY][newX].(*Emitter); ok {
		n.Arm()
		n.Trig(g.Key, g.Scale, g.pulse) // Trigger the emitter if it exists at the new position
	}

	g.nodes[y][x] = nil // Clear the original position
}

// Resize changes the size of the grid and preserves existing nodes within the new dimensions.
func (g *Grid) Resize(newWidth, newHeight int) {
	newNodes := make([][]Node, newHeight)
	for i := range newNodes {
		newNodes[i] = make([]Node, newWidth)
	}

	minWidth := g.Width
	if newWidth < g.Width {
		minWidth = newWidth // Determine the smaller of the current and new widths
	}

	minHeight := g.Height
	if newHeight < g.Height {
		minHeight = newHeight // Determine the smaller of the current and new heights
	}

	for y := 0; y < minHeight; y++ {
		for x := 0; x < minWidth; x++ {
			newNodes[y][x] = g.nodes[y][x] // Copy existing nodes to the resized grid
		}
	}

	g.Width = newWidth
	g.Height = newHeight
	g.nodes = newNodes // Update the grid with the resized node array
}

// outOfBounds checks if the specified coordinates are outside the grid dimensions.
func (g Grid) outOfBounds(x, y int) bool {
	return x >= g.Width || y >= g.Height || x < 0 || y < 0
}
