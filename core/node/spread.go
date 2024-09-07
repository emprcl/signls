package node

import (
	"fmt"

	"cykl/core/common"
	"cykl/core/music"
	"cykl/midi"
)

// SpreadEmitter is a type of emitter that generates signals in a spread pattern.
// It implements the behavior for spreading signals across multiple directions.
type SpreadEmitter struct{}

// NewSpreadEmitter creates a new Emitter configured with the SpreadEmitter behavior.
// It initializes the emitter with the given MIDI interface and starting direction.
func NewSpreadEmitter(midi midi.Midi, direction common.Direction) *Emitter {
	return &Emitter{
		direction: direction,
		note:      music.NewNote(midi),
		behavior:  &SpreadEmitter{},
	}
}

// EmitDirections determines the direction in which the emitter should propagate signals.
// For SpreadEmitter, it returns the same direction as provided without modification.
func (e *SpreadEmitter) EmitDirections(dir common.Direction, inDir common.Direction, pulse uint64) common.Direction {
	return dir
}

// ArmedOnStart indicates whether the emitter is armed and ready to emit at the start.
// For SpreadEmitter, it returns false, meaning it is not armed initially.
func (e *SpreadEmitter) ArmedOnStart() bool {
	return false
}

func (e *SpreadEmitter) Copy() EmitterBehavior {
	return &SpreadEmitter{}
}

// Symbol returns a string representation of the emitter's symbol.
func (e *SpreadEmitter) Symbol(dir common.Direction) string {
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
