package core

import (
	"cykl/midi"
)

type noteBehavior uint8

const (
	defaultKey      uint8 = 60
	defaultChannel  uint8 = 0
	defaultVelocity uint8 = 100
	defaultLength   uint8 = uint8(pulsesPerStep)

	maxKey      uint8 = 127
	maxVelocity uint8 = 127
	maxLength   uint8 = 127 // 127 is infinity
	minLength   uint8 = 1
	maxChannel  uint8 = 15

	silence noteBehavior = iota
	fixed
)

type Note struct {
	midi     midi.Midi
	behavior noteBehavior // TODO: implement
	Channel  uint8
	Key      uint8
	Velocity uint8
	Length   uint8

	nextKey   uint8
	pulse     uint64
	triggered bool
}

func NewNote(midi midi.Midi) *Note {
	return &Note{
		midi:     midi,
		Channel:  defaultChannel,
		Key:      defaultKey,
		Velocity: defaultVelocity,
		Length:   defaultLength,
	}
}

func (n Note) IsValid() bool {
	return n.Key == 0
}

func (n Note) KeyName() string {
	if n.nextKey > 0 {
		return midi.Note(n.nextKey)
	}
	return midi.Note(n.Key)
}

func (n Note) KeyValue() uint8 {
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

func (n *Note) Play() {
	if n.nextKey > 0 {
		n.Stop()
		n.Key = n.nextKey
	}
	n.midi.NoteOn(n.Channel, n.Key, n.Velocity)
	n.triggered = true
	n.pulse = 0
	n.nextKey = 0
}

func (n *Note) Stop() {
	n.midi.NoteOff(n.Channel, n.Key)
	n.triggered = false
	n.pulse = 0
}

func (n *Note) SetKey(key uint8) {
	if key > maxKey {
		n.nextKey = 0
		return
	}
	if key < 0 {
		n.nextKey = 127
	}
	n.nextKey = key
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
