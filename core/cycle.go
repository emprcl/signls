package core

import (
	"cykl/midi"
	"fmt"
)

type CycleEmitter struct {
	next int
}

func NewCycleEmitter(midi midi.Midi, direction Direction) *Emitter {
	return &Emitter{
		direction: direction,
		note:      NewNote(midi),
		behavior:  &CycleEmitter{},
	}
}

func (e *CycleEmitter) EmitDirections(dir Direction, pulse uint64) Direction {
	if dir.Count() == 0 {
		return NONE
	}
	d := e.next % dir.Count()
	e.next = (e.next + 1) % dir.Count()
	return dir.Decompose()[d]
}

func (e *CycleEmitter) ArmedOnStart() bool {
	return false
}

func (e *CycleEmitter) Symbol(dir Direction) string {
	return fmt.Sprintf("%s%s", "C", dir.Symbol())
}

func (e *CycleEmitter) Name() string {
	return "cycle"
}

func (e *CycleEmitter) Color() string {
	return "063"
}
