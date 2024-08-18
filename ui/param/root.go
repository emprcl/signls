package param

import (
	"cykl/core"
)

const (
	maxKey int = 127
)

type Root struct {
	grid *core.Grid
}

func (r Root) Name() string {
	return "root"
}

func (r Root) Display() string {
	return r.grid.Key.Name()
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

func (r Root) Left() {}

func (r Root) Right() {}

func (r Root) Set(value int) {
	if value < 0 || value > maxKey {
		return
	}
	r.grid.SetKey(core.Key(value))
}
