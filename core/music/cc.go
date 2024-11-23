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
)

type ControlType int

const (
	NONEControlType ControlType = iota
	CCControlType
	ATControlType // After Touch
)

var (
	AllControlTypes = []ControlType{
		NONEControlType,
		CCControlType,
		ATControlType,
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

func (c *CC) SetController(controller uint8) {
	if controller < minController || controller > maxController {
		return
	}
	c.Controller = controller
}

func (c CC) Send(channel uint8) {
	switch c.Type {
	case NONEControlType:
		return
	case CCControlType:
		c.midi.ControlChange(channel, c.Controller, uint8(c.Value.Computed()))
	case ATControlType:
		c.midi.AfterTouch(channel, uint8(c.Value.Computed()))
	}
}
