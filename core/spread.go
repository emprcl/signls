package core

import (
	"cykl/midi"
	"fmt"
)

type SpreadEmitter struct {
	note      *Note
	direction Direction
	pulse     uint64
	armed     bool
	triggered bool
	muted     bool
}

func NewSimpleEmitter(midi midi.Midi, direction Direction) *SpreadEmitter {
	return &SpreadEmitter{
		direction: direction,
		note:      NewNote(midi),
	}
}

func (e *SpreadEmitter) Copy() Node {
	return &SpreadEmitter{
		direction: e.direction,
		note:      e.note,
	}
}

func (e *SpreadEmitter) Activated() bool {
	return e.armed || e.triggered
}

func (e *SpreadEmitter) Note() *Note {
	return e.note
}

func (e *SpreadEmitter) Arm() {
	e.armed = true
}

func (e *SpreadEmitter) SetMute(mute bool) {
	e.muted = mute
}

func (e *SpreadEmitter) Muted() bool {
	return e.muted
}

func (e *SpreadEmitter) Trig(pulse uint64) {
	if !e.armed {
		return
	}
	if !e.muted {
		e.note.Play()
	}
	e.triggered = true
	e.armed = false
	e.pulse = pulse
	// TODO: handle length and note off
}

func (e *SpreadEmitter) Emit(g *Grid, x, y int) {
	if e.updated(g.Pulse) || !e.triggered {
		return
	}
	for _, dir := range e.direction.Decompose() {
		g.Emit(x, y, dir)
	}
	e.triggered = false
	e.pulse = g.Pulse
}

func (e *SpreadEmitter) Direction() Direction {
	return e.direction
}

func (e *SpreadEmitter) SetDirection(dir Direction) {
	if e.direction.Contains(dir) {
		e.direction = e.direction.Remove(dir)
		return
	}
	e.direction = e.direction.Add(dir)
}

func (e *SpreadEmitter) Symbol() string {
	return fmt.Sprintf("%s%s", "S", e.Direction().Symbol())
}

func (e *SpreadEmitter) Name() string {
	return "spread"
}

func (e *SpreadEmitter) Color() string {
	return "177"
}

func (e *SpreadEmitter) Reset() {
	e.pulse = 0
	e.armed = false
	e.triggered = false
}

func (e *SpreadEmitter) updated(pulse uint64) bool {
	return e.pulse == pulse
}
