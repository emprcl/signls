package core

// EmitterBehavior defines the behavior of different types of emitters.
// This interface is implemented by different emitter types (e.g., BangEmitter).
type EmitterBehavior interface {
	// ArmedOnStart indicates if the emitter is armed when the grid starts.
	ArmedOnStart() bool

	// EmitDirections determines which directions the emitter will emit signals
	// based on its current direction and the pulse count.
	EmitDirections(dir Direction, pulse uint64) Direction

	// Symbol returns a string representation of the emitter, potentially
	// taking its direction into account for visualization.
	Symbol(dir Direction) string

	// Name returns the name of the emitter type.
	Name() string

	// Color returns the color code associated with the emitter.
	Color() string
}

// Emitter represents a node that emits signals when triggered. It contains
// information about its behavior, direction, associated note, and state.
type Emitter struct {
	behavior EmitterBehavior // The specific behavior of this emitter.

	direction Direction // The direction(s) in which this emitter is facing.
	note      *Note     // The musical note associated with this emitter.

	pulse     uint64 // The last pulse when the emitter was triggered.
	armed     bool   // Whether the emitter is armed and ready to trigger.
	triggered bool   // Whether the emitter has been triggered.
	muted     bool   // Whether the emitter is muted.
}

// Copy creates a deep copy of the emitter, returning it as a Node interface.
// It clones the associated note and keeps the same behavior and direction.
func (e *Emitter) Copy() Node {
	newNote := *e.note // Deep copy the note to maintain state separately.
	return &Emitter{
		behavior:  e.behavior,
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
func (e *Emitter) Note() *Note {
	return e.note
}

// Arm sets the emitter to an armed state, meaning it is ready to trigger.
func (e *Emitter) Arm() {
	e.armed = true
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
func (e *Emitter) Trig(key Key, scale Scale, pulse uint64) {
	if !e.updated(pulse) {
		e.note.Tick() // Move the note's internal clock forward.
	}
	if !e.armed {
		return // If not armed, don't do anything.
	}
	if !e.muted {
		e.note.Play(key, scale) // Play the note if not muted.
	}
	e.triggered = true // Mark as triggered.
	e.armed = false    // Disarm the emitter.
	e.pulse = pulse    // Update the pulse to the current one.
}

// Emit handles the signal emission process. If the emitter was triggered on the current pulse,
// it calculates the directions to emit signals and updates the grid accordingly.
func (e *Emitter) Emit(g *Grid, x, y int) {
	if e.updated(g.pulse) || !e.triggered {
		return // If the emitter hasn't been triggered, do nothing.
	}
	directions := e.behavior.EmitDirections(e.direction, g.pulse) // Get emission directions.
	for _, dir := range directions.Decompose() {
		g.Emit(x, y, dir) // Emit signals in each direction.
	}
	e.triggered = false // Reset the triggered state after emitting.
	e.pulse = g.pulse   // Update the pulse.
}

// Direction returns the current direction(s) the emitter is facing.
func (e *Emitter) Direction() Direction {
	return e.direction
}

// SetDirection adds or removes a direction from the emitter's current direction(s).
// If the direction is already set, it removes it; otherwise, it adds it.
func (e *Emitter) SetDirection(dir Direction) {
	if e.direction.Contains(dir) {
		e.direction = e.direction.Remove(dir) // Remove the direction if it exists.
		return
	}
	e.direction = e.direction.Add(dir) // Add the direction if it doesn't exist.
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
	e.pulse = 0                         // Reset pulse counter.
	e.armed = e.behavior.ArmedOnStart() // Arm or disarm based on behavior.
	e.triggered = false                 // Reset triggered state.
	e.Note().Stop()                     // Stop any playing note.
}

// updated checks if the emitter was updated on the given pulse.
func (e *Emitter) updated(pulse uint64) bool {
	return e.pulse == pulse
}
