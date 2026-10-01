package param

import (
	"signls/core/field"
	"signls/midi"
)

type InputParamEdit struct {
	grid *field.Grid
}

func (t InputParamEdit) Name() string {
	return "param edit"
}

func (t InputParamEdit) Help() string {
	return ""
}

func (t InputParamEdit) Display() string {
	if t.grid.MidiEditAllParams {
		return "all"
	}
	return "selected"
}

func (t InputParamEdit) Value() int {
	return 0
}

func (t InputParamEdit) AltValue() int {
	return 0
}

func (t InputParamEdit) Up() {
	t.grid.SetMidiEditAllParams(true)
}

func (t InputParamEdit) Down() {
	t.grid.SetMidiEditAllParams(false)
}

func (t InputParamEdit) Left() {}

func (t InputParamEdit) Right() {}

func (t InputParamEdit) AltUp() {}

func (t InputParamEdit) AltDown() {}

func (t InputParamEdit) AltLeft() {}

func (t InputParamEdit) AltRight() {}

func (t InputParamEdit) Set(value int) {}

func (t InputParamEdit) SetAlt(value int) {}

func (t InputParamEdit) SetEditValue(input string) {}

func (t InputParamEdit) SetFromMidiIn(msg midi.InMessage) {}
