package node

import (
	"signls/core/common"
	"signls/core/music"
	"signls/midi"
)

type CycleEmitter struct {
	next int
}

func NewCycleEmitter(midi midi.Midi, direction common.Direction) *Emitter {
	return &Emitter{
		direction: direction,
		note:      music.NewNote(midi),
		behavior:  &CycleEmitter{},
	}
}

func (e *CycleEmitter) EmitDirections(dir common.Direction, inDir common.Direction, pulse uint64) common.Direction {
	if dir.Count() == 0 {
		return common.NONE
	}
	d := e.next % dir.Count()
	e.next = (e.next + 1) % dir.Count()
	return dir.Decompose()[d]
}

func (e *CycleEmitter) ShouldPropagate() bool {
	return false
}

func (e *CycleEmitter) ArmedOnStart() bool {
	return false
}

func (e *CycleEmitter) Copy() common.EmitterBehavior {
	return &CycleEmitter{
		next: e.next,
	}
}

func (e *CycleEmitter) Symbol() string {
	return "C"
}

func (e *CycleEmitter) Name() string {
	return "cycle"
}

func (e *CycleEmitter) Color() string {
	return "63"
}

func (e *CycleEmitter) Reset() {
	e.next = 0
}
