package core

import (
	"cykl/midi"
)

type noteBehavior uint8

const (
	defaultKey      uint8 = 60
	defaultChannel  uint8 = 0
	defaultVelocity uint8 = 100
	defaultLength   uint8 = 1

	maxKey uint8 = 127

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

func (n Note) Name() string {
	return midi.Note(n.Key)
}

func (n Note) Play() {
	n.midi.NoteOn(n.Channel, n.Key, n.Velocity)
}

func (n *Note) SetKey(key uint8) {
	if n.Key > maxKey {
		n.Key = 0
		return
	}
	n.Key = key
}
