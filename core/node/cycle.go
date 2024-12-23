package node

import (
	"math"
	"signls/core/common"
	"signls/core/music"
	"signls/midi"
)

type CycleEmitter struct {
	repeat *common.ControlValue[int]
	count  int
	next   int
}

func NewCycleEmitter(midi midi.Midi, device midi.Device, direction common.Direction) *Emitter {
	return &Emitter{
		direction: direction,
		note:      music.NewNote(midi, device),
		behavior: &CycleEmitter{
			repeat: common.NewControlValue[int](0, 0, math.MaxInt32),
		},
	}
}

func (e *CycleEmitter) EmitDirections(dir common.Direction, inDir common.Direction, pulse uint64) common.Direction {
	if dir.Count() == 0 {
		return common.NONE
	}
	d := e.next % dir.Count()
	if e.count < e.repeat.Last() {
		e.count++
		return dir.Decompose()[d]
	}
	e.Repeat().Computed()
	e.count = 0
	e.next = (e.next + 1) % dir.Count()
	return dir.Decompose()[d]
}

func (e *CycleEmitter) Repeat() *common.ControlValue[int] {
	return e.repeat
}

func (e *CycleEmitter) ShouldPropagate() bool {
	return false
}

func (e *CycleEmitter) ArmedOnStart() bool {
	return false
}

func (e *CycleEmitter) Copy() common.EmitterBehavior {
	newRepeat := *e.repeat
	return &CycleEmitter{
		next:   e.next,
		repeat: &newRepeat,
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
