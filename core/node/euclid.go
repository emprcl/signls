package node

import (
	"cykl/core/common"
	"cykl/core/music"
	"cykl/midi"
	"fmt"
)

const (
	defaultSteps    = 16
	defaultTriggers = 4
	defaultOffset   = 0
)

type EuclidEmitter struct {
	direction common.Direction
	note      *music.Note

	Steps    int
	Triggers int
	Offset   int

	step int

	pulse      uint64
	armed      bool
	triggered  bool
	hasTrigged bool
	muted      bool
}

func NewEuclidEmitter(midi midi.Midi, direction common.Direction) *EuclidEmitter {
	return &EuclidEmitter{
		Steps:     defaultSteps,
		Triggers:  defaultTriggers,
		Offset:    defaultOffset,
		direction: direction,
		note:      music.NewNote(midi),
	}
}

func (e *EuclidEmitter) Copy(dx, dy int) common.Node {
	newNote := *e.note
	return &EuclidEmitter{
		direction: e.direction,
		armed:     e.armed,
		note:      &newNote,
		muted:     e.muted,
		Steps:     e.Steps,
		Triggers:  e.Triggers,
		Offset:    e.Offset,
	}
}

func (e *EuclidEmitter) Activated() bool {
	return e.hasTrigged
}

func (e *EuclidEmitter) Note() *music.Note {
	return e.note
}

func (e *EuclidEmitter) Arm() {
	e.armed = true
}

func (e *EuclidEmitter) SetMute(mute bool) {
	e.note.Stop()
	e.muted = mute
}

func (e *EuclidEmitter) Muted() bool {
	return e.muted
}

func (e *EuclidEmitter) Trig(key music.Key, scale music.Scale, inDir common.Direction, pulse uint64) {
	if !e.armed {
		return
	}
	if !e.muted {
		e.note.Play(key, scale)
	}
	e.triggered = true
	e.armed = false
}

func (e *EuclidEmitter) Emit(pulse uint64) []common.Direction {
	if !e.triggered {
		return []common.Direction{}
	}
	e.triggered = false
	return e.direction.Decompose()
}

func (e *EuclidEmitter) Tick() {
	e.patternTrigger()
	e.note.Tick()
	e.pulse++
}

func (e *EuclidEmitter) patternTrigger() {
	if e.pulse%uint64(common.PulsesPerStep) != 0 {
		return
	}
	pattern := generateEuclideanPattern(e.Steps, e.Triggers)
	adjusetedStep := (e.step + e.Offset) % e.Steps
	if e.hasTrigged {
		e.hasTrigged = false
	}
	if pattern[adjusetedStep] {
		e.hasTrigged = true
		e.armed = true
	}
	e.step = (e.step + 1) % e.Steps
}

func generateEuclideanPattern(steps, triggers int) []bool {
	pattern := make([]bool, steps)
	bucket := 0
	for i := 0; i < steps; i++ {
		bucket += triggers
		if bucket >= steps {
			bucket -= steps
			pattern[i] = true
		} else {
			pattern[i] = false
		}
	}
	return pattern
}

func (e *EuclidEmitter) Direction() common.Direction {
	return e.direction
}

func (e *EuclidEmitter) SetDirection(dir common.Direction) {
	if e.direction.Contains(dir) {
		e.direction = e.direction.Remove(dir)
		return
	}
	e.direction = e.direction.Add(dir)
}

func (e *EuclidEmitter) Symbol() string {
	return fmt.Sprintf("%s%s", "E", e.direction.Symbol())
}

func (e *EuclidEmitter) Name() string {
	return "euclid"
}

func (e *EuclidEmitter) Color() string {
	return "39"
}

func (e *EuclidEmitter) Reset() {
	e.pulse = 0
	e.triggered = false
	e.armed = false
	e.step = 0
	e.Note().Stop()
}
