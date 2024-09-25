package node

import (
	"fmt"

	"cykl/core/common"
	"cykl/core/music"
	"cykl/midi"
)

type PassEmitter struct{}

func NewPassEmitter(midi midi.Midi, direction common.Direction) *Emitter {
	return &Emitter{
		direction: direction,
		note:      music.NewNote(midi),
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
	return fmt.Sprintf("%s%s", "P", "⇌")
}

func (e *PassEmitter) Name() string {
	return "pass"
}

func (e *PassEmitter) Color() string {
	return "35"
}

func (e *PassEmitter) Reset() {}
