package core

import (
	"cykl/midi"
)

const (
	defaultKey      Key   = 60
	defaultChannel  uint8 = 0
	defaultVelocity uint8 = 100
	defaultLength   uint8 = uint8(pulsesPerStep)

	maxKey      Key   = 127
	maxVelocity uint8 = 127
	maxLength   uint8 = 127 // 127 is infinity
	minLength   uint8 = 1
	maxChannel  uint8 = 15
)

type Note struct {
	Behavior NoteBehavior

	midi midi.Midi

	Key      Key
	Interval int
	Channel  uint8
	Velocity uint8
	Length   uint8

	nextKey   Key
	pulse     uint64
	triggered bool
}

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

func (n Note) KeyName() string {
	if n.nextKey > 0 {
		return midi.Note(uint8(n.nextKey))
	}
	return midi.Note(uint8(n.Key))
}

func (n Note) KeyValue() Key {
	if n.nextKey > 0 {
		return n.nextKey
	}
	return n.Key
}

func (n *Note) Tick() {
	if !n.triggered {
		return
	}
	n.pulse++
	if n.Length < maxLength && n.pulse >= uint64(n.Length) {
		n.Stop()
	}
}

func (n *Note) Transpose(root Key, scale Scale) {
	n.SetKey(n.Key.Transpose(root, scale, n.Interval), root)
}

func (n *Note) Play(key Key, scale Scale) {
	n.Behavior.Play(n, key, scale)
	n.triggered = true
	n.pulse = 0
	n.nextKey = 0
}

func (n *Note) Stop() {
	n.midi.NoteOff(n.Channel, uint8(n.Key))
	n.triggered = false
	n.pulse = 0
}

func (n *Note) SetKey(key Key, root Key) {
	if key > maxKey {
		n.nextKey = Key(0)
	} else if key < 0 {
		n.nextKey = maxKey
	} else {
		n.nextKey = Key(key)
	}
	n.Interval = n.nextKey.AllSemitonesFrom(root)
	if !n.triggered {
		n.Key = n.nextKey
		n.nextKey = 0
	}
}

func (n *Note) SetVelocity(velocity uint8) {
	if velocity > maxVelocity {
		return
	}
	n.Velocity = velocity
}

func (n *Note) SetLength(length uint8) {
	if length > maxLength || length < minLength {
		return
	}
	n.Length = length
}

func (n *Note) SetChannel(channel uint8) {
	if channel > maxChannel {
		return
	}
	n.Channel = channel
}

func (n *Note) ClockDivision() (int, int) {
	return pulsesPerStep, stepsPerQuarterNote
}
