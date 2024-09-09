package param

import (
	"fmt"

	"cykl/core/common"
	"cykl/core/music"
)

type Channel struct {
	nodes []common.Node
}

func (c Channel) Name() string {
	return "cha"
}

func (c Channel) Display() string {
	return fmt.Sprintf("%d", c.nodes[0].(music.Audible).Note().Channel+1)
}

func (c Channel) Value() int {
	return int(c.nodes[0].(music.Audible).Note().Channel)
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
	for _, n := range c.nodes {
		n.(music.Audible).Note().SetChannel(uint8(value))
	}
}
