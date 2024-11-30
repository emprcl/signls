package param

import (
	"fmt"
	"signls/core/common"
	"signls/core/music"
	"signls/ui/util"
)

const (
	scaleCmdIndex = 3
)

type ScaleCmd struct {
	nodes []common.Node
}

func (s ScaleCmd) Name() string {
	return "scale"
}

func (s ScaleCmd) Display() string {
	if !s.nodes[0].(music.Audible).Note().MetaCommands[scaleCmdIndex].Active() {
		return "⨯"
	}
	if s.nodes[0].(music.Audible).Note().MetaCommands[scaleCmdIndex].Value().RandomAmount() != 0 {
		return util.Normalize(
			fmt.Sprintf(
				"%s%+d\u033c",
				s.nodes[0].(music.Audible).Note().MetaCommands[scaleCmdIndex].Display(),
				s.nodes[0].(music.Audible).Note().MetaCommands[scaleCmdIndex].Value().RandomAmount(),
			),
		)
	}
	return s.nodes[0].(music.Audible).Note().MetaCommands[scaleCmdIndex].Display()
}

func (s ScaleCmd) Value() int {
	return s.nodes[0].(music.Audible).Note().MetaCommands[scaleCmdIndex].Value().Value()
}

func (s ScaleCmd) AltValue() int {
	return s.nodes[0].(music.Audible).Note().MetaCommands[scaleCmdIndex].Value().RandomAmount()
}

func (s ScaleCmd) Up() {
	s.Set(s.Value() + 1)
}

func (s ScaleCmd) Down() {
	s.Set(s.Value() - 1)
}

func (s ScaleCmd) Left() {
	s.SetAlt(s.AltValue() - 1)
}

func (s ScaleCmd) Right() {
	s.SetAlt(s.AltValue() + 1)
}

func (s ScaleCmd) AltUp() {}

func (s ScaleCmd) AltDown() {}

func (s ScaleCmd) AltLeft() {
	active := s.nodes[0].(music.Audible).Note().MetaCommands[scaleCmdIndex].Active()
	for _, n := range s.nodes {
		n.(music.Audible).Note().MetaCommands[scaleCmdIndex].SetActive(!active)
	}
}

func (s ScaleCmd) AltRight() {
	active := s.nodes[0].(music.Audible).Note().MetaCommands[scaleCmdIndex].Active()
	for _, n := range s.nodes {
		n.(music.Audible).Note().MetaCommands[scaleCmdIndex].SetActive(!active)
	}
}

func (s ScaleCmd) Set(value int) {
	for _, n := range s.nodes {
		n.(music.Audible).Note().MetaCommands[scaleCmdIndex].Value().Set(value)
	}
}

func (s ScaleCmd) SetAlt(value int) {
	for _, n := range s.nodes {
		n.(music.Audible).Note().MetaCommands[scaleCmdIndex].Value().SetRandomAmount(value)
	}
}

func (s ScaleCmd) SetEditValue(input string) {}
