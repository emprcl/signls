package core

import (
	"cykl/midi"
	"fmt"
)

// CycleEmitter defines an emitter behavior that cycles through multiple directions
// each time it's triggered. It tracks the next direction to emit using the `next` field.
type CycleEmitter struct {
	next int // Keeps track of the next direction to emit.
}

// NewCycleEmitter creates and returns a new Emitter instance with the CycleEmitter behavior.
// It initializes the emitter with the provided MIDI interface and direction.
func NewCycleEmitter(midi midi.Midi, direction Direction) *Emitter {
	return &Emitter{
		direction: direction,       // The set of directions this emitter will cycle through.
		note:      NewNote(midi),   // Create a new Note instance associated with this emitter.
		behavior:  &CycleEmitter{}, // Set the behavior to CycleEmitter.
	}
}

// EmitDirections cycles through the provided directions, returning the next direction in sequence.
// If there are no directions, it returns NONE. The cycling is done in a round-robin manner.
func (e *CycleEmitter) EmitDirections(dir Direction, pulse uint64) Direction {
	if dir.Count() == 0 {
		return NONE // Return NONE if no directions are available.
	}
	d := e.next % dir.Count()           // Determine the next direction index.
	e.next = (e.next + 1) % dir.Count() // Update the index to the next direction for future emissions.
	return dir.Decompose()[d]           // Return the decomposed direction corresponding to the current index.
}

// ArmedOnStart returns false, indicating that the CycleEmitter is not armed when the grid starts.
func (e *CycleEmitter) ArmedOnStart() bool {
	return false
}

// Symbol returns the visual representation of the emitter on the grid.
func (e *CycleEmitter) Symbol(dir Direction) string {
	return fmt.Sprintf("%s%s", "C", dir.Symbol()) // "C" for CycleEmitter plus direction symbol.
}

// Name returns the name of this emitter behavior.
func (e *CycleEmitter) Name() string {
	return "cycle"
}

// Color returns the color code associated with this emitter behavior.
func (e *CycleEmitter) Color() string {
	return "063"
}
