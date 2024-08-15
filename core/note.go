package core

import (
	"cykl/midi"
)

type NoteBehavior uint8

const (
	Silence NoteBehavior = iota
	Fixed
)

type Note struct {
	Behavior NoteBehavior
	Channel  uint8
	Key      uint8
	Velocity uint8
	Length   uint8
}

func (n Note) IsValid() bool {
	return n.Key == 0
}

func (n Note) Name() string {
	return midi.Note(n.Key)
}
