package param

import (
	"signls/core/field"
)

type ClockSend struct {
	grid *field.Grid
}

func (c ClockSend) Name() string {
	return "clock"
}

func (c ClockSend) Display() string {
	if c.grid.SendClock {
		return "on"
	}
	return "off"
}

func (c ClockSend) Value() int {
	return 0
}

func (c ClockSend) AltValue() int {
	return 0
}

func (c ClockSend) Up() {
	c.grid.SendClock = true
}

func (c ClockSend) Down() {
	c.grid.SendClock = false
}

func (c ClockSend) Left() {}

func (c ClockSend) Right() {}

func (c ClockSend) AltUp() {}

func (c ClockSend) AltDown() {}

func (c ClockSend) AltLeft() {}

func (c ClockSend) AltRight() {}

func (c ClockSend) Set(value int) {}

func (c ClockSend) SetAlt(value int) {}

func (c ClockSend) SetEditValue(input string) {}
