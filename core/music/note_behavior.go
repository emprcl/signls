package music

import "math/rand"

// NoteBehavior defines an interface for various note behaviors.
// Implementations of this interface define how a note is played.
type NoteBehavior interface {
	Name() string
	Play(n *Note, root Key, scale Scale)
	Symbol() string
}

// AllNoteBehaviors returns a slice containing all available note behaviors.
// This is useful for iterating through or selecting different behaviors.
func AllNoteBehaviors() []NoteBehavior {
	return []NoteBehavior{
		FixedNote{},
		SilentNote{},
		RandomNote{},
	}
}

// SilentNote is a behavior that does not play any sound.
// It implements the NoteBehavior interface.
type SilentNote struct{}

// Play is the implementation of the SilentNote behavior, which does nothing.
func (b SilentNote) Play(n *Note, root Key, scale Scale) {}

// Name returns the name of the SilentNote behavior.
func (b SilentNote) Name() string {
	return "silent"
}

func (b SilentNote) Symbol() string {
	return "̥"
}

// FixedNote is a behavior where the note plays at a fixed pitch.
// It implements the NoteBehavior interface.
type FixedNote struct{}

// Play is the implementation of the FixedNote behavior, which plays the note
// at its set key. If there is a transposition, it stops the note and applies the new key.
func (b FixedNote) Play(n *Note, root Key, scale Scale) {
	if n.nextKey > 0 {
		n.Stop()
		n.Key = n.nextKey
	}
	n.Transpose(root, scale)
	n.midi.NoteOn(n.Channel, uint8(n.Key), n.Velocity)
}

// Name returns the name of the FixedNote behavior.
func (b FixedNote) Name() string {
	return "fixed"
}

func (b FixedNote) Symbol() string {
	return ""
}

// RandomNote is a behavior where the note plays at a random pitch.
// It implements the NoteBehavior interface.
type RandomNote struct{}

// Play is the implementation of the RandomNote behavior, which selects
// a random key within an octave and transposes it according to the scale and root key.
func (b RandomNote) Play(n *Note, root Key, scale Scale) {
	n.Key = root + Key(rand.Intn(12))
	interval := n.Key.AllSemitonesFrom(root)
	n.Key = n.Key.Transpose(root, scale, interval)
	n.midi.NoteOn(n.Channel, uint8(n.Key), n.Velocity)
}

// Name returns the name of the RandomNote behavior.
func (b RandomNote) Name() string {
	return "random"
}

func (b RandomNote) Symbol() string {
	return "̰"
}
