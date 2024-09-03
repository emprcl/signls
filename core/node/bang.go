package node

import (
	"fmt"

	"cykl/core/common"
	"cykl/core/music"
	"cykl/midi"
)

// BangEmitter defines a simple emitter behavior where it triggers and emits
// signals in a specific direction when activated. This emitter is always armed at the start.
type BangEmitter struct{}

// NewBangEmitter creates and returns a new Emitter instance with the BangEmitter behavior.
// It initializes the emitter with the provided MIDI interface, direction, and armed state.
func NewBangEmitter(midi midi.Midi, direction common.Direction, armed bool) *Emitter {
	return &Emitter{
		direction: direction,           // The direction this emitter will emit signals.
		armed:     armed,               // Whether the emitter is armed at initialization.
		note:      music.NewNote(midi), // Create a new Note instance associated with this emitter.
		behavior:  &BangEmitter{},      // Set the behavior to BangEmitter.
	}
}

// ArmedOnStart returns true, indicating that the BangEmitter is always armed when the grid starts.
func (e *BangEmitter) ArmedOnStart() bool {
	return true
}

// EmitDirections returns the current direction. For the BangEmitter, it always
// emits in the direction it's facing without any modification.
func (e *BangEmitter) EmitDirections(dir common.Direction, pulse uint64) common.Direction {
	return dir // Emit in the same direction it's facing.
}

// Symbol returns the visual representation of the emitter on the grid.
// It concatenates "B" with the direction's symbol for a complete identifier.
func (e *BangEmitter) Symbol(dir common.Direction) string {
	return fmt.Sprintf("%s%s", "B", dir.Symbol()) // "B" for BangEmitter plus direction symbol.
}

// Name returns the name of this emitter behavior, which is "bang".
func (e *BangEmitter) Name() string {
	return "bang"
}

// Color returns the color code associated with this emitter behavior.
// In this case, the color code is "165".
func (e *BangEmitter) Color() string {
	return "165"
}
