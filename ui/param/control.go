package param

import (
	"fmt"
	"strconv"

	"signls/core/common"
	"signls/core/music"
	"signls/ui/util"
)

type Control struct {
	index int
	nodes []common.Node
}

func (c Control) Name() string {
	switch c.nodes[0].(music.Audible).Note().Controls[c.index].Type {
	case music.ControlChangeControlType:
		return fmt.Sprintf("cc%d", c.nodes[0].(music.Audible).Note().Controls[c.index].Controller)
	case music.AfterTouchControlType:
		return "at"
	case music.PitchBendControlType:
		return "pb"
	case music.ProgramChangeControlType:
		return "pc"
	default:
		return "cc"
	}
}

func (c Control) Display() string {
	if c.nodes[0].(music.Audible).Note().Controls[c.index].Type == music.SilentControlType {
		return "⨯"
	}
	if c.nodes[0].(music.Audible).Note().Controls[c.index].Value.RandomAmount() != 0 {
		return util.Normalize(
			fmt.Sprintf(
				"%d%+d\u033c",
				c.nodes[0].(music.Audible).Note().Controls[c.index].Value.Value(),
				c.nodes[0].(music.Audible).Note().Controls[c.index].Value.RandomAmount(),
			),
		)
	}
	return fmt.Sprintf("%d", c.nodes[0].(music.Audible).Note().Controls[c.index].Value.Value())
}

func (c Control) Value() int {
	return int(c.nodes[0].(music.Audible).Note().Controls[c.index].Value.Value())
}

func (c Control) AltValue() int {
	return 0
}

func (c Control) Up() {
	c.Set(c.Value() + 1)
}

func (c Control) Down() {
	c.Set(c.Value() - 1)
}

func (c Control) Left() {
	c.SetAlt(c.nodes[0].(music.Audible).Note().Controls[c.index].Value.RandomAmount() - 1)
}

func (c Control) Right() {
	c.SetAlt(c.nodes[0].(music.Audible).Note().Controls[c.index].Value.RandomAmount() + 1)
}

func (c Control) AltUp() {
	c.SetController(c.nodes[0].(music.Audible).Note().Controls[c.index].Controller + 1)
}

func (c Control) AltDown() {
	c.SetController(c.nodes[0].(music.Audible).Note().Controls[c.index].Controller - 1)
}

func (c Control) AltLeft() {
	newMode := util.Mod((int(c.nodes[0].(music.Audible).Note().Controls[c.index].Type) - 1), len(music.AllControlTypes))
	for _, n := range c.nodes {
		n.(music.Audible).Note().Controls[c.index].SetType(newMode)
	}
}

func (c Control) AltRight() {
	newMode := util.Mod((int(c.nodes[0].(music.Audible).Note().Controls[c.index].Type) + 1), len(music.AllControlTypes))
	for _, n := range c.nodes {
		n.(music.Audible).Note().Controls[c.index].SetType(newMode)
	}
}

func (c Control) Set(value int) {
	for _, n := range c.nodes {
		n.(music.Audible).Note().Controls[c.index].Value.Set(uint8(value))
	}
}

func (c Control) SetAlt(value int) {
	for _, n := range c.nodes {
		n.(music.Audible).Note().Controls[c.index].Value.SetRandomAmount(value)
	}
}

func (c Control) SetController(value uint8) {
	for _, n := range c.nodes {
		n.(music.Audible).Note().Controls[c.index].SetController(value)
	}
}

func (c Control) SetEditValue(input string) {
	value, err := strconv.Atoi(input)
	if err != nil {
		return
	}
	c.Set(value - 1)
}
