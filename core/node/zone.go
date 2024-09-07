package node

import (
	"cykl/core/common"
	"cykl/core/music"
	"cykl/midi"
)

type ZoneEmitter struct{}

func NewZoneEmitter(midi midi.Midi, direction common.Direction) *Emitter {
	return &Emitter{
		direction: direction,
		note:      music.NewNote(midi),
		behavior:  &ZoneEmitter{},
	}
}

func (e *ZoneEmitter) EmitDirections(dir common.Direction, inDir common.Direction, pulse uint64) common.Direction {
	return common.NONE
}

func (e *ZoneEmitter) ShouldPropagate() bool {
	return true
}

func (e *ZoneEmitter) ArmedOnStart() bool {
	return false
}

func (e *ZoneEmitter) Copy() EmitterBehavior {
	return &ZoneEmitter{}
}

func (e *ZoneEmitter) Symbol(dir common.Direction) string {
	return "Z╬"
}

func (e *ZoneEmitter) Name() string {
	return "zone"
}

func (e *ZoneEmitter) Color() string {
	return "204"
}
