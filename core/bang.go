package core

import (
	"cykl/midi"
	"fmt"
)

type BangEmitter struct{}

func NewBangEmitter(midi midi.Midi, direction Direction, armed bool) *BaseEmitter {
	return &BaseEmitter{
		direction: direction,
		armed:     armed,
		note:      NewNote(midi),
		behavior:  &BangEmitter{},
	}
}

func (e *BangEmitter) ArmedOnStart() bool {
	return true
}

func (e *BangEmitter) EmitDirections(dir Direction, pulse uint64) Direction {
	return dir
}

func (e *BangEmitter) Symbol(dir Direction) string {
	return fmt.Sprintf("%s%s", "B", dir.Symbol())
}

func (e *BangEmitter) Name() string {
	return "bang"
}

func (e *BangEmitter) Color() string {
	return "165"
}
