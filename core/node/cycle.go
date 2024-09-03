package node

import (
	"fmt"

	"cykl/core/common"
	"cykl/core/music"
	"cykl/midi"
)

// CycleEmitter defines an emitter behavior that cycles through multiple directions
// each time it's triggered. It tracks the next direction to emit using the `next` field.
type CycleEmitter struct {
	next int // Keeps track of the next direction to emit.
}

// NewCycleEmitter creates and returns a new Emitter instance with the CycleEmitter behavior.
// It initializes the emitter with the provided MIDI interface and direction.
func NewCycleEmitter(midi midi.Midi, direction common.Direction) *Emitter {
	return &Emitter{
		direction: direction,
		note:      music.NewNote(midi),
		behavior:  &CycleEmitter{},
	}
}

// EmitDirections cycles through the provided directions, returning the next direction in sequence.
// If there are no directions, it returns NONE. The cycling is done in a round-robin manner.
func (e *CycleEmitter) EmitDirections(dir common.Direction, pulse uint64) common.Direction {
	if dir.Count() == 0 {
		return common.NONE
	}
	d := e.next % dir.Count()
	e.next = (e.next + 1) % dir.Count()
	return dir.Decompose()[d]
}

// ArmedOnStart returns false, indicating that the CycleEmitter is not armed when the grid starts.
func (e *CycleEmitter) ArmedOnStart() bool {
	return false
}

// Symbol returns the visual representation of the emitter on the grid.
func (e *CycleEmitter) Symbol(dir common.Direction) string {
	return fmt.Sprintf("%s%s", "C", dir.Symbol())
}

// Name returns the name of this emitter behavior.
func (e *CycleEmitter) Name() string {
	return "cycle"
}

// Color returns the color code associated with this emitter behavior.
func (e *CycleEmitter) Color() string {
	return "063"
}
