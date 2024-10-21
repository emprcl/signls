package node

import (
	"signls/core/common"
	"signls/core/music"
	"signls/midi"
)

type RelayEmitter struct{}

func NewRelayEmitter(midi midi.Midi, direction common.Direction) *Emitter {
	return &Emitter{
		direction: direction,
		note:      music.NewNote(midi),
		behavior:  &RelayEmitter{},
	}
}

func (e *RelayEmitter) EmitDirections(dir common.Direction, inDir common.Direction, pulse uint64) common.Direction {
	return dir
}

func (e *RelayEmitter) ShouldPropagate() bool {
	return false
}

func (e *RelayEmitter) ArmedOnStart() bool {
	return false
}

func (e *RelayEmitter) Copy() common.EmitterBehavior {
	return &RelayEmitter{}
}

func (e *RelayEmitter) Symbol() string {
	return "R"
}

func (e *RelayEmitter) Name() string {
	return "relay"
}

func (e *RelayEmitter) Color() string {
	return "56"
}

func (e *RelayEmitter) Reset() {}
