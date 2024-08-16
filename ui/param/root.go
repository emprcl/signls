package param

import (
	"cykl/core"
	"cykl/midi"
)

type Root struct {
	grid *core.Grid
}

func (r Root) Name() string {
	return "root"
}

func (r Root) Display() string {
	return midi.Note(r.grid.Key)
}

func (r Root) Value() int {
	return int(r.grid.Key)
}

func (r Root) Increment() {
	r.Set(r.Value() + 1)
}

func (r Root) Decrement() {
	r.Set(r.Value() - 1)
}

func (r Root) Set(value int) {
	if value > 127 {
		return
	}
	r.grid.Key = uint8(value)
}
