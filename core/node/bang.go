package node

import (
	"signls/core/common"
	"signls/core/music"
	"signls/midi"
)

type BangEmitter struct{}

func NewBangEmitter(midi midi.Midi, device *midi.Device, direction common.Direction, armed bool) *Emitter {
	return &Emitter{
		direction: direction,
		armed:     armed,
		note:      music.NewNote(midi, device),
		behavior:  &BangEmitter{},
	}
}

func (e *BangEmitter) ArmedOnStart() bool {
	return true
}

func (e *BangEmitter) Copy() common.EmitterBehavior {
	return &BangEmitter{}
}

func (e *BangEmitter) EmitDirections(dir common.Direction, inDir common.Direction, pulse uint64) common.Direction {
	return dir
}

func (e *BangEmitter) ShouldPropagate() bool {
	return false
}

func (e *BangEmitter) Symbol() string {
	return "B"
}

func (e *BangEmitter) Name() string {
	return "bang"
}

func (e *BangEmitter) Color() string {
	return "165"
}

func (e *BangEmitter) Reset() {}
