package param

import (
	"fmt"

	"cykl/core/common"
	"cykl/core/node"
)

type Velocity struct {
	nodes []common.Node
}

func (v Velocity) Name() string {
	return "vel"
}

func (v Velocity) Display() string {
	return fmt.Sprintf("%d", v.nodes[0].(*node.BaseEmitter).Note().Velocity)
}

func (v Velocity) Value() int {
	return int(v.nodes[0].(*node.BaseEmitter).Note().Velocity)
}

func (v Velocity) Increment() {
	v.Set(v.Value() + 1)
}

func (v Velocity) Decrement() {
	v.Set(v.Value() - 1)
}

func (v Velocity) Left() {}

func (v Velocity) Right() {}

func (v Velocity) Set(value int) {
	for _, n := range v.nodes {
		n.(*node.BaseEmitter).Note().SetVelocity(uint8(value))
	}
}
