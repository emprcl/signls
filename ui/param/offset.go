package param

import (
	"cykl/core/common"
	"cykl/core/node"
	"fmt"
)

type Offset struct {
	nodes []common.Node
}

func (o Offset) Name() string {
	return "offset"
}

func (o Offset) Display() string {
	return fmt.Sprintf("%d", o.Value())
}

func (o Offset) Value() int {
	return o.nodes[0].(*node.EuclidEmitter).Offset
}

func (o Offset) Increment() {
	o.Set(o.Value() + 1)
}

func (o Offset) Decrement() {
	o.Set(o.Value() - 1)
}

func (o Offset) Left() {}

func (o Offset) Right() {}

func (o Offset) Set(value int) {
	if value < 0 {
		return
	}
	for _, n := range o.nodes {
		if value > n.(*node.EuclidEmitter).Steps {
			continue
		}
		n.(*node.EuclidEmitter).Offset = value
	}
}
