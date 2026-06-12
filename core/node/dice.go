package node

import (
	"math"
	"math/rand"
	"time"

	"signls/core/common"
	"signls/core/music"
	"signls/midi"
)

type DiceEmitter struct {
	rand   *rand.Rand
	repeat *common.ControlValue[int]
	last   int
	count  int
}

func NewDiceEmitter(midi midi.Midi, device *midi.Device, direction common.Direction) *Emitter {
	source := rand.NewSource(time.Now().UnixNano())
	return &Emitter{
		direction: direction,
		note:      music.NewNote(midi, device),
		behavior: &DiceEmitter{
			rand:   rand.New(source),
			repeat: common.NewControlValue[int](0, 0, math.MaxInt32),
		},
	}
}

func (e *DiceEmitter) EmitDirections(dir common.Direction, inDir common.Direction, pulse uint64) common.Direction {
	if dir.Count() == 0 {
		return common.NONE
	}
	if e.count < e.repeat.Last() {
		e.count++
		return dir.Decompose()[e.last]
	}
	e.repeat.Computed()
	e.count = 0
	e.last = e.rand.Intn(dir.Count())
	return dir.Decompose()[e.last]
}

func (e *DiceEmitter) Repeat() *common.ControlValue[int] {
	return e.repeat
}

func (e *DiceEmitter) ShouldPropagate() bool {
	return false
}

func (e *DiceEmitter) ArmedOnStart() bool {
	return false
}

func (e *DiceEmitter) Copy() common.EmitterBehavior {
	source := rand.NewSource(time.Now().UnixNano())
	newRepeat := *e.repeat
	return &DiceEmitter{
		rand:   rand.New(source),
		repeat: &newRepeat,
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
