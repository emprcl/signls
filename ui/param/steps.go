package param

import (
	"fmt"

	"cykl/core/common"
	"cykl/core/node"
)

const (
	minSteps = 1
	maxSteps = 128
)

type Steps struct {
	nodes []common.Node
}

func (s Steps) Name() string {
	return "stp"
}

func (s Steps) Display() string {
	return fmt.Sprintf("%d", s.Value())
}

func (s Steps) Value() int {
	return s.nodes[0].(*node.EuclidEmitter).Steps
}

func (s Steps) AltValue() int {
	return 0
}

func (s Steps) Up() {
	s.Set(s.Value() + 1)
}

func (s Steps) Down() {
	s.Set(s.Value() - 1)
}

func (s Steps) Left() {}

func (s Steps) Right() {}

func (s Steps) AltUp() {}

func (s Steps) AltDown() {}

func (s Steps) AltLeft() {}

func (s Steps) AltRight() {}

func (s Steps) Set(value int) {
	if value < minSteps || value >= maxSteps {
		return
	}
	for _, n := range s.nodes {
		n.(*node.EuclidEmitter).Steps = value
	}
}

func (s Steps) SetAlt(value int) {}
