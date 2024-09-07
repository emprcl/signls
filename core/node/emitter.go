package node

import (
	"cykl/core/common"
	"cykl/core/music"
)

// EmitterBehavior defines the behavior of different types of emitters.
type EmitterBehavior interface {
	// ArmedOnStart indicates if the emitter is armed when the grid starts.
	ArmedOnStart() bool

	// Copy makes a copy of the behavior.
	Copy() EmitterBehavior

	// EmitDirections determines which directions the emitter will emit signals
	// based on its current direction, the incoming direction, and the pulse count.
	EmitDirections(dir common.Direction, inDir common.Direction, pulse uint64) common.Direction

	// ShouldPropagate indicates if triggers should be propagated to direct
	// neighbors.
	ShouldPropagate() bool

	// Symbol returns a string representation of the emitter, potentially
	// taking its direction into account for visualization.
	Symbol(dir common.Direction) string

	// Name returns the name of the emitter type.
	Name() string

	// Color returns the color code associated with the emitter.
	Color() string
}

// Emitter represents a node that emits signals when triggered. It contains
// information about its behavior, direction, associated note, and state.
type Emitter struct {
	behavior EmitterBehavior // The specific behavior of this emitter.

	direction         common.Direction // The signal emitting direction(s)
	incomingDirection common.Direction // The direction from where the emitter has been triggered
	note              *music.Note      // The musical note associated with this emitter.

	pulse     uint64 // The last pulse when the emitter was triggered.
	armed     bool   // Whether the emitter is armed and ready to trigger.
	triggered bool   // Whether the emitter has been triggered.
	muted     bool   // Whether the emitter is muted.
}

// Copy creates a deep copy of the emitter, returning it as a Node interface.
// It clones the associated note and keeps the same behavior and direction.
func (e *Emitter) Copy(dx, dy int) common.Node {
	newNote := *e.note
	return &Emitter{
		behavior:  e.behavior.Copy(),
		direction: e.direction,
		armed:     e.armed,
		note:      &newNote,
	}
}

// Activated checks if the emitter is currently active, meaning it's either
// armed or has been triggered.
func (e *Emitter) Activated() bool {
	return e.armed || e.triggered
}

// Note returns the pointer to the Note associated with the emitter.
func (e *Emitter) Note() *music.Note {
	return e.note
}

// Arm sets the emitter to an armed state, meaning it is ready to trigger.
func (e *Emitter) Arm() {
	e.armed = true
}

func (e *Emitter) Behavior() EmitterBehavior {
	return e.behavior
}

func (e *Emitter) SetBehavior(behavior EmitterBehavior) {
	e.behavior = behavior
}

// SetMute mutes or unmutes the emitter. If muted, it stops any currently
// playing note.
func (e *Emitter) SetMute(mute bool) {
	e.note.Stop() // Stop the note if we're muting.
	e.muted = mute
}

// Muted returns whether the emitter is currently muted.
func (e *Emitter) Muted() bool {
	return e.muted
}

// Trig triggers the emitter, playing its note if it is armed and not muted.
// It also updates the pulse to the current one, and disarms the emitter.
func (e *Emitter) Trig(key music.Key, scale music.Scale, inDir common.Direction, pulse uint64) {
	if !e.updated(pulse) {
		e.note.Tick() // Move the note's internal clock forward.
	}
	if !e.armed {
		return
	}
	if !e.muted {
		e.note.Play(key, scale)
	}
	e.incomingDirection = inDir
	e.triggered = true
	e.armed = false
	e.pulse = pulse
}

// Emit returns the directions to emits for a given pulse.
func (e *Emitter) Emit(pulse uint64) []common.Direction {
	if e.updated(pulse) || !e.triggered {
		return []common.Direction{}
	}
	e.triggered = false
	e.pulse = pulse
	return e.behavior.EmitDirections(e.direction, e.incomingDirection, pulse).Decompose()
}

func (e *Emitter) Tick() {
	e.note.Tick()
}

// Direction returns the current direction(s) the emitter is facing.
func (e *Emitter) Direction() common.Direction {
	return e.direction
}

// SetDirection adds or removes a direction from the emitter's current direction(s).
// If the direction is already set, it removes it; otherwise, it adds it.
func (e *Emitter) SetDirection(dir common.Direction) {
	if e.direction.Contains(dir) {
		e.direction = e.direction.Remove(dir)
		return
	}
	e.direction = e.direction.Add(dir)
}

// Symbol returns a string representation of the emitter, typically used for visualization.
func (e *Emitter) Symbol() string {
	return e.behavior.Symbol(e.direction)
}

// Name returns the name of the emitter type.
func (e *Emitter) Name() string {
	return e.behavior.Name()
}

// Color returns the color associated with the emitter.
func (e *Emitter) Color() string {
	return e.behavior.Color()
}

// Reset restores the emitter to its initial state, resetting the pulse count,
// disarming the emitter, and stopping any playing notes.
func (e *Emitter) Reset() {
	e.pulse = 0
	e.armed = e.behavior.ArmedOnStart()
	e.triggered = false
	e.Note().Stop()
}

// updated checks if the emitter was updated on the given pulse.
func (e *Emitter) updated(pulse uint64) bool {
	return e.pulse == pulse
}
