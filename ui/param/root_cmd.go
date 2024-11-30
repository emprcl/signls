package param

import (
	"fmt"
	"signls/core/common"
	"signls/core/music"
	"signls/ui/util"
)

const (
	cmdIndex = 0
)

type RootCmd struct {
	nodes []common.Node
}

func (r RootCmd) Name() string {
	return "root"
}

func (r RootCmd) Display() string {
	if !r.nodes[0].(music.Audible).Note().MetaCommands[cmdIndex].Active() {
		return "⨯"
	}
	if r.nodes[0].(music.Audible).Note().MetaCommands[cmdIndex].Value().RandomAmount() != 0 {
		return util.Normalize(
			fmt.Sprintf(
				"%s%+d\u033c",
				r.nodes[0].(music.Audible).Note().MetaCommands[cmdIndex].Display(),
				r.nodes[0].(music.Audible).Note().MetaCommands[cmdIndex].Value().RandomAmount(),
			),
		)
	}
	return r.nodes[0].(music.Audible).Note().MetaCommands[cmdIndex].Display()
}

func (r RootCmd) Value() int {
	return r.nodes[0].(music.Audible).Note().MetaCommands[cmdIndex].Value().Value()
}

func (r RootCmd) AltValue() int {
	return r.nodes[0].(music.Audible).Note().MetaCommands[cmdIndex].Value().RandomAmount()
}

func (r RootCmd) Up() {
	r.Set(r.Value() + 1)
}

func (r RootCmd) Down() {
	r.Set(r.Value() - 1)
}

func (r RootCmd) Left() {
	r.SetAlt(r.AltValue() - 1)
}

func (r RootCmd) Right() {
	r.SetAlt(r.AltValue() + 1)
}

func (r RootCmd) AltUp() {}

func (r RootCmd) AltDown() {}

func (r RootCmd) AltLeft() {
	active := r.nodes[0].(music.Audible).Note().MetaCommands[cmdIndex].Active()
	for _, n := range r.nodes {
		n.(music.Audible).Note().MetaCommands[cmdIndex].SetActive(!active)
	}
}

func (r RootCmd) AltRight() {
	active := r.nodes[0].(music.Audible).Note().MetaCommands[cmdIndex].Active()
	for _, n := range r.nodes {
		n.(music.Audible).Note().MetaCommands[cmdIndex].SetActive(!active)
	}
}

func (r RootCmd) Set(value int) {
	if value < 0 || value > maxKey {
		return
	}
	for _, n := range r.nodes {
		n.(music.Audible).Note().MetaCommands[cmdIndex].Value().Set(value)
	}
}

func (r RootCmd) SetAlt(value int) {
	for _, n := range r.nodes {
		n.(music.Audible).Note().MetaCommands[cmdIndex].Value().SetRandomAmount(value)
	}
}

func (r RootCmd) SetEditValue(input string) {}
