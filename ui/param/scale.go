package param

import (
	"signls/core/field"
	"signls/core/theory"
)

type Scale struct {
	grid   *field.Grid
	scales []theory.Scale
}

func (s Scale) Name() string {
	return "scale"
}

func (s Scale) Help() string {
	return ""
}

func (s Scale) Display() string {
	return s.grid.Scale.Name()
}

func (s Scale) Value() int {
	return int(s.grid.Scale)
}

func (s Scale) AltValue() int {
	return 0
}

func (s Scale) Up() {
	s.Set(s.scaleIndex() + 1)
}

func (s Scale) Down() {
	s.Set(s.scaleIndex() - 1)
}

func (s Scale) Left() {}

func (s Scale) Right() {}

func (s Scale) AltUp() {}

func (s Scale) AltDown() {}

func (s Scale) AltLeft() {}

func (s Scale) AltRight() {}

func (s Scale) Set(value int) {
	if value < 0 {
		value = len(s.scales) - 1
	} else if value >= len(s.scales) {
		value = 0
	}
	s.grid.SetScale(s.scales[value])
}

func (s Scale) SetAlt(value int) {}

func (s Scale) scaleIndex() int {
	for i := 0; i < len(s.scales); i++ {
		if s.grid.Scale == s.scales[i] {
			return i
		}
	}
	return 0
}

func (s Scale) SetEditValue(input string) {}
