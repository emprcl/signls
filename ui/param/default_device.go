package param

import (
	"signls/core/field"
)

type DefaultDevice struct {
	grid *field.Grid
}

func (d DefaultDevice) Name() string {
	return "device"
}

func (d DefaultDevice) Help() string {
	return ""
}

func (d DefaultDevice) Display() string {
	// TODO: handle disconnected as well
	return d.grid.MidiDevice().Name
}

func (d DefaultDevice) Value() int {
	return d.grid.MidiDevice().ID
}

func (d DefaultDevice) AltValue() int {
	return 0
}

func (d DefaultDevice) Up() {
	d.grid.SetMidiDevice(d.grid.Midi().GetDevice(d.Value() + 1))
}

func (d DefaultDevice) Down() {
	d.grid.SetMidiDevice(d.grid.Midi().GetDevice(d.Value() - 1))
}

func (d DefaultDevice) Left() {}

func (d DefaultDevice) Right() {}

func (d DefaultDevice) AltUp() {}

func (d DefaultDevice) AltDown() {}

func (d DefaultDevice) AltLeft() {}

func (d DefaultDevice) AltRight() {}

func (d DefaultDevice) Set(value int) {}

func (d DefaultDevice) SetAlt(value int) {}

func (d DefaultDevice) SetEditValue(input string) {}
