package param

import (
	"fmt"
	"math"

	"cykl/core/common"
	"cykl/core/node"
)

type Threshold struct {
	nodes []common.Node
}

func (t Threshold) Name() string {
	return "thd"
}

func (t Threshold) Display() string {
	return fmt.Sprintf("%d", t.Value())
}

func (t Threshold) Value() int {
	return t.nodes[0].(*node.Emitter).Behavior().(*node.TollEmitter).Threshold
}

func (t Threshold) AltValue() int {
	return 0
}

func (t Threshold) Up() {
	t.Set(t.Value() + 1)
}

func (t Threshold) Down() {
	t.Set(t.Value() - 1)
}

func (t Threshold) Left() {}

func (t Threshold) Right() {}

func (t Threshold) AltUp() {}

func (t Threshold) AltDown() {}

func (t Threshold) AltLeft() {}

func (t Threshold) AltRight() {}

func (t Threshold) Set(value int) {
	if value < 0 || value >= math.MaxInt32 {
		return
	}
	for _, n := range t.nodes {
		n.(*node.Emitter).Behavior().(*node.TollEmitter).Threshold = value
	}
}

func (t Threshold) SetAlt(value int) {}
