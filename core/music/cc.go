package music

import (
	"signls/core/common"
	"signls/midi"
)

const (
	defaultController uint8 = 0
	minController     uint8 = 0
	maxController     uint8 = 119

	defaultControlValue uint8 = 0
	minControlValue     uint8 = 0
	maxControlValue     uint8 = 127

	defaultPitchBendValue = 64
	minPitchBendValue     = -128
	maxPitchBendValue     = 128
)

type ControlType int

const (
	SilentControlType ControlType = iota
	ControlChangeControlType
	AfterTouchControlType
	PitchBendControlType
	ProgramChangeControlType
)

var (
	AllControlTypes = []ControlType{
		SilentControlType,
		ControlChangeControlType,
		AfterTouchControlType,
		PitchBendControlType,
		ProgramChangeControlType,
	}
)

type CC struct {
	midi midi.Midi

	Type       ControlType
	Controller uint8
	Value      *common.ControlValue[uint8]
}

func NewCC(midi midi.Midi, controlType ControlType) *CC {
	return &CC{
		midi:       midi,
		Type:       controlType,
		Controller: defaultController,
		Value:      common.NewControlValue[uint8](defaultControlValue, minControlValue, maxControlValue),
	}
}

func (c CC) Copy() *CC {
	newValue := *c.Value
	return &CC{
		midi:       c.midi,
		Type:       c.Type,
		Controller: c.Controller,
		Value:      &newValue,
	}
}

func (c *CC) SetController(controller uint8) {
	if controller < minController || controller > maxController {
		return
	}
	c.Controller = controller
}

func (c *CC) SetType(t int) {
	c.Type = ControlType(t)
	if c.Type == PitchBendControlType {

	}
}

func (c CC) Send(channel uint8) {
	switch c.Type {
	case SilentControlType:
		return
	case ControlChangeControlType:
		c.midi.ControlChange(channel, c.Controller, uint8(c.Value.Computed()))
	case AfterTouchControlType:
		c.midi.AfterTouch(channel, uint8(c.Value.Computed()))
	case PitchBendControlType:
		c.midi.Pitchbend(
			channel,
			int16(remap(
				int(c.Value.Computed()),
				int(minControlValue),
				int(maxControlValue),
				minPitchBendValue,
				maxPitchBendValue,
			)),
		)
	}
}

func remap(value, oldMin, oldMax, newMin, newMax int) int {
	if value < oldMin {
		value = oldMin
	} else if value > oldMax {
		value = oldMax
	}

	return newMin + (value-oldMin)*(newMax-newMin)/(oldMax-oldMin)
}
