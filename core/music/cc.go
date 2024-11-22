package music

import (
	"math/rand"
	"signls/core/common"
	"signls/midi"
	"time"
)

type ControlType int

const (
	NONEType ControlType = iota
	CCType
	RPNType
	NRPNType
)

type CC struct {
	midi midi.Midi

	rand *rand.Rand

	Type       ControlType
	Controller uint8
	Value      common.ControlValue[uint8]
}

func NewCC(midi midi.Midi, controlType ControlType) *CC {
	source := rand.NewSource(time.Now().UnixNano())
	return &CC{
		midi: midi,
		rand: rand.New(source),
		Type: controlType,
	}
}

func (c CC) Send(channel uint8) {
	switch c.Type {
	case NONEType:
		return
	case CCType:
		c.midi.ControlChange(channel, c.Controller, c.Value.Computed())
	}
}
