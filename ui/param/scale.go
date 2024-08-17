package param

import (
	"cykl/core"
)

type Scale struct {
	grid   *core.Grid
	scales []core.Scale
}

func (s Scale) Name() string {
	return "scale"
}

func (s Scale) Display() string {
	return s.grid.Scale.Name()
}

func (s Scale) Value() int {
	return int(s.grid.Scale)
}

func (s Scale) Increment() {
	s.Set(s.scaleIndex() + 1)
}

func (s Scale) Decrement() {
	s.Set(s.scaleIndex() - 1)
}

func (s Scale) Set(value int) {
	if value < 0 {
		value = len(s.scales) - 1
	} else if value >= len(s.scales) {
		value = 0
	}
	s.grid.SetScale(s.scales[value])
}

func (s Scale) scaleIndex() int {
	for i := 0; i < len(s.scales); i++ {
		if s.grid.Scale == s.scales[i] {
			return i
		}
	}
	return 0
}
