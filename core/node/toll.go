package node

import (
	"cykl/core/common"
	"cykl/core/music"
	"cykl/midi"
)

type TollEmitter struct {
	Threshold int
	count     int
}

func NewTollEmitter(midi midi.Midi, direction common.Direction) *Emitter {
	return &Emitter{
		direction: direction,
		note:      music.NewNote(midi),
		behavior:  &TollEmitter{},
	}
}

func (e *TollEmitter) EmitDirections(dir common.Direction, inDir common.Direction, pulse uint64) common.Direction {
	e.count++
	if e.count < e.Threshold {
		return common.NONE
	}
	e.count = 0
	return dir
}

func (e *TollEmitter) ArmedOnStart() bool {
	return false
}

func (e *TollEmitter) ShouldPropagate() bool {
	return false
}

func (e *TollEmitter) Copy() common.EmitterBehavior {
	return &TollEmitter{
		Threshold: e.Threshold,
	}
}

func (e *TollEmitter) Symbol() string {
	return "T"
}

func (e *TollEmitter) Name() string {
	return "toll"
}

func (e *TollEmitter) Color() string {
	return "197"
}

func (e *TollEmitter) Reset() {
	e.count = 0
}
