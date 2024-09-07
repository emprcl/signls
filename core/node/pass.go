package node

import (
	"fmt"
	"math/rand"

	"cykl/core/common"
	"cykl/core/music"
	"cykl/midi"
)

type PassEmitter struct {
	rand *rand.Rand
}

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

func (e *PassEmitter) Symbol(dir common.Direction) string {
	return fmt.Sprintf("%s%s", "P", "⇌")
}

func (e *PassEmitter) Name() string {
	return "pass"
}

func (e *PassEmitter) Color() string {
	return "35"
}
