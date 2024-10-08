package param

import (
	"fmt"

	"cykl/core/common"
	"cykl/core/node"
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
	return fmt.Sprintf("%d", t.Value())
}

func (t Triggers) Value() int {
	return t.nodes[0].(*node.EuclidEmitter).Triggers
}

func (t Triggers) AltValue() int {
	return 0
}

func (t Triggers) Increment() {
	t.Set(t.Value() + 1)
}

func (t Triggers) Decrement() {
	t.Set(t.Value() - 1)
}

func (t Triggers) Left() {}

func (t Triggers) Right() {}

func (t Triggers) Set(value int) {
	if value < minTriggers {
		return
	}
	for _, n := range t.nodes {
		if value > n.(*node.EuclidEmitter).Steps {
			continue
		}
		n.(*node.EuclidEmitter).Triggers = value
	}
}

func (t Triggers) SetAlt(value int) {}

func (t Triggers) ChangeAltMode() {}
