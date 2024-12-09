package param

import (
	"fmt"
	"math"
	"strconv"

	"signls/core/common"
	"signls/core/node"

	"signls/ui/util"
)

type Repeat struct {
	nodes []common.Node
}

func (r Repeat) Name() string {
	return "rpt"
}

func (r Repeat) Help() string {
	return ""
}

func (r Repeat) Display() string {
	if r.control().RandomAmount() != 0 {
		return util.Normalize(
			fmt.Sprintf(
				"%d%+d\u033c",
				r.control().Value(),
				r.control().RandomAmount(),
			),
		)
	}
	return fmt.Sprintf("%d", r.Value())
}

func (r Repeat) control() *common.ControlValue[int] {
	return r.nodes[0].(*node.Emitter).Behavior().(common.Repeatable).Repeat()
}

func (r Repeat) Value() int {
	return r.control().Value()
}

func (r Repeat) AltValue() int {
	return r.control().RandomAmount()
}

func (r Repeat) Up() {
	r.Set(r.Value() + 1)
}

func (r Repeat) Down() {
	r.Set(r.Value() - 1)
}

func (r Repeat) Left() {
	r.SetAlt(r.AltValue() - 1)
}

func (r Repeat) Right() {
	r.SetAlt(r.AltValue() + 1)
}

func (r Repeat) AltUp() {}

func (r Repeat) AltDown() {}

func (r Repeat) AltLeft() {}

func (r Repeat) AltRight() {}

func (r Repeat) Set(value int) {
	if value < 0 || value >= math.MaxInt32 {
		return
	}
	for _, n := range r.nodes {
		n.(*node.Emitter).Behavior().(common.Repeatable).Repeat().Set(value)
	}
}

func (r Repeat) SetAlt(value int) {
	for _, n := range r.nodes {
		n.(*node.Emitter).Behavior().(common.Repeatable).Repeat().SetRandomAmount(value)
	}
}

func (r Repeat) SetEditValue(input string) {
	value, err := strconv.Atoi(input)
	if err != nil {
		return
	}
	r.Set(value)
}
