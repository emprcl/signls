package node

import (
	"math/rand"
	"time"

	"signls/core/common"
	"signls/core/music"
	"signls/midi"
)

type DiceEmitter struct {
	rand *rand.Rand
}

func NewDiceEmitter(midi midi.Midi, direction common.Direction) *Emitter {
	source := rand.NewSource(time.Now().UnixNano())
	return &Emitter{
		direction: direction,
		note:      music.NewNote(midi),
		behavior: &DiceEmitter{
			rand: rand.New(source),
		},
	}
}

func (e *DiceEmitter) EmitDirections(dir common.Direction, inDir common.Direction, pulse uint64) common.Direction {
	if dir.Count() == 0 {
		return common.NONE
	}
	d := e.rand.Intn(dir.Count())
	return dir.Decompose()[d]
}

func (e *DiceEmitter) ShouldPropagate() bool {
	return false
}

func (e *DiceEmitter) ArmedOnStart() bool {
	return false
}

func (e *DiceEmitter) Copy() common.EmitterBehavior {
	source := rand.NewSource(time.Now().UnixNano())
	return &DiceEmitter{
		rand: rand.New(source),
	}
}

func (e *DiceEmitter) Symbol() string {
	return "D"
}

func (e *DiceEmitter) Name() string {
	return "dice"
}

func (e *DiceEmitter) Color() string {
	return "33"
}

func (e *DiceEmitter) Reset() {}
