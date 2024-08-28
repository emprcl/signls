package param

import (
	"cykl/core"
	"fmt"
)

type Velocity struct {
	nodes []core.Node
}

func (v Velocity) Name() string {
	return "vel"
}

func (v Velocity) Display() string {
	return fmt.Sprintf("%d", v.nodes[0].(*core.Emitter).Note().Velocity)
}

func (v Velocity) Value() int {
	return int(v.nodes[0].(*core.Emitter).Note().Velocity)
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
	for _, node := range v.nodes {
		node.(*core.Emitter).Note().SetVelocity(uint8(value))
	}
}
