package node

import (
	"fmt"
	"unicode/utf8"

	"signls/core/common"
	"signls/core/music"
	"signls/core/theory"
)

type Emitter struct {
	behavior common.EmitterBehavior

	direction         common.Direction
	incomingDirection common.Direction
	incomingDivision  int
	note              *music.Note

	pulse     uint64
	armed     bool
	triggered bool
	retrig    bool
	muted     bool
}

func (e *Emitter) Copy(dx, dy int) common.Node {
	newNote := e.note.Copy()
	return &Emitter{
		behavior:  e.behavior.Copy(),
		direction: e.direction,
		armed:     e.armed,
		note:      newNote,
		muted:     e.muted,
	}
}

func (e *Emitter) Activated() bool {
	return e.armed || e.triggered
}

func (e *Emitter) Note() *music.Note {
	return e.note
}

func (e *Emitter) Arm() {
	e.armed = true
}

func (e *Emitter) Behavior() common.EmitterBehavior {
	return e.behavior
}

func (e *Emitter) SetBehavior(behavior common.EmitterBehavior) {
	e.behavior = behavior
}

func (e *Emitter) SetMute(mute bool) {
	e.note.Stop()
	e.muted = mute
}

func (e *Emitter) Muted() bool {
	return e.muted
}

func (e *Emitter) Trig(key theory.Key, scale theory.Scale, inDir common.Direction, inDiv int, pulse uint64) {
	if !e.updated(pulse) {
		e.note.Tick()
	}
	if !e.armed {
		return
	}
	if !e.muted {
		e.note.TransposeAndPlay(key, scale)
	}
	if !e.updated(pulse) && e.triggered {
		e.retrig = true
	} else {
		e.pulse = pulse
	}
	e.incomingDirection = inDir
	e.incomingDivision = inDiv
	e.triggered = true
	e.armed = false
}

func (e *Emitter) Emit(pulse uint64) []common.Direction {
	if e.updated(pulse) || !e.triggered {
		return []common.Direction{}
	}
	if e.retrig {
		e.retrig = false
	} else {
		e.triggered = false
	}
	e.pulse = pulse
	return e.behavior.EmitDirections(e.direction, e.incomingDirection, pulse).Decompose()
}

func (e *Emitter) Tick() {
	e.note.Tick()
}

func (e *Emitter) Direction() common.Direction {
	return e.direction
}

func (e *Emitter) Division() int {
	return e.incomingDivision
}

func (e *Emitter) SetDirection(dir common.Direction) {
	if e.direction.Contains(dir) {
		e.direction = e.direction.Remove(dir)
		return
	}
	e.direction = e.direction.Add(dir)
}

func (e *Emitter) Symbol() string {
	symbol := e.behavior.Symbol()
	if utf8.RuneCountInString(symbol) >= 2 {
		return fmt.Sprintf("%.2s", symbol)
	}
	if _, ok := e.behavior.(*PassEmitter); ok {
		return fmt.Sprintf("%s%s%s", symbol, e.note.Symbol(), "⇌")
	}
	if e.note == nil {
		return fmt.Sprintf("%s%s", symbol, e.direction.Symbol())
	}
	return fmt.Sprintf("%s%s%s", symbol, e.note.Symbol(), e.direction.Symbol())
}

func (e *Emitter) Name() string {
	return e.behavior.Name()
}

func (e *Emitter) Color() string {
	return e.behavior.Color()
}

func (e *Emitter) Reset() {
	e.pulse = 0
	e.armed = e.behavior.ArmedOnStart()
	e.triggered = false
	e.Note().Stop()
	e.behavior.Reset()
}

func (e *Emitter) updated(pulse uint64) bool {
	return e.pulse == pulse
}
