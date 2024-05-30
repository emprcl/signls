package midi

import (
	"fmt"

	gomidi "gitlab.com/gomidi/midi/v2"
)

const (
	minCC      = 0
	maxCC      = 127
	minPitch   = -8192
	maxPitch   = 8192
	resetPitch = 0
)

type msgType uint8

const (
	controlChange msgType = iota
	programChange
	pitchBend // TODO: handle better parameter set
	afterTouch
)

// Controllable defines struct that can hold midi controls.
type Controllable interface {
	Device() int
	Channel() uint8
	Control(nb int) Control
	SetControl(nb int, value int16)
	IsActiveControl(control int) bool
}

// Control represents a midi control.
// It can represent 4 types of message (see msgType):
//   - Control Change (midi cc)
//   - Program Change // TODO: send all 3 program related cc messages on change
//   - Pitchbend
//   - Aftertouch
type Control struct {
	midi Midi

	// We keep a reference to the parent item to known in which device and
	// channel we need to send the message.
	parent  Controllable
	msgType msgType

	// The controller value is used only for controlChange type.
	controller uint8
	value      int16
}

// NewControls create every possible midi controls:
//   - 127 Control Changes
//   - 1 Program Change
//   - 1 Pichbend
//   - 1 Aftertouch
func NewControls(midi Midi, parent Controllable) []Control {
	controls := []Control{
		{
			midi:    midi,
			parent:  parent,
			msgType: programChange,
		},
		{
			midi:    midi,
			parent:  parent,
			msgType: pitchBend,
		},
		{
			midi:    midi,
			parent:  parent,
			msgType: afterTouch,
		},
	}

	for i := 0; i <= 127; i++ {
		controls = append(controls, Control{
			midi:       midi,
			parent:     parent,
			msgType:    controlChange,
			controller: uint8(i),
		})
	}
	return controls
}

// Value returns the control value.
func (c Control) Value() int16 {
	return c.value
}

// String returns the string representation of the control value.
func (c Control) String() string {
	return fmt.Sprintf("%d", c.value)
}

// Name returns the name of the control.
func (c Control) Name() string {
	switch c.msgType {
	case programChange:
		return "Program"
	case pitchBend:
		return "Pitchbend"
	case afterTouch:
		return "After Touch"
	default:
		name := gomidi.ControlChangeName[c.controller]
		return fmt.Sprintf("CC%d %s", c.controller, name)
	}
}

// Set sets the control value.
func (c *Control) Set(value int16) {
	switch c.msgType {
	case pitchBend:
		if value < minPitch || value > maxPitch {
			return
		}
	default:
		if value < minCC || value > maxCC {
			return
		}
	}
	c.value = value
}

// Send sends the actual midi messages.
func (c Control) Send() {
	switch c.msgType {
	case controlChange:
		c.midi.ControlChange(
			c.parent.Device(),
			c.parent.Channel(),
			c.controller,
			uint8(c.value),
		)
	case programChange:
		c.midi.ProgramChange(
			c.parent.Device(),
			c.parent.Channel(),
			uint8(c.value),
		)
	case pitchBend:
		c.midi.Pitchbend(c.parent.Device(), c.parent.Channel(), c.value)
	case afterTouch:
		c.midi.AfterTouch(c.parent.Device(), c.parent.Channel(), uint8(c.value))
	}
}
