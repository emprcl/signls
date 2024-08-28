package param

import (
	"cykl/core"
	"fmt"
)

type Channel struct {
	nodes []core.Node
}

func (c Channel) Name() string {
	return "cha"
}

func (c Channel) Display() string {
	return fmt.Sprintf("%d", c.nodes[0].(*core.Emitter).Note().Channel+1)
}

func (c Channel) Value() int {
	return int(c.nodes[0].(*core.Emitter).Note().Channel)
}

func (c Channel) Increment() {
	c.Set(c.Value() + 1)
}

func (c Channel) Decrement() {
	c.Set(c.Value() - 1)
}

func (c Channel) Left() {}

func (c Channel) Right() {}

func (c Channel) Set(value int) {
	for _, node := range c.nodes {
		node.(*core.Emitter).Note().SetChannel(uint8(value))
	}
}
