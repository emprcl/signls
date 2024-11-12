package node

import (
	"math"

	"signls/core/common"
	"signls/core/music"
	"signls/midi"
)

type TollEmitter struct {
	Threshold *common.ControlValue[int]
	count     int
}

func NewTollEmitter(midi midi.Midi, device int, direction common.Direction) *Emitter {
	return &Emitter{
		direction: direction,
		note:      music.NewNote(midi, device),
		behavior: &TollEmitter{
			Threshold: common.NewControlValue[int](1, 1, math.MaxInt32),
		},
	}
}

func (e *TollEmitter) EmitDirections(dir common.Direction, inDir common.Direction, pulse uint64) common.Direction {
	e.count++
	if e.count < e.Threshold.Last() {
		return common.NONE
	}
	e.count = 0
	e.Threshold.Computed()
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
	return "39"
}

func (e *TollEmitter) Reset() {
	e.count = 0
}
