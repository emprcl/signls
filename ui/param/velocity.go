package param

import (
	"fmt"

	"signls/core/common"
	"signls/core/music"
)

type Velocity struct {
	nodes []common.Node
}

func (v Velocity) Name() string {
	return "vel"
}

func (v Velocity) Display() string {
	if v.nodes[0].(music.Audible).Note().Velocity.RandomAmount() != 0 {
		return fmt.Sprintf(
			"%d%+d\u033c",
			v.nodes[0].(music.Audible).Note().Velocity.Value(),
			v.nodes[0].(music.Audible).Note().Velocity.RandomAmount(),
		)
	}
	return fmt.Sprintf("%d", v.nodes[0].(music.Audible).Note().Velocity.Value())
}

func (v Velocity) Value() int {
	return int(v.nodes[0].(music.Audible).Note().Velocity.Value())
}

func (v Velocity) AltValue() int {
	return 0
}

func (v Velocity) Up() {
	v.Set(v.Value() + 1)
}

func (v Velocity) Down() {
	v.Set(v.Value() - 1)
}

func (v Velocity) Left() {
	v.SetAlt(v.nodes[0].(music.Audible).Note().Velocity.RandomAmount() - 1)
}

func (v Velocity) Right() {
	v.SetAlt(v.nodes[0].(music.Audible).Note().Velocity.RandomAmount() + 1)
}

func (v Velocity) AltUp() {}

func (v Velocity) AltDown() {}

func (v Velocity) AltLeft() {}

func (v Velocity) AltRight() {}

func (v Velocity) Set(value int) {
	for _, n := range v.nodes {
		n.(music.Audible).Note().SetVelocity(uint8(value))
	}
}

func (v Velocity) SetAlt(value int) {
	for _, n := range v.nodes {
		n.(music.Audible).Note().Velocity.SetRandomAmount(value)
	}
}
