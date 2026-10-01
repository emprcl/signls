package param

import (
	"fmt"

	"signls/core/field"
	"signls/midi"
)

type InputDevice struct {
	grid *field.Grid
}

func (d InputDevice) Name() string {
	return "in device"
}

func (d InputDevice) Help() string {
	if !d.grid.MidiInputDevice().Enabled() {
		return ""
	} else if d.grid.MidiInputDevice().Fallback {
		return fmt.Sprintf("disconnected: %s", d.grid.MidiInputDevice().Name)
	}
	return d.grid.MidiInputDevice().Name
}

func (d InputDevice) Display() string {
	if d.grid.MidiInputDevice().Fallback {
		return "??"
	}
	return fmt.Sprintf("%d", d.grid.MidiInputDevice().ID)
}

func (d InputDevice) Value() int {
	return int(d.grid.MidiInputDevice().ID)
}

func (d InputDevice) AltValue() int {
	return 0
}

func (d InputDevice) Up() {
	d.grid.SetMidiInputDevice(d.grid.MidiInputDevice().Next())
}

func (d InputDevice) Down() {
	d.grid.SetMidiInputDevice(d.grid.MidiInputDevice().Prev())
}

func (d InputDevice) Left() {}

func (d InputDevice) Right() {}

func (d InputDevice) AltUp() {}

func (d InputDevice) AltDown() {}

func (d InputDevice) AltLeft() {}

func (d InputDevice) AltRight() {}

func (d InputDevice) Set(value int) {}

func (d InputDevice) SetAlt(value int) {}

func (d InputDevice) SetEditValue(input string) {}

func (d InputDevice) SetFromMidiIn(msg midi.InMessage) {}
