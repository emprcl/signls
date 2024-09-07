package node

import (
	"cykl/core/common"
	"cykl/core/music"
	"cykl/midi"
	"fmt"
)

type EuclidEmitter struct {
	direction common.Direction
	note      *music.Note

	pulse     uint64
	armed     bool
	triggered bool
	muted     bool
}

func NewEuclidEmitter(midi midi.Midi, direction common.Direction) *EuclidEmitter {
	return &EuclidEmitter{
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
	if !e.updated(pulse) {
		e.note.Tick()
	}
	if !e.armed {
		return
	}
	if !e.muted {
		e.note.Play(key, scale)
	}
	e.triggered = true
	e.armed = false
	e.pulse = pulse
}

func (e *EuclidEmitter) Emit(pulse uint64) []common.Direction {
	if e.updated(pulse) || !e.triggered {
		return []common.Direction{}
	}
	e.triggered = false
	e.pulse = pulse
	return e.direction.Decompose()
}

func (e *EuclidEmitter) Tick() {
	e.note.Tick()
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
	e.Note().Stop()
}

func (e *EuclidEmitter) updated(pulse uint64) bool {
	return e.pulse == pulse
}
