package param

import (
	"signls/core/field"
	"signls/core/music"
)

const (
	maxKey int = 127
)

type Root struct {
	grid *field.Grid
}

func (r Root) Name() string {
	return "root"
}

func (r Root) Display() string {
	return "key"
}

func (r Root) Value() int {
	return int(r.grid.Key)
}

func (r Root) AltValue() int {
	return 0
}

func (r Root) Up() {
	r.Set(r.Value() + 1)
}

func (r Root) Down() {
	r.Set(r.Value() - 1)
}

func (r Root) Left() {}

func (r Root) Right() {}

func (r Root) AltUp() {}

func (r Root) AltDown() {}

func (r Root) AltLeft() {}

func (r Root) AltRight() {}

func (r Root) Set(value int) {
	if value < 0 || value > maxKey {
		return
	}
	r.grid.SetKey(music.Key(value))
}

func (r Root) SetAlt(value int) {}

func (r Root) SetEditValue(input string) {}
