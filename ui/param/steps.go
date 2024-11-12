package param

import (
	"fmt"
	"strconv"

	"signls/core/common"
	"signls/core/node"

	"signls/ui/util"
)

type Steps struct {
	nodes []common.Node
}

func (s Steps) Name() string {
	return "stp"
}

func (s Steps) Display() string {
	if s.nodes[0].(*node.EuclidEmitter).Steps.RandomAmount() != 0 {
		return util.Normalize(
			fmt.Sprintf(
				"%d%+d\u033c",
				s.nodes[0].(*node.EuclidEmitter).Steps.Value(),
				s.nodes[0].(*node.EuclidEmitter).Steps.RandomAmount(),
			),
		)
	}
	return fmt.Sprintf("%d", s.Value())
}

func (s Steps) Value() int {
	return s.nodes[0].(*node.EuclidEmitter).Steps.Value()
}

func (s Steps) AltValue() int {
	return s.nodes[0].(*node.EuclidEmitter).Steps.RandomAmount()
}

func (s Steps) Up() {
	s.Set(s.Value() + 1)
}

func (s Steps) Down() {
	s.Set(s.Value() - 1)
}

func (s Steps) Left() {
	s.SetAlt(s.nodes[0].(*node.EuclidEmitter).Steps.RandomAmount() - 1)
}

func (s Steps) Right() {
	s.SetAlt(s.nodes[0].(*node.EuclidEmitter).Steps.RandomAmount() + 1)
}

func (s Steps) AltUp() {}

func (s Steps) AltDown() {}

func (s Steps) AltLeft() {}

func (s Steps) AltRight() {}

func (s Steps) Set(value int) {
	for _, n := range s.nodes {
		n.(*node.EuclidEmitter).Steps.Set(value)
		if n.(*node.EuclidEmitter).Offset.Value() > value {
			n.(*node.EuclidEmitter).Offset.Set(value)
		}

		if n.(*node.EuclidEmitter).Triggers.Value() > value {
			n.(*node.EuclidEmitter).Triggers.Set(value)
		}
	}
}

func (s Steps) SetAlt(value int) {
	for _, n := range s.nodes {
		n.(*node.EuclidEmitter).Steps.SetRandomAmount(value)
	}
}

func (s Steps) SetEditValue(input string) {
	value, err := strconv.Atoi(input)
	if err != nil {
		return
	}
	s.Set(value)
}
