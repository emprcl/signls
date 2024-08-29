package core

import (
	"fmt"

	"cykl/midi"
)

// SpreadEmitter is a type of emitter that generates signals in a spread pattern.
// It implements the behavior for spreading signals across multiple directions.
type SpreadEmitter struct{}

// NewSpreadEmitter creates a new Emitter configured with the SpreadEmitter behavior.
// It initializes the emitter with the given MIDI interface and starting direction.
func NewSpreadEmitter(midi midi.Midi, direction Direction) *Emitter {
	return &Emitter{
		direction: direction,
		note:      NewNote(midi),
		behavior:  &SpreadEmitter{},
	}
}

// EmitDirections determines the direction in which the emitter should propagate signals.
// For SpreadEmitter, it returns the same direction as provided without modification.
func (e *SpreadEmitter) EmitDirections(dir Direction, pulse uint64) Direction {
	return dir
}

// ArmedOnStart indicates whether the emitter is armed and ready to emit at the start.
// For SpreadEmitter, it returns false, meaning it is not armed initially.
func (e *SpreadEmitter) ArmedOnStart() bool {
	return false
}

// Symbol returns a string representation of the emitter's symbol.
func (e *SpreadEmitter) Symbol(dir Direction) string {
	return fmt.Sprintf("%s%s", "S", dir.Symbol())
}

// Name returns the name of the emitter.
func (e *SpreadEmitter) Name() string {
	return "spread"
}

// Color returns a string representing the color code for the emitter.
func (e *SpreadEmitter) Color() string {
	return "177"
}
