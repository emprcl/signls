package param

import (
	"fmt"

	"signls/core/common"
	"signls/core/music"
	"signls/midi"
)

type Device struct {
	nodes []common.Node
}

func (d Device) Name() string {
	return "dvc"
}

func (d Device) Help() string {
	if !d.nodes[0].(music.Audible).Note().Device.Enabled {
		return ""
	} else if d.nodes[0].(music.Audible).Note().Device.Device.Fallback {
		return fmt.Sprintf("disconnected: %s", d.nodes[0].(music.Audible).Note().Device.Name())
	}
	return d.nodes[0].(music.Audible).Note().Device.Name()
}

func (d Device) Display() string {
	if !d.nodes[0].(music.Audible).Note().Device.Enabled {
		return "⨯"
	} else if d.nodes[0].(music.Audible).Note().Device.Device.Fallback {
		return "??"
	}
	return fmt.Sprintf("%d", d.nodes[0].(music.Audible).Note().Device.Device.ID)
}

func (d Device) Value() int {
	return int(d.nodes[0].(music.Audible).Note().Device.Get())
}

func (d Device) AltValue() int {
	return 0
}

func (d Device) Up() {
	d.Set(int(d.nodes[0].(music.Audible).Note().Device.Device.Next().ID))
}

func (d Device) Down() {
	d.Set(int(d.nodes[0].(music.Audible).Note().Device.Device.Prev().ID))
}

func (d Device) Left() {}

func (d Device) Right() {}

func (d Device) AltUp() {}

func (d Device) AltDown() {}

func (d Device) AltLeft() {
	d.SetAlt(0)
}

func (d Device) AltRight() {
	d.SetAlt(0)
}

func (d Device) Set(value int) {
	if !d.nodes[0].(music.Audible).Note().Device.Enabled {
		return
	}
	device := d.nodes[0].(music.Audible).Note().Midi().GetOutDevice(midi.MidiOutput(value))
	for _, n := range d.nodes {
		n.(music.Audible).Note().Device.Device = device
	}
}

func (d Device) SetAlt(value int) {
	enabled := !d.nodes[0].(music.Audible).Note().Device.Enabled
	for _, n := range d.nodes {
		if enabled {
			n.(music.Audible).Note().Device.Device = d.nodes[0].(music.Audible).Note().Midi().GetOutDevice(midi.MidiOutput(value))
		}
		n.(music.Audible).Note().Device.Enabled = enabled
	}
}

func (c Device) SetEditValue(input string) {}

func (c Device) SetFromMidiIn(msg midi.InMessage) {}
