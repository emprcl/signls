package music

import (
	"cykl/core/common"
	"cykl/midi"
	"math/rand"
	"time"
)

// Constants defining default values for note properties and their limits.
const (
	defaultKey      Key   = 60 // Default MIDI key value (Middle C).
	defaultChannel  uint8 = 0
	defaultVelocity uint8 = 100
	defaultLength   uint8 = uint8(common.PulsesPerStep)

	maxKey         Key   = 127
	maxVelocity    uint8 = 127
	maxLength      uint8 = 127 // Maximum length of the note (127 is treated as infinity).
	minLength      uint8 = 1
	maxChannel     uint8 = 15
	maxProbability uint8 = 100
)

// Note represents a musical note with various properties such as key, velocity, and length.
// It includes behavior that defines how the note plays and interacts with other elements.
type Note struct {
	midi midi.Midi

	rand *rand.Rand

	Key         KeyBehavior
	Channel     *MidiParam
	Velocity    *MidiParam
	Length      *MidiParam
	Probability uint8

	pulse     uint64 // Internal pulse counter to manage note length.
	triggered bool   // Indicates whether the note is currently playing.
}

// NewNote initializes a new Note with default settings and the provided MIDI interface.
func NewNote(midi midi.Midi) *Note {
	source := rand.NewSource(time.Now().UnixNano())
	return &Note{
		midi:        midi,
		rand:        rand.New(source),
		Key:         &FixedKey{key: defaultKey},
		Channel:     NewMidiParam(defaultChannel, 0, maxChannel),
		Velocity:    NewMidiParam(defaultVelocity, 0, maxVelocity),
		Length:      NewMidiParam(defaultLength, minLength, maxLength),
		Probability: maxProbability,
	}
}

// Copy creates a copy of the note.
func (n Note) Copy() *Note {
	newChannel := *n.Channel
	newVelocity := *n.Velocity
	newLength := *n.Length
	source := rand.NewSource(time.Now().UnixNano())
	return &Note{
		midi:        n.midi,
		rand:        rand.New(source),
		Key:         n.Key,
		Channel:     &newChannel,
		Velocity:    &newVelocity,
		Length:      &newLength,
		Probability: n.Probability,
	}
}

// KeyName returns the MIDI note name based on the current or next key.
func (n Note) KeyName() string {
	return midi.Note(uint8(n.Key.Value()))
}

// KeyValue returns the MIDI key value of the current or next key.
func (n Note) KeyValue() Key {
	return n.Key.Value()
}

// Tick advances the internal pulse counter and stops the note if it exceeds its length.
func (n *Note) Tick() {
	if !n.triggered {
		return
	}
	n.pulse++

	// Stop the note if its duration is complete.
	if n.Length.Last() < maxLength && n.pulse >= uint64(n.Length.Last()) {
		n.Stop()
	}
}

// Play triggers the note with a specific root and scale, resetting internal state.
func (n *Note) Play(root Key, scale Scale) {
	if n.Probability == maxProbability ||
		uint8(rand.Int31n((100))) < n.Probability {
		if n.Key.IsChanging() {
			n.Stop()
		}
		n.Key.Transpose(root, scale)
		n.midi.NoteOn(n.Channel.Computed(), uint8(n.Key.Computed()), n.Velocity.Computed())
	}
	n.triggered = true
	n.pulse = 0
}

// Stop sends a MIDI Note Off message and resets the triggered state.
func (n *Note) Stop() {
	n.midi.NoteOff(n.Channel.Last(), uint8(n.Key.Value()))
	n.triggered = false
	n.pulse = 0
}

// SetKey updates the note's key and calculates the interval relative to the root.
func (n *Note) SetKey(key Key, root Key) {
	n.Key.SetNext(key, root)
	if !n.triggered {
		n.Key.Set(n.Key.Value())
	}
}

// SetVelocity updates the velocity of the note, ensuring it is within valid limits.
func (n *Note) SetVelocity(velocity uint8) {
	n.Velocity.Set(velocity)
}

// SetLength updates the length of the note, ensuring it is within valid limits.
func (n *Note) SetLength(length uint8) {
	n.Length.Set(length)
}

// SetChannel updates the MIDI channel of the note, ensuring it is within valid limits.
func (n *Note) SetChannel(channel uint8) {
	n.Channel.Set(channel)
}

// ClockDivision returns the pulses per step and steps per quarter note,
// which might be used for timing or synchronization purposes.
func (n *Note) ClockDivision() (int, int) {
	return common.PulsesPerStep, common.StepsPerQuarterNote
}
