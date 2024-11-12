package param

import (
	"fmt"
	"strconv"

	"signls/core/common"
	"signls/core/music"
	"signls/ui/util"
)

type Device struct {
	nodes []common.Node
}

func (d Device) Name() string {
	return "dvc"
}

func (d Device) Display() string {
	if d.nodes[0].(music.Audible).Note().Device.RandomAmount() != 0 {
		return util.Normalize(
			fmt.Sprintf(
				"%d%+d\u033c",
				d.nodes[0].(music.Audible).Note().Device.Value(),
				d.nodes[0].(music.Audible).Note().Device.RandomAmount(),
			),
		)
	}
	return fmt.Sprintf("%s", d.nodes[0].(music.Audible).Note().DeviceName())
}

func (d Device) Value() int {
	return int(d.nodes[0].(music.Audible).Note().Device.Value())
}

func (d Device) AltValue() int {
	return 0
}

func (d Device) Up() {
	d.Set(d.Value() + 1)
}

func (d Device) Down() {
	d.Set(d.Value() - 1)
}

func (d Device) Left() {
	d.SetAlt(d.nodes[0].(music.Audible).Note().Device.RandomAmount() - 1)
}

func (d Device) Right() {
	d.SetAlt(d.nodes[0].(music.Audible).Note().Device.RandomAmount() + 1)
}

func (d Device) AltUp() {}

func (d Device) AltDown() {}

func (d Device) AltLeft() {}

func (d Device) AltRight() {}

func (d Device) Set(value int) {
	for _, n := range d.nodes {
		n.(music.Audible).Note().SetDevice(value)
	}
}

func (d Device) SetAlt(value int) {
	for _, n := range d.nodes {
		n.(music.Audible).Note().Device.SetRandomAmount(value)
	}
}

func (c Device) SetEditValue(input string) {
	value, err := strconv.Atoi(input)
	if err != nil {
		return
	}
	c.Set(value - 1)
}
