package node

import (
	"signls/core/common"
	"signls/core/music"
	"signls/midi"
)

type SpreadEmitter struct{}

func NewSpreadEmitter(midi midi.Midi, device *midi.Device, direction common.Direction) *Emitter {
	return &Emitter{
		direction: direction,
		note:      music.NewNote(midi, device),
		behavior:  &SpreadEmitter{},
	}
}

func (e *SpreadEmitter) EmitDirections(dir common.Direction, inDir common.Direction, pulse uint64) common.Direction {
	return dir
}

func (e *SpreadEmitter) ShouldPropagate() bool {
	return false
}

func (e *SpreadEmitter) ArmedOnStart() bool {
	return false
}

func (e *SpreadEmitter) Copy() common.EmitterBehavior {
	return &SpreadEmitter{}
}

func (e *SpreadEmitter) Symbol() string {
	return "S"
}

func (e *SpreadEmitter) Name() string {
	return "spread"
}

func (e *SpreadEmitter) Color() string {
	return "56"
}

func (e *SpreadEmitter) Reset() {}
