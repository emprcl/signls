package core

import (
	"cykl/midi"
	"fmt"
)

type SequencialEmitter struct {
	next int
}

func NewSequencialEmitter(midi midi.Midi, direction Direction) *BaseEmitter {
	return &BaseEmitter{
		direction: direction,
		note:      NewNote(midi),
		behavior:  SequencialEmitter{},
	}
}

func (e SequencialEmitter) EmitDirections(dir Direction, pulse uint64) Direction {
	if dir.Count() == 0 {
		return NONE
	}
	d := e.next
	e.next = (e.next + 1) % dir.Count()
	return dir.Decompose()[d]
}

func (e SequencialEmitter) ArmedOnStart() bool {
	return false
}

func (e SequencialEmitter) Symbol(dir Direction) string {
	return fmt.Sprintf("%s%s", "S", dir.Symbol())
}

func (e SequencialEmitter) Name() string {
	return "sequ"
}

func (e SequencialEmitter) Color() string {
	return "063"
}
