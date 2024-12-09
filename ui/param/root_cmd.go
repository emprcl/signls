package param

import (
	"fmt"
	"signls/core/common"
	"signls/core/music"
	"signls/ui/util"
)

const (
	rootCmdIndex = 2
)

type RootCmd struct {
	nodes []common.Node
}

func (r RootCmd) Name() string {
	return "root"
}

func (r RootCmd) Help() string {
	return ""
}

func (r RootCmd) Display() string {
	if !r.nodes[0].(music.Audible).Note().MetaCommands[rootCmdIndex].Active() {
		return "⨯"
	}
	if r.nodes[0].(music.Audible).Note().MetaCommands[rootCmdIndex].Value().RandomAmount() != 0 {
		return util.Normalize(
			fmt.Sprintf(
				"%s%+d\u033c",
				r.nodes[0].(music.Audible).Note().MetaCommands[rootCmdIndex].Display(),
				r.nodes[0].(music.Audible).Note().MetaCommands[rootCmdIndex].Value().RandomAmount(),
			),
		)
	}
	return r.nodes[0].(music.Audible).Note().MetaCommands[rootCmdIndex].Display()
}

func (r RootCmd) Value() int {
	return r.nodes[0].(music.Audible).Note().MetaCommands[rootCmdIndex].Value().Value()
}

func (r RootCmd) AltValue() int {
	return r.nodes[0].(music.Audible).Note().MetaCommands[rootCmdIndex].Value().RandomAmount()
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
	active := r.nodes[0].(music.Audible).Note().MetaCommands[rootCmdIndex].Active()
	for _, n := range r.nodes {
		n.(music.Audible).Note().MetaCommands[rootCmdIndex].SetActive(!active)
	}
}

func (r RootCmd) AltRight() {
	active := r.nodes[0].(music.Audible).Note().MetaCommands[rootCmdIndex].Active()
	for _, n := range r.nodes {
		n.(music.Audible).Note().MetaCommands[rootCmdIndex].SetActive(!active)
	}
}

func (r RootCmd) Set(value int) {
	for _, n := range r.nodes {
		n.(music.Audible).Note().MetaCommands[rootCmdIndex].Value().Set(value)
	}
}

func (r RootCmd) SetAlt(value int) {
	for _, n := range r.nodes {
		n.(music.Audible).Note().MetaCommands[rootCmdIndex].Value().SetRandomAmount(value)
	}
}

func (r RootCmd) SetEditValue(input string) {
	midiKey, err := music.ConvertNoteToMIDI(input)
	if err != nil {
		return
	}
	r.Set(midiKey)
}
