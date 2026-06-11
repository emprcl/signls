package param

import (
	"fmt"
	"signls/core/common"
	"signls/core/music"
	"signls/ui/util"
	"strconv"
)

const (
	tempoCmdIndex = 0
)

type TempoCmd struct {
	nodes []common.Node
}

func (t TempoCmd) Name() string {
	return "tempo"
}

func (t TempoCmd) Help() string {
	return ""
}

func (t TempoCmd) Display() string {
	if !t.nodes[0].(music.Audible).Note().MetaCommands[tempoCmdIndex].Active() {
		return "⨯"
	}
	if t.nodes[0].(music.Audible).Note().MetaCommands[tempoCmdIndex].Value().RandomAmount() != 0 {
		return util.Normalize(
			fmt.Sprintf(
				"%s%+d\u033c",
				t.nodes[0].(music.Audible).Note().MetaCommands[tempoCmdIndex].Display(),
				t.nodes[0].(music.Audible).Note().MetaCommands[tempoCmdIndex].Value().RandomAmount(),
			),
		)
	}
	return t.nodes[0].(music.Audible).Note().MetaCommands[tempoCmdIndex].Display()
}

func (t TempoCmd) Value() int {
	return t.nodes[0].(music.Audible).Note().MetaCommands[tempoCmdIndex].Value().Value()
}

func (t TempoCmd) AltValue() int {
	return t.nodes[0].(music.Audible).Note().MetaCommands[tempoCmdIndex].Value().RandomAmount()
}

func (t TempoCmd) Up() {
	t.Set(t.Value() + 1)
}

func (t TempoCmd) Down() {
	t.Set(t.Value() - 1)
}

func (t TempoCmd) Left() {
	t.SetAlt(t.AltValue() - 1)
}

func (t TempoCmd) Right() {
	t.SetAlt(t.AltValue() + 1)
}

func (t TempoCmd) AltUp() {}

func (t TempoCmd) AltDown() {}

func (t TempoCmd) AltLeft() {
	active := t.nodes[0].(music.Audible).Note().MetaCommands[tempoCmdIndex].Active()
	for _, n := range t.nodes {
		n.(music.Audible).Note().MetaCommands[tempoCmdIndex].SetActive(!active)
	}
}

func (t TempoCmd) AltRight() {
	active := t.nodes[0].(music.Audible).Note().MetaCommands[tempoCmdIndex].Active()
	for _, n := range t.nodes {
		n.(music.Audible).Note().MetaCommands[tempoCmdIndex].SetActive(!active)
	}
}

func (t TempoCmd) Set(value int) {
	for _, n := range t.nodes {
		n.(music.Audible).Note().MetaCommands[tempoCmdIndex].Value().Set(value)
	}
}

func (t TempoCmd) SetAlt(value int) {
	for _, n := range t.nodes {
		n.(music.Audible).Note().MetaCommands[tempoCmdIndex].Value().SetRandomAmount(value)
	}
}

func (t TempoCmd) SetEditValue(input string) {
	value, err := strconv.Atoi(input)
	if err != nil {
		return
	}
	t.Set(value)
}
