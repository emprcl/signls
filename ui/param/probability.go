package param

import (
	"fmt"

	"cykl/core/common"
	"cykl/core/music"
)

const (
	maxProbability = 100
)

type Probability struct {
	nodes []common.Node
}

func (p Probability) Name() string {
	return "prob"
}

func (p Probability) Display() string {
	return fmt.Sprintf("%d", p.Value())
}

func (p Probability) Value() int {
	return int(p.nodes[0].(music.Audible).Note().Probability)
}

func (p Probability) Increment() {
	p.Set(p.Value() + 1)
}

func (p Probability) Decrement() {
	p.Set(p.Value() - 1)
}

func (p Probability) Left() {}

func (p Probability) Right() {}

func (p Probability) Set(value int) {
	if value < 0 || value > maxProbability {
		return
	}
	for _, n := range p.nodes {
		n.(music.Audible).Note().Probability = uint8(value)
	}
}
