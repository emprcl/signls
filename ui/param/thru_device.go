package param

import (
	"fmt"

	"signls/core/field"
	"signls/midi"
)

type ThruDevice struct {
	grid *field.Grid
}

func (d ThruDevice) Name() string {
	return "thru device"
}

func (d ThruDevice) Help() string {
	if !d.grid.MidiThruDevice().Enabled() {
		return ""
	} else if d.grid.MidiThruDevice().Fallback {
		return fmt.Sprintf("disconnected: %s", d.grid.MidiThruDevice().Name)
	}
	return d.grid.MidiThruDevice().Name
}

func (d ThruDevice) Display() string {
	if d.grid.MidiThruDevice().Fallback {
		return "??"
	}
	return fmt.Sprintf("%d", d.grid.MidiThruDevice().ID)
}

func (d ThruDevice) Value() int {
	return int(d.grid.MidiThruDevice().ID)
}

func (d ThruDevice) AltValue() int {
	return 0
}

func (d ThruDevice) Up() {
	d.grid.SetMidiThruDevice(d.grid.MidiThruDevice().Next())
}

func (d ThruDevice) Down() {
	d.grid.SetMidiThruDevice(d.grid.MidiThruDevice().Prev())
}

func (d ThruDevice) Left() {}

func (d ThruDevice) Right() {}

func (d ThruDevice) AltUp() {}

func (d ThruDevice) AltDown() {}

func (d ThruDevice) AltLeft() {}

func (d ThruDevice) AltRight() {}

func (d ThruDevice) Set(value int) {}

func (d ThruDevice) SetAlt(value int) {}

func (d ThruDevice) SetEditValue(input string) {}

func (d ThruDevice) SetFromMidiIn(msg midi.InMessage) {}
