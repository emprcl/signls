package node

import (
	"fmt"
	"math/rand"
	"time"

	"cykl/core/common"
	"cykl/core/music"
	"cykl/midi"
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

func (e *DiceEmitter) ArmedOnStart() bool {
	return false
}

func (e *DiceEmitter) Symbol(dir common.Direction) string {
	return fmt.Sprintf("%s%s", "D", dir.Symbol())
}

func (e *DiceEmitter) Name() string {
	return "dice"
}

func (e *DiceEmitter) Color() string {
	return "27"
}
