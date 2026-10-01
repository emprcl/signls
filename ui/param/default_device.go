package param

import (
	"fmt"

	"signls/core/field"
	"signls/midi"
)

type DefaultDevice struct {
	grid *field.Grid
}

func (d DefaultDevice) Name() string {
	return "out device"
}

func (d DefaultDevice) Help() string {
	if !d.grid.MidiOutputDevice().Enabled() {
		return ""
	} else if d.grid.MidiOutputDevice().Fallback {
		return fmt.Sprintf("disconnected: %s", d.grid.MidiOutputDevice().Name)
	}
	return d.grid.MidiOutputDevice().Name
}

func (d DefaultDevice) Display() string {
	if d.grid.MidiOutputDevice().Fallback {
		return "??"
	}
	return fmt.Sprintf("%d", d.grid.MidiOutputDevice().ID)
}

func (d DefaultDevice) Value() int {
	return int(d.grid.MidiOutputDevice().ID)
}

func (d DefaultDevice) AltValue() int {
	return 0
}

func (d DefaultDevice) Up() {
	d.grid.SetMidiOutputDevice(d.grid.MidiOutputDevice().Next())
}

func (d DefaultDevice) Down() {
	d.grid.SetMidiOutputDevice(d.grid.MidiOutputDevice().Prev())
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

func (d DefaultDevice) SetFromMidiIn(msg midi.InMessage) {}
