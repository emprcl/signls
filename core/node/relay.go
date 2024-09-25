package node

import (
	"fmt"

	"cykl/core/common"
	"cykl/core/music"
	"cykl/midi"
)

// RelayEmitter is a type of emitter that generates signals in a spread pattern.
// It implements the behavior for spreading signals across multiple directions.
type RelayEmitter struct{}

// NewRelayEmitter creates a new Emitter configured with the RelayEmitter behavior.
// It initializes the emitter with the given MIDI interface and starting direction.
func NewRelayEmitter(midi midi.Midi, direction common.Direction) *Emitter {
	return &Emitter{
		direction: direction,
		note:      music.NewNote(midi),
		behavior:  &RelayEmitter{},
	}
}

// EmitDirections determines the direction in which the emitter should propagate signals.
// For RelayEmitter, it returns the same direction as provided without modification.
func (e *RelayEmitter) EmitDirections(dir common.Direction, inDir common.Direction, pulse uint64) common.Direction {
	return dir
}

// ShouldPropagate indicates if triggers should be propagated to direct neighbors.
func (e *RelayEmitter) ShouldPropagate() bool {
	return false
}

// ArmedOnStart indicates whether the emitter is armed and ready to emit at the start.
// For RelayEmitter, it returns false, meaning it is not armed initially.
func (e *RelayEmitter) ArmedOnStart() bool {
	return false
}

func (e *RelayEmitter) Copy() common.EmitterBehavior {
	return &RelayEmitter{}
}

// Symbol returns a string representation of the emitter's symbol.
func (e *RelayEmitter) Symbol(dir common.Direction) string {
	return fmt.Sprintf("%s%s", "R", dir.Symbol())
}

// Name returns the name of the emitter.
func (e *RelayEmitter) Name() string {
	return "relay"
}

// Color returns a string representing the color code for the emitter.
func (e *RelayEmitter) Color() string {
	return "177"
}

func (e *RelayEmitter) Reset() {}
