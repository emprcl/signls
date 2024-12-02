package param

import (
	"fmt"
	"strconv"
	"strings"

	"signls/core/common"
	"signls/core/music"
	"signls/midi"
	"signls/ui/util"
)

type CC struct {
	index int
	nodes []common.Node
}

func (c CC) Name() string {
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

func (c CC) Help() string {
	switch c.nodes[0].(music.Audible).Note().Controls[c.index].Type {
	case music.ControlChangeControlType:
		return strings.ToLower(midi.CC(c.nodes[0].(music.Audible).Note().Controls[c.index].Controller))
	case music.AfterTouchControlType:
		return "after touch"
	case music.PitchBendControlType:
		return "pitch bend"
	case music.ProgramChangeControlType:
		return "program change"
	default:
		return ""
	}
}

func (c CC) Display() string {
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

func (c CC) Value() int {
	return int(c.nodes[0].(music.Audible).Note().Controls[c.index].Value.Value())
}

func (c CC) AltValue() int {
	return 0
}

func (c CC) Up() {
	c.Set(c.Value() + 1)
}

func (c CC) Down() {
	c.Set(c.Value() - 1)
}

func (c CC) Left() {
	c.SetAlt(c.nodes[0].(music.Audible).Note().Controls[c.index].Value.RandomAmount() - 1)
}

func (c CC) Right() {
	c.SetAlt(c.nodes[0].(music.Audible).Note().Controls[c.index].Value.RandomAmount() + 1)
}

func (c CC) AltUp() {
	c.SetController(c.nodes[0].(music.Audible).Note().Controls[c.index].Controller + 1)
}

func (c CC) AltDown() {
	c.SetController(c.nodes[0].(music.Audible).Note().Controls[c.index].Controller - 1)
}

func (c CC) AltLeft() {
	newMode := util.Mod((int(c.nodes[0].(music.Audible).Note().Controls[c.index].Type) - 1), len(music.AllControlTypes))
	for _, n := range c.nodes {
		n.(music.Audible).Note().Controls[c.index].SetType(newMode)
	}
}

func (c CC) AltRight() {
	newMode := util.Mod((int(c.nodes[0].(music.Audible).Note().Controls[c.index].Type) + 1), len(music.AllControlTypes))
	for _, n := range c.nodes {
		n.(music.Audible).Note().Controls[c.index].SetType(newMode)
	}
}

func (c CC) Set(value int) {
	for _, n := range c.nodes {
		n.(music.Audible).Note().Controls[c.index].Value.Set(uint8(value))
	}
}

func (c CC) SetAlt(value int) {
	for _, n := range c.nodes {
		n.(music.Audible).Note().Controls[c.index].Value.SetRandomAmount(value)
	}
}

func (c CC) SetController(value uint8) {
	for _, n := range c.nodes {
		n.(music.Audible).Note().Controls[c.index].SetController(value)
	}
}

func (c CC) SetEditValue(input string) {
	value, err := strconv.Atoi(input)
	if err != nil {
		return
	}
	c.Set(value)
}
