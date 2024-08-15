package core

import (
	"cykl/midi"
	"fmt"
)

type BangEmitter struct {
	note      *Note
	direction Direction
	pulse     uint64
	armed     bool
	triggered bool
	muted     bool
}

func NewInitEmitter(midi midi.Midi, direction Direction) *BangEmitter {
	return &BangEmitter{
		direction: direction,
		armed:     true,
		note:      NewNote(midi),
	}
}

func (e *BangEmitter) Copy() Node {
	return &BangEmitter{
		direction: e.direction,
		armed:     true,
		note:      e.note,
	}
}

func (e *BangEmitter) Activated() bool {
	return e.armed || e.triggered
}

func (e *BangEmitter) Note() *Note {
	return e.note
}

func (e *BangEmitter) Arm() {
	e.armed = true
}

func (e *BangEmitter) SetMute(mute bool) {
	e.muted = mute
}

func (e *BangEmitter) Muted() bool {
	return e.muted
}

func (e *BangEmitter) Trig(pulse uint64) {
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

func (e *BangEmitter) Emit(g *Grid, x, y int) {
	if e.updated(g.Pulse) || !e.triggered {
		return
	}
	for _, dir := range e.direction.Decompose() {
		g.Emit(x, y, dir)
	}
	e.triggered = false
	e.pulse = g.Pulse
}

func (e *BangEmitter) Direction() Direction {
	return e.direction
}

func (e *BangEmitter) SetDirection(dir Direction) {
	if e.direction.Contains(dir) {
		e.direction = e.direction.Remove(dir)
		return
	}
	e.direction = e.direction.Add(dir)
}

func (e *BangEmitter) Symbol() string {
	return fmt.Sprintf("%s%s", "B", e.Direction().Symbol())
}

func (e *BangEmitter) Name() string {
	return "bang"
}

func (e *BangEmitter) Color() string {
	return "165"
}

func (e *BangEmitter) Reset() {
	e.pulse = 0
	e.armed = true
	e.triggered = false
}

func (e *BangEmitter) updated(pulse uint64) bool {
	return e.pulse == pulse
}
