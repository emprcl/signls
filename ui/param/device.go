package param

import (
	"fmt"

	"signls/core/common"
	"signls/core/music"
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
	}
	return d.nodes[0].(music.Audible).Note().Device.Name()
}

func (d Device) Display() string {
	if !d.nodes[0].(music.Audible).Note().Device.Enabled {
		return "⨯"
	}
	return fmt.Sprintf("%d", d.nodes[0].(music.Audible).Note().Device.Device.ID)
}

func (d Device) Value() int {
	return d.nodes[0].(music.Audible).Note().Device.Get()
}

func (d Device) AltValue() int {
	return 0
}

func (d Device) Up() {
	d.Set(d.nodes[0].(music.Audible).Note().Device.Get() + 1)
}

func (d Device) Down() {
	d.Set(d.nodes[0].(music.Audible).Note().Device.Get() - 1)
}

func (d Device) Left() {
	for _, n := range d.nodes {
		n.(music.Audible).Note().Device.Enabled = !n.(music.Audible).Note().Device.Enabled
	}
}

func (d Device) Right() {
	for _, n := range d.nodes {
		n.(music.Audible).Note().Device.Enabled = !n.(music.Audible).Note().Device.Enabled
	}
}

func (d Device) AltUp() {}

func (d Device) AltDown() {}

func (d Device) AltLeft() {}

func (d Device) AltRight() {}

func (d Device) Set(value int) {
	if !d.nodes[0].(music.Audible).Note().Device.Enabled {
		return
	}
	device := d.nodes[0].(music.Audible).Note().Midi().GetDevice(value)
	for _, n := range d.nodes {
		n.(music.Audible).Note().Device.Device = device
	}
}

func (d Device) SetAlt(value int) {}

func (c Device) SetEditValue(input string) {}
