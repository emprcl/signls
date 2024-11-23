package param

import (
	"signls/core/field"
)

type TransportSend struct {
	grid *field.Grid
}

func (t TransportSend) Name() string {
	return "transport send"
}

func (t TransportSend) Display() string {
	if t.grid.SendTransport {
		return "on"
	}
	return "off"
}

func (t TransportSend) Value() int {
	return 0
}

func (t TransportSend) AltValue() int {
	return 0
}

func (t TransportSend) Up() {
	t.grid.SendTransport = true
}

func (t TransportSend) Down() {
	t.grid.SendTransport = false
}

func (t TransportSend) Left() {}

func (t TransportSend) Right() {}

func (t TransportSend) AltUp() {}

func (t TransportSend) AltDown() {}

func (t TransportSend) AltLeft() {}

func (t TransportSend) AltRight() {}

func (t TransportSend) Set(value int) {}

func (t TransportSend) SetAlt(value int) {}

func (t TransportSend) SetEditValue(input string) {}
