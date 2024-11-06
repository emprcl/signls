package param

import (
	"fmt"

	"signls/core/common"
	"signls/core/node"
	"signls/ui/util"
)

const (
	minTriggers = 1
)

type Triggers struct {
	nodes []common.Node
}

func (t Triggers) Name() string {
	return "trg"
}

func (t Triggers) Display() string {
	if t.nodes[0].(*node.EuclidEmitter).Triggers.RandomAmount() != 0 {
		return util.Normalize(
			fmt.Sprintf(
				"%d%+d\u033c",
				t.nodes[0].(*node.EuclidEmitter).Triggers.Value(),
				t.nodes[0].(*node.EuclidEmitter).Triggers.RandomAmount(),
			),
		)
	}
	return fmt.Sprintf("%d", t.Value())
}

func (t Triggers) Value() int {
	return t.nodes[0].(*node.EuclidEmitter).Triggers.Value()
}

func (t Triggers) AltValue() int {
	return t.nodes[0].(*node.EuclidEmitter).Triggers.RandomAmount()
}

func (t Triggers) Up() {
	t.Set(t.Value() + 1)
}

func (t Triggers) Down() {
	t.Set(t.Value() - 1)
}

func (t Triggers) Left() {
	t.SetAlt(t.nodes[0].(*node.EuclidEmitter).Triggers.RandomAmount() - 1)
}

func (t Triggers) Right() {
	t.SetAlt(t.nodes[0].(*node.EuclidEmitter).Triggers.RandomAmount() + 1)
}

func (t Triggers) AltUp() {}

func (t Triggers) AltDown() {}

func (t Triggers) AltLeft() {}

func (t Triggers) AltRight() {}

func (t Triggers) Set(value int) {
	if value < minTriggers {
		return
	}
	for _, n := range t.nodes {
		if value > n.(*node.EuclidEmitter).Steps.Value() {
			continue
		}
		n.(*node.EuclidEmitter).Triggers.Set(value)
	}
}

func (t Triggers) SetAlt(value int) {
	for _, n := range t.nodes {
		n.(*node.EuclidEmitter).Triggers.SetRandomAmount(value)
	}
}
