package node

import (
	"cykl/core/common"
	"cykl/core/music"
	"cykl/midi"
	"fmt"
)

const (
	defaultSteps = 16
	minSteps     = 1
	maxSteps     = 128

	defaultTriggers = 4
	defaultOffset   = 0
)

type EuclidEmitter struct {
	direction common.Direction
	note      *music.Note

	Steps    *common.ControlValue[int]
	Triggers *common.ControlValue[int]
	Offset   *common.ControlValue[int]

	step int

	pulse     uint64
	ticks     uint64
	armed     bool
	triggered bool
	retrig    bool
	muted     bool
}

func NewEuclidEmitter(midi midi.Midi, direction common.Direction) *EuclidEmitter {
	return &EuclidEmitter{
		Steps:     common.NewControlValue[int](defaultSteps, minSteps, maxSteps),
		Triggers:  common.NewControlValue[int](defaultTriggers, minSteps, maxSteps),
		Offset:    common.NewControlValue[int](defaultOffset, defaultOffset, maxSteps),
		direction: direction,
		armed:     true,
		note:      music.NewNote(midi),
	}
}

func (e *EuclidEmitter) Copy(dx, dy int) common.Node {
	newSteps := *e.Steps
	newTriggers := *e.Triggers
	newOffset := *e.Offset
	newNote := e.note.Copy()
	return &EuclidEmitter{
		direction: e.direction,
		armed:     e.armed,
		note:      newNote,
		muted:     e.muted,
		Steps:     &newSteps,
		Triggers:  &newTriggers,
		Offset:    &newOffset,
	}
}

func (e *EuclidEmitter) Activated() bool {
	return e.armed || e.triggered
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
	if e.triggered {
		e.retrig = true
	} else {
		e.pulse = pulse
	}
	e.triggered = true
	e.armed = false
}

func (e *EuclidEmitter) Emit(pulse uint64) []common.Direction {
	if e.updated(pulse) || !e.triggered {
		return []common.Direction{}
	}
	if e.retrig {
		e.retrig = false
	} else {
		e.triggered = false
	}
	e.pulse = pulse
	return e.direction.Decompose()
}

func (e *EuclidEmitter) Tick() {
	e.patternTrigger()
	e.note.Tick()
	e.ticks++
}

func (e *EuclidEmitter) patternTrigger() {
	if e.ticks%uint64(common.PulsesPerStep) != 0 {
		return
	}

	if e.ticks%uint64(common.PulsesPerStep*e.Steps.Value()) == 0 {
		newSteps := e.Steps.Computed()
		e.Triggers.SetMax(newSteps)
		e.Offset.SetMax(newSteps)
		e.Triggers.Computed()
		e.Offset.Computed()
	}

	pattern := generateEuclideanPattern(e.Steps.Last(), e.Triggers.Last())
	adjusetedStep := (e.step + e.Offset.Last()) % e.Steps.Last()
	if pattern[adjusetedStep] {
		e.armed = true
	}
	e.step = (e.step + 1) % e.Steps.Last()
}

func generateEuclideanPattern(steps, triggers int) []bool {
	pattern := make([]bool, steps)
	bucket := 0
	pattern[0] = true
	for i := 1; i < steps; i++ {
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
	return fmt.Sprintf("%s%s%s", "E", e.note.Key.Symbol(), e.direction.Symbol())
}

func (e *EuclidEmitter) Name() string {
	return "euclid"
}

func (e *EuclidEmitter) Color() string {
	return "162"
}

func (e *EuclidEmitter) Reset() {
	e.pulse = 0
	e.ticks = 0
	e.triggered = false
	e.armed = true
	e.retrig = false
	e.step = 0
	e.Note().Stop()
}

// updated checks if the emitter was updated on the given pulse.
func (e *EuclidEmitter) updated(pulse uint64) bool {
	return e.pulse == pulse
}
