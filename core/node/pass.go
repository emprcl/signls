package node

import (
	"signls/core/common"
	"signls/core/music"
	"signls/midi"
)

type PassEmitter struct{}

func NewPassEmitter(midi midi.Midi, device int, direction common.Direction) *Emitter {
	return &Emitter{
		direction: direction,
		note:      music.NewNote(midi, device),
		behavior:  &PassEmitter{},
	}
}

func (e *PassEmitter) EmitDirections(dir common.Direction, inDir common.Direction, pulse uint64) common.Direction {
	return inDir
}

func (e *PassEmitter) ArmedOnStart() bool {
	return false
}

func (e *PassEmitter) ShouldPropagate() bool {
	return false
}

func (e *PassEmitter) Copy() common.EmitterBehavior {
	return &PassEmitter{}
}

func (e *PassEmitter) Symbol() string {
	return "P"
}

func (e *PassEmitter) Name() string {
	return "pass"
}

func (e *PassEmitter) Color() string {
	return "35"
}

func (e *PassEmitter) Reset() {}
