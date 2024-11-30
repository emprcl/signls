package param

import (
	"fmt"
	"signls/core/common"
	"signls/core/music"
	"signls/ui/util"
)

const (
	bankCmdIndex = 1
)

type BankCmd struct {
	nodes []common.Node
}

func (b BankCmd) Name() string {
	return "bank"
}

func (b BankCmd) Display() string {
	if !b.nodes[0].(music.Audible).Note().MetaCommands[bankCmdIndex].Active() {
		return "⨯"
	}
	if b.nodes[0].(music.Audible).Note().MetaCommands[bankCmdIndex].Value().RandomAmount() != 0 {
		return util.Normalize(
			fmt.Sprintf(
				"%s%+d\u033c",
				b.nodes[0].(music.Audible).Note().MetaCommands[bankCmdIndex].Display(),
				b.nodes[0].(music.Audible).Note().MetaCommands[bankCmdIndex].Value().RandomAmount(),
			),
		)
	}
	return b.nodes[0].(music.Audible).Note().MetaCommands[bankCmdIndex].Display()
}

func (b BankCmd) Value() int {
	return b.nodes[0].(music.Audible).Note().MetaCommands[bankCmdIndex].Value().Value()
}

func (b BankCmd) AltValue() int {
	return b.nodes[0].(music.Audible).Note().MetaCommands[bankCmdIndex].Value().RandomAmount()
}

func (b BankCmd) Up() {
	b.Set(b.Value() + 1)
}

func (b BankCmd) Down() {
	b.Set(b.Value() - 1)
}

func (b BankCmd) Left() {
	b.SetAlt(b.AltValue() - 1)
}

func (b BankCmd) Right() {
	b.SetAlt(b.AltValue() + 1)
}

func (b BankCmd) AltUp() {}

func (b BankCmd) AltDown() {}

func (b BankCmd) AltLeft() {
	active := b.nodes[0].(music.Audible).Note().MetaCommands[bankCmdIndex].Active()
	for _, n := range b.nodes {
		n.(music.Audible).Note().MetaCommands[bankCmdIndex].SetActive(!active)
	}
}

func (b BankCmd) AltRight() {
	active := b.nodes[0].(music.Audible).Note().MetaCommands[bankCmdIndex].Active()
	for _, n := range b.nodes {
		n.(music.Audible).Note().MetaCommands[bankCmdIndex].SetActive(!active)
	}
}

func (b BankCmd) Set(value int) {
	for _, n := range b.nodes {
		n.(music.Audible).Note().MetaCommands[bankCmdIndex].Value().Set(value)
	}
}

func (b BankCmd) SetAlt(value int) {
	for _, n := range b.nodes {
		n.(music.Audible).Note().MetaCommands[bankCmdIndex].Value().SetRandomAmount(value)
	}
}

func (b BankCmd) SetEditValue(input string) {}
