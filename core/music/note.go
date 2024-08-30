package music

import (
	"cykl/core/common"
	"cykl/midi"
)

// Constants defining default values for note properties and their limits.
const (
	defaultKey      Key   = 60 // Default MIDI key value (Middle C).
	defaultChannel  uint8 = 0
	defaultVelocity uint8 = 100
	defaultLength   uint8 = uint8(common.PulsesPerStep)

	maxKey      Key   = 127
	maxVelocity uint8 = 127
	maxLength   uint8 = 127 // Maximum length of the note (127 is treated as infinity).
	minLength   uint8 = 1
	maxChannel  uint8 = 15
)

// Note represents a musical note with various properties such as key, velocity, and length.
// It includes behavior that defines how the note plays and interacts with other elements.
type Note struct {
	Behavior NoteBehavior // Interface for defining different behaviors for a note.

	midi midi.Midi // Interface to interact with MIDI devices.

	Key      Key   // The MIDI key (pitch) of the note.
	Interval int   // Interval relative to a scale's root note.
	Channel  uint8 // The MIDI channel on which the note is played.
	Velocity uint8 // The velocity (volume) of the note.
	Length   uint8 // The duration the note is held for.

	nextKey   Key    // Next key to be played, used for transposition.
	pulse     uint64 // Internal pulse counter to manage note length.
	triggered bool   // Indicates whether the note is currently playing.
}

// NewNote initializes a new Note with default settings and the provided MIDI interface.
func NewNote(midi midi.Midi) *Note {
	return &Note{
		Behavior: FixedNote{},
		midi:     midi,
		Channel:  defaultChannel,
		Key:      defaultKey,
		Velocity: defaultVelocity,
		Length:   defaultLength,
	}
}

// KeyName returns the MIDI note name based on the current or next key.
func (n Note) KeyName() string {
	if n.nextKey > 0 {
		return midi.Note(uint8(n.nextKey))
	}
	return midi.Note(uint8(n.Key))
}

// KeyValue returns the MIDI key value of the current or next key.
func (n Note) KeyValue() Key {
	if n.nextKey > 0 {
		return n.nextKey
	}
	return n.Key
}

// Tick advances the internal pulse counter and stops the note if it exceeds its length.
func (n *Note) Tick() {
	if !n.triggered {
		return
	}
	n.pulse++

	// Stop the note if its duration is complete.
	if n.Length < maxLength && n.pulse >= uint64(n.Length) {
		n.Stop()
	}
}

// Transpose changes the note's key according to a scale and root note.
func (n *Note) Transpose(root Key, scale Scale) {
	n.SetKey(n.Key.Transpose(root, scale, n.Interval), root)
}

// Play triggers the note with a specific key and scale, resetting internal state.
func (n *Note) Play(key Key, scale Scale) {
	n.Behavior.Play(n, key, scale)
	n.triggered = true
	n.pulse = 0
	n.nextKey = 0
}

// Stop sends a MIDI Note Off message and resets the triggered state.
func (n *Note) Stop() {
	n.midi.NoteOff(n.Channel, uint8(n.Key))
	n.triggered = false
	n.pulse = 0
}

// SetKey updates the note's key and calculates the interval relative to the root.
func (n *Note) SetKey(key Key, root Key) {
	if key > maxKey {
		n.nextKey = Key(0)
	} else {
		n.nextKey = key
	}
	n.Interval = n.nextKey.AllSemitonesFrom(root)
	if !n.triggered {
		n.Key = n.nextKey
		n.nextKey = 0
	}
}

// SetVelocity updates the velocity of the note, ensuring it is within valid limits.
func (n *Note) SetVelocity(velocity uint8) {
	if velocity > maxVelocity {
		return
	}
	n.Velocity = velocity
}

// SetLength updates the length of the note, ensuring it is within valid limits.
func (n *Note) SetLength(length uint8) {
	if length > maxLength || length < minLength {
		return
	}
	n.Length = length
}

// SetChannel updates the MIDI channel of the note, ensuring it is within valid limits.
func (n *Note) SetChannel(channel uint8) {
	if channel > maxChannel {
		return
	}
	n.Channel = channel
}

// ClockDivision returns the pulses per step and steps per quarter note,
// which might be used for timing or synchronization purposes.
func (n *Note) ClockDivision() (int, int) {
	return common.PulsesPerStep, common.StepsPerQuarterNote
}
