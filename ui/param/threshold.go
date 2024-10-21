package param

import (
	"fmt"
	"math"

	"signls/core/common"
	"signls/core/node"
)

type Threshold struct {
	nodes []common.Node
}

func (t Threshold) Name() string {
	return "thd"
}

func (t Threshold) Display() string {
	if t.control().RandomAmount() != 0 {
		return fmt.Sprintf(
			"%d%+d\u033c",
			t.control().Value(),
			t.control().RandomAmount(),
		)
	}
	return fmt.Sprintf("%d", t.Value())
}

func (t Threshold) control() *common.ControlValue[int] {
	return t.nodes[0].(*node.Emitter).Behavior().(*node.TollEmitter).Threshold
}

func (t Threshold) Value() int {
	return t.control().Value()
}

func (t Threshold) AltValue() int {
	return t.control().RandomAmount()
}

func (t Threshold) Up() {
	t.Set(t.Value() + 1)
}

func (t Threshold) Down() {
	t.Set(t.Value() - 1)
}

func (t Threshold) Left() {
	t.SetAlt(t.AltValue() - 1)
}

func (t Threshold) Right() {
	t.SetAlt(t.AltValue() + 1)
}

func (t Threshold) AltUp() {}

func (t Threshold) AltDown() {}

func (t Threshold) AltLeft() {}

func (t Threshold) AltRight() {}

func (t Threshold) Set(value int) {
	if value < 0 || value >= math.MaxInt32 {
		return
	}
	for _, n := range t.nodes {
		n.(*node.Emitter).Behavior().(*node.TollEmitter).Threshold.Set(value)
	}
}

func (t Threshold) SetAlt(value int) {
	for _, n := range t.nodes {
		n.(*node.Emitter).Behavior().(*node.TollEmitter).Threshold.SetRandomAmount(value)
	}
}
